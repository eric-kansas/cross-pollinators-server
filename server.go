package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/eric-kansas/cross-pollinators-server/database"
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
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "The id of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(models.User); ok {
						return user.ID, nil
					}
					return "Sad", nil
				},
			},
			"full_name": &graphql.Field{
				Type:        graphql.String,
				Description: "Users full name",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(models.User); ok {
						return user.Username, nil
					}
					return "Sad", nil
				},
			},
			"avatar_url": &graphql.Field{
				Type:        graphql.String,
				Description: "Avatar url",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(models.User); ok {
						return user.Username, nil
					}
					return "Sad", nil
				},
			},
			"organization": &graphql.Field{
				Type:        graphql.String,
				Description: "Organization user is a part of",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(models.User); ok {
						return user.Organization, nil
					}
					return "Sad", nil
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
				Type:        graphql.ID,
				Description: "Project ID",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if project, ok := p.Source.(models.Project); ok {
						return project.ID, nil
					}
					return "Sad", nil
				},
			},
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "Text content of the project",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if project, ok := p.Source.(models.Project); ok {
						return project.Name, nil
					}
					return "Sad", nil
				},
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "Description of the project",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if project, ok := p.Source.(models.Project); ok {
						return project.Description, nil
					}
					return "No description", nil
				},
			},
			"objective": &graphql.Field{
				Type:        graphql.String,
				Description: "Text objective of the project",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if project, ok := p.Source.(models.Project); ok {
						return project.Objective, nil
					}
					return "No objective... Lets talk!", nil
				},
			},
			"author": &graphql.Field{
				Type:        userType,
				Description: "Author of the post",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if project, ok := p.Source.(models.Project); ok {
						user, err := database.GetUser(project.UserID)
						if err == nil {
							return user, nil
						}
					}
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
		"hello": &graphql.Field{
			Type:        graphql.String,
			Description: "Hello world",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
		"discover": &graphql.Field{
			Type:        graphql.NewList(projectType),
			Description: "List of projects to discover",
			Args: graphql.FieldConfigArgument{
				"first": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 10,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				first := 10
				if v, ok := p.Args["first"].(int); ok {
					first = v
				}
				log.Printf("first: %d", first)

				return database.GetProjects("testing12345", first)
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
