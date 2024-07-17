package repository

type UserRepository interface {
	SetUserUUID(tokenUUID string, userUUID string, expiresIn int64) error
	GetUserUUID(tokenUUID string) (string, error)
	DelUserUUID(tokenUUID string, accessTokenUUID string) (int64, error)
}
