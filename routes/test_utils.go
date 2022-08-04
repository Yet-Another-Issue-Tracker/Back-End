package routes

import (
	"log"

	"github.com/spf13/viper"
)

func GetConfig(path string) (EnvConfiguration, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {

		return EnvConfiguration{}, err
	}

	var config EnvConfiguration

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error reading env configuration")
		return EnvConfiguration{}, err
	}
	log.Printf("Connected %s", config.DATABASE_HOST)

	return config, nil
}
