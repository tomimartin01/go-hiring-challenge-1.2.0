package filters

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCategoryFilter_Parse(t *testing.T) {
	tests := map[string]struct {
		queryParams      map[string]string
		expectedCategory string
	}{
		"empty category name": {
			queryParams:      map[string]string{},
			expectedCategory: "",
		},
		"valid category name": {
			queryParams:      map[string]string{"category_name": "Electronics"},
			expectedCategory: "Electronics",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := require.New(t)

			values := url.Values{}
			for key, val := range tt.queryParams {
				values.Set(key, val)
			}

			filter := &CategoryFilter{}
			err := filter.Parse(values)

			req.NoError(err)
			req.Equal(tt.expectedCategory, filter.CategoryName)
		})
	}
}
