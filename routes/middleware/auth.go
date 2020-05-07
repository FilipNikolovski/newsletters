package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/csrf"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
	log "github.com/sirupsen/logrus"
)

// Authorization header prefixes.
const (
	APIKeyAuth = "Api-Key"
)

// SetUser fetches the token and then from the token fetches the user entity
// and sets it to the context.
func SetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var authHeader = c.GetHeader("Authorization")

		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 {
				c.Next()
				return
			}

			if parts[0] == APIKeyAuth {
				key, err := storage.GetAPIKey(c, parts[1])
				if err != nil {
					log.WithError(err).Error("unable to fetch api key")
					c.Next()
					return
				}

				c.Set("user", &key.User)
			}

			// When using api keys it's ok to skip the csrf token
			// since we are not using cookies to authenticate the user
			c.Request = csrf.UnsafeSkipCheck(c.Request)
			c.Next()
			return
		}

		session := sessions.Default(c)
		v := session.Get("sess_id")
		if v == nil {
			c.Next()
			return
		}
		sessID := v.(string)
		s, err := storage.GetSession(c, sessID)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user", &s.User)

		c.Next()
	}
}

// GetUser returns the user set in the context
func GetUser(c *gin.Context) *entities.User {
	val, ok := c.Get("user")
	if !ok {
		return nil
	}

	user, ok := val.(*entities.User)
	if !ok {
		return nil
	}

	return user
}

// Authorized is a middleware that checks if the user is authorized to do the
// requested action.
func Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
			c.Abort()
			return
		}
		_, ok = val.(*entities.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
