package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/models"
	"github.com/nimble-link/backend/pkg/linkutils"
	"github.com/nimble-link/backend/services/authentication"
)

type ShortLinkInput struct {
	OriginalURL string `json:"original_url" binding:"required"`
	Alias       string `json:"alias"`
	Password    string `json:"password"`
}

func CreateLink(c *gin.Context) {
	var input ShortLinkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := authentication.GetCurrentUserFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID := user.ID
	originalURL := input.OriginalURL

	if alias := input.Alias; alias != "" {
		if linkutils.IsDuplicateAlias(alias) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate Alias"})
		} else {
			saveLink(c, originalURL, alias, input.Password, userID)
		}
	} else {
		alias := linkutils.GenerateHashFrom(originalURL, userID)

		for linkutils.IsDuplicateAlias(alias) {
			alias = linkutils.GenerateHashFrom(originalURL, userID)
		}

		saveLink(c, originalURL, alias, input.Password, userID)
	}
}

func saveLink(c *gin.Context, OriginalURL string, Alias string, Password string, UserID uint) {
	link := models.Link{OriginalURL: OriginalURL, Alias: Alias, Password: Password, UserID: UserID}
	if errors := link.Save(); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": link})
}
