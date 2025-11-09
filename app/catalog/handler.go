package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/filters"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total_products"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Product struct {
	Code       string     `json:"code"`
	Categories []Category `json:"categories"`
	Price      float64    `json:"price"`
}

type CatalogHandler struct {
	repo models.ProductsRepositoryInterface
}

func NewCatalogHandler(r models.ProductsRepositoryInterface) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	pageFilters := &filters.PaginationFilter{}
	if err := pageFilters.Parse(r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	categoryFilters := &filters.CategoryFilter{}
	if err := categoryFilters.Parse(r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	productFilters := &filters.ProductFilter{}
	if err := productFilters.Parse(r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, total, err := h.repo.GetAllProducts(pageFilters, categoryFilters, productFilters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	products := make([]Product, len(res))
	for i, p := range res {
		categories := make([]Category, len(p.Categories))
		for j, c := range p.Categories {
			categories[j] = Category{
				Code: c.Code,
				Name: c.Name,
			}
		}
		products[i] = Product{
			Code:       p.Code,
			Categories: categories,
			Price:      p.Price.InexactFloat64(),
		}
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Products: products,
		Total:    total,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
