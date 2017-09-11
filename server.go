package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/eric-kansas/cross-pollinators-server/server/api"
	"github.com/eric-kansas/cross-pollinators-server/server/configs"
)

var httpServer = &http.Server{
	Addr:         configs.Data.Addr,
	ReadTimeout:  1 * time.Second,
	WriteTimeout: 1 * time.Second,
	IdleTimeout:  1 * time.Second,
}

func init() {
	fmt.Printf("Server started: Version %s \n", "alpha-0.3.0")
	fmt.Printf("Server running go version %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	configs.Initialize()
}

func main() {
	setupAPI()

	log.Fatal(http.ListenAndServe(configs.Data.Addr, nil))
}

func setupAPI() {
	mux := http.NewServeMux()

	// TODO: FIX USE MUX
	http.HandleFunc("/healthz", api.HealthzHandler)
	http.HandleFunc("/login", api.LoginHandler)
	http.HandleFunc("/register", api.RegisterHandler)

	// define GraphQL schema using relay library helpers
	fields := graphql.Fields{
		"user": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, _ := graphql.NewSchema(schemaConfig)

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})
	authGraphQL := api.AuthWrapper(h)

	http.Handle("/graphql", authGraphQL)

	httpServer.Handler = mux
}
