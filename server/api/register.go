package api

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

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
