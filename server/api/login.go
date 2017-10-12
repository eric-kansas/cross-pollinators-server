package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eric-kansas/cross-pollinators-server/database"
	"github.com/eric-kansas/cross-pollinators-server/database/models"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		logError(w, ErrNotPostRequest)
		return
	}

	err := validLoginReq(req)
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

	email := req.Form["email"][0]
	password := []byte(req.Form["password"][0])

	user := models.User{}
	db.Where(&models.User{Email: email}).First(&user)
	if user.ID == 0 {
		logError(w, ErrEmailNotFound)
		return
	}

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(user.Password, password)
	if err != nil {
		logError(w, err)
		return
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": req.Form["email"],
		"id":    user.ID,
		"nbf":   time.Date(2017, 6, 20, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		fmt.Fprintf(w, "Failed to sign token error: %+v", err)
		return
	}

	fmt.Println("tokenString:", tokenString)

	cookie1 := &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().UTC().Add(time.Hour * time.Duration(1)),
		HttpOnly: false,
	}
	http.SetCookie(w, cookie1)
	fmt.Fprintf(w, "Success")
}

func validLoginReq(req *http.Request) error {
	req.ParseForm()

	if len(req.Form["email"]) == 0 || len(req.Form["email"][0]) == 0 {
		return ErrNoEmailProvided
	}

	if len(req.Form["password"]) == 0 || len(req.Form["password"][0]) == 0 {
		return ErrNoPasswordProvided
	}
	return nil
}
