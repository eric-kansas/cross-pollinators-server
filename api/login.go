package api

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		fmt.Fprintf(w, "Failed login not a POST")
		return
	}

	req.ParseForm()
	password := []byte(req.Form["password"][0])

	// get HashedPassword form data based keyed off user name
	hashedPassword := []byte("hashed-from-database")
	// Comparing the password with the hash
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		fmt.Fprintf(w, "Failed comparing of hashed passwords: %+v", err)
		return
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Form["username"],
		"nbf":      time.Date(2017, 6, 20, 12, 0, 0, 0, time.UTC).Unix(),
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
}
