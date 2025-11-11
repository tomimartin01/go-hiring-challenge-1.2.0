package models

type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"not null"`
}

type CategoryDto struct {
	Code string
	Name string
}

func (c *Category) TableName() string {
	return "categories"
}

func NewCategoryToCreate(code, name string) CategoryDto {
	return CategoryDto{
		Code: code,
		Name: name,
	}
}
