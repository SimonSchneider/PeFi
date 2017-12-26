package main

import (
	"github.com/simonschneider/pefi/services/pefi/http"
	"github.com/simonschneider/pefi/services/pefi/postgres"
	//"github.com/simonschneider/pefi/services/pefi/redis"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	log.Info("starting")
	log.Warn("warnign")
	config := &postgres.Config{
		User:    "postgres",
		DbName:  "pefi",
		Sslmode: "disable",
	}
	service, err := postgres.NewAccountService(config)
	if err != nil {
		log.Error("Error connecting to db: ", err)
		os.Exit(1)
	}
	//service := redis.NewAccountService(pservice)

	handler := http.NewAccountHandler(service)

	http.AttachAndStart(handler)
}
