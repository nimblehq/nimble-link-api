package models

import (
	"github.com/nimble-link/backend/database"
)

type Link struct {
	BaseModel
	OriginalURL string `json:"original_url" gorm:"column:original_url;not null"`
	Alias       string `json:"alias" gorm:"column:alias;type:varchar(255);unique;not null"`
	Password    string `json:"password" gorm:"column:password;type:varchar(255);;default:NULL"`
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
