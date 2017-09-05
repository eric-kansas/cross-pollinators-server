package models

import (
	"log"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	gorm.Model
	Username  string
	Email     string
	Password  []byte     // bcrypted password
	Interests []Interest // User has many interests
	Projects  []Project  // User has many Projects
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	log.Print("Before Create user")
	scope.SetColumn("ID", uuid.NewV4())
	return nil
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
