package sqlutil

import (
	"strings"

	"github.com/huandu/go-sqlbuilder"
)

func parseIn(value string) []interface{} {
	values := strings.Split(value, ",")
	result := make([]interface{}, len(values))
	for i := range values {
		result[i] = values[i]
	}
	return result
}

func parseFilter(sb *sqlbuilder.SelectBuilder, key string, value interface{}) {
	if strings.Contains(key, ".") {
		split := strings.Split(key, ".")
		parsedKey := split[0]
		compare := split[1]
		switch compare {
		case "in":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				sb.Where(sb.In(parsedKey, values...))
			}
		case "notin":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				sb.Where(sb.NotIn(parsedKey, values...))
			}
		case "not":
			sb.Where(sb.NotEqual(parsedKey, value))
		case "gt":
			sb.Where(sb.GreaterThan(parsedKey, value))
		case "gte":
			sb.Where(sb.GreaterEqualThan(parsedKey, value))
		case "lt":
			sb.Where(sb.LessThan(parsedKey, value))
		case "lte":
			sb.Where(sb.LessEqualThan(parsedKey, value))
		case "like":
			sb.Where(sb.Like(parsedKey, value))
		case "null":
			valueBool, ok := value.(bool)
			if ok {
				if valueBool {
					sb.Where(sb.IsNull(key))
				} else {
					sb.Where(sb.IsNotNull(key))
				}
			}
		}
	} else {
		switch value.(type) {
		case nil:
			sb.Where(sb.IsNull(key))
		default:
			sb.Where(sb.Equal(key, value))
		}
	}
}

// FindQuery returns compiled SELECT string and args.
func FindQuery(tableName string, options *FindOptions) (string, []interface{}) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.Flavor(options.Flavor))
	sb.Select(options.Fields...).From(tableName)
	for key, value := range options.Filters {
		parseFilter(sb, key, value)
	}
	return sb.Build()
}

// FindAllQuery returns compiled SELECT string and args.
func FindAllQuery(tableName string, options *FindAllOptions) (string, []interface{}) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.Flavor(options.Flavor))
	sb.Select(options.Fields...).From(tableName).Limit(options.Limit).Offset(options.Offset)
	for key, value := range options.Filters {
		parseFilter(sb, key, value)
	}
	if options.OrderBy != "" {
		sb.OrderBy(options.OrderBy)
	}
	return sb.Build()
}

// InsertQuery returns compiled INSERT string and args.
func InsertQuery(flavor Flavor, tag, tableName string, structValue interface{}) (string, []interface{}) {
	theStruct := sqlbuilder.NewStruct(structValue).For(sqlbuilder.Flavor(flavor))
	ib := theStruct.InsertIntoForTag(tableName, tag, structValue)
	ib.SQL("RETURNING *")
	return ib.Build()
}

// UpdateQuery returns compiled UPDATE string and args.
func UpdateQuery(flavor Flavor, tag, tableName string, id interface{}, structValue interface{}) (string, []interface{}) {
	theStruct := sqlbuilder.NewStruct(structValue).For(sqlbuilder.Flavor(flavor))
	ub := theStruct.UpdateForTag(tableName, tag, structValue)
	ub.Where(ub.Equal("id", id))
	return ub.Build()
}

// DeleteQuery returns compiled DELETE string and args.
func DeleteQuery(flavor Flavor, tableName string, id interface{}) (string, []interface{}) {
	db := sqlbuilder.NewDeleteBuilder()
	db.SetFlavor(sqlbuilder.Flavor(flavor))
	db.DeleteFrom(tableName)
	db.Where(db.Equal("id", id))
	return db.Build()
}
