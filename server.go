package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/eric-kansas/cross-pollinators-server/database/models"
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

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "A user in cross pollinators",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The id of the droid.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},
		},
	})

	tagEnum := graphql.NewEnum(graphql.EnumConfig{
		Name:        "Episode",
		Description: "One of the films in the Star Wars Trilogy",
		Values: graphql.EnumValueConfigMap{
			"Tag1": &graphql.EnumValueConfig{
				Value:       1,
				Description: "Tag 1.",
			},
			"Tag2": &graphql.EnumValueConfig{
				Value:       2,
				Description: "Tag 2.",
			},
			"Tag3": &graphql.EnumValueConfig{
				Value:       3,
				Description: "Tag 3.",
			},
		},
	})

	/*
		id: ID!
		name: String!
		header: String!
		sub_header: String!
		body: String!
		header_img_url: String!
		author: Author!
		followed_by: [User]!
		tags: [Tag]!
	*/
	projectType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Project",
		Description: "Project object",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.Int,
				Description: "Project ID",
			},
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "Text content of the project",
			},
			"header": &graphql.Field{
				Type:        graphql.String,
				Description: "Text content of the project",
			},
			"sub_header": &graphql.Field{
				Type:        graphql.String,
				Description: "Text content of the project",
			},
			"body": &graphql.Field{
				Type:        graphql.String,
				Description: "Text content of the project",
			},
			"author": &graphql.Field{
				Type:        userType,
				Description: "Author of the post",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},
			"followed_by": &graphql.Field{
				Type:        graphql.NewList(userType),
				Description: "List of users who liked the project",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},
			"tags": &graphql.Field{
				Type:        graphql.NewList(tagEnum),
				Description: "List of users who liked the project",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},
		},
	})

	/*
	   discover (user: user-id) {
	   	projects: (first: 10) [
	   		project {
	   			header_img_url:
	   			header:
	   			sub_header:
	   			body:
	   			author: {
	   				avatar_url:
	   				username:
	   				organization:
	   			}
	   			tags: ["tag1","tag2"]
	   		},
	   		...
	   	]
	   }
	*/

	// define GraphQL schema using relay library helpers
	fields := graphql.Fields{
		"discover": &graphql.Field{
			Type:        projectType,
			Description: "List of projects to discover",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// get value of `limit` arg
				/*limit := 10
				if v, ok := p.Args["limit"].(int); ok {
					limit = v
				}
				return data.GetPosts(limit)
				*/
				return models.Project{}, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	}

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
