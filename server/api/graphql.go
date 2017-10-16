package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/eric-kansas/cross-pollinators-server/server/schema"
	"github.com/graphql-go/handler"
)

func GraphQLHander(auth bool) http.Handler {
	graphQLHander := handler.New(&handler.Config{
		Schema:   &schema.Root,
		Pretty:   true,
		GraphiQL: true,
	})

	if auth {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, err := GetUserID(r)
			if err != nil {
				log.Printf("Failed to verify user: %+v", err)
				fmt.Fprintf(w, "Failed to verify user: %+v", err)
				return
			}

			ctx := context.WithValue(context.Background(), "username", username)
			graphQLHander.ContextHandler(ctx, w, r)
		})
	}

	return graphQLHander
}
