package models

import (
	"github.com/jinzhu/gorm"
	"github.com/nimble-link/backend/database"
)

type User struct {
	gorm.Model
	Sub     string `json:"sub" gorm:"column:sub;type:varchar(255);not null"`
	Picture string `json:"picture" gorm:"column:picture;type:varchar(255);not null"`
	Email   string `json:"email" gorm:"column:email;type:varchar(255);unique;default:NULL" sql:"index"`
	Links 	[]Link
}

func FindUserByEmail(userEmail string) *User {
	var user = new(User)

	database.DB.Where("email = ?", userEmail).First(user)

	return user
}

func (user *User) Save() {
	database.DB.Create(user)
}
