package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/database"
	"github.com/nimble-link/backend/models"
	"github.com/nimble-link/backend/services/authentication"
)

func Logout(c *gin.Context) {
	user, err := authentication.GetCurrentUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var tokens []models.Token

	database.DB.Model(user).Related(&tokens)
	for _, token := range tokens {
		token.Delete()
	}

	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}
