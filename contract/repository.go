package contract

import (
	"context"
)

type (
	Repository interface {
		With(relations ...string) Repository
		Where(query any, args ...any) Repository
		Unscoped() Repository
		Limit(limit int) Repository
		Offset(offset int) Repository
		OrderBy(column, direction string) Repository

		Find(ctx context.Context, id any) (Model, error)
		FindOrFail(ctx context.Context, id any) (Model, error)
		First(ctx context.Context) (Model, error)
		FirstOrFail(ctx context.Context) (Model, error)
		Get(ctx context.Context) ([]Model, error)
		Pluck(ctx context.Context, column string, dest any) error

		Create(ctx context.Context, models ...Model) error
		CreateInBatches(ctx context.Context, models []Model, batchSize int) error
		Update(ctx context.Context, models ...Model) error
		Delete(ctx context.Context, models ...Model) error
		ForceDelete(ctx context.Context, models ...Model) error

		FirstOrCreate(ctx context.Context, condition Model, create ...Model) (Model, error)
		UpdateOrCreate(ctx context.Context, condition Model, values any) (Model, error)

		// QueryBuilder provides access to the fluent query builder interface
		QueryBuilder() QueryBuilder
	}
)
