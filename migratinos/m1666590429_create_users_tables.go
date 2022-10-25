package migratinos

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"restful-gingorm/models"
)

func m1666590429CreateUsersTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1666590429",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("users")
		},
	}
}
