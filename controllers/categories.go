package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"net/http"
	"restful-gingorm/models"
	"strconv"
)

type Categories struct {
	DB *gorm.DB
}

type categoryResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Article []struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	} `json:"article"`
}

type allCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type createCategoryForm struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc" binding:"required"`
}

type updateCategoryForm struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (cg *Categories) FindAll(c *gin.Context) {
	var categories []models.Category
	cg.DB.Order("id desc").Find(&categories)

	serializedCategories := []allCategoryResponse{}
	copier.Copy(&serializedCategories, &categories)

	c.JSON(http.StatusOK, gin.H{
		"categories": serializedCategories,
	})
}

func (cg *Categories) FindOne(c *gin.Context) {
	category, err := cg.findCategoryByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"FindOne's error message": err.Error()})
		return
	}

	serializedCategory := categoryResponse{}
	copier.Copy(&serializedCategory, category)

	c.JSON(http.StatusOK, gin.H{
		"category": serializedCategory,
	})
}

func (cg *Categories) Create(c *gin.Context) {
	var form createCategoryForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message error": err.Error()})
		return
	}

	var category models.Category
	copier.Copy(&category, &form)

	result := cg.DB.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message error": result.Error.Error()})
		return
	}

	serializedCategory := categoryResponse{}
	copier.Copy(&serializedCategory, &category)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Created category successfully",
		"category": serializedCategory,
	})
}

func (cg *Categories) Update(c *gin.Context) {
	var form updateCategoryForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message error": err.Error()})
		return
	}

	category, err := cg.findCategoryByID(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message error": err.Error()})
		return
	}

	var updateCategory models.Category
	copier.Copy(&updateCategory, &form)

	result := cg.DB.Model(&category).Updates(&updateCategory)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error message when update": result.Error.Error()})
		return
	}

	serializedCategory := categoryResponse{}
	copier.Copy(&serializedCategory, category)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Update category successfully",
		"articles": serializedCategory,
	})
}

func (cg *Categories) Delete(c *gin.Context) {
	categoryId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error message": "Category id is not a number",
		})
	}

	cg.DB.Unscoped().Delete(&models.Category{}, categoryId)

	c.JSON(http.StatusOK, gin.H{
		"message": "Delete category successfully",
	})
}

func (cg *Categories) findCategoryByID(c *gin.Context) (*models.Category, error) {
	var category models.Category
	id := c.Param("id")

	result := cg.DB.Preload("Article").First(&category, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &category, nil
}
