package sqlutil

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type player struct {
	ID   int    `db:"id" fieldtag:"insert"`
	Name string `db:"name" fieldtag:"insert,update"`
}

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Ronaldinho Gaúcho")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM players WHERE id = $1`)).
		WithArgs("R10").
		WillReturnRows(rows)

	options := NewFindOptions(PostgreSQLFlavor).WithFilter("id", "R10")
	p := player{}
	err = Get(context.Background(), db, "players", options, &p)
	assert.Nil(t, err)
	assert.Equal(t, 1, p.ID)
	assert.Equal(t, "Ronaldinho Gaúcho", p.Name)
}

func TestSelect(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Ronaldinho Gaúcho").AddRow(2, "Ronaldo Fenômeno")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM players LIMIT 10 OFFSET 0 FOR UPDATE SKIP LOCKED`)).WillReturnRows(rows)

	options := NewFindAllOptions(PostgreSQLFlavor).WithLimit(10).WithOffset(0).WithForUpdate("SKIP LOCKED")
	p := []*player{}
	err = Select(context.Background(), db, "players", options, &p)
	assert.Nil(t, err)
	assert.Len(t, p, 2)
	assert.Equal(t, 1, p[0].ID)
	assert.Equal(t, "Ronaldinho Gaúcho", p[0].Name)
	assert.Equal(t, 2, p[1].ID)
	assert.Equal(t, "Ronaldo Fenômeno", p[1].Name)
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	//rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Ronaldinho Gaúcho")
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO players`)).
		WithArgs(1, "Ronaldinho Gaúcho").
		WillReturnResult(sqlmock.NewResult(1, 1))

	p := player{
		ID:   1,
		Name: "Ronaldinho Gaúcho",
	}
	err = Insert(context.Background(), db, PostgreSQLFlavor, "", "players", &p)
	assert.Nil(t, err)
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE players SET id = $1, name = $2 WHERE id = $3`)).
		WithArgs(1, "Ronaldinho Gaúcho", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	p := player{
		ID:   1,
		Name: "Ronaldinho Gaúcho",
	}
	err = Update(context.Background(), db, PostgreSQLFlavor, "", "players", p.ID, &p)
	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM players WHERE id = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	p := player{
		ID:   1,
		Name: "Ronaldinho Gaúcho",
	}
	err = Delete(context.Background(), db, PostgreSQLFlavor, "players", p.ID)
	assert.Nil(t, err)
}
