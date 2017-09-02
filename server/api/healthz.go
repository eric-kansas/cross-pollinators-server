package api

import (
	"fmt"
	"net/http"
)

func HealthzHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Success")
}
