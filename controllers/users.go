package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"net/http"
	"os"
	"restful-gingorm/config"
	"restful-gingorm/models"
	"strconv"
	"strings"
)

type Users struct {
	DB *gorm.DB
}

type createUserForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type updateUserForm struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=8"`
	Name     string `json:"name"`
}

type userResponse struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Role   string `json:"role"`
}

type usersPaging struct {
	Items  []userResponse `json:"items"`
	Paging *pagingResult  `json:"paging"`
}

func (u *Users) FindAll(c *gin.Context) {
	var users []models.User
	query := u.DB.Order("id desc").Find(&users)

	term := c.Query("term")
	if term != "" {
		query = query.Where("name LIKE ?", "%"+term+"%")
	}

	pagination := pagination{ctx: c, query: query, records: &users}
	paging := pagination.paginate()

	serializedUsers := []userResponse{}
	copier.Copy(&serializedUsers, &users)
	c.JSON(http.StatusOK, gin.H{
		"users": usersPaging{Items: serializedUsers, Paging: paging},
	})
}

func (u *Users) FindOne(c *gin.Context) {
	user, err := u.findUserByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	c.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) Create(c *gin.Context) {
	var form createUserForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()

	if err := u.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	c.JSON(http.StatusCreated, gin.H{"user": serializedUser})
}

func (u *Users) Update(c *gin.Context) {
	var form updateUserForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := u.findUserByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if form.Password != "" {
		user.Password = user.GenerateEncryptedPassword()
	}

	var updateUserData models.User
	copier.Copy(&updateUserData, &form)

	if err := u.DB.Model(&user).Updates(&updateUserData).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	c.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) Delete(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))

	u.DB.Unscoped().Delete(&models.User{}, uint(userId))

	c.JSON(http.StatusOK, gin.H{"message": "A user has already been deleted from database"})
}

func (u *Users) Promote(c *gin.Context) {
	user, err := u.findUserByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Promote()
	u.DB.Save(user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	c.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) Demote(c *gin.Context) {
	user, err := u.findUserByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Demote()
	u.DB.Save(user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	c.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) findUserByID(c *gin.Context) (*models.User, error) {
	id := c.Param("id")
	var user models.User
	if err := u.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func setUserImage(ctx *gin.Context, user *models.User) error {
	file, _ := ctx.FormFile("avatar")
	if file == nil {
		return nil
	}

	if user.Avatar != "" {
		user.Avatar = strings.Replace(user.Avatar, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + user.Avatar)
	}

	path := "uploads/users/" + strconv.Itoa(int(user.ID))
	os.MkdirAll(path, os.ModePerm)
	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}

	db := config.GetDB()
	user.Avatar = os.Getenv("HOST") + "/" + filename
	db.Save(user)

	return nil
}
