package models

import (
	"github.com/mytheresa/go-hiring-challenge/app/filters"
	"gorm.io/gorm"
)

type ProductsRepositoryInterface interface {
	GetAllProducts(filters *filters.PaginationFilter) ([]Product, int64, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductsRepositoryInterface {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(filters *filters.PaginationFilter) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.Model(&Product{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Variants").
		Preload("Categories").
		Offset(filters.Offset).
		Limit(filters.Limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}
