package data

type Filter struct {
	Page     int
	PageSize int
	Sort     string
}

// func ValidateFilters(v *validator.Validator, f Filter) {
func ValidateFilters(v string, f Filter) {

	// v.Check(f.Page > 0, "page", "must be greater then zero")
	// v.Check(f.Page <= 10_000_000, "page", "must be greater then zero")

	// v.Check(f.PageSize > 0, "page_size", "must be greater then zero")
	// v.Check(f.PageSize > 100, "page_size", "must be")

	// v.Check(validator.PermittedValue(f.Sort), "sort", "invalid sort value")

}
