package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/eric-kansas/cross-pollinators-server/database/models"
	"github.com/jinzhu/gorm"
)

// Errors
var (
	ErrNotPostRequest       = errors.New("Request method should be POST")
	ErrNoUsernameProvided   = errors.New("No username provided")
	ErrNoPasswordProvided   = errors.New("No password provided")
	ErrNoEmailProvided      = errors.New("No email provided")
	ErrFailedToConnectToDB  = errors.New("Failed to connect to database")
	ErrUsernameAlreadyTaken = errors.New("Username is already taken")
	ErrEmailAlreadyTaken    = errors.New("Email is already taken")
	ErrUsernameNotFound     = errors.New("Username was not found")
)

// Helper functions
func checkUsernameExists(db *gorm.DB, req *http.Request) error {
	var user models.User
	db.Where(models.User{Username: req.Form["username"][0]}).First(&user)
	if user.ID != 0 {
		return ErrUsernameAlreadyTaken
	}
	return nil
}

func checkEmailExists(db *gorm.DB, req *http.Request) error {
	var user models.User
	db.Where(models.User{Email: req.Form["email"][0]}).First(&user)
	if user.ID != 0 {
		return ErrEmailAlreadyTaken
	}
	return nil
}

func logError(w http.ResponseWriter, err error) {
	log.Printf("Failed to register: %+v", err)
	fmt.Fprintf(w, "Failed to register: %+v", err)
}
