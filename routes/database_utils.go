package routes

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getConnectionString(config EnvConfiguration) string {
	return fmt.Sprintf("host=%s user=%s password=%s port=5432 dbname=%s",
		config.DATABASE_HOST,
		config.DATABASE_USERNAME,
		config.DATABASE_PASSWORD,
		config.DATABASE_NAME,
	)
}

func ConnectDatabase(config EnvConfiguration) (*gorm.DB, error) {
	dsn := getConnectionString(config)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
		return &gorm.DB{}, err
	}
	return db, nil
}

func SetupAndResetDatabase(database *gorm.DB) {
	database.AutoMigrate(&Project{})
	database.Exec("DELETE FROM projects")
	database.Exec("ALTER SEQUENCE projects_id_seq RESTART WITH 1")
}
