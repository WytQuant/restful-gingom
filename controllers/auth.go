package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"restful-gingorm/models"
)

type Auth struct {
	DB *gorm.DB
}

type authForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type updateProfileForm struct {
	Email  string                `form:"email"`
	Name   string                `form:"name"`
	Avatar *multipart.FileHeader `form:"avatar"`
}

type authResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func (a *Auth) GetProfile(c *gin.Context) {
	sub, _ := c.Get("sub")
	user := sub.(*models.User)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	c.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (a *Auth) Signup(c *gin.Context) {
	var form authForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()
	if err := a.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser authResponse
	copier.Copy(&serializedUser, &user)
	c.JSON(http.StatusCreated, gin.H{"user": serializedUser})
}

func (a *Auth) UpdateProfile(c *gin.Context) {
	var form updateProfileForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error message": err.Error()})
		return
	}

	sub, _ := c.Get("sub")
	user := sub.(*models.User)

	setUserImage(c, user)

	var updateUserData models.User
	copier.Copy(&updateUserData, &form)

	result := a.DB.Model(user).Updates(&updateUserData)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error message": result.Error.Error()})
		return
	}

	var serialziedUser userResponse
	copier.Copy(&serialziedUser, user)
	c.JSON(http.StatusOK, gin.H{
		"user": serialziedUser,
	})
}
