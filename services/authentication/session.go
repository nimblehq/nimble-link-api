package authentication

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/models"
)

func GetCurrentUserFromContext(c *gin.Context) (*models.User, error) {
	err := errors.New("Failed to get current user from context")
	savedUser, exists := c.Get("current_user")
	if !exists {
		return nil, err
	}

	user, ok := savedUser.(*models.User)
	if !ok {
		return nil, err
	}

	return user, nil
}
