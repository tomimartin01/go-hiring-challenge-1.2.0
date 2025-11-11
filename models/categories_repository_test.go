package models

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCategoryRepository_Create(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		req := require.New(t)
		db, mock, err := sqlmock.New()
		req.NoError(err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockCategory := CategoryDto{
			Code: "CATEGORY_1",
			Name: "Category 1",
		}

		mock.ExpectQuery(`SELECT *`).
			WithArgs(mockCategory.Code, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name"}).
				AddRow(1, mockCategory.Code, mockCategory.Name))

		mock.ExpectQuery(`INSERT INTO "categories"`).
			WithArgs(mockCategory.Code, mockCategory.Name).
			WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name"}).
				AddRow(1, mockCategory.Code, mockCategory.Name))

		repo := NewCategoryRepository(gormDB)

		createdCategory, err := repo.Create(mockCategory)
		req.NoError(err)

		req.Equal(mockCategory.Code, createdCategory.Code)
		req.Equal(mockCategory.Name, createdCategory.Name)
	})
	t.Run("DB Error", func(t *testing.T) {
		req := require.New(t)
		db, mock, err := sqlmock.New()
		req.NoError(err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		mock.ExpectQuery(`SELECT *`).
			WithArgs("CATEGORY_1", 1).
			WillReturnError(errors.New("error"))

		repo := NewCategoryRepository(gormDB)

		mockCategory := CategoryDto{
			Code: "CATEGORY_1",
			Name: "Category 1",
		}

		createdCategory, err := repo.Create(mockCategory)
		req.Error(err)

		req.EqualError(err, "error")
		req.Equal(Category{}, createdCategory)
	})
}

func TestCategoryRepository_GetAllCategories(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		req := require.New(t)
		db, mock, err := sqlmock.New()
		req.NoError(err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mockCategories := []Category{
			{
				ID:   1,
				Code: "CATEGORY_1",
				Name: "Test Category 1",
			},
			{
				ID:   2,
				Code: "CATEGORY_2",
				Name: "Test Category 2",
			},
		}

		mock.ExpectQuery(`SELECT *`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name"}).
				AddRow(mockCategories[0].ID, mockCategories[0].Code, mockCategories[0].Name).
				AddRow(mockCategories[1].ID, mockCategories[1].Code, mockCategories[1].Name))

		repo := NewCategoryRepository(gormDB)

		categories, err := repo.GetAllCategories()
		req.NoError(err)

		req.Equal(2, len(categories))
		req.Equal(mockCategories[0].Code, categories[0].Code)
		req.Equal(mockCategories[0].Name, categories[0].Name)
		req.Equal(mockCategories[1].Code, categories[1].Code)
		req.Equal(mockCategories[1].Name, categories[1].Name)
	})

	t.Run("DB Error", func(t *testing.T) {
		req := require.New(t)
		db, mock, err := sqlmock.New()
		req.NoError(err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		mock.ExpectQuery(`SELECT *`).
			WillReturnError(errors.New("error"))

		repo := NewCategoryRepository(gormDB)

		categories, err := repo.GetAllCategories()
		req.Error(err)
		req.EqualError(err, "error")
		req.Equal(0, len(categories))
	})
}
