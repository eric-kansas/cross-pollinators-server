package api

import (
	"fmt"
	"log"
	"net/http"
)

func HealthzHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Health Handler")
	fmt.Fprintf(w, "Success")
}
