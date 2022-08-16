package internal

import (
	"fmt"
	"log"
	"strings"

	"issue-service/app/issue-api/cfg"
	models "issue-service/app/issue-api/routes/models"

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

func dropTables(database *gorm.DB) {
	database.Migrator().DropTable(&models.Project{})
	database.Migrator().DropTable(&models.Sprint{})
}

func automigrateDatabaseSchema(database *gorm.DB) {
	database.AutoMigrate(&models.Project{})
	database.AutoMigrate(&models.Sprint{})
}

func resetSequenceId(database *gorm.DB) {
	database.Exec("ALTER SEQUENCE projects_id_seq RESTART WITH 1")
	database.Exec("ALTER SEQUENCE sprints_id_seq RESTART WITH 1")
}

func SetupAndResetDatabase(database *gorm.DB) {
	dropTables(database)
	automigrateDatabaseSchema(database)
	resetSequenceId(database)
}

func IsDuplicateKeyError(databaseError error) bool {
	return strings.Contains(databaseError.Error(), DUPLICATE_KEY_ERROR)
}
