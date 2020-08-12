package models

import (
	"time"

	"github.com/nimble-link/backend/database"
	"github.com/nimble-link/backend/pkg/token_generator"
)

type User struct {
	BaseModel
	Sub     string  `json:"sub" gorm:"column:sub;type:varchar(255);not null"`
	Picture string  `json:"picture" gorm:"column:picture;type:varchar(255);not null"`
	Email   string  `json:"email" gorm:"column:email;type:varchar(255);unique;default:NULL" sql:"index"`
	Tokens  []Token `json:"-"`
	Links   []Link  `json:"-"`
}

func FindUserByEmail(userEmail string) *User {
	var user = new(User)

	database.DB.Where("email = ?", userEmail).First(user)

	return user
}

func FindUserById(id uint) *User {
	var user = new(User)

	database.DB.Where("id = ?", id).First(user)

	return user
}

func (user *User) Save() {
	database.DB.Create(user)
}

func (user *User) GenerateAccessToken() *Token {
	token := &Token{
		User:                  *user,
		AccessToken:           token_generator.Generate(),
		RefreshToken:          token_generator.Generate(),
		AccessTokenExpiresAt:  time.Now().Add(time.Hour * 12),
		RefreshTokenExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	database.DB.Create(token)

	return token
}
