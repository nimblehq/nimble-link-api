package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/database"
	"github.com/nimble-link/backend/models"
	"github.com/nimble-link/backend/services/authentication"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type oauth2Form struct {
	Code         string `form:"code"`
	CodeVerifier string `form:"code_verifier"`
}

var conf *oauth2.Config

func OAuth2Handler(c *gin.Context) {
	var form oauth2Form

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	savedUser, err := exchangeCode(form.Code, form.CodeVerifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	token := savedUser.GenerateAccessToken()

	c.JSON(http.StatusOK, token)
}

func Logout(c *gin.Context) {
	user, err := authentication.GetCurrentUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var tokens []models.Token

	database.DB.Model(user).Related(&tokens)
	for _, token := range tokens {
		database.DB.Unscoped().Delete(token) // Delete record permanently
	}

	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func exchangeCode(code string, codeVerifier string) (*models.User, error) {
	if conf == nil {
		conf = &oauth2.Config{
			ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OAUTH2_REDIRECT_URL"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		}
	}

	codeVerifierOption := oauth2.SetAuthURLParam("code_verifier", codeVerifier)
	token, err := conf.Exchange(oauth2.NoContext, code, codeVerifierOption)
	if err != nil {
		return nil, err
	}
	client := conf.Client(oauth2.NoContext, token)

	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer userInfo.Body.Close()

	data, _ := ioutil.ReadAll(userInfo.Body)
	user := new(models.User)
	json.Unmarshal(data, user)

	savedUser := models.FindUserByEmail(user.Email)
	if savedUser.ID == 0 {
		user.Save()
		savedUser = user
	}

	return savedUser, nil
}
