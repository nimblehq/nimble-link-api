package models

import (
	"time"

	"github.com/nimble-link/backend/database"
)

type Token struct {
	BaseModel
	AccessToken           string    `json:"access_token" gorm:"column:access_token;"`
	RefreshToken          string    `json:"refresh_token" gorm:"column:refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at" gorm:"column:access_token_expires_at;not_null"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at" gorm:"column:refresh_token_expires_at;not_null"`
	User                  User      `json:"-"`
	UserID                uint      `json:"user_id"`
}

func (token *Token) IsValid() bool {
	return token.AccessTokenExpiresAt.Unix() > time.Now().Unix() && token.DeletedAt == nil
}

func (token *Token) Delete() {
	database.DB.Delete(token)
}

func FindByAccessToken(accessToken string) *Token {
	token := new(Token)
	database.DB.Where("access_token = ?", accessToken).First(token)

	return token
}

func FindByRefreshToken(refreshToken string) *Token {
	token := new(Token)
	database.DB.Where("refresh_token = ?", refreshToken).First(token)

	return token
}

func (token *Token) IsRefreshable() bool {
	return token.RefreshTokenExpiresAt.Unix() > time.Now().Unix() && token.DeletedAt == nil
}
