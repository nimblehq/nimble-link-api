package models

import (
	"github.com/jinzhu/gorm"
	"github.com/nimble-link/backend/database"
)

type Link struct {
	gorm.Model
	OriginalURL string `json:"omitempty" gorm:"column:original_url;not null"`
	CustomAlias string `gorm:"column:custom_alias;type:varchar(255);unique;not null"`
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

func FindLinkByCustomAlias(customAlias string) *Link {
	var link = new(Link)
	database.DB.Where("custom_alias = ?", customAlias).First(link)
	return link
}
