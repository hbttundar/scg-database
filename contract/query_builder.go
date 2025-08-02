package contract

import (
	"context"
)

type (
	// QueryBuilder provides a fluent interface for building database queries
	// Each adapter should implement this interface to provide database-specific query building
	QueryBuilder interface {
		// Query building methods
		Select(columns ...string) QueryBuilder
		Where(condition string, args ...any) QueryBuilder
		WhereIn(column string, values []any) QueryBuilder
		WhereNotIn(column string, values []any) QueryBuilder
		WhereNull(column string) QueryBuilder
		WhereNotNull(column string) QueryBuilder
		WhereBetween(column string, start, end any) QueryBuilder
		OrWhere(condition string, args ...any) QueryBuilder

		// Join methods
		Join(table, condition string) QueryBuilder
		LeftJoin(table, condition string) QueryBuilder
		RightJoin(table, condition string) QueryBuilder
		InnerJoin(table, condition string) QueryBuilder

		// Ordering and grouping
		OrderBy(column, direction string) QueryBuilder
		GroupBy(columns ...string) QueryBuilder
		Having(condition string, args ...any) QueryBuilder

		// Limiting and pagination
		Limit(limit int) QueryBuilder
		Offset(offset int) QueryBuilder

		// Relationships
		With(relations ...string) QueryBuilder
		WithCount(relations ...string) QueryBuilder

		// Scopes and advanced features
		Scoped() QueryBuilder
		Unscoped() QueryBuilder

		// Execution methods
		Find(ctx context.Context, dest any) error
		First(ctx context.Context, dest any) error
		Get(ctx context.Context, dest any) error
		Count(ctx context.Context) (int64, error)
		Exists(ctx context.Context) (bool, error)

		// Mutation methods
		Create(ctx context.Context, value any) error
		Update(ctx context.Context, values any) error
		Delete(ctx context.Context) error

		// Raw query methods
		Raw(sql string, args ...any) QueryBuilder
		Exec(ctx context.Context, sql string, args ...any) error

		// Utility methods
		ToSQL() (string, []any, error)
		Clone() QueryBuilder
		Reset() QueryBuilder
	}

	// QueryBuilderFactory creates QueryBuilder instances for specific models and connections
	QueryBuilderFactory interface {
		NewQueryBuilder(model Model, connection any) QueryBuilder
		Name() string
	}

	// QueryBuilderRegistry manages query builder factories for different adapters
	QueryBuilderRegistry interface {
		Register(adapterName string, factory QueryBuilderFactory)
		Get(adapterName string) (QueryBuilderFactory, error)
		List() []string
	}
)
