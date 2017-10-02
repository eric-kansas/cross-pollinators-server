package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username     string
	Email        string
	Organization string
	Password     []byte // bcrypted password
	AvatarURL    string
	Interests    []Interest // User has many interests
	Projects     []Project  // User has many Projects
	//Followers []User
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
