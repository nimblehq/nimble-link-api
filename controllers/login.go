package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/models"
	"github.com/nimble-link/backend/pkg/uuid"
	"github.com/nimble-link/backend/services/authentication"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var conf *oauth2.Config

func Login(c *gin.Context) {
	conf = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH2_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	state, err := uuid.RandToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	session := sessions.Default(c)
	session.Set("state", state)
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	url := getLoginURL(state)

	c.Redirect(http.StatusSeeOther, url)
}

func LoginCallback(c *gin.Context) {
	session := sessions.Default(c)

	retrievedState := session.Get("state")
	queryState := c.Request.URL.Query().Get("state")

	if retrievedState != queryState {
		c.JSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	code := c.Request.URL.Query().Get("code")
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	client := conf.Client(oauth2.NoContext, token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer userInfo.Body.Close()

	data, _ := ioutil.ReadAll(userInfo.Body)
	var user models.User
	json.Unmarshal(data, &user)

	savedUser := models.FindUserByEmail(user.Email)
	if savedUser.ID == 0 {
		user.Save()
		savedUser = &user
	}
	err = authentication.Login(c, savedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}
