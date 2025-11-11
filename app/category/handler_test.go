package category

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	models_mocks "github.com/mytheresa/go-hiring-challenge/models/mocks"
	repoMocks "github.com/mytheresa/go-hiring-challenge/models/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCategoryHandler_HandleGet(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		req := require.New(t)
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockCategories := []models.Category{
			{
				ID:   1,
				Code: "CATEGORY_1",
				Name: "Category 1",
			},
		}

		request, err := http.NewRequest("GET", "/categories", nil)
		req.NoError(err)
		recorder := httptest.NewRecorder()

		mockRepo := repoMocks.NewMockCategoriesRepositoryInterface(ctl)
		mockRepo.EXPECT().GetAllCategories().Return(
			mockCategories, nil)
		handler := NewCategoryHandler(mockRepo)

		handlerFunc := http.HandlerFunc(handler.HandleGet)
		handlerFunc.ServeHTTP(recorder, request)

		req.Equal(http.StatusOK, recorder.Code)
		req.Equal(http.Header{"Content-Type": []string{"application/json"}}, recorder.Header())

		responseCategories := make([]models.Category, 1)
		err = json.Unmarshal(recorder.Body.Bytes(), &responseCategories)
		req.NoError(err)

		req.Equal(1, len(responseCategories))
		req.Equal(mockCategories[0].Code, responseCategories[0].Code)
		req.Equal(mockCategories[0].Name, responseCategories[0].Name)
	})

	t.Run("DB Error", func(t *testing.T) {
		req := require.New(t)
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockRepo := models_mocks.NewMockCategoriesRepositoryInterface(ctl)
		mockRepo.EXPECT().GetAllCategories().Return(nil, errors.New("error"))

		request, err := http.NewRequest("GET", "/categories", nil)
		req.NoError(err)
		recorder := httptest.NewRecorder()

		handler := NewCategoryHandler(mockRepo)
		handlerFunc := http.HandlerFunc(handler.HandleGet)
		handlerFunc.ServeHTTP(recorder, request)

		req.Equal(http.StatusInternalServerError, recorder.Code)
		req.Equal(http.Header{"Content-Type": []string{"application/json"}}, recorder.Header())

		responseError := make(map[string]string)
		err = json.Unmarshal(recorder.Body.Bytes(), &responseError)
		req.NoError(err)

		req.Equal(api.InternalServerError, responseError["error"])
	})
}

func TestCategoryHandler_HandleCreate(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		req := require.New(t)
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockCategoryDto := models.CategoryDto{
			Code: "CATEGORY_1",
			Name: "Category 1",
		}

		mockCategory := models.Category{
			ID:   1,
			Code: mockCategoryDto.Code,
			Name: mockCategoryDto.Name,
		}

		mockCategoryDtoJson, err := json.Marshal(mockCategoryDto)
		req.NoError(err)

		request, err := http.NewRequest("POST", "/categories", bytes.NewBuffer(mockCategoryDtoJson))
		req.NoError(err)
		recorder := httptest.NewRecorder()

		mockRepo := repoMocks.NewMockCategoriesRepositoryInterface(ctl)
		mockRepo.EXPECT().Create(mockCategoryDto).Return(mockCategory, nil)
		handler := NewCategoryHandler(mockRepo)
		handlerFunc := http.HandlerFunc(handler.HandleCreate)
		handlerFunc.ServeHTTP(recorder, request)

		req.Equal(http.StatusOK, recorder.Code)
		req.Equal(http.Header{"Content-Type": []string{"application/json"}}, recorder.Header())

		responseCategory := models.Category{}
		err = json.Unmarshal(recorder.Body.Bytes(), &responseCategory)
		req.NoError(err)

		req.Equal(mockCategory.Code, responseCategory.Code)
		req.Equal(mockCategory.Name, responseCategory.Name)
	})

	t.Run("DB error", func(t *testing.T) {
		req := require.New(t)
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockCategoryDto := models.CategoryDto{
			Code: "CATEGORY_1",
			Name: "Category 1",
		}

		mockCategoryDtoJson, err := json.Marshal(mockCategoryDto)
		req.NoError(err)

		request, err := http.NewRequest("POST", "/categories", bytes.NewBuffer(mockCategoryDtoJson))
		req.NoError(err)
		recorder := httptest.NewRecorder()

		mockRepo := repoMocks.NewMockCategoriesRepositoryInterface(ctl)
		mockRepo.EXPECT().Create(mockCategoryDto).Return(models.Category{}, errors.New("error"))
		handler := NewCategoryHandler(mockRepo)
		handlerFunc := http.HandlerFunc(handler.HandleCreate)
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

		request, err := http.NewRequest("POST", "/categories", bytes.NewBuffer(nil))
		req.NoError(err)
		recorder := httptest.NewRecorder()

		mockRepo := repoMocks.NewMockCategoriesRepositoryInterface(ctl)
		handler := NewCategoryHandler(mockRepo)
		handlerFunc := http.HandlerFunc(handler.HandleCreate)
		handlerFunc.ServeHTTP(recorder, request)

		req.Equal(http.StatusBadRequest, recorder.Code)
		req.Equal(http.Header{"Content-Type": []string{"application/json"}}, recorder.Header())

		responseError := make(map[string]string)
		err = json.Unmarshal(recorder.Body.Bytes(), &responseError)
		req.NoError(err)

		req.Equal(api.BadRequestError, responseError["error"])
	})
}
