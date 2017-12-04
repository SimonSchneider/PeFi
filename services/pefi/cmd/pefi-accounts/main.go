package main

import (
	"github.com/simonschneider/pefi/services/pefi/http"
	"github.com/simonschneider/pefi/services/pefi/postgres"
	"github.com/simonschneider/pefi/services/pefi/redis"
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

	service := redis.NewAccountService(postgres.NewAccountService())

	handler := http.NewAccountHandler(service)

	http.AttachAndStart(handler)
}
