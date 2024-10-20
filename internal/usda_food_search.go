package food

import "strings"

type Search_Options struct {
	Query      string `url:"query"`
	DataType   string `url:"dataType,omitempty"`
	PageSize   int    `url:"pageSize,omitempty"`
	PageNumber int    `url:"pageNumber,omitempty"`
	SortBy     string `url:"sortBy,omitempty"`
	SortOrder  string `url:"sortOrder,omitempty"`
	BrandOwner string `url:"brandOwner,omitempty"`
}

type SearchOptions func(*Search_Options)

func NewSearch(options ...SearchOptions) *Search_Options {
	search_struct := &Search_Options{}
	for _, opt := range options {
		opt(search_struct)
	}
	return search_struct
}

func SearchQuery(query string) func(*Search_Options) {
	return func(so *Search_Options) {
		so.Query = query
	}
}

func SearchDataType(dataType []string) func(*Search_Options) {
	formatted_dataType := strings.Join(dataType, ",")
	return func(so *Search_Options) {
		so.DataType = formatted_dataType
	}
}

func SearchPageSize(pageSize int) func(*Search_Options) {
	return func(so *Search_Options) {
		so.PageSize = pageSize
	}
}

func SearchPageNumber(pageNumber int) func(*Search_Options) {
	return func(so *Search_Options) {
		so.PageNumber = pageNumber
	}
}

func SearchSortBy(sortBy string) func(*Search_Options) {
	return func(so *Search_Options) {
		so.SortBy = sortBy
	}
}

func SearchSortOrder(sortOrder string) func(*Search_Options) {
	return func(so *Search_Options) {
		so.SortOrder = sortOrder
	}
}

func SearchBrandOwner(brandOwner string) func(*Search_Options) {
	return func(so *Search_Options) {
		so.BrandOwner = brandOwner
	}
}
