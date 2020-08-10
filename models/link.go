package models

import (
	"github.com/jinzhu/gorm"
	"github.com/nimble-link/backend/database"
)

type Link struct {
	gorm.Model
	OriginalURL string `gorm:"column:original_url;not null"`
	Alias       string `gorm:"column:alias;type:varchar(255);unique;not null"`
	Password    string `gorm:"column:password;type:varchar(255);;default:NULL"`
	UserID      uint   `json:"-"`
}

func (link *Link) Save() []error {
	errors := database.DB.Create(link).GetErrors()
	if errors != nil && len(errors) > 0 {
		return errors
	}
	return nil
}

func FindLinkByAlias(alias string) *Link {
	var link = new(Link)
	database.DB.Where("alias = ?", alias).First(link)
	return link
}
