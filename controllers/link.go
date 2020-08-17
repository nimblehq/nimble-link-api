package controllers

import (
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

	user, err := authentication.GetCurrentUserFromContext(c)
	if err != nil || user == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID := user.ID
	originalURL := input.OriginalURL

	if alias := input.Alias; alias != "" {
		if linkutils.IsDuplicateAlias(alias) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This alias is already used"})
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

	database.DB.Delete(link)

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
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL"))
		return
	}

	link.IncrementCounter()

	c.Redirect(http.StatusTemporaryRedirect, link.OriginalURL)
}

func GetLinkWithPassword(c *gin.Context) {
	alias := c.Param("alias")
	password := c.PostForm("password")

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

	if link.Password != password {
		c.JSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
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

	if input.Password == "" && input.Alias == "" {
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

	if input.Password != "" {
		link.Password = input.Password
	}

	if input.Alias != "" {
		if linkutils.IsDuplicateAlias(input.Alias) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This alias is already used"})
			return
		}
		link.Alias = input.Alias
	}

	database.DB.Save(link)

	c.JSON(http.StatusOK, link)
}

func saveLink(c *gin.Context, OriginalURL string, Alias string, Password string, UserID uint) {
	link := models.Link{OriginalURL: OriginalURL, Alias: Alias, Password: Password, UserID: UserID}
	if errors := link.Save(); errors != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": errors})
		return
	}

	c.JSON(http.StatusOK, link)
}
