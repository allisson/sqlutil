package sqlutil

import (
	"context"
	"database/sql"

	"github.com/georgysavva/scany/sqlscan"
)

// Querier is a abstraction over *sql.DB/*sql.Conn/*sql.Tx.
type Querier interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Get is a high-level function that calls FindQuery and scany sqlscan.Get function.
func Get(ctx context.Context, db Querier, tableName string, options *FindOptions, dst interface{}) error {
	sqlQuery, args := FindQuery(tableName, options)
	return sqlscan.Get(ctx, db, dst, sqlQuery, args...)
}

// Select is a high-level function that calls FindAllQuery and scany sqlscan.Select function.
func Select(ctx context.Context, db Querier, tableName string, options *FindAllOptions, dst interface{}) error {
	sqlQuery, args := FindAllQuery(tableName, options)
	return sqlscan.Select(ctx, db, dst, sqlQuery, args...)
}

// Insert is a high-level function that calls InsertQuery and db.ExecContext.
func Insert(ctx context.Context, db Querier, flavor Flavor, tag, tableName string, structValue interface{}) error {
	sqlQuery, args := InsertQuery(flavor, tag, tableName, structValue)
	return sqlscan.Get(ctx, db, structValue, sqlQuery, args...)
}

// Update is a high-level function that calls UpdateQuery and db.ExecContext.
func Update(ctx context.Context, db Querier, flavor Flavor, tag, tableName string, id interface{}, structValue interface{}) error {
	sqlQuery, args := UpdateQuery(flavor, tag, tableName, id, structValue)
	_, err := db.ExecContext(ctx, sqlQuery, args...)
	return err
}

// Delete is a high-level function that calls DeleteQuery and db.ExecContext.
func Delete(ctx context.Context, db Querier, flavor Flavor, tableName string, id interface{}) error {
	sqlQuery, args := DeleteQuery(flavor, tableName, id)
	_, err := db.ExecContext(ctx, sqlQuery, args...)
	return err
}
