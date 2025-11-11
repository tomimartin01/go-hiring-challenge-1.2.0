package models

import (
	"gorm.io/gorm"
)

//go:generate mockgen -source=categories_repository.go -destination=mocks/categories_repository_mock.go -package=models_mocks
type CategoriesRepositoryInterface interface {
	Create(dto CategoryDto) (Category, error)
	GetAllCategories() ([]Category, error)
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoriesRepositoryInterface {
	return &CategoryRepository{
		db: db,
	}
}

func (r CategoryRepository) Create(dto CategoryDto) (Category, error) {
	newCategory := Category{
		Code: dto.Code,
		Name: dto.Name,
	}

	result := r.db.FirstOrCreate(&newCategory, Category{Code: newCategory.Code})
	if result.Error != nil {
		return Category{}, result.Error
	}

	return newCategory, nil
}

func (r CategoryRepository) GetAllCategories() ([]Category, error) {
	var categories []Category

	if err := r.db.Model(&Category{}).
		Find(&categories).Error; err != nil {
		return categories, err
	}

	return categories, nil
}
