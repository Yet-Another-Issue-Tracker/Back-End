package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"issue-service/app/issue-api/cfg"
	models "issue-service/app/issue-api/routes/makes/models"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func GetConfig(path string) (cfg.EnvConfiguration, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {

		return cfg.EnvConfiguration{}, err
	}

	var config cfg.EnvConfiguration

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error reading env configuration")
		return cfg.EnvConfiguration{}, err
	}
	log.Printf("Connected %s", config.DATABASE_HOST)

	return config, nil
}

func ReturnErrorResponse(err error, w http.ResponseWriter) {
	var errorResponse *models.ErrorResponse
	errors.As(err, &errorResponse)

	jsonResponse, _ := json.Marshal(err)

	http.Error(w, string(jsonResponse), errorResponse.ErrorCode)
}

func ValidateRequest(inputRequest interface{}) error {
	var validationError string = ""
	validate := validator.New()
	err := validate.Struct(inputRequest)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Validation error, field: %s, tag: %s", err.Namespace(), err.Tag())
			if validationError == "" {
				validationError = fmt.Sprintf("%s%s", validationError, errorMessage)
			} else {
				validationError = fmt.Sprintf("%s\n%s", validationError, errorMessage)
			}
		}

		return &models.ErrorResponse{
			ErrorMessage: validationError,
			ErrorCode:    400,
		}
	}
	return nil
}
