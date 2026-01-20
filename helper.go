package sqlutil

import (
	"context"
	"database/sql"

	"github.com/allisson/sqlquery"
	"github.com/georgysavva/scany/v2/sqlscan"
)

type (
	Flavor         = sqlquery.Flavor
	FindOptions    = sqlquery.FindOptions
	FindAllOptions = sqlquery.FindAllOptions
	UpdateOptions  = sqlquery.UpdateOptions
	DeleteOptions  = sqlquery.DeleteOptions
)

var (
	MySQLFlavor      = sqlquery.MySQLFlavor
	PostgreSQLFlavor = sqlquery.PostgreSQLFlavor
	SQLiteFlavor     = sqlquery.SQLiteFlavor
)

// NewFindOptions creates and returns a new FindOptions instance for the specified SQL flavor.
// FindOptions is used to build queries for finding a single record with support for field selection,
// filtering, and various comparison operators.
func NewFindOptions(flavor Flavor) *FindOptions {
	return sqlquery.NewFindOptions(flavor)
}

// NewFindAllOptions creates and returns a new FindAllOptions instance for the specified SQL flavor.
// FindAllOptions is used to build queries for finding multiple records with support for field selection,
// filtering, pagination (limit/offset), ordering, and row locking (FOR UPDATE).
func NewFindAllOptions(flavor Flavor) *FindAllOptions {
	return sqlquery.NewFindAllOptions(flavor)
}

// NewUpdateOptions creates and returns a new UpdateOptions instance for the specified SQL flavor.
// UpdateOptions is used to build UPDATE queries with support for filtering and setting multiple fields.
func NewUpdateOptions(flavor Flavor) *UpdateOptions {
	return sqlquery.NewUpdateOptions(flavor)
}

// NewDeleteOptions creates and returns a new DeleteOptions instance for the specified SQL flavor.
// DeleteOptions is used to build DELETE queries with support for filtering conditions.
func NewDeleteOptions(flavor Flavor) *DeleteOptions {
	return sqlquery.NewDeleteOptions(flavor)
}

// Querier is an abstraction over *sql.DB, *sql.Conn, and *sql.Tx.
// It provides a common interface for executing queries and commands against a database,
// allowing functions to work with any of these types.
type Querier interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Get retrieves a single record from the database and scans it into dst.
// It builds a SELECT query using sqlquery.FindQuery with the provided options,
// then uses scany's sqlscan.Get to scan the result into the destination struct.
// Returns an error if no rows are found or if scanning fails.
func Get(ctx context.Context, db Querier, tableName string, options *FindOptions, dst interface{}) error {
	sqlQuery, args := sqlquery.FindQuery(tableName, options)
	return sqlscan.Get(ctx, db, dst, sqlQuery, args...)
}

// Select retrieves multiple records from the database and scans them into dst.
// It builds a SELECT query using sqlquery.FindAllQuery with the provided options,
// then uses scany's sqlscan.Select to scan the results into the destination slice.
// The dst parameter should be a pointer to a slice of structs.
func Select(ctx context.Context, db Querier, tableName string, options *FindAllOptions, dst interface{}) error {
	sqlQuery, args := sqlquery.FindAllQuery(tableName, options)
	return sqlscan.Select(ctx, db, dst, sqlQuery, args...)
}

// Insert inserts a new record into the database table.
// It builds an INSERT query using sqlquery.InsertQuery based on the struct fields
// that match the specified tag (e.g., "insert"), then executes the query.
// The structValue parameter should be a pointer to a struct with appropriate field tags.
func Insert(ctx context.Context, db Querier, flavor Flavor, tag, tableName string, structValue interface{}) error {
	sqlQuery, args := sqlquery.InsertQuery(flavor, tag, tableName, structValue)
	_, err := db.ExecContext(ctx, sqlQuery, args...)
	return err
}

// Update updates an existing record in the database table by its ID.
// It builds an UPDATE query using sqlquery.UpdateQuery based on the struct fields
// that match the specified tag (e.g., "update"), adding a WHERE clause for the given ID.
// The structValue parameter should be a pointer to a struct with appropriate field tags.
func Update(ctx context.Context, db Querier, flavor Flavor, tag, tableName string, id interface{}, structValue interface{}) error {
	sqlQuery, args := sqlquery.UpdateQuery(flavor, tag, tableName, id, structValue)
	_, err := db.ExecContext(ctx, sqlQuery, args...)
	return err
}

// Delete deletes a record from the database table by its ID.
// It builds a DELETE query using sqlquery.DeleteQuery with a WHERE clause
// matching the given ID, then executes the query.
func Delete(ctx context.Context, db Querier, flavor Flavor, tableName string, id interface{}) error {
	sqlQuery, args := sqlquery.DeleteQuery(flavor, tableName, id)
	_, err := db.ExecContext(ctx, sqlQuery, args...)
	return err
}

// UpdateWithOptions updates records in the database table based on custom criteria.
// It builds an UPDATE query using sqlquery.UpdateWithOptionsQuery with the provided options,
// which allow for complex WHERE clauses and multiple field updates.
// This provides more flexibility than the simple Update function.
func UpdateWithOptions(ctx context.Context, db Querier, flavor Flavor, tableName string, options *UpdateOptions) error {
	sqlQuery, args := sqlquery.UpdateWithOptionsQuery(tableName, options)
	_, err := db.ExecContext(ctx, sqlQuery, args...)
	return err
}

// DeleteWithOptions deletes records from the database table based on custom criteria.
// It builds a DELETE query using sqlquery.DeleteWithOptionsQuery with the provided options,
// which allow for complex WHERE clauses to match multiple records.
// This provides more flexibility than the simple Delete function.
func DeleteWithOptions(ctx context.Context, db Querier, flavor Flavor, tableName string, options *DeleteOptions) error {
	sqlQuery, args := sqlquery.DeleteWithOptionsQuery(tableName, options)
	_, err := db.ExecContext(ctx, sqlQuery, args...)
	return err
}
