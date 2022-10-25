package middleware

import (
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
	"restful-gingorm/models"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("sub")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error message": "Unauthorized",
			})
			return
		}

		enforcer := casbin.NewCachedEnforcer("config/acl_model.conf", "config/policy.csv")
		ok = enforcer.Enforce(user.(*models.User), c.Request.URL.Path, c.Request.Method)

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error message": "You are not allowed to access this resource",
			})
			return
		}

		c.Next()
	}
}
