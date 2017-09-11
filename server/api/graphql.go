package api

import (
	"fmt"
	"net/http"
)

func GraphQLHandler(w http.ResponseWriter, req *http.Request) {
	err := verifyUser(req)
	if err != nil {
		fmt.Fprintf(w, "Failed to verify user: %+v", err)
	}

}
