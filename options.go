package sqlutil

// Supported flavors.
const (
	MySQLFlavor Flavor = iota + 1
	PostgreSQLFlavor
	SQLiteFlavor
)

// Flavor is the flag to control the format of compiled sql.
type Flavor int

// FindOptions provides configuration for FindQuery function.
type FindOptions struct {
	Flavor  Flavor
	Fields  []string
	Filters map[string]interface{}
}

// WithFields is a helper function to construct functional options that sets Fields field.
func (f *FindOptions) WithFields(fields []string) *FindOptions {
	copy := *f
	copy.Fields = fields
	return &copy
}

// WithFilter is a helper function to construct functional options that sets Filters field.
func (f *FindOptions) WithFilter(field string, value interface{}) *FindOptions {
	copy := *f
	copy.Filters[field] = value
	return &copy
}

// NewFindOptions returns a FindOptions.
func NewFindOptions(flavor Flavor) *FindOptions {
	return &FindOptions{
		Fields:  []string{"*"},
		Flavor:  flavor,
		Filters: make(map[string]interface{}),
	}
}

// FindAllOptions provides configuration for FindAllQuery function.
type FindAllOptions struct {
	Flavor  Flavor
	Fields  []string
	Filters map[string]interface{}
	Limit   int
	Offset  int
	OrderBy string
}

// WithFields is a helper function to construct functional options that sets Fields field.
func (f *FindAllOptions) WithFields(fields []string) *FindAllOptions {
	copy := *f
	copy.Fields = fields
	return &copy
}

// WithFilter is a helper function to construct functional options that sets Filters field.
func (f *FindAllOptions) WithFilter(field string, value interface{}) *FindAllOptions {
	copy := *f
	copy.Filters[field] = value
	return &copy
}

// WithLimit is a helper function to construct functional options that sets Limit field.
func (f *FindAllOptions) WithLimit(limit int) *FindAllOptions {
	copy := *f
	copy.Limit = limit
	return &copy
}

// WithOffset is a helper function to construct functional options that sets Offset field.
func (f *FindAllOptions) WithOffset(offset int) *FindAllOptions {
	copy := *f
	copy.Offset = offset
	return &copy
}

// WithOrderBy is a helper function to construct functional options that sets OrderBy field.
func (f *FindAllOptions) WithOrderBy(orderBy string) *FindAllOptions {
	copy := *f
	copy.OrderBy = orderBy
	return &copy
}

// NewFindAllOptions returns a FindAllOptions.
func NewFindAllOptions(flavor Flavor) *FindAllOptions {
	return &FindAllOptions{
		Fields:  []string{"*"},
		Flavor:  flavor,
		Filters: make(map[string]interface{}),
	}
}
