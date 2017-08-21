package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/eric-kansas/cross-pollinators-server/configs"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

var hmacSampleSecret = []byte("my_secret_key")

const (
	dbUser = "kansas"
	dbPass = "pass1234"
	dbName = "cross-pollinators-db"
)

func init() {
	configs.Initialize()
	connectDB()
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
	mux.HandleFunc("/register", RegisterHandler)
	mux.HandleFunc("/dostuff", DoTheThingsHandler)
	httpServer.Handler = mux

	log.Fatal(httpServer.ListenAndServe())
}

func HealthzHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Success")
}

func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		fmt.Fprintf(w, "Failed login not a POST")
		return
	}

	req.ParseForm()
	fmt.Println("username:", req.Form["username"])
	fmt.Println("password:", req.Form["password"])
	password := []byte(req.Form["password"][0])

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintf(w, "Failed to hash password: %+v", err)
		return
	}

	// Save username with hashed passedword to data base
	fmt.Println(string(hashedPassword))
}

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

func connectDB() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPass, dbName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Printf("Failed to open connection to postgres database: %+v \n", err)
		return
	}
	log.Printf("Cross Pollinators Service connected to DB!")

	age := 21
	_, err = db.Query("SELECT name FROM users WHERE age = $1", age)

	defer db.Close()
}
