package models

import (
	"errors"

	"github.com/mytheresa/go-hiring-challenge/app/filters"
	"gorm.io/gorm"
)

//go:generate mockgen -source=products_repository.go -destination=mocks/products_repository_mock.go -package=models_mocks
type ProductsRepositoryInterface interface {
	GetAllProducts(filters *filters.PaginationFilter, categoryFilters *filters.CategoryFilter, productFilters *filters.ProductFilter) ([]Product, int64, error)
	GetProductDetails(productCode string) (Product, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductsRepositoryInterface {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(pageFilters *filters.PaginationFilter, categoryFilters *filters.CategoryFilter, productFilters *filters.ProductFilter) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.Model(&Product{})
	if productFilters.PriceLessThan != nil {
		query = query.Where("products.price <= ?", productFilters.PriceLessThan)
	}

	if categoryFilters.CategoryName != "" {
		query = query.
			Joins("JOIN product_categories ON products.id = product_categories.product_id").
			Joins("JOIN categories ON product_categories.category_id = categories.id").
			Where("categories.name = ?", categoryFilters.CategoryName)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Variants").
		Preload("Categories").
		Offset(pageFilters.Offset).
		Limit(pageFilters.Limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductsRepository) GetProductDetails(productCode string) (Product, error) {
	var product Product
	result := r.db.Model(&Product{}).
		Where("products.code = ?", productCode).
		Preload("Variants").
		Preload("Categories").
		Joins("JOIN product_categories ON products.id = product_categories.product_id").
		Joins("JOIN categories ON product_categories.category_id = categories.id").
		First(&product)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Product{}, errors.New(ProductNotFoundError)
		}
		return Product{}, result.Error
	}

	return product, nil
}
