package models

import (
	"github.com/jinzhu/gorm"
)

type FutureProgram struct {
	gorm.Model
	Name   string `json:"name"`
	Number int    `json:"number"`
}
