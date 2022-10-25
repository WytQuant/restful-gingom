package migratinos

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"restful-gingorm/models"
)

func m1666014124CreateCategoriesTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1666577721",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Category{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("categories")
		},
	}
}
