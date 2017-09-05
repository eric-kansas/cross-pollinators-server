package api

import (
	"fmt"
	"net/http"

	"github.com/eric-kansas/cross-pollinators-server/database"
	"github.com/eric-kansas/cross-pollinators-server/database/models"
	"github.com/jinzhu/gorm"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		logError(w, ErrNotPostRequest)
		return
	}

	err := validRegisterReq(req)
	if err != nil {
		logError(w, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		logError(w, ErrFailedToConnectToDB)
		return
	}
	defer db.Close()

	err = checkUserExistance(db, req)
	if err != nil {
		logError(w, err)
		return
	}

	// Hash password
	password := []byte(req.Form["password"][0])

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		logError(w, err)
		return
	}

	user := models.User{
		Username: req.Form["username"][0],
		Email:    req.Form["email"][0],
		Password: hashedPassword, // bcrypted password
	}

	// Save username with hashed passedword to data base
	if db.NewRecord(user) {
		db.Create(&user)
	}
	fmt.Fprintf(w, "Success")
}

func validRegisterReq(req *http.Request) error {
	req.ParseForm()

	if len(req.Form["username"]) == 0 || len(req.Form["username"][0]) == 0 {
		return ErrNoUsernameProvided
	}

	if len(req.Form["email"]) == 0 || len(req.Form["email"][0]) == 0 {
		return ErrNoEmailProvided
	}

	if len(req.Form["password"]) == 0 || len(req.Form["password"][0]) == 0 {
		return ErrNoPasswordProvided
	}
	return nil
}

func checkUserExistance(db *gorm.DB, req *http.Request) error {
	var user models.User
	db.Where(models.User{Username: req.Form["username"][0]}).First(&user)
	if user.ID != 0 {
		return ErrUsernameAlreadyTaken
	}

	user = models.User{}

	db.Where(models.User{Email: req.Form["email"][0]}).First(&user)
	if user.ID != 0 {
		return ErrEmailAlreadyTaken
	}
	return nil
}
