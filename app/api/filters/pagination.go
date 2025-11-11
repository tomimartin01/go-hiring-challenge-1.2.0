package filters

import (
	"errors"
	"net/url"
	"strconv"
)

type PaginationFilter struct {
	Offset int
	Limit  int
}

const (
	minLimitValue        = 1
	maxLimitValue        = 100
	defaultOffsetValue   = 0
	defaultLimitValue    = 10
	offsetParamName      = "offset"
	limitParamName       = "limit"
	ErrInvalidParamValue = "INVALID_PARAM_VALUE"
)

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
