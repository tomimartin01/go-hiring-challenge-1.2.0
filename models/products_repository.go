package models

import (
	"github.com/mytheresa/go-hiring-challenge/app/filters"
	"gorm.io/gorm"
)

type ProductsRepositoryInterface interface {
	GetAllProducts(filters *filters.PaginationFilter, categoryFilters *filters.CategoryFilter, productFilters *filters.ProductFilter) ([]Product, int64, error)
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
