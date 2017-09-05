package database

import (
	"fmt"
	"log"

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
