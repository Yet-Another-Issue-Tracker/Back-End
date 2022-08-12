package internal

import (
	"fmt"
	"log"
	"strings"

	"issue-service/app/issue-api/cfg"
	models "issue-service/app/issue-api/routes/makes/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DUPLICATE_KEY_ERROR = "duplicate key value violates unique constraint"

func getConnectionString(config cfg.EnvConfiguration) string {
	return fmt.Sprintf("host=%s user=%s password=%s port=5432 dbname=%s",
		config.DATABASE_HOST,
		config.DATABASE_USERNAME,
		config.DATABASE_PASSWORD,
		config.DATABASE_NAME,
	)
}

func ConnectDatabase(config cfg.EnvConfiguration) (*gorm.DB, error) {
	dsn := getConnectionString(config)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
		return &gorm.DB{}, err
	}
	return db, nil
}

func SetupAndResetDatabase(database *gorm.DB) {
	database.AutoMigrate(&models.Project{})
	database.Exec("DELETE FROM projects")
	database.Exec("ALTER SEQUENCE projects_id_seq RESTART WITH 1")
}

func IsDuplicateKeyError(databaseError error) bool {
	return strings.Contains(databaseError.Error(), DUPLICATE_KEY_ERROR)
}
