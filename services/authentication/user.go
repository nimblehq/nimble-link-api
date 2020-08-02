package authentication

import (
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/models"
)

func Login(c *gin.Context, user *models.User) error {
	session := sessions.Default(c)

	session.Set("current_user", user.Email)
	return session.Save()
}

func SaveCurrentUserToContext(c *gin.Context) *models.User {
	user, err := getCurrentUserFromSession(c)
	if err != nil {
		return nil
	}

	c.Keys["current_user"] = user

	return user
}

func getCurrentUserFromSession(c *gin.Context) (*models.User, error) {
	session := sessions.Default(c)

	userEmail := session.Get("current_user")
	if userEmail == nil {
		return nil, errors.New("Invalid session")
	}
	user := models.FindUserByEmail(userEmail.(string))

	return user, nil
}
