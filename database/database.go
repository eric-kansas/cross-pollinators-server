package database

import (
	"fmt"
	"log"

	"github.com/eric-kansas/cross-pollinators-server/database/models"
	"github.com/eric-kansas/cross-pollinators-server/server/configs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Connect() (*gorm.DB, error) {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", "db", configs.Data.DbUser, configs.Data.DbPass, configs.Data.DbName)
	db, err := gorm.Open("postgres", dbinfo)

	if err != nil {
		log.Printf("Failed to open connection to postgres database: %+v \n", err)
		return nil, err
	}

	log.Printf("Connected to Cross Pollinators DB!")
	return db, nil
}
func GetProject() (models.Project, error) {
	return models.Project{
		Name: "Project 1",
	}, nil
}

func GetProjects(username string, amount int) ([]models.Project, error) {

	// TODO: DO database call

	db, err := Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	projects := []models.Project{}

	db.Find(&projects)
	if len(projects) == 0 {
		return nil, fmt.Errorf("failed to find projects")
	}
	/*
		var projects []models.Project

		projects = append(projects,
			models.Project{
				Name:        "Project 1",
				Description: "Description 1",
			},
			models.Project{
				Name:        "Project 2",
				Description: "Description 2",
			},
			models.Project{
				Name:        "Project 3",
				Description: "Description 3",
			},
		)*/
	return projects, nil
}
