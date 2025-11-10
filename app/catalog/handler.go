package catalog

import (
	"encoding/json"
	"fmt"
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

type ProductDetails struct {
	Code       string     `json:"code"`
	Categories []Category `json:"categories"`
	Price      float64    `json:"price"`
	Variants   []Variant  `json:"variants"`
}

type Variant struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductDetailsResponse struct {
	Product ProductDetails `json:"product"`
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

func (h *CatalogHandler) HandleGetProductDetails(w http.ResponseWriter, r *http.Request) {
	productCode := r.PathValue("code")
	if productCode == "" {
		http.Error(w, "Product code is required", http.StatusBadRequest)
		return
	}

	product, err := h.repo.GetProductDetails(productCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	variants := make([]Variant, len(product.Variants))
	for i, variant := range product.Variants {
		variants[i] = Variant{
			Code: variant.SKU,
			Name: variant.Name,
			Price: product.Price.InexactFloat64(),
		}
		// We need to inherit the price from the product
		if variant.Price.IsZero() {
			variants[i].Price = product.Price.InexactFloat64()
		}
	}
	categories := make([]Category, len(product.Categories))
	for i, category := range product.Categories {
		categories[i] = Category{
			Code: category.Code,
			Name: category.Name,
		}
	}
	productDetails := ProductDetails{
		Code:       product.Code,
		Categories: categories,
		Price:      product.Price.InexactFloat64(),
		Variants:   variants,
	}

	response := ProductDetailsResponse{
		Product: productDetails,
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	fmt.Println(product)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
