package catalog

import (
	"log"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
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

const (
	prodCodeSubPath = "code"
)

func NewCatalogHandler(r models.ProductsRepositoryInterface) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	pageFilters := &filters.PaginationFilter{}
	if err := pageFilters.Parse(r.URL.Query()); err != nil {
		log.Printf("Error parsing page filters: %v", err.Error())
		api.ErrorResponse(w, http.StatusBadRequest, api.InvalidParamValueError)
		return

	}

	categoryFilters := &filters.CategoryFilter{}
	if err := categoryFilters.Parse(r.URL.Query()); err != nil {
		log.Printf("Error parsing category filters: %v", err.Error())
		api.ErrorResponse(w, http.StatusBadRequest, api.InvalidParamValueError)
		return
	}

	productFilters := &filters.ProductFilter{}
	if err := productFilters.Parse(r.URL.Query()); err != nil {
		log.Printf("Error parsing product filters: %v", err.Error())
		api.ErrorResponse(w, http.StatusBadRequest, api.InvalidParamValueError)
		return
	}

	res, total, err := h.repo.GetAllProducts(pageFilters, categoryFilters, productFilters)
	if err != nil {
		log.Printf("Error getting all products: %v", err.Error())
		api.ErrorResponse(w, http.StatusInternalServerError, api.InternalServerError)
		return
	}

	response := Response{
		Products: mapToProducts(res),
		Total:    total,
	}

	api.OKResponse(w, response)
}

func (h *CatalogHandler) HandleGetProductDetails(w http.ResponseWriter, r *http.Request) {
	productCode := r.PathValue(prodCodeSubPath)
	if productCode == "" {
		log.Printf("Product code is required")
		api.ErrorResponse(w, http.StatusBadRequest, api.BadRequestError)
		return
	}

	product, err := h.repo.GetProductDetails(productCode)
	if err != nil {
		if err.Error() == models.ProductNotFoundError {
			log.Printf("Product not found: %v", err.Error())
			api.ErrorResponse(w, http.StatusNotFound, api.NotFoundError)
			return
		}

		log.Printf("Error getting product details: %v", err.Error())
		api.ErrorResponse(w, http.StatusInternalServerError, api.InternalServerError)
		return
	}

	response := ProductDetailsResponse{
		Product: mapToProductDetails(product),
	}

	api.OKResponse(w, response)
}

func mapToProducts(modelsProducts []models.Product) []Product {
	products := make([]Product, len(modelsProducts))
	for i, p := range modelsProducts {
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
	return products
}

func mapToProductDetails(product models.Product) ProductDetails {
	variants := make([]Variant, len(product.Variants))
	for i, variant := range product.Variants {
		variants[i] = Variant{
			Code:  variant.SKU,
			Name:  variant.Name,
			Price: variant.Price.InexactFloat64(),
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

	return ProductDetails{
		Code:       product.Code,
		Categories: categories,
		Price:      product.Price.InexactFloat64(),
		Variants:   variants,
	}
}
