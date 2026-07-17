package redis

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	redislib "github.com/redis/go-redis/v9"
)

const (
	sessionTimeout = 5 * time.Second
	resultOK       = int64(0)
	resultMissing  = int64(1)
	resultReused   = int64(2)
)

var (
	createSessionScript = redislib.NewScript(`
if redis.call('EXISTS', KEYS[1], KEYS[2], KEYS[3]) > 0 then return 4 end
redis.call('SET', KEYS[3], 'active', 'EX', ARGV[3])
redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[2])
redis.call('HSET', KEYS[2], 'user_uuid', ARGV[1], 'state', 'active')
redis.call('EXPIRE', KEYS[2], ARGV[3])
redis.call('SADD', KEYS[4], KEYS[1], KEYS[2])
redis.call('EXPIRE', KEYS[4], ARGV[3])
return 0
`)
	rotateSessionScript = redislib.NewScript(`
local function revoke_family()
  redis.call('SET', KEYS[4], 'revoked', 'KEEPTTL')
  local members = redis.call('SMEMBERS', KEYS[5])
  for _, key in ipairs(members) do
    if string.find(key, ':refresh:', 1, true) then
      if redis.call('EXISTS', key) == 1 then redis.call('HSET', key, 'state', 'revoked') end
    else
      redis.call('DEL', key)
    end
  end
end
if redis.call('EXISTS', KEYS[1]) == 0 then return 1 end
local state = redis.call('HGET', KEYS[1], 'state')
local family_state = redis.call('GET', KEYS[4])
if state ~= 'active' or family_state ~= 'active' then
  revoke_family()
  return 2
end
if redis.call('HGET', KEYS[1], 'user_uuid') ~= ARGV[1] then return 3 end
if redis.call('EXISTS', KEYS[2], KEYS[3]) > 0 then return 4 end
redis.call('HSET', KEYS[1], 'state', 'rotated', 'replaced_by', KEYS[3])
redis.call('SET', KEYS[2], ARGV[1], 'EX', ARGV[2])
redis.call('HSET', KEYS[3], 'user_uuid', ARGV[1], 'state', 'active')
redis.call('EXPIRE', KEYS[3], ARGV[3])
redis.call('SADD', KEYS[5], KEYS[2], KEYS[3])
local family_ttl = redis.call('TTL', KEYS[4])
if family_ttl < tonumber(ARGV[3]) then
  redis.call('EXPIRE', KEYS[4], ARGV[3])
  redis.call('EXPIRE', KEYS[5], ARGV[3])
end
return 0
`)
	revokeSessionScript = redislib.NewScript(`
if redis.call('EXISTS', KEYS[1]) == 0 then return 1 end
redis.call('SET', KEYS[2], 'revoked', 'KEEPTTL')
local members = redis.call('SMEMBERS', KEYS[3])
for _, key in ipairs(members) do
  if string.find(key, ':refresh:', 1, true) then
    if redis.call('EXISTS', key) == 1 then redis.call('HSET', key, 'state', 'revoked') end
  else
    redis.call('DEL', key)
  end
end
return 0
`)
)

type TokenRepository struct{ client *redislib.Client }

func NewTokenRepository(client *redislib.Client) repository.TokenRepository {
	return &TokenRepository{client: client}
}

func (r TokenRepository) Create(ctx context.Context, session entity.TokenSession) error {
	accessTTL, refreshTTL, err := sessionTTLs(session)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, sessionTimeout)
	defer cancel()
	keys := sessionKeys(session.FamilyUUID, session.AccessTokenUUID, session.RefreshTokenUUID)
	result, err := createSessionScript.Run(ctx, r.client, keys, session.UserUUID, accessTTL, refreshTTL).Int64()
	if err != nil {
		return fmt.Errorf("create token session: %w", err)
	}
	if result != resultOK {
		return errors.New("create token session: session key already exists")
	}
	return nil
}

func (r TokenRepository) Rotate(ctx context.Context, currentRefreshTokenUUID string, replacement entity.TokenSession) error {
	accessTTL, refreshTTL, err := sessionTTLs(replacement)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, sessionTimeout)
	defer cancel()
	prefix := sessionPrefix(replacement.FamilyUUID)
	keys := []string{
		prefix + "refresh:" + currentRefreshTokenUUID,
		prefix + "access:" + replacement.AccessTokenUUID,
		prefix + "refresh:" + replacement.RefreshTokenUUID,
		prefix + "family",
		prefix + "members",
	}
	result, err := rotateSessionScript.Run(ctx, r.client, keys, replacement.UserUUID, accessTTL, refreshTTL).Int64()
	if err != nil {
		return fmt.Errorf("rotate refresh session: %w", err)
	}
	switch result {
	case resultOK:
		return nil
	case resultMissing:
		return apperror.ErrSessionNotFound
	case resultReused:
		return apperror.ErrRefreshTokenReused
	default:
		return errors.New("rotate refresh session: inconsistent session data")
	}
}

func (r TokenRepository) GetUserUUID(ctx context.Context, familyUUID, accessTokenUUID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	result, err := r.client.Get(ctx, sessionPrefix(familyUUID)+"access:"+accessTokenUUID).Result()
	if errors.Is(err, redislib.Nil) {
		return "", apperror.ErrSessionNotFound
	}
	if err != nil {
		return "", fmt.Errorf("get access session: %w", err)
	}
	return result, nil
}

func (r TokenRepository) Revoke(ctx context.Context, familyUUID, refreshTokenUUID string) error {
	ctx, cancel := context.WithTimeout(ctx, sessionTimeout)
	defer cancel()
	prefix := sessionPrefix(familyUUID)
	result, err := revokeSessionScript.Run(ctx, r.client, []string{
		prefix + "refresh:" + refreshTokenUUID,
		prefix + "family",
		prefix + "members",
	}).Int64()
	if err != nil {
		return fmt.Errorf("revoke token session: %w", err)
	}
	if result == resultMissing {
		return apperror.ErrSessionNotFound
	}
	return nil
}

func sessionTTLs(session entity.TokenSession) (int64, int64, error) {
	accessTTL := int64(math.Ceil(time.Until(session.AccessExpiresAt).Seconds()))
	refreshTTL := int64(math.Ceil(time.Until(session.RefreshExpiresAt).Seconds()))
	if accessTTL <= 0 || refreshTTL <= 0 {
		return 0, 0, errors.New("token session expiry must be in the future")
	}
	return accessTTL, refreshTTL, nil
}

func sessionKeys(familyUUID, accessTokenUUID, refreshTokenUUID string) []string {
	prefix := sessionPrefix(familyUUID)
	return []string{
		prefix + "access:" + accessTokenUUID,
		prefix + "refresh:" + refreshTokenUUID,
		prefix + "family",
		prefix + "members",
	}
}

func sessionPrefix(familyUUID string) string { return "session:{" + familyUUID + "}:" }
