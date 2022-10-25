package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Dashboard struct {
	DB *gorm.DB
}

type dashboardArticle struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Image   string `json:"image"`
}

type dashboardResponse struct {
	LatestArticles []dashboardArticle `json:"latestArticles"`
	UsersCount     []struct {
		Role  string `json:"role"`
		Count uint   `json:"count"`
	} `json:"usersCount"`
	CategoriesCount uint `json:"categoriesCount"`
	ArticlesCount   uint `json:"articlesCount"`
}

func (d *Dashboard) GetInfo(ctx *gin.Context) {
	res := dashboardResponse{}
	var articlesCount int64
	var categoriesCount int64
	d.DB.Table("articles").Order("id desc").Limit(5).Find(&res.LatestArticles)
	d.DB.Table("articles").Count(&articlesCount)
	d.DB.Table("categories").Count(&categoriesCount)
	d.DB.Table("users").Select("role, count(*) as count").Group("role").Scan(&res.UsersCount)

	res.ArticlesCount = uint(articlesCount)
	res.CategoriesCount = uint(categoriesCount)

	ctx.JSON(http.StatusOK, gin.H{"dashboard": &res})
}
