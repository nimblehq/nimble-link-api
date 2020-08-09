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
	CustomAlias string `json:"custom_alias"`
	Password    string `json:"password"`
}

func CreateLink(c *gin.Context) {

	var input ShortLinkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// get user info
	user, err := authentication.GetCurrentUserFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID := user.ID
	originalURL := input.OriginalURL

	if customAlias := input.CustomAlias; customAlias != "" {
		// save user's CustomAlias
		if linkutils.IsDuplicateCustomAlias(customAlias) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate Custom Alias"})
		} else {
			saveLink(c, originalURL, customAlias, input.Password, userID)
		}
	} else {
		// generate random CustomAlias
		customAlias := linkutils.GenerateHashFrom(originalURL, userID)

		for linkutils.IsDuplicateCustomAlias(customAlias) {
			customAlias = linkutils.GenerateHashFrom(originalURL, userID)
		}

		saveLink(c, originalURL, customAlias, input.Password, userID)
	}
}

func saveLink(c *gin.Context, OriginalURL string, CustomAlias string, Password string, UserID uint) {
	link := models.Link{OriginalURL: OriginalURL, CustomAlias: CustomAlias, Password: Password, UserID: UserID}
	if errors := link.Save(); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": link})
}
