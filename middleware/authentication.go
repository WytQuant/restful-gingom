package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"restful-gingorm/config"
	"restful-gingorm/models"
)

type login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

var identityKey = "sub"

func Authenticate() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key: []byte(os.Getenv("SECRET_KEY")),

		IdentityKey: identityKey,

		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",

		IdentityHandler: func(c *gin.Context) interface{} {
			var user models.User
			claims := jwt.ExtractClaims(c)
			id := claims[identityKey]

			db := config.GetDB()
			if err := db.First(&user, uint(id.(float64))).Error; err != nil {
				return nil
			}
			return &user
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var form login
			var user models.User

			if err := c.ShouldBindJSON(&form); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			db := config.GetDB()
			if err := db.Where("email = ?", form.Email).First(&user).Error; err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &user, nil
		},

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				claims := jwt.MapClaims{
					identityKey: v.ID,
				}

				return claims
			}
			return jwt.MapClaims{}
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"error message": message,
			})
		},
	})

	if err != nil {
		log.Fatalln("JWT Error:" + err.Error())
	}

	return authMiddleware
}