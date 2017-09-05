package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/eric-kansas/cross-pollinators-server/database/models"
	"github.com/eric-kansas/cross-pollinators-server/server/configs"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		fmt.Fprintf(w, "Failed login not a POST")
		return
	}

	req.ParseForm()
	if len(req.Form["username"]) == 0 {
		fmt.Fprintf(w, "Failed: no username provided")
		return
	}

	db, err := connectToDB()

	if err != nil {
		fmt.Fprintf(w, "Failure: failed to connect to database")
		return
	}

	var user models.User
	db.Not("username = ?", req.Form["username"]).First(&user)
	if user.ID == 0 {
		fmt.Fprintf(w, "Failure: username already in use")
		return
	}

	if len(req.Form["email"]) == 0 {
		fmt.Fprintf(w, "Failed: no email provided")
		return
	}

	if len(req.Form["password"]) == 0 {
		fmt.Fprintf(w, "Failed: no password provided")
		return
	}

	password := []byte(req.Form["password"][0])

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintf(w, "Failed to hash password: %+v", err)
		return
	}

	user = models.User{
		Username: req.Form["username"][0],
		Email:    req.Form["email"][0],
		Password: hashedPassword, // bcrypted password
		// Interests:   nil, // User has many interests
		// Projects:    nil,  // User has many Projects
	}

	// Save username with hashed passedword to data base
	log.Printf("user: %+v \n", user)

	if db.NewRecord(user) {
		db.Create(&user)
	}
	fmt.Fprintf(w, "Success")
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
