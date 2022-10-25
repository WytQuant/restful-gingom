package migratinos

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"log"
	"restful-gingorm/config"
)

func Migrate() {
	db := config.GetDB()
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		m1666014124CreateArticlesTable(),
		m1666014124CreateCategoriesTable(),
		m1666584582AddCategoryIDToArticles(),
		m1666590429CreateUsersTable(),
		m1666584582AddUserIDToArticles(),
	})

	if err := m.Migrate(); err != nil {
		log.Fatalln("Failed to migrate: ", err)
	}
}
