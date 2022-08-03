package routes

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func getConfig() (EnvConfiguration, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return EnvConfiguration{}, err
	}

	var config EnvConfiguration

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error reading env configuration")
		return EnvConfiguration{}, err
	}

	return config, nil
}

func getConnectionString(config EnvConfiguration) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s",
		config.DATABASE_USERNAME,
		config.DATABASE_PASSWORD,
		config.DATABASE_CONNECTION_STRING,
	)
}

func InitDatabase(config EnvConfiguration) (*gorm.DB, error) {
	dsn := getConnectionString(config)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
		return &gorm.DB{}, err
	}
	return db, nil
}
