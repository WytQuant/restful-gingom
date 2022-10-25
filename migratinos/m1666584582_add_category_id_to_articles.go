package migratinos

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"restful-gingorm/models"
)

func m1666584582AddCategoryIDToArticles() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1666584582",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(&models.Article{})

			var articles []models.Article
			tx.Unscoped().Find(&articles)
			for _, article := range articles {
				article.CategoryID = 3
				tx.Save(&article)
			}

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&models.Article{}, "category_id")
		},
	}
}
