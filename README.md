# sqlutil
[![Build Status](https://github.com/allisson/sqlutil/workflows/Release/badge.svg)](https://github.com/allisson/sqlutil/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/allisson/sqlutil)](https://goreportcard.com/report/github.com/allisson/sqlutil)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/allisson/sqlutil)

A collection of helpers to deal with database.

Example:

```golang
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/allisson/sqlutil"
	_ "github.com/lib/pq"
)

type Player struct {
	ID   int    `db:"id"`
	Name string `db:"name" fieldtag:"insert,update"`
	Age  int    `db:"age" fieldtag:"insert,update"`
}

func main() {
	// Run a database with docker: docker run --name test --restart unless-stopped -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=sqlutil -p 5432:5432 -d postgres:14-alpine
	// Connect to database
	db, err := sql.Open("postgres", "postgres://user:password@localhost/sqlutil?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS players(
			id SERIAL PRIMARY KEY,
			name VARCHAR NOT NULL,
			age INTEGER NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert players
	r9 := Player{
		Name: "Ronaldo Fenômeno",
		Age:  44,
	}
	r10 := Player{
		Name: "Ronaldinho Gaúcho",
		Age:  41,
	}
	flavour := sqlutil.PostgreSQLFlavor
	tag := "insert" // will use fields with fieldtag:"insert"
	ctx := context.Background()
	if err := sqlutil.Insert(ctx, db, sqlutil.PostgreSQLFlavor, tag, "players", &r9); err != nil {
		log.Fatal(err)
	}
	if err := sqlutil.Insert(ctx, db, sqlutil.PostgreSQLFlavor, tag, "players", &r10); err != nil {
		log.Fatal(err)
	}

	// Get player
	findOptions := sqlutil.NewFindOptions(flavour).WithFilter("name", r10.Name)
	if err := sqlutil.Get(ctx, db, "players", findOptions, &r10); err != nil {
		log.Fatal(err)
	}
	findOptions = sqlutil.NewFindOptions(flavour).WithFilter("name", r9.Name)
	if err := sqlutil.Get(ctx, db, "players", findOptions, &r9); err != nil {
		log.Fatal(err)
	}

	// Select players
	players := []*Player{}
	findAllOptions := sqlutil.NewFindAllOptions(flavour).WithLimit(10).WithOffset(0).WithOrderBy("name asc")
	if err := sqlutil.Select(ctx, db, "players", findAllOptions, &players); err != nil {
		log.Fatal(err)
	}
	for _, p := range players {
		fmt.Printf("%#v\n", p)
	}

	// Update player
	tag = "update" // will use fields with fieldtag:"update"
	r10.Name = "Ronaldinho Bruxo"
	if err := sqlutil.Update(ctx, db, sqlutil.PostgreSQLFlavor, tag, "players", r10.ID, &r10); err != nil {
		log.Fatal(err)
	}

	// Delete player
	if err := sqlutil.Delete(ctx, db, sqlutil.PostgreSQLFlavor, "players", r9.ID); err != nil {
		log.Fatal(err)
	}
}
```

Options for FindOptions and FindAllOptions:

```golang
package main

import (
	"github.com/allisson/sqlutil"
	_ "github.com/lib/pq"
)

func main() {
	findOptions := sqlutil.NewFindOptions(sqlutil.PostgreSQLFlavor).
		WithFields([]string{"id", "name"}). // Return only id and name fields
		WithFilter("id", 1).                // WHERE id = 1
		WithFilter("id", nil).              // WHERE id IS NULL
		WithFilter("id.in", "1,2,3").       // WHERE id IN (1, 2, 3)
		WithFilter("id.notin", "1,2,3").    // WHERE id NOT IN ($1, $2, $3)
		WithFilter("id.not", 1).            // WHERE id <> 1
		WithFilter("id.gt", 1).             // WHERE id > 1
		WithFilter("id.gte", 1).            // WHERE id >= 1
		WithFilter("id.lt", 1).             // WHERE id < 1
		WithFilter("id.lte", 1).            // WHERE id <= 1
		WithFilter("id.like", 1).           // WHERE id LIKE 1
		WithFilter("id.null", true).        // WHERE id.null IS NULL
		WithFilter("id.null", false)        // WHERE id.null IS NOT NULL

	findAllOptions := sqlutil.NewFindAllOptions(sqlutil.PostgreSQLFlavor).
		WithFields([]string{"id", "name"}). // Return only id and name fields
		WithFilter("id", 1).                // WHERE id = 1
		WithFilter("id", nil).              // WHERE id IS NULL
		WithFilter("id.in", "1,2,3").       // WHERE id IN (1, 2, 3)
		WithFilter("id.notin", "1,2,3").    // WHERE id NOT IN ($1, $2, $3)
		WithFilter("id.not", 1).            // WHERE id <> 1
		WithFilter("id.gt", 1).             // WHERE id > 1
		WithFilter("id.gte", 1).            // WHERE id >= 1
		WithFilter("id.lt", 1).             // WHERE id < 1
		WithFilter("id.lte", 1).            // WHERE id <= 1
		WithFilter("id.like", 1).           // WHERE id LIKE 1
		WithFilter("id.null", true).        // WHERE id.null IS NULL
		WithFilter("id.null", false).       // WHERE id.null IS NOT NULL
		WithLimit(10).                      // LIMIT 10
		WithOffset(0).                      // OFFSET 0
		WithOrderBy("name asc")             // ORDER BY name asc
}
```
