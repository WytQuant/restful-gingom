package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"os"
	"restful-gingorm/models"
	"strconv"
	"strings"
)

type Article struct {
	DB *gorm.DB
}

type createArticleForm struct {
	Title      string                `form:"title" binding:"required"`
	Body       string                `form:"body" binding:"required"`
	Excerpt    string                `form:"excerpt" binding:"required"`
	Image      *multipart.FileHeader `form:"image" binding:"required"`
	CategoryID uint                  `form:"categoryId" binding:"required"`
}

type updateArticleForm struct {
	Title      string                `form:"title"`
	Body       string                `form:"body"`
	Excerpt    string                `form:"excerpt"`
	Image      *multipart.FileHeader `form:"image"`
	CategoryID uint                  `form:"categoryId"`
}

type articleResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Excerpt    string `json:"excerpt"`
	Body       string `json:"body"`
	Image      string `json:"image"`
	CategoryID uint   `json:"categoryId"`
	Category   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`

	User struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"user"`
}

type createdOrUpdatedResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Excerpt    string `json:"excerpt"`
	Body       string `json:"body"`
	Image      string `json:"image"`
	CategoryID uint   `json:"categoryId"`
	UserID     uint   `json:"userId"`
}

type articlesPaging struct {
	Items  []articleResponse `json:"item"`
	Paging *pagingResult     `json:"paging"`
}

func (a *Article) FindAll(c *gin.Context) {
	var articles []models.Article

	query := a.DB.Preload("User").Preload("Category").Order("id desc")

	categoryId := c.Query("categoryId")
	if categoryId != "" {
		query = query.Where("category_id = ?", categoryId)
	}

	term := c.Query("term")
	if term != "" {
		query = query.Where("title LIKE ?", "%"+term+"%")
	}

	pagination := pagination{ctx: c, query: query, records: &articles}
	paging := pagination.paginate()

	serializedArticles := []articleResponse{}
	copier.Copy(&serializedArticles, &articles)

	c.JSON(http.StatusOK, gin.H{
		"articles": articlesPaging{
			Items:  serializedArticles,
			Paging: paging,
		},
	})
}

func (a *Article) FindOne(c *gin.Context) {
	article, err := a.findArticleByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"FindOne|error message": err.Error(),
		})
		return
	}

	serializedArticle := articleResponse{}
	copier.Copy(&serializedArticle, article)

	c.JSON(http.StatusOK, gin.H{
		"article": serializedArticle,
	})
}

func (a *Article) Create(c *gin.Context) {
	var form createArticleForm

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error message": err.Error(),
		})
		return
	}

	// form => articles
	var article models.Article
	user, _ := c.Get("sub")
	if err := copier.Copy(&article, &form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error message": err.Error(),
		})
		return
	}

	article.User = *user.(*models.User)

	// articles => db
	result := a.DB.Create(&article)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error message": result.Error.Error(),
		})
		return
	}

	if err := a.setArticleImage(c, &article); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error message": result.Error.Error(),
		})
		return
	}

	serializedArticle := createdOrUpdatedResponse{}
	copier.Copy(&serializedArticle, &article)

	c.JSON(http.StatusCreated, gin.H{
		"article": serializedArticle,
	})

}

func (a *Article) Update(c *gin.Context) {
	var form updateArticleForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error message": err.Error()})
		return
	}

	article, err := a.findArticleByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error message": err.Error()})
		return
	}

	var updateArticle models.Article
	copier.Copy(&updateArticle, &form)

	result := a.DB.Model(&article).Updates(&updateArticle)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error message when update": result.Error.Error()})
		return
	}

	a.setArticleImage(c, article)

	var serializedArticle createdOrUpdatedResponse
	copier.Copy(&serializedArticle, article)

	c.JSON(http.StatusOK, gin.H{
		"articles": serializedArticle,
	})
}

func (a *Article) Delete(c *gin.Context) {
	articleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error message": "Articles id is not a number",
		})
	}

	a.DB.Unscoped().Delete(&models.Article{}, articleId)

	c.JSON(http.StatusOK, gin.H{
		"message": "Delete article successfully",
	})
}

func (a *Article) setArticleImage(c *gin.Context, article *models.Article) error {
	file, err := c.FormFile("image")
	if err != nil || file == nil {
		return err
	}

	if article.Image != "" {
		article.Image = strings.Replace(article.Image, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + article.Image)
	}

	path := "uploads/articles/" + strconv.Itoa(int(article.ID))
	os.MkdirAll(path, 0755)
	fileName := path + "/" + file.Filename
	if err := c.SaveUploadedFile(file, fileName); err != nil {
		return err
	}

	article.Image = os.Getenv("HOST") + "/" + fileName
	a.DB.Save(article)

	return nil
}

func (a *Article) findArticleByID(c *gin.Context) (*models.Article, error) {
	var article models.Article
	id := c.Param("id")

	result := a.DB.Preload("User").Preload("Category").First(&article, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &article, nil
}
