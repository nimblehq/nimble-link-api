package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/database"
	"github.com/nimble-link/backend/models"
	"github.com/nimble-link/backend/pkg/linkutils"
	"github.com/nimble-link/backend/services/authentication"
)

type ShortLinkInput struct {
	OriginalURL string `json:"original_url" binding:"required"`
	Alias       string `json:"alias"`
	Password    string `json:"password"`
}

type UpdateLinkInput struct {
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

func CreateLink(c *gin.Context) {
	var input ShortLinkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	originalURL := input.OriginalURL

	if !linkutils.IsValidUrl(originalURL) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid URL"})
		return
	}

	user, err := authentication.GetCurrentUserFromContext(c)
	if err != nil || user == nil {
		generateRandomLinkAndSave(c, originalURL, "", 0)
		return
	}

	userID := user.ID
	password := input.Password

	if alias := input.Alias; alias != "" {
		if linkutils.IsDuplicateAlias(alias) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This alias is already used"})
		} else {
			saveLink(c, originalURL, alias, password, userID)
		}
	} else {
		generateRandomLinkAndSave(c, originalURL, password, userID)
	}
}

func GetLinks(c *gin.Context) {
	user, err := authentication.GetCurrentUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var links []models.Link
	database.DB.Model(user).Related(&links)

	c.JSON(http.StatusOK, links)
}

func DeleteLink(c *gin.Context) {
	user, _ := authentication.GetCurrentUserFromContext(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	var link = new(models.Link)
	database.DB.Where(&models.Link{
		UserID: user.ID,
		BaseModel: models.BaseModel{
			ID: uint(id),
		},
	}).First(link)

	if link.ID == 0 {
		c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	database.DB.Unscoped().Delete(link) // Delete record permanently

	c.JSON(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

func GetLink(c *gin.Context) {
	alias := c.Param("alias")
	if alias == "" {
		c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	var link = new(models.Link)
	database.DB.Where(&models.Link{
		Alias: alias,
	}).First(link)

	if link.ID == 0 {
		c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if link.Password != "" {
		c.Redirect(
			http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/go?alias=%s", os.Getenv("FRONTEND_URL"), link.Alias),
		)
		return
	}

	link.IncrementCounter()

	c.Redirect(http.StatusTemporaryRedirect, link.OriginalURL)
}

type GetLinkWithPasswordInput struct {
	Password string `json:"password" binding:"required"`
}

func GetLinkWithPassword(c *gin.Context) {
	var input GetLinkWithPasswordInput
	alias := c.Param("alias")

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	if alias == "" {
		c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	var link = new(models.Link)
	database.DB.Where(&models.Link{
		Alias: alias,
	}).First(link)

	if link.ID == 0 {
		c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if link.Password != input.Password {
		c.JSON(http.StatusForbidden, http.StatusText(http.StatusUnauthorized))
		return
	}

	link.IncrementCounter()

	c.JSON(http.StatusOK, link)
}

func UpdateLink(c *gin.Context) {
	var input UpdateLinkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	if input.Alias == "" {
		c.JSON(http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
		return
	}

	user, _ := authentication.GetCurrentUserFromContext(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	var link = new(models.Link)
	database.DB.Where(&models.Link{
		UserID: user.ID,
		BaseModel: models.BaseModel{
			ID: uint(id),
		},
	}).First(link)

	if link.ID == 0 {
		c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if input.Password != link.Password {
		link.Password = input.Password
	}

	if input.Alias != link.Alias {
		if linkutils.IsDuplicateAlias(input.Alias) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This alias is already used"})
			return
		}
		link.Alias = input.Alias
	}

	database.DB.Save(link)

	c.JSON(http.StatusOK, link)
}

func generateRandomLinkAndSave(c *gin.Context, originalURL string, password string, userID uint) {
	alias := linkutils.GenerateHashFrom(originalURL, userID)

	for linkutils.IsDuplicateAlias(alias) {
		alias = linkutils.GenerateHashFrom(originalURL, userID)
	}

	saveLink(c, originalURL, alias, password, userID)
}

func saveLink(c *gin.Context, originalURL string, alias string, password string, userID uint) {
	link := models.Link{OriginalURL: originalURL, Alias: alias, Password: password, UserID: userID}
	if errors := link.Save(); errors != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": errors})
		return
	}

	c.JSON(http.StatusOK, link)
}
