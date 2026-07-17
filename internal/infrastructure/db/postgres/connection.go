package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	errMsgCreatingConnectionPool = "unable to create connection pool"
	errMsgPingDatabase           = "unable to ping PostgreSQL database"
)

// Open creates and verifies a PostgreSQL connection pool. The caller owns the
// returned pool and must close it.
func Open(ctx context.Context, postgresDBConfig *entity.PostgresDBConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%s",
		postgresDBConfig.Username, postgresDBConfig.Password,
		postgresDBConfig.Host, postgresDBConfig.Port,
		postgresDBConfig.Name, postgresDBConfig.SslMode,
		postgresDBConfig.PoolMaxConns,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsgCreatingConnectionPool, err)
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsgCreatingConnectionPool, err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := dbpool.Ping(pingCtx); err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("%s: %w", errMsgPingDatabase, err)
	}

	return dbpool, nil
}
