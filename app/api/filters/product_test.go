package filters

import (
	"net/url"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestProductFilter_Parse(t *testing.T) {
	tests := map[string]struct {
		queryParams       map[string]string
		expectedPriceLess *decimal.Decimal
		expectError       bool
	}{
		"no params sets nil": {
			queryParams:       map[string]string{},
			expectedPriceLess: nil,
			expectError:       false,
		},
		"valid param": {
			queryParams:       map[string]string{"price_less_than": "99.99"},
			expectedPriceLess: decimalPtr("99.99"),
			expectError:       false,
		},
		"valid integer": {
			queryParams:       map[string]string{"price_less_than": "100"},
			expectedPriceLess: decimalPtr("100"),
			expectError:       false,
		},
		"zero value": {
			queryParams:       map[string]string{"price_less_than": "0"},
			expectedPriceLess: decimalPtr("0"),
			expectError:       false,
		},
		"negative value": {
			queryParams:       map[string]string{"price_less_than": "-10.50"},
			expectedPriceLess: decimalPtr("-10.50"),
			expectError:       false,
		},
		"invalid param": {
			queryParams: map[string]string{"price_less_than": "abc"},
			expectError: true,
		},
		"invalid with symbols": {
			queryParams: map[string]string{"price_less_than": "$99.99"},
			expectError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := require.New(t)

			values := url.Values{}
			for key, val := range tt.queryParams {
				values.Set(key, val)
			}

			filter := &ProductFilter{}
			err := filter.Parse(values)

			if tt.expectError {
				req.Error(err)
				req.EqualError(err, ErrInvalidParamValue)
			} else {
				req.NoError(err)
				if tt.expectedPriceLess == nil {
					req.Nil(filter.PriceLessThan)
				} else {
					req.NotNil(filter.PriceLessThan)
					req.True(tt.expectedPriceLess.Equal(*filter.PriceLessThan))
				}
			}
		})
	}
}

func decimalPtr(value string) *decimal.Decimal {
	d, _ := decimal.NewFromString(value)
	return &d
}
