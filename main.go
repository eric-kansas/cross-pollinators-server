package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/eric-kansas/cross-pollinators-server/configs"
)

var hmacSampleSecret = []byte("my_secret_key")

func init() {
	configs.Initialize()
}

func main() {
	fmt.Printf("Server started with go version %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	setupServer()
}

func setupServer() {
	httpServer := &http.Server{
		Addr:         configs.Data.Addr,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  1 * time.Second,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", HealthzHandler)
	mux.HandleFunc("/login", LoginHandler)
	mux.HandleFunc("/stuff", DoTheThingsHandler)
	httpServer.Handler = mux

	log.Fatal(httpServer.ListenAndServe())
}

func HealthzHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Success")
}

func LoginHandler(w http.ResponseWriter, req *http.Request) {

	//DO AUTH

	// if success

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		fmt.Fprintf(w, "Failed to sign token error: %+v", err)
		return
	}

	cookie1 := &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().UTC().Add(time.Hour * time.Duration(1)),
		HttpOnly: false,
	}
	http.SetCookie(w, cookie1)
}

func DoTheThingsHandler(w http.ResponseWriter, req *http.Request) {
	// Get token
	var authCookie, err = req.Cookie("auth_token")
	if err != nil || authCookie == nil || authCookie.Value == "" {
		fmt.Fprintf(w, "Failed to find auth_token error: %+v", err)
		return
	}
	authToken := authCookie.Value

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		fmt.Fprintf(w, "Here #1")
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Fprintf(w, "Here #3")
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		fmt.Fprintf(w, "Here #2")
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	fmt.Fprintf(w, "Here #4")
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Fprintf(w, "Here #5")
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Fprintf(w, "Here #6 %+v", err)
		fmt.Println(err)
	}
}
