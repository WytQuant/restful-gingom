package seed

import (
	"github.com/bxcodec/faker/v3"
	"math/rand"
	"restful-gingorm/config"
	"restful-gingorm/migratinos"
	"restful-gingorm/models"
	"strconv"
)

func Load() {
	db := config.GetDB()

	// Cleaning database
	db.Migrator().DropTable("users", "articles", "categories", "migrations")
	migratinos.Migrate()

	//Add admin
	admin := models.User{
		Email:    "admin@wyt.com",
		Password: "passw0rd",
		Name:     "Wyt Admin",
		Role:     "Admin",
		Avatar:   "https://i.pravatar.cc/100",
	}

	admin.Password = admin.GenerateEncryptedPassword()
	db.Create(&admin)

	// Add normal users
	numOfUsers := 50
	users := make([]models.User, 0, numOfUsers)
	userRoles := [2]string{"Editor", "Member"}

	for i := 1; i <= numOfUsers; i++ {
		user := models.User{
			Name:     faker.Name(),
			Email:    faker.Email(),
			Password: "normalUserPassword",
			Avatar:   "https://i.pravatar.cc/100?" + strconv.Itoa(i),
			Role:     userRoles[rand.Intn(2)],
		}

		user.Password = user.GenerateEncryptedPassword()
		db.Create(&user)
		users = append(users, user)
	}

	// Add categories
	numOfCategories := 20
	categories := make([]models.Category, 0, numOfCategories)

	for i := 1; i <= numOfCategories; i++ {
		category := models.Category{
			Name: faker.Word(),
			Desc: faker.Paragraph(),
		}

		db.Create(&category)
		categories = append(categories, category)
	}

	// Add articles
	numOfArticles := 50
	articles := make([]models.Article, 0, numOfArticles)

	for i := 1; i <= numOfArticles; i++ {
		article := models.Article{
			Title:      faker.Sentence(),
			Excerpt:    faker.Sentence(),
			Body:       faker.Paragraph(),
			Image:      "https://source.unsplash.com/random/300x200?" + strconv.Itoa(i),
			CategoryID: uint(rand.Intn(numOfCategories) + 1),
			UserID:     uint(rand.Intn(numOfUsers) + 1),
		}

		db.Create(&article)
		articles = append(articles, article)
	}

}
