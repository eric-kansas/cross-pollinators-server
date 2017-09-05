package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/eric-kansas/cross-pollinators-server/database/models"
	"github.com/eric-kansas/cross-pollinators-server/server/configs"

	"golang.org/x/crypto/bcrypt"
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
)

func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	err := validateReq(req)
	if err != nil {
		log.Printf("Failed to register: %s", err)
		fmt.Fprintf(w, "Failed to register: %s", err)
		return
	}

	db, err := connectToDB()
	if err != nil {
		log.Printf("Failed to register: %s", ErrFailedToConnectToDB)
		fmt.Fprintf(w, "Failed to register: %s", ErrFailedToConnectToDB)
		return
	}

	err = checkUserExistance(db, req)
	if err != nil {
		log.Printf("Failed to register: %s", err)
		fmt.Fprintf(w, "Failed to register: %s", err)
		return
	}

	// Hash password
	password := []byte(req.Form["password"][0])

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %+v", err)
		fmt.Fprintf(w, "Failed to hash password: %+v", err)
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

func validateReq(req *http.Request) error {
	if req.Method != "POST" {
		return ErrNotPostRequest
	}

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

func connectToDB() (*gorm.DB, error) {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", "db", configs.Data.DbUser, configs.Data.DbPass, configs.Data.DbName)
	db, err := gorm.Open("postgres", dbinfo)

	if err != nil {
		log.Printf("Failed to open connection to postgres database: %+v \n", err)
		return nil, err
	}

	log.Printf("Connected to Cross Pollinators DB!")
	return db, nil
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
