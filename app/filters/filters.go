package filters

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/shopspring/decimal"
)

const (
	minLimitValue          = 1
	maxLimitValue          = 100
	defaultOffsetValue     = 0
	defaultLimitValue      = 10
	offsetParamName        = "offset"
	limitParamName         = "limit"
	categoryNameParamName  = "category_name"
	priceLessThanParamName = "price_less_than"
	ErrInvalidParamValue   = "INVALID_PARAM_VALUE"
)

type FilterBuilder interface {
	Parse(values url.Values) error
}

type PaginationFilter struct {
	Offset int
	Limit  int
}

type CategoryFilter struct {
	CategoryName string
}

type ProductFilter struct {
	PriceLessThan *decimal.Decimal
}

func (p *ProductFilter) Parse(values url.Values) error {
	p.PriceLessThan = nil
	if values.Get(string(priceLessThanParamName)) != "" {
		priceValue, err := decimal.NewFromString(values.Get(string(priceLessThanParamName)))
		if err != nil {
			return errors.New(ErrInvalidParamValue)
		}
		p.PriceLessThan = &priceValue
	}
	return nil
}

func (c *CategoryFilter) Parse(values url.Values) error {
	if values.Get(string(categoryNameParamName)) != "" {
		c.CategoryName = values.Get(string(categoryNameParamName))
	}
	return nil
}

func (p *PaginationFilter) Parse(values url.Values) error {
	p.Offset = defaultOffsetValue
	if values.Get(string(offsetParamName)) != "" {
		offsetValue, err := strconv.Atoi(values.Get(string(offsetParamName)))
		if err != nil {
			return errors.New(ErrInvalidParamValue)
		}
		p.Offset = offsetValue
	}

	p.Limit = defaultLimitValue
	if values.Get(string(limitParamName)) != "" {
		limitValue, err := strconv.Atoi(values.Get(string(limitParamName)))
		if err != nil {
			return errors.New(ErrInvalidParamValue)
		}
		p.Limit = validateLimit(limitValue)
	}

	return nil
}

func validateLimit(limit int) int {
	if limit < minLimitValue {
		return minLimitValue
	}
	if limit > maxLimitValue {
		return maxLimitValue
	}
	return limit
}
