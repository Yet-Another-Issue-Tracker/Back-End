package main

import (
	"fmt"
	"issue-service/app/issue-api/cfg"
	"issue-service/app/issue-api/webserver"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func getLogLevel(logLevel string) log.Level {
	switch logLevel {
	case "trace":
		return log.TraceLevel
	case "info":
		return log.InfoLevel
	case "debug":
		return log.DebugLevel
	case "warning":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}

func initLogging(logLevel string) {
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	log.SetLevel(getLogLevel(logLevel))
}

func main() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.WithField("error", err.Error()).Fatal("Error reading env configuration")
		return
	}

	var config cfg.EnvConfiguration

	if err := viper.Unmarshal(&config); err != nil {
		log.WithField("error", err.Error()).Fatal("Error unmarshaling env configuration")
		return
	}
	initLogging(config.LOG_LEVEL)
	router := webserver.NewRouter(config)

	log.Info(fmt.Sprintf("Server starting on port: %s", config.HTTP_PORT))

	log.Info(http.ListenAndServe(fmt.Sprintf(":%s", config.HTTP_PORT), router))
}
