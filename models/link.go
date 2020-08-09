package models

import (
	"github.com/jinzhu/gorm"
	"github.com/nimble-link/backend/database"
)

type Link struct {
	gorm.Model
	Source     	string `json:"source" gorm:"column:source;not null"`
	Destination	string `json:"destination" gorm:"column:destination;type:varchar(255);unique;not null"`
	Password   	string `gorm:"column:password;type:varchar(255);;default:NULL"`
	UserID  		uint
}
