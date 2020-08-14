package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/services/authentication"
)

func GetUserInfo(c *gin.Context) {
	user, err := authentication.GetCurrentUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}
