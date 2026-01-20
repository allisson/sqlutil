# sqlutil

[![Build Status](https://github.com/allisson/sqlutil/workflows/Release/badge.svg)](https://github.com/allisson/sqlutil/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/allisson/sqlutil)](https://goreportcard.com/report/github.com/allisson/sqlutil)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/allisson/sqlutil)

A powerful and flexible collection of Go helpers for database operations. sqlutil simplifies common database tasks with a clean, chainable API that works with PostgreSQL, MySQL, and SQLite.

## Features

- Clean, intuitive API for CRUD operations
- Support for multiple SQL flavors (PostgreSQL, MySQL, SQLite)
- Flexible filtering with comparison operators (=, <>, >, <, >=, <=, LIKE, IN, NOT IN)
- Pagination and ordering support
- Field selection to retrieve only needed columns
- Row locking support (FOR UPDATE)
- Works with `*sql.DB`, `*sql.Conn`, and `*sql.Tx`
- Built on top of [sqlquery](https://github.com/allisson/sqlquery) and [scany](https://github.com/georgysavva/scany)

## Installation

```bash
go get github.com/allisson/sqlutil
```

## Table of Contents

- [Quick Start](#quick-start)
- [Database Setup](#database-setup)
- [Basic CRUD Operations](#basic-crud-operations)
  - [Insert](#insert)
  - [Get (Single Record)](#get-single-record)
  - [Select (Multiple Records)](#select-multiple-records)
  - [Update](#update)
  - [Delete](#delete)
- [Advanced Queries](#advanced-queries)
  - [Filtering](#filtering)
  - [Field Selection](#field-selection)
  - [Pagination](#pagination)
  - [Ordering](#ordering)
  - [Row Locking](#row-locking)
- [Bulk Operations](#bulk-operations)
- [Using with Transactions](#using-with-transactions)
- [Multiple SQL Flavors](#multiple-sql-flavors)

## Quick Start

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/allisson/sqlutil"
	_ "github.com/lib/pq"
)

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name" fieldtag:"insert,update"`
	Email string `db:"email" fieldtag:"insert,update"`
	Age   int    `db:"age" fieldtag:"insert,update"`
}

func main() {
	// Connect to database
	db, err := sql.Open("postgres", "postgres://user:password@localhost/mydb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	flavor := sqlutil.PostgreSQLFlavor

	// Insert a user
	user := User{Name: "Alice", Email: "alice@example.com", Age: 30}
	if err := sqlutil.Insert(ctx, db, flavor, "insert", "users", &user); err != nil {
		log.Fatal(err)
	}

	// Get a user
	var alice User
	opts := sqlutil.NewFindOptions(flavor).WithFilter("email", "alice@example.com")
	if err := sqlutil.Get(ctx, db, "users", opts, &alice); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found user: %+v\n", alice)

	// Update the user
	alice.Age = 31
	if err := sqlutil.Update(ctx, db, flavor, "update", "users", alice.ID, &alice); err != nil {
		log.Fatal(err)
	}

	// Delete the user
	if err := sqlutil.Delete(ctx, db, flavor, "users", alice.ID); err != nil {
		log.Fatal(err)
	}
}
```

## Database Setup

For the examples below, we'll use this table structure:

```sql
CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	age INTEGER NOT NULL,
	country VARCHAR(100),
	active BOOLEAN DEFAULT true,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

Run a PostgreSQL database with Docker:

```bash
docker run --name sqlutil-postgres \
	-e POSTGRES_USER=user \
	-e POSTGRES_PASSWORD=password \
	-e POSTGRES_DB=sqlutil \
	-p 5432:5432 \
	-d postgres:14-alpine
```

## Basic CRUD Operations

### Insert

Insert a single record into the database:

```go
package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/allisson/sqlutil"
	_ "github.com/lib/pq"
)

type User struct {
	ID        int    `db:"id"`
	Name      string `db:"name" fieldtag:"insert,update"`
	Email     string `db:"email" fieldtag:"insert,update"`
	Age       int    `db:"age" fieldtag:"insert,update"`
	Country   string `db:"country" fieldtag:"insert,update"`
	Active    bool   `db:"active" fieldtag:"insert,update"`
}

func main() {
	db, _ := sql.Open("postgres", "postgres://user:password@localhost/sqlutil?sslmode=disable")
	defer db.Close()

	ctx := context.Background()
	flavor := sqlutil.PostgreSQLFlavor

	// Insert a new user
	newUser := User{
		Name:    "Bob Smith",
		Email:   "bob@example.com",
		Age:     28,
		Country: "USA",
		Active:  true,
	}

	// The "insert" tag tells sqlutil to use fields marked with fieldtag:"insert"
	if err := sqlutil.Insert(ctx, db, flavor, "insert", "users", &newUser); err != nil {
		log.Fatal(err)
	}

	log.Println("User inserted successfully")
}
```

### Get (Single Record)

Retrieve a single record from the database:

```go
// Get user by ID
var user User
opts := sqlutil.NewFindOptions(flavor).WithFilter("id", 1)
if err := sqlutil.Get(ctx, db, "users", opts, &user); err != nil {
	log.Fatal(err)
}
fmt.Printf("User: %+v\n", user)

// Get user by email
var userByEmail User
opts = sqlutil.NewFindOptions(flavor).WithFilter("email", "bob@example.com")
if err := sqlutil.Get(ctx, db, "users", opts, &userByEmail); err != nil {
	log.Fatal(err)
}

// Get user with specific fields only
var partialUser struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}
opts = sqlutil.NewFindOptions(flavor).
	WithFields([]string{"id", "name", "email"}).
	WithFilter("id", 1)
if err := sqlutil.Get(ctx, db, "users", opts, &partialUser); err != nil {
	log.Fatal(err)
}
```

### Select (Multiple Records)

Retrieve multiple records from the database:

```go
// Get all users
var users []User
opts := sqlutil.NewFindAllOptions(flavor)
if err := sqlutil.Select(ctx, db, "users", opts, &users); err != nil {
	log.Fatal(err)
}
fmt.Printf("Found %d users\n", len(users))

// Get active users only
var activeUsers []User
opts = sqlutil.NewFindAllOptions(flavor).WithFilter("active", true)
if err := sqlutil.Select(ctx, db, "users", opts, &activeUsers); err != nil {
	log.Fatal(err)
}

// Get users with pagination
var pagedUsers []User
opts = sqlutil.NewFindAllOptions(flavor).
	WithLimit(10).
	WithOffset(0).
	WithOrderBy("created_at DESC")
if err := sqlutil.Select(ctx, db, "users", opts, &pagedUsers); err != nil {
	log.Fatal(err)
}
```

### Update

Update existing records:

```go
// Update user by ID
user := User{
	ID:      1,
	Name:    "Bob Smith Jr.",
	Email:   "bob.jr@example.com",
	Age:     29,
	Country: "USA",
	Active:  true,
}

// The "update" tag tells sqlutil to use fields marked with fieldtag:"update"
if err := sqlutil.Update(ctx, db, flavor, "update", "users", user.ID, &user); err != nil {
	log.Fatal(err)
}

// Update with custom options (more flexible)
updateOpts := sqlutil.NewUpdateOptions(flavor).
	WithSet("age", 30).
	WithSet("country", "Canada").
	WithFilter("id", 1)

if err := sqlutil.UpdateWithOptions(ctx, db, flavor, "users", updateOpts); err != nil {
	log.Fatal(err)
}

// Bulk update - update all inactive users
updateOpts = sqlutil.NewUpdateOptions(flavor).
	WithSet("active", false).
	WithFilter("age.lt", 18)

if err := sqlutil.UpdateWithOptions(ctx, db, flavor, "users", updateOpts); err != nil {
	log.Fatal(err)
}
```

### Delete

Delete records from the database:

```go
// Delete user by ID
if err := sqlutil.Delete(ctx, db, flavor, "users", 1); err != nil {
	log.Fatal(err)
}

// Delete with custom options (more flexible)
deleteOpts := sqlutil.NewDeleteOptions(flavor).WithFilter("active", false)
if err := sqlutil.DeleteWithOptions(ctx, db, flavor, "users", deleteOpts); err != nil {
	log.Fatal(err)
}

// Delete users older than 65
deleteOpts = sqlutil.NewDeleteOptions(flavor).WithFilter("age.gt", 65)
if err := sqlutil.DeleteWithOptions(ctx, db, flavor, "users", deleteOpts); err != nil {
	log.Fatal(err)
}

// Delete users from specific country
deleteOpts = sqlutil.NewDeleteOptions(flavor).WithFilter("country", "USA")
if err := sqlutil.DeleteWithOptions(ctx, db, flavor, "users", deleteOpts); err != nil {
	log.Fatal(err)
}
```

## Advanced Queries

### Filtering

sqlutil supports a wide range of filter operators:

```go
ctx := context.Background()
flavor := sqlutil.PostgreSQLFlavor

// Equality
opts := sqlutil.NewFindAllOptions(flavor).WithFilter("age", 30)
// WHERE age = 30

// Not equal
opts = sqlutil.NewFindAllOptions(flavor).WithFilter("age.not", 30)
// WHERE age <> 30

// Greater than
opts = sqlutil.NewFindAllOptions(flavor).WithFilter("age.gt", 18)
// WHERE age > 18

// Greater than or equal
opts = sqlutil.NewFindAllOptions(flavor).WithFilter("age.gte", 18)
// WHERE age >= 18

// Less than
opts = sqlutil.NewFindAllOptions(flavor).WithFilter("age.lt", 65)
// WHERE age < 65

// Less than or equal
opts = sqlutil.NewFindAllOptions(flavor).WithFilter("age.lte", 65)
// WHERE age <= 65

// LIKE operator
opts = sqlutil.NewFindAllOptions(flavor).WithFilter("name.like", "%Smith%")
// WHERE name LIKE '%Smith%'

// IN operator
opts = sqlutil.NewFindAllOptions(flavor).WithFilter("country.in", "USA,Canada,Mexico")
// WHERE country IN ('USA', 'Canada', 'Mexico')

// NOT IN operator
opts = sqlutil.NewFindAllOptions(flavor).WithFilter("country.notin", "USA,Canada")
// WHERE country NOT IN ('USA', 'Canada')

// IS NULL
opts = sqlutil.NewFindAllOptions(flavor).WithFilter("country.null", true)
// WHERE country IS NULL

// IS NOT NULL
opts = sqlutil.NewFindAllOptions(flavor).WithFilter("country.null", false)
// WHERE country IS NOT NULL

// Multiple filters (AND condition)
opts = sqlutil.NewFindAllOptions(flavor).
	WithFilter("age.gte", 18).
	WithFilter("age.lte", 65).
	WithFilter("active", true).
	WithFilter("country", "USA")
// WHERE age >= 18 AND age <= 65 AND active = true AND country = 'USA'

// Complex example: Find active users aged 25-40 from specific countries
var users []User
opts = sqlutil.NewFindAllOptions(flavor).
	WithFilter("active", true).
	WithFilter("age.gte", 25).
	WithFilter("age.lte", 40).
	WithFilter("country.in", "USA,Canada,UK")

if err := sqlutil.Select(ctx, db, "users", opts, &users); err != nil {
	log.Fatal(err)
}
```

### Field Selection

Retrieve only the fields you need:

```go
// Define a struct with only the fields you need
type UserBasic struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

// Select only specific fields
var users []UserBasic
opts := sqlutil.NewFindAllOptions(flavor).
	WithFields([]string{"id", "name", "email"}).
	WithFilter("active", true)

if err := sqlutil.Select(ctx, db, "users", opts, &users); err != nil {
	log.Fatal(err)
}
// SELECT id, name, email FROM users WHERE active = true

// Get single user with limited fields
var userBasic UserBasic
opts := sqlutil.NewFindOptions(flavor).
	WithFields([]string{"id", "name", "email"}).
	WithFilter("id", 1)

if err := sqlutil.Get(ctx, db, "users", opts, &userBasic); err != nil {
	log.Fatal(err)
}
```

### Pagination

Implement pagination for large result sets:

```go
// Page 1: First 10 users
var page1 []User
opts := sqlutil.NewFindAllOptions(flavor).
	WithLimit(10).
	WithOffset(0).
	WithOrderBy("id ASC")

if err := sqlutil.Select(ctx, db, "users", opts, &page1); err != nil {
	log.Fatal(err)
}

// Page 2: Next 10 users
var page2 []User
opts = sqlutil.NewFindAllOptions(flavor).
	WithLimit(10).
	WithOffset(10).
	WithOrderBy("id ASC")

if err := sqlutil.Select(ctx, db, "users", opts, &page2); err != nil {
	log.Fatal(err)
}

// Pagination helper function
func GetUserPage(ctx context.Context, db *sql.DB, page, pageSize int) ([]User, error) {
	var users []User
	offset := (page - 1) * pageSize
	
	opts := sqlutil.NewFindAllOptions(flavor).
		WithLimit(pageSize).
		WithOffset(offset).
		WithOrderBy("id ASC")
	
	err := sqlutil.Select(ctx, db, "users", opts, &users)
	return users, err
}

// Usage
page1Users, err := GetUserPage(ctx, db, 1, 20) // First page, 20 items
page2Users, err := GetUserPage(ctx, db, 2, 20) // Second page, 20 items
```

### Ordering

Sort results by one or multiple columns:

```go
// Order by single column ascending
opts := sqlutil.NewFindAllOptions(flavor).WithOrderBy("name ASC")

// Order by single column descending
opts = sqlutil.NewFindAllOptions(flavor).WithOrderBy("created_at DESC")

// Order by multiple columns
opts = sqlutil.NewFindAllOptions(flavor).WithOrderBy("country ASC, age DESC")

// Complex example: Get top 10 youngest active users from USA
var youngUsers []User
opts = sqlutil.NewFindAllOptions(flavor).
	WithFilter("active", true).
	WithFilter("country", "USA").
	WithOrderBy("age ASC").
	WithLimit(10)

if err := sqlutil.Select(ctx, db, "users", opts, &youngUsers); err != nil {
	log.Fatal(err)
}
```

### Row Locking

Use row locking for concurrent operations:

```go
// FOR UPDATE - locks rows for update
opts := sqlutil.NewFindAllOptions(flavor).
	WithFilter("active", true).
	WithForUpdate("")

// FOR UPDATE SKIP LOCKED - skip rows that are already locked
opts = sqlutil.NewFindAllOptions(flavor).
	WithFilter("active", true).
	WithForUpdate("SKIP LOCKED")

// FOR UPDATE NOWAIT - return immediately if rows are locked
opts = sqlutil.NewFindAllOptions(flavor).
	WithFilter("active", true).
	WithForUpdate("NOWAIT")

// Example: Process queue items with row locking
tx, err := db.BeginTx(ctx, nil)
if err != nil {
	log.Fatal(err)
}
defer tx.Rollback()

var queueItems []QueueItem
opts := sqlutil.NewFindAllOptions(flavor).
	WithFilter("status", "pending").
	WithOrderBy("created_at ASC").
	WithLimit(10).
	WithForUpdate("SKIP LOCKED")

if err := sqlutil.Select(ctx, tx, "queue", opts, &queueItems); err != nil {
	log.Fatal(err)
}

// Process items...
for _, item := range queueItems {
	// Process item
	updateOpts := sqlutil.NewUpdateOptions(flavor).
		WithSet("status", "processed").
		WithFilter("id", item.ID)
	
	if err := sqlutil.UpdateWithOptions(ctx, tx, flavor, "queue", updateOpts); err != nil {
		log.Fatal(err)
	}
}

tx.Commit()
```

## Bulk Operations

Efficiently handle multiple operations:

```go
// Insert multiple users
users := []User{
	{Name: "Alice", Email: "alice@example.com", Age: 30, Country: "USA"},
	{Name: "Bob", Email: "bob@example.com", Age: 25, Country: "Canada"},
	{Name: "Charlie", Email: "charlie@example.com", Age: 35, Country: "UK"},
}

for _, user := range users {
	if err := sqlutil.Insert(ctx, db, flavor, "insert", "users", &user); err != nil {
		log.Fatal(err)
	}
}

// Bulk update with filter
updateOpts := sqlutil.NewUpdateOptions(flavor).
	WithSet("active", false).
	WithFilter("country.in", "USA,Canada,Mexico")

if err := sqlutil.UpdateWithOptions(ctx, db, flavor, "users", updateOpts); err != nil {
	log.Fatal(err)
}

// Bulk delete with filter
deleteOpts := sqlutil.NewDeleteOptions(flavor).
	WithFilter("active", false).
	WithFilter("age.lt", 18)

if err := sqlutil.DeleteWithOptions(ctx, db, flavor, "users", deleteOpts); err != nil {
	log.Fatal(err)
}
```

## Using with Transactions

sqlutil works seamlessly with database transactions:

```go
// Start a transaction
tx, err := db.BeginTx(ctx, nil)
if err != nil {
	log.Fatal(err)
}
defer tx.Rollback() // Rollback if not committed

// Insert user
user := User{Name: "David", Email: "david@example.com", Age: 40, Country: "France"}
if err := sqlutil.Insert(ctx, tx, flavor, "insert", "users", &user); err != nil {
	log.Fatal(err)
}

// Get the inserted user
var insertedUser User
opts := sqlutil.NewFindOptions(flavor).WithFilter("email", "david@example.com")
if err := sqlutil.Get(ctx, tx, "users", opts, &insertedUser); err != nil {
	log.Fatal(err)
}

// Update the user
insertedUser.Age = 41
if err := sqlutil.Update(ctx, tx, flavor, "update", "users", insertedUser.ID, &insertedUser); err != nil {
	log.Fatal(err)
}

// Commit the transaction
if err := tx.Commit(); err != nil {
	log.Fatal(err)
}

// Example: Transfer operation with rollback
func TransferUser(ctx context.Context, db *sql.DB, userID int, newCountry string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get user
	var user User
	opts := sqlutil.NewFindOptions(flavor).WithFilter("id", userID)
	if err := sqlutil.Get(ctx, tx, "users", opts, &user); err != nil {
		return err
	}

	// Update country
	updateOpts := sqlutil.NewUpdateOptions(flavor).
		WithSet("country", newCountry).
		WithFilter("id", userID)
	
	if err := sqlutil.UpdateWithOptions(ctx, tx, flavor, "users", updateOpts); err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}
```

## Multiple SQL Flavors

sqlutil supports PostgreSQL, MySQL, and SQLite:

### PostgreSQL

```go
import (
	"github.com/allisson/sqlutil"
	_ "github.com/lib/pq"
)

db, _ := sql.Open("postgres", "postgres://user:password@localhost/mydb?sslmode=disable")
flavor := sqlutil.PostgreSQLFlavor

opts := sqlutil.NewFindOptions(flavor).WithFilter("id", 1)
```

### MySQL

```go
import (
	"github.com/allisson/sqlutil"
	_ "github.com/go-sql-driver/mysql"
)

db, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/mydb")
flavor := sqlutil.MySQLFlavor

opts := sqlutil.NewFindOptions(flavor).WithFilter("id", 1)
```

### SQLite

```go
import (
	"github.com/allisson/sqlutil"
	_ "github.com/mattn/go-sqlite3"
)

db, _ := sql.Open("sqlite3", "./mydb.db")
flavor := sqlutil.SQLiteFlavor

opts := sqlutil.NewFindOptions(flavor).WithFilter("id", 1)
```

### Cross-Database Compatibility

Write database-agnostic code by passing the flavor as a parameter:

```go
type UserRepository struct {
	db     *sql.DB
	flavor sqlutil.Flavor
}

func NewUserRepository(db *sql.DB, flavor sqlutil.Flavor) *UserRepository {
	return &UserRepository{db: db, flavor: flavor}
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	var user User
	opts := sqlutil.NewFindOptions(r.flavor).WithFilter("id", id)
	err := sqlutil.Get(ctx, r.db, "users", opts, &user)
	return &user, err
}

func (r *UserRepository) GetActiveUsers(ctx context.Context) ([]User, error) {
	var users []User
	opts := sqlutil.NewFindAllOptions(r.flavor).
		WithFilter("active", true).
		WithOrderBy("name ASC")
	err := sqlutil.Select(ctx, r.db, "users", opts, &users)
	return users, err
}

// Works with any database flavor
pgRepo := NewUserRepository(pgDB, sqlutil.PostgreSQLFlavor)
mysqlRepo := NewUserRepository(mysqlDB, sqlutil.MySQLFlavor)
sqliteRepo := NewUserRepository(sqliteDB, sqlutil.SQLiteFlavor)
```

## Complete Example

Here's a complete example demonstrating various features:

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/allisson/sqlutil"
	_ "github.com/lib/pq"
)

type User struct {
	ID        int    `db:"id"`
	Name      string `db:"name" fieldtag:"insert,update"`
	Email     string `db:"email" fieldtag:"insert,update"`
	Age       int    `db:"age" fieldtag:"insert,update"`
	Country   string `db:"country" fieldtag:"insert,update"`
	Active    bool   `db:"active" fieldtag:"insert,update"`
}

func main() {
	// Connect to database
	db, err := sql.Open("postgres", "postgres://user:password@localhost/sqlutil?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	flavor := sqlutil.PostgreSQLFlavor

	// Create table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			age INTEGER NOT NULL,
			country VARCHAR(100),
			active BOOLEAN DEFAULT true
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert users
	users := []User{
		{Name: "Alice Johnson", Email: "alice@example.com", Age: 30, Country: "USA", Active: true},
		{Name: "Bob Smith", Email: "bob@example.com", Age: 25, Country: "Canada", Active: true},
		{Name: "Charlie Brown", Email: "charlie@example.com", Age: 35, Country: "UK", Active: true},
		{Name: "David Lee", Email: "david@example.com", Age: 28, Country: "USA", Active: false},
	}

	for _, user := range users {
		if err := sqlutil.Insert(ctx, db, flavor, "insert", "users", &user); err != nil {
			log.Printf("Error inserting user: %v", err)
		}
	}

	// Get single user by email
	var alice User
	opts := sqlutil.NewFindOptions(flavor).WithFilter("email", "alice@example.com")
	if err := sqlutil.Get(ctx, db, "users", opts, &alice); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found user: %+v\n", alice)

	// Get all active users from USA
	var activeUSAUsers []User
	opts2 := sqlutil.NewFindAllOptions(flavor).
		WithFilter("active", true).
		WithFilter("country", "USA").
		WithOrderBy("name ASC")
	
	if err := sqlutil.Select(ctx, db, "users", opts2, &activeUSAUsers); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Active USA users: %d\n", len(activeUSAUsers))

	// Get users with age between 25 and 30
	var youngUsers []User
	opts3 := sqlutil.NewFindAllOptions(flavor).
		WithFilter("age.gte", 25).
		WithFilter("age.lte", 30).
		WithOrderBy("age ASC")
	
	if err := sqlutil.Select(ctx, db, "users", opts3, &youngUsers); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Users aged 25-30: %d\n", len(youngUsers))

	// Update user
	alice.Age = 31
	alice.Country = "Canada"
	if err := sqlutil.Update(ctx, db, flavor, "update", "users", alice.ID, &alice); err != nil {
		log.Fatal(err)
	}

	// Bulk update: deactivate all users older than 30
	updateOpts := sqlutil.NewUpdateOptions(flavor).
		WithSet("active", false).
		WithFilter("age.gt", 30)
	
	if err := sqlutil.UpdateWithOptions(ctx, db, flavor, "users", updateOpts); err != nil {
		log.Fatal(err)
	}

	// Get users with pagination
	page1 := []User{}
	paginationOpts := sqlutil.NewFindAllOptions(flavor).
		WithLimit(2).
		WithOffset(0).
		WithOrderBy("name ASC")
	
	if err := sqlutil.Select(ctx, db, "users", paginationOpts, &page1); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Page 1: %d users\n", len(page1))

	// Delete inactive users
	deleteOpts := sqlutil.NewDeleteOptions(flavor).WithFilter("active", false)
	if err := sqlutil.DeleteWithOptions(ctx, db, flavor, "users", deleteOpts); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inactive users deleted")
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Related Projects

- [sqlquery](https://github.com/allisson/sqlquery) - SQL query builder for Go
- [scany](https://github.com/georgysavva/scany) - Library for scanning data from a database into Go structs
