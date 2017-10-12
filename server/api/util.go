package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/eric-kansas/cross-pollinators-server/database/models"
	"github.com/jinzhu/gorm"
)

// Move to config / .env
var hmacSampleSecret = []byte("my_secret_key")

// Errors
var (
	ErrNotPostRequest       = errors.New("Request method should be POST")
	ErrNoPasswordProvided   = errors.New("No password provided")
	ErrNoEmailProvided      = errors.New("No email provided")
	ErrFailedToConnectToDB  = errors.New("Failed to connect to database")
	ErrUsernameAlreadyTaken = errors.New("Username is already taken")
	ErrEmailAlreadyTaken    = errors.New("Email is already taken")
	ErrEmailNotFound        = errors.New("Email was not found")
	ErrParsingAuthToken     = errors.New("Error parsing auth token")
)

func AuthWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := GetUserID(r)
		if err != nil {
			fmt.Fprintf(w, "Failed to verify user: %+v", err)
			return
		}

		h.ServeHTTP(w, r) // call original
	})
}

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

func GetUserID(req *http.Request) (string, error) {
	// Get token
	var authCookie, err = req.Cookie("auth_token")
	if err != nil || authCookie == nil || authCookie.Value == "" {
		return "", err
	}
	authToken := authCookie.Value

	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if err != nil {
		return "", ErrParsingAuthToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["username"], claims["id"], claims["nbf"])
		return claims["id"].(string), nil
	}
	return "", err
}

func logError(w http.ResponseWriter, err error) {
	log.Printf("Failed with error: %+v", err)
	fmt.Fprintf(w, "Failed with error: %+v", err)
}
