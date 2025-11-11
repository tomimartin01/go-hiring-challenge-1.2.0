package filters

import "net/url"

const (
	categoryNameParamName = "category_name"
)

type CategoryFilter struct {
	CategoryName string
}

func (c *CategoryFilter) Parse(values url.Values) error {
	if values.Get(string(categoryNameParamName)) != "" {
		c.CategoryName = values.Get(string(categoryNameParamName))
	}
	return nil
}
