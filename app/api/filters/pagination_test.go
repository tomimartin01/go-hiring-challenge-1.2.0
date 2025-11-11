package filters

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaginationFilter_Parse(t *testing.T) {
	tests := map[string]struct {
		queryParams    map[string]string
		expectedLimit  int
		expectedOffset int
		expectError    bool
	}{
		"defaults when no params": {
			queryParams:    map[string]string{},
			expectedLimit:  10,
			expectedOffset: 0,
			expectError:    false,
		},
		"valid offset and limit": {
			queryParams:    map[string]string{"offset": "20", "limit": "50"},
			expectedLimit:  50,
			expectedOffset: 20,
			expectError:    false,
		},
		"limit to 0 should return 1": {
			queryParams:    map[string]string{"limit": "0"},
			expectedLimit:  1,
			expectedOffset: 0,
			expectError:    false,
		},
		"limit to 500 should return 100": {
			queryParams:    map[string]string{"limit": "500"},
			expectedLimit:  100,
			expectedOffset: 0,
			expectError:    false,
		},
		"negative limit clamped to 1": {
			queryParams:    map[string]string{"limit": "-10"},
			expectedLimit:  1,
			expectedOffset: 0,
			expectError:    false,
		},
		"invalid offset returns error": {
			queryParams: map[string]string{"offset": "abc"},
			expectError: true,
		},
		"invalid limit returns error": {
			queryParams: map[string]string{"limit": "xyz"},
			expectError: true,
		},
		"offset zero": {
			queryParams:    map[string]string{"offset": "0"},
			expectedLimit:  10,
			expectedOffset: 0,
			expectError:    false,
		},
		"large offset value": {
			queryParams:    map[string]string{"offset": "10000"},
			expectedLimit:  10,
			expectedOffset: 10000,
			expectError:    false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := require.New(t)

			values := url.Values{}
			for key, val := range tt.queryParams {
				values.Set(key, val)
			}

			filter := &PaginationFilter{}
			err := filter.Parse(values)

			if tt.expectError {
				req.Error(err)
				req.EqualError(err, ErrInvalidParamValue)
			} else {
				req.NoError(err)
				req.Equal(tt.expectedLimit, filter.Limit)
				req.Equal(tt.expectedOffset, filter.Offset)
			}
		})
	}
}
