package authentication

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/database"
	"github.com/nimble-link/backend/models"
)

func SaveCurrentUserToContext(c *gin.Context) {
	user, err := getCurrentUserFromAccessToken(c)
	if err != nil {
		return
	}

	c.Set("current_user", user)
}

func getCurrentUserFromAccessToken(c *gin.Context) (*models.User, error) {
	err := errors.New("invalid Bearer token")
	authorizationHeader := c.Request.Header.Get("Authorization")
	re := regexp.MustCompile(`(^Bearer)\s(.+)$`)
	out := re.FindAllStringSubmatch(authorizationHeader, -1)

	if len(out) == 0 || len(out[0]) != 3 {
		return nil, err
	}

	accessToken := out[0][2]
	token := models.FindByAccessToken(accessToken)

	if token.ID == 0 || !token.IsValid() {
		return nil, errors.New("invalid access token")
	}

	user := new(models.User)

	database.DB.Model(token).Related(user)

	return user, nil
}
