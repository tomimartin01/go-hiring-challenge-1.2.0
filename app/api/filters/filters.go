package filters

import (
	"net/url"
)

type FilterBuilder interface {
	Parse(values url.Values) error
}
