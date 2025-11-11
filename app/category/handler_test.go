package category

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCategoryHandler_HandleGet(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		req := require.New(t)

		mux := http.NewServeMux()

		mux.HandleFunc("GET /categories", func(w http.ResponseWriter, r *http.Request) {
			handler := NewCategoryHandler(nil)
			handler.HandleGet(w, r)
		})

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, &http.Request{
			Method: "GET",
			URL:    &url.URL{Path: "/categories"},
		})

		req.Equal(http.StatusOK, recorder.Code)
		req.Equal(http.Header{"Content-Type": []string{"application/json"}}, recorder.Header())
		req.Equal(`[]`, recorder.Body.String())

	})
}
