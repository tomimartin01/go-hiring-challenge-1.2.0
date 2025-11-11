package category

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type Category struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type CreateRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoryHandler struct {
	repo models.CategoriesRepositoryInterface
}

func NewCategoryHandler(r models.CategoriesRepositoryInterface) CategoryHandler {
	return CategoryHandler{
		repo: r,
	}
}

func (h CategoryHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	categoryDto := models.NewCategoryToCreate(req.Code, req.Name)
	newCategory, err := h.repo.Create(categoryDto)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := Category{
		ID:   newCategory.ID,
		Code: newCategory.Code,
		Name: newCategory.Name,
	}

	api.OKResponse(w, res)
}

func (h CategoryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	res, err := h.repo.GetAllCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, mapToCategories(res))
}

func mapToCategories(modelsCategories []models.Category) []Category {
	categories := make([]Category, len(modelsCategories))
	for i, category := range modelsCategories {
		categories[i] = Category{
			Code: category.Code,
			Name: category.Name,
		}
	}
	return categories
}
