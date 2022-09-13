package sqlutil

import (
	"context"
	"database/sql"

	"github.com/allisson/sqlquery"
	"github.com/georgysavva/scany/sqlscan"
)

type Flavor = sqlquery.Flavor
type FindOptions = sqlquery.FindOptions
type FindAllOptions = sqlquery.FindAllOptions

var (
	MySQLFlavor      = sqlquery.MySQLFlavor
	PostgreSQLFlavor = sqlquery.PostgreSQLFlavor
	SQLiteFlavor     = sqlquery.SQLiteFlavor
)

// NewFindOptions returns a FindOptions.
func NewFindOptions(flavor Flavor) *FindOptions {
	return sqlquery.NewFindOptions(flavor)
}

// NewFindAllOptions returns a FindAllOptions.
func NewFindAllOptions(flavor Flavor) *FindAllOptions {
	return sqlquery.NewFindAllOptions(flavor)
}

// Querier is a abstraction over *sql.DB/*sql.Conn/*sql.Tx.
type Querier interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Get is a high-level function that calls sqlquery.FindQuery and scany sqlscan.Get function.
func Get(ctx context.Context, db Querier, tableName string, options *FindOptions, dst interface{}) error {
	sqlQuery, args := sqlquery.FindQuery(tableName, options)
	return sqlscan.Get(ctx, db, dst, sqlQuery, args...)
}

// Select is a high-level function that calls sqlquery.FindAllQuery and scany sqlscan.Select function.
func Select(ctx context.Context, db Querier, tableName string, options *FindAllOptions, dst interface{}) error {
	sqlQuery, args := sqlquery.FindAllQuery(tableName, options)
	return sqlscan.Select(ctx, db, dst, sqlQuery, args...)
}

// Insert is a high-level function that calls sqlquery.InsertQuery and db.ExecContext.
func Insert(ctx context.Context, db Querier, flavor Flavor, tag, tableName string, structValue interface{}) error {
	sqlQuery, args := sqlquery.InsertQuery(flavor, tag, tableName, structValue)
	_, err := db.ExecContext(ctx, sqlQuery, args...)
	return err
}

// Update is a high-level function that calls sqlquery.pdateQuery and db.ExecContext.
func Update(ctx context.Context, db Querier, flavor Flavor, tag, tableName string, id interface{}, structValue interface{}) error {
	sqlQuery, args := sqlquery.UpdateQuery(flavor, tag, tableName, id, structValue)
	_, err := db.ExecContext(ctx, sqlQuery, args...)
	return err
}

// Delete is a high-level function that calls sqlquery.DeleteQuery and db.ExecContext.
func Delete(ctx context.Context, db Querier, flavor Flavor, tableName string, id interface{}) error {
	sqlQuery, args := sqlquery.DeleteQuery(flavor, tableName, id)
	_, err := db.ExecContext(ctx, sqlQuery, args...)
	return err
}
