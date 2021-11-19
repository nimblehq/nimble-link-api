package controllers

import (
	"context"
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
	"google.golang.org/api/idtoken"
)

type oauth2CodeForm struct {
	Code         string `form:"code"`
	CodeVerifier string `form:"code_verifier"`
}

type oauth2IdTokenForm struct {
	IdToken string `form:"id_token"`
}

var conf *oauth2.Config

func OAuth2CodeHandler(c *gin.Context) {
	var form oauth2CodeForm

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

func OAuth2IdTokenHandler(c *gin.Context) {
	var form oauth2IdTokenForm

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	savedUser, err := verifyIdToken(form.IdToken)

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

func verifyIdToken(idToken string) (*models.User, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, os.Getenv("OAUTH2_IOS_CLIENT_ID"))
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(payload.Claims)
	if err != nil {
		return nil, err
	}

	user := new(models.User)

	err = json.Unmarshal(data, user)
	if err != nil {
		return nil, err
	}

	savedUser := models.FindUserByEmail(user.Email)
	if savedUser.ID == 0 {
		user.Save()
		savedUser = user
	}

	return savedUser, err
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

	err = json.Unmarshal(data, user)
	if err != nil {
		return nil, err
	}

	savedUser := models.FindUserByEmail(user.Email)
	if savedUser.ID == 0 {
		user.Save()
		savedUser = user
	}

	return savedUser, nil
}
