package routes

import (
	"github.com/gin-gonic/gin"
	"restful-gingorm/config"
	"restful-gingorm/controllers"
	"restful-gingorm/middleware"
)

func Serve(r *gin.Engine) {
	db := config.GetDB()
	v1 := r.Group("/api/v1")
	authenticate := middleware.Authenticate().MiddlewareFunc()
	authorize := middleware.Authorize()

	authController := controllers.Auth{DB: db}
	authGroup := v1.Group("auth")
	{
		authGroup.POST("/sign-up", authController.Signup)
		authGroup.POST("/sign-in", middleware.Authenticate().LoginHandler)
		authGroup.GET("/profile", authenticate, authController.GetProfile)
		authGroup.PATCH("profile", authenticate, authController.UpdateProfile)
	}

	// Users Route
	usersController := controllers.Users{DB: db}
	usersGroup := v1.Group("users")
	usersGroup.Use(authenticate, authorize)
	{
		usersGroup.GET("", usersController.FindAll)
		usersGroup.POST("", usersController.Create)
		usersGroup.GET("/:id", usersController.FindOne)
		usersGroup.PATCH("/:id", usersController.Update)
		usersGroup.DELETE("/:id", usersController.Delete)
		usersGroup.PATCH("/:id/promote", usersController.Promote)
		usersGroup.PATCH("/:id/demote", usersController.Demote)
	}

	// Articles Route
	articlesController := controllers.Article{DB: db}
	articlesGroup := v1.Group("/articles")
	articlesGroup.GET("", articlesController.FindAll)
	articlesGroup.GET("/:id", articlesController.FindOne)
	articlesGroup.Use(authenticate, authorize)
	{
		articlesGroup.PATCH("/:id", articlesController.Update)
		articlesGroup.POST("", authenticate, articlesController.Create)
		articlesGroup.DELETE("/:id", articlesController.Delete)
	}

	// Categories Route
	categoriesController := controllers.Categories{DB: db}
	categoriesGroup := v1.Group("categories")
	categoriesGroup.GET("", categoriesController.FindAll)
	categoriesGroup.GET("/:id", categoriesController.FindOne)
	categoriesGroup.Use(authenticate, authorize)
	{
		categoriesGroup.PATCH("/:id", categoriesController.Update)
		categoriesGroup.POST("", categoriesController.Create)
		categoriesGroup.DELETE("/:id", categoriesController.Delete)
	}

	dashboardController := controllers.Dashboard{DB: db}
	dashboardGroup := v1.Group("dashboard")
	dashboardGroup.Use(authenticate, authorize)
	{
		dashboardGroup.GET("", dashboardController.GetInfo)
	}
}
