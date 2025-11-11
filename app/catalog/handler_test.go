package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	repoMocks "github.com/mytheresa/go-hiring-challenge/models/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCatalogHandler_HandleGetProductDetails(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		req := require.New(t)
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockProduct := models.Product{
			ID:    1,
			Code:  "PROD001",
			Price: decimal.NewFromFloat(10.99),
			Variants: []models.Variant{
				{
					ID:        1,
					ProductID: 1,
					Name:      "Variant A",
					SKU:       "SKU001A",
					Price:     decimal.NewFromFloat(11.99),
				},
			},
			Categories: []models.Category{
				{
					ID:   1,
					Code: "CATEGORY_1",
					Name: "Category 1",
				},
			},
		}

		request, err := http.NewRequest("GET", "/catalog/PROD001", nil)
		req.NoError(err)
		request.SetPathValue("code", "PROD001")

		recorder := httptest.NewRecorder()

		mockRepo := repoMocks.NewMockProductsRepositoryInterface(ctl)
		mockRepo.EXPECT().GetProductDetails("PROD001").Return(
			mockProduct, nil)
		handler := NewCatalogHandler(mockRepo)

		handlerFunc := http.HandlerFunc(handler.HandleGetProductDetails)
		handlerFunc.ServeHTTP(recorder, request)

		req.Equal(http.StatusOK, recorder.Code)
		req.Equal(http.Header{"Content-Type": []string{"application/json"}}, recorder.Header())

		responseProduct := ProductDetailsResponse{}
		err = json.Unmarshal(recorder.Body.Bytes(), &responseProduct)
		req.NoError(err)

		req.Equal(mockProduct.Code, responseProduct.Product.Code)
		req.Equal(10.99, responseProduct.Product.Price)

		req.Len(responseProduct.Product.Variants, 1)
		req.Equal(mockProduct.Variants[0].SKU, responseProduct.Product.Variants[0].Code)
		req.Equal(mockProduct.Variants[0].Name, responseProduct.Product.Variants[0].Name)
		req.Equal(mockProduct.Variants[0].Price.InexactFloat64(), responseProduct.Product.Variants[0].Price) // Variant price as float

		req.Len(responseProduct.Product.Categories, 1)
		req.Equal(mockProduct.Categories[0].Code, responseProduct.Product.Categories[0].Code)
		req.Equal(mockProduct.Categories[0].Name, responseProduct.Product.Categories[0].Name)
	})

	t.Run("Ok, no variant price", func(t *testing.T) {
		req := require.New(t)
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockProduct := models.Product{
			ID:    1,
			Code:  "PROD001",
			Price: decimal.NewFromFloat(10.99),
			Variants: []models.Variant{
				{
					ID:        1,
					ProductID: 1,
					Name:      "Variant A",
					SKU:       "SKU001A",
				},
			},
			Categories: []models.Category{
				{
					ID:   1,
					Code: "CATEGORY_1",
					Name: "Category 1",
				},
			},
		}

		request, err := http.NewRequest("GET", "/catalog/PROD001", nil)
		req.NoError(err)
		request.SetPathValue("code", "PROD001")

		recorder := httptest.NewRecorder()

		mockRepo := repoMocks.NewMockProductsRepositoryInterface(ctl)
		mockRepo.EXPECT().GetProductDetails("PROD001").Return(
			mockProduct, nil)
		handler := NewCatalogHandler(mockRepo)

		handlerFunc := http.HandlerFunc(handler.HandleGetProductDetails)
		handlerFunc.ServeHTTP(recorder, request)

		req.Equal(http.StatusOK, recorder.Code)
		req.Equal(http.Header{"Content-Type": []string{"application/json"}}, recorder.Header())

		responseProduct := ProductDetailsResponse{}
		err = json.Unmarshal(recorder.Body.Bytes(), &responseProduct)
		req.NoError(err)

		req.Equal(mockProduct.Code, responseProduct.Product.Code)

		// Check that the price is the product price and not the variant price
		req.Equal(mockProduct.Price.InexactFloat64(), responseProduct.Product.Price)
		req.Equal(mockProduct.Price.InexactFloat64(), responseProduct.Product.Variants[0].Price)
	})

	t.Run("DB Error", func(t *testing.T) {
		req := require.New(t)
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockRepo := repoMocks.NewMockProductsRepositoryInterface(ctl)
		mockRepo.EXPECT().GetProductDetails("PROD001").Return(models.Product{}, errors.New("error"))
		handler := NewCatalogHandler(mockRepo)

		handlerFunc := http.HandlerFunc(handler.HandleGetProductDetails)
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/catalog/PROD001", nil)
		req.NoError(err)
		request.SetPathValue("code", "PROD001")
		handlerFunc.ServeHTTP(recorder, request)

		req.Equal(http.StatusInternalServerError, recorder.Code)
		req.Equal(http.Header{"Content-Type": []string{"application/json"}}, recorder.Header())

		responseError := make(map[string]string)
		err = json.Unmarshal(recorder.Body.Bytes(), &responseError)
		req.NoError(err)

		req.Equal(api.InternalServerError, responseError["error"])
	})

	t.Run("Bad request", func(t *testing.T) {
		req := require.New(t)
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockRepo := repoMocks.NewMockProductsRepositoryInterface(ctl)
		handler := NewCatalogHandler(mockRepo)
		handlerFunc := http.HandlerFunc(handler.HandleGetProductDetails)
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/catalog/PROD001", nil)
		req.NoError(err)
		handlerFunc.ServeHTTP(recorder, request)

		req.Equal(http.StatusBadRequest, recorder.Code)
		req.Equal(http.Header{"Content-Type": []string{"application/json"}}, recorder.Header())

		responseError := make(map[string]string)
		err = json.Unmarshal(recorder.Body.Bytes(), &responseError)
		req.NoError(err)

		req.Equal(api.BadRequestError, responseError["error"])
	})

	t.Run("No product found", func(t *testing.T) {
		req := require.New(t)
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockRepo := repoMocks.NewMockProductsRepositoryInterface(ctl)
		mockRepo.EXPECT().GetProductDetails("PROD009").Return(models.Product{}, errors.New(models.ProductNotFoundError))
		handler := NewCatalogHandler(mockRepo)

		handlerFunc := http.HandlerFunc(handler.HandleGetProductDetails)
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/catalog/PROD009", nil)
		request.SetPathValue("code", "PROD009")
		req.NoError(err)
		handlerFunc.ServeHTTP(recorder, request)

		req.Equal(http.StatusNotFound, recorder.Code)
		req.Equal(http.Header{"Content-Type": []string{"application/json"}}, recorder.Header())

		responseError := make(map[string]string)
		err = json.Unmarshal(recorder.Body.Bytes(), &responseError)
		req.NoError(err)

		req.Equal(api.NotFoundError, responseError["error"])
	})
}

func TestCatalogHandler_HandleGet(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		req := require.New(t)
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockProducts := []models.Product{
			{
				ID:    1,
				Code:  "PROD001",
				Price: decimal.NewFromFloat(10.99),
			},
		}
		mockTotal := len(mockProducts)

		mockRepo := repoMocks.NewMockProductsRepositoryInterface(ctl)
		mockRepo.EXPECT().GetAllProducts(gomock.Any()).Return(mockProducts, int64(mockTotal), nil)
		handler := NewCatalogHandler(mockRepo)
		handlerFunc := http.HandlerFunc(handler.HandleGet)
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/catalog", nil)
		req.NoError(err)
		handlerFunc.ServeHTTP(recorder, request)

		req.Equal(http.StatusOK, recorder.Code)
		req.Equal(http.Header{"Content-Type": []string{"application/json"}}, recorder.Header())

		response := Response{}
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		req.NoError(err)

		req.Equal(mockTotal, len(response.Products))
		req.Equal(int64(mockTotal), response.Total)
		req.Equal(mockProducts[0].Code, response.Products[0].Code)
		req.Equal(mockProducts[0].Price.InexactFloat64(), response.Products[0].Price)
	})

	t.Run("DB Error", func(t *testing.T) {
		req := require.New(t)
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockRepo := repoMocks.NewMockProductsRepositoryInterface(ctl)
		mockRepo.EXPECT().GetAllProducts(gomock.Any()).Return([]models.Product{}, int64(0), errors.New("error"))
		handler := NewCatalogHandler(mockRepo)
		handlerFunc := http.HandlerFunc(handler.HandleGet)
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/catalog", nil)
		req.NoError(err)
		handlerFunc.ServeHTTP(recorder, request)

		req.Equal(http.StatusInternalServerError, recorder.Code)
		req.Equal(http.Header{"Content-Type": []string{"application/json"}}, recorder.Header())

		responseError := make(map[string]string)
		err = json.Unmarshal(recorder.Body.Bytes(), &responseError)
		req.NoError(err)

		req.Equal(api.InternalServerError, responseError["error"])
	})
}
