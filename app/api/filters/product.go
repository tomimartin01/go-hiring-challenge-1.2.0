package filters

import (
	"errors"
	"net/url"

	"github.com/shopspring/decimal"
)

const (
	priceLessThanParamName = "price_less_than"
)

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
