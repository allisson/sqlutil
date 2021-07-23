package sqlutil

import (
	"testing"

	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

func TestParseFilter(t *testing.T) {
	var tests = []struct {
		kind         string
		key          string
		value        interface{}
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{"equals", "id", 1, `SELECT * FROM test_table WHERE id = $1`, []interface{}{1}},
		{"equals nil", "id", nil, `SELECT * FROM test_table WHERE id IS NULL`, []interface{}(nil)},
		{"in", "id.in", "1,2,3", `SELECT * FROM test_table WHERE id IN ($1, $2, $3)`, []interface{}{"1", "2", "3"}},
		{"notin", "id.notin", "1,2,3", `SELECT * FROM test_table WHERE id NOT IN ($1, $2, $3)`, []interface{}{"1", "2", "3"}},
		{"not", "id.not", 1, `SELECT * FROM test_table WHERE id <> $1`, []interface{}{1}},
		{"gt", "id.gt", 1, `SELECT * FROM test_table WHERE id > $1`, []interface{}{1}},
		{"gte", "id.gte", 1, `SELECT * FROM test_table WHERE id >= $1`, []interface{}{1}},
		{"lt", "id.lt", 1, `SELECT * FROM test_table WHERE id < $1`, []interface{}{1}},
		{"lte", "id.lte", 1, `SELECT * FROM test_table WHERE id <= $1`, []interface{}{1}},
		{"like", "id.like", 1, `SELECT * FROM test_table WHERE id LIKE $1`, []interface{}{1}},
		{"null true", "id.null", true, `SELECT * FROM test_table WHERE id.null IS NULL`, []interface{}(nil)},
		{"null false", "id.null", false, `SELECT * FROM test_table WHERE id.null IS NOT NULL`, []interface{}(nil)},
	}
	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			sb := sqlbuilder.NewSelectBuilder()
			sb.SetFlavor(sqlbuilder.Flavor(PostgreSQLFlavor))
			sb.Select("*").From("test_table")
			parseFilter(sb, tt.key, tt.value)
			sqlQuery, args := sb.Build()
			assert.Equal(t, tt.expectedSQL, sqlQuery)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestFindQuery(t *testing.T) {
	expectedSQLQuery := `SELECT * FROM test_table WHERE id = $1`
	expectedArgs := []interface{}{1}
	options := NewFindOptions(PostgreSQLFlavor).WithFilter("id", 1)
	sqlQuery, args := FindQuery("test_table", options)
	assert.Equal(t, expectedSQLQuery, sqlQuery)
	assert.Equal(t, expectedArgs, args)
}

func TestFindAllQuery(t *testing.T) {
	expectedSQLQuery := `SELECT * FROM test_table WHERE id = $1 ORDER BY id asc LIMIT 50 OFFSET 10`
	expectedArgs := []interface{}{1}
	options := NewFindAllOptions(PostgreSQLFlavor).
		WithFilter("id", 1).
		WithLimit(50).
		WithOffset(10).
		WithOrderBy("id asc")
	sqlQuery, args := FindAllQuery("test_table", options)
	assert.Equal(t, expectedSQLQuery, sqlQuery)
	assert.Equal(t, expectedArgs, args)
}

type player struct {
	ID   int    `db:"id" fieldtag:"insert"`
	Name string `db:"name" fieldtag:"insert,update"`
}

func TestInsertQuery(t *testing.T) {
	expectedSQLQuery := `INSERT INTO players (id, name) VALUES ($1, $2)`
	expectedArgs := []interface{}{1, "Ronaldinho 10"}
	r10 := player{ID: 1, Name: "Ronaldinho 10"}
	sqlQuery, args := InsertQuery(PostgreSQLFlavor, "insert", "players", &r10)
	assert.Equal(t, expectedSQLQuery, sqlQuery)
	assert.Equal(t, expectedArgs, args)
}

func TestUpdateQuery(t *testing.T) {
	expectedSQLQuery := `UPDATE players SET name = $1 WHERE id = $2`
	expectedArgs := []interface{}{"Ronaldinho Bruxo", 1}
	r10 := player{ID: 1, Name: "Ronaldinho Bruxo"}
	sqlQuery, args := UpdateQuery(PostgreSQLFlavor, "update", "players", r10.ID, &r10)
	assert.Equal(t, expectedSQLQuery, sqlQuery)
	assert.Equal(t, expectedArgs, args)
}

func TestDeleteQuery(t *testing.T) {
	expectedSQLQuery := `DELETE FROM players WHERE id = $1`
	expectedArgs := []interface{}{1}
	sqlQuery, args := DeleteQuery(PostgreSQLFlavor, "players", 1)
	assert.Equal(t, expectedSQLQuery, sqlQuery)
	assert.Equal(t, expectedArgs, args)
}
