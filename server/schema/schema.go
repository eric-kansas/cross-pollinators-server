package schema

import (
	"log"

	"github.com/eric-kansas/cross-pollinators-server/database"
	"github.com/eric-kansas/cross-pollinators-server/database/models"
	"github.com/graphql-go/graphql"
)

var Root graphql.Schema

func init() {
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

				log.Printf("resolve: %v+", p.Context)

				return database.GetProjects("testing12345", first)
			},
		},

		"following": &graphql.Field{
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

				return database.GetProjects("testing12345", first)
			},
		},

		"user": &graphql.Field{
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

				return database.GetProjects("testing12345", first)
			},
		},
	}

	rootQuery := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	}

	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	Root, _ = graphql.NewSchema(schemaConfig)
}
