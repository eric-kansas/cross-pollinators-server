package api

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

var hmacSampleSecret = []byte("my_secret_key")

func DoTheThingsHandler(w http.ResponseWriter, req *http.Request) {
	err := verifyUser(req)
	if err != nil {
		fmt.Fprintf(w, "Failed to verify user: %+v", err)
	}
}

func verifyUser(req *http.Request) error {
	// Get token
	var authCookie, err = req.Cookie("auth_token")
	if err != nil || authCookie == nil || authCookie.Value == "" {
		return err
	}
	authToken := authCookie.Value

	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["username"], claims["password"], claims["nbf"])
	} else {
		return err
	}
	return nil
}
