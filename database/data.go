package database

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email       string // user has one Email
	Password    string // bcrypted password
	Description string
	Interests   []Intrest // User has many interests
	Projects    []Project // User has many Projects
	Location    string
}

type Email struct {
	OwnerID    User
	Email      string
	Subscribed bool
}

type Intrest struct {
	gorm.Model
	Name string
}

type Project struct {
	gorm.Model
	Name        string
	Description string
	Objective   string
	Location    string
	Photos      []byte // figure out images
	Category    string
	SubCategory string
	OwnerID     User //FK: user
	OrgName     string
	Tags        string
}

type Comment struct {
	gorm.Model
	ID        uint
	Message   string
	CreatedBy User    //FK: user
	Project   Project //FK: Project
}
