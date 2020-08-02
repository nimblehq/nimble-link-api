package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	session := sessions.Default(c)

	session.Delete("state")
	session.Delete("current_user")

	session.Save()
	c.Keys["current_user"] = nil

	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}
