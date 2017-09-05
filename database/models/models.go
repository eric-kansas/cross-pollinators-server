package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username  string
	Email     string
	Password  []byte     // bcrypted password
	Interests []Interest // User has many interests
	Projects  []Project  // User has many Projects
}

type Interest struct {
	gorm.Model
	Name   string
	UserID uint
}

type Project struct {
	gorm.Model
	Name        string
	Description string
	Objective   string
	Location    string
	Category    string
	SubCategory string
	Tags        string
	UserID      uint
}

type Comment struct {
	gorm.Model
	Message   string
	CreatedBy User    //FK: user
	Project   Project //FK: Project
}
