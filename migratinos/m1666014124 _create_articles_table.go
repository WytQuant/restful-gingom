package migratinos

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"restful-gingorm/models"
)

func m1666014124CreateArticlesTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1666014124",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Article{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("articles")
		},
	}
}
