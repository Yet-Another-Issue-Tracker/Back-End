package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"issue-service/app/issue-api/cfg"
	models "issue-service/app/issue-api/routes/models"
	"math/rand"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetRandomStringName(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

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

	return config, nil
}

func LogAndReturnErrorResponse(err error, w http.ResponseWriter) {
	var errorResponse *models.ErrorResponse
	errors.As(err, &errorResponse)
	log.WithField("error", errorResponse.ErrorMessage)

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

func GetCreateResponseBody(id uint) ([]byte, error) {
	response := models.CreateResponse{Id: fmt.Sprint(id)}

	responseBody, err := json.Marshal(response)
	if err != nil {
		return []byte{}, &models.ErrorResponse{
			ErrorMessage: "Error mashaling the response",
			ErrorCode:    500,
		}
	}
	return responseBody, nil
}
