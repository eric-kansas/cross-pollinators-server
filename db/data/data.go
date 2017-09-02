package data

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	ID          uint
	Email       string
	Password    string
	ProfilePic  []byte
	Description string
	Interests   string
	Location    string
}

type Project struct {
	gorm.Model
	ID          uint
	Name        string
	Description string
	Objective   string
	Location    string
	Photos      []byte
	Category    string
	SubCategory string
	OwnerID     string //FK: user
	OrgName     string
	Tags        string
}

type Comment struct {
	gorm.Model
	ID        uint
	CreatedBy string //FK: user
	Message   string
	Projectid string //FK: Project
}

type Story struct {
	gorm.Model
	StoryID      string
	StoryContent string
	CreatedBy    string // FK: User
}
