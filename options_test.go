package sqlutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindOptions(t *testing.T) {
	t.Run("NewFindOptions", func(t *testing.T) {
		options := NewFindOptions(PostgreSQLFlavor)
		assert.Equal(t, PostgreSQLFlavor, options.Flavor)
		assert.Equal(t, []string{"*"}, options.Fields)
		assert.Equal(t, map[string]interface{}{}, options.Filters)
	})

	t.Run("NewFindOptions with modifiers", func(t *testing.T) {
		options := NewFindOptions(PostgreSQLFlavor).
			WithFields([]string{"id"}).
			WithFilter("key1", "value1").
			WithFilter("key2", "value2")
		assert.Equal(t, PostgreSQLFlavor, options.Flavor)
		assert.Equal(t, []string{"id"}, options.Fields)
		assert.Equal(t, map[string]interface{}{"key1": "value1", "key2": "value2"}, options.Filters)
	})
}

func TestFindAllOptions(t *testing.T) {
	t.Run("NewFindAllOptions", func(t *testing.T) {
		options := NewFindAllOptions(PostgreSQLFlavor)
		assert.Equal(t, PostgreSQLFlavor, options.Flavor)
		assert.Equal(t, []string{"*"}, options.Fields)
		assert.Equal(t, map[string]interface{}{}, options.Filters)
		assert.Equal(t, 0, options.Limit)
		assert.Equal(t, 0, options.Offset)
		assert.Equal(t, "", options.OrderBy)
	})

	t.Run("NewFindAllOptions with modifiers", func(t *testing.T) {
		options := NewFindAllOptions(PostgreSQLFlavor).
			WithFields([]string{"id"}).
			WithFilter("key1", "value1").
			WithFilter("key2", "value2").
			WithLimit(10).
			WithOffset(10).
			WithOrderBy("column asc")
		assert.Equal(t, PostgreSQLFlavor, options.Flavor)
		assert.Equal(t, []string{"id"}, options.Fields)
		assert.Equal(t, map[string]interface{}{"key1": "value1", "key2": "value2"}, options.Filters)
		assert.Equal(t, 10, options.Limit)
		assert.Equal(t, 10, options.Offset)
		assert.Equal(t, "column asc", options.OrderBy)
	})
}
