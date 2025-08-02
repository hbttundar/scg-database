package contract

import (
	"context"
	"database/sql"
)

type (
	Connection interface {
		GetConnection() any
		Ping(ctx context.Context) error
		Close() error
		NewRepository(model Model) (Repository, error)
		Transaction(ctx context.Context, fn func(txConnection Connection) error) error
		Select(ctx context.Context, query string, bindings ...any) ([]map[string]any, error)
		Statement(ctx context.Context, query string, bindings ...any) (sql.Result, error)
	}
)
