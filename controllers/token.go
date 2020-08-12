package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/database"
	"github.com/nimble-link/backend/models"
	"github.com/nimble-link/backend/pkg/token_generator"
)

type refreshTokenForm struct {
	RefreshToken string `json:"refresh_token"`
}

func RefreshToken(c *gin.Context) {
	var form refreshTokenForm

	if err := c.ShouldBindJSON(&form); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	token := models.FindByRefreshToken(form.RefreshToken)

	if token.ID == 0 || !token.IsRefreshable() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token.AccessToken = token_generator.Generate()
	token.AccessTokenExpiresAt = time.Now().Add(time.Hour * 12)

	database.DB.Save(token)

	c.JSON(http.StatusOK, token)
}
