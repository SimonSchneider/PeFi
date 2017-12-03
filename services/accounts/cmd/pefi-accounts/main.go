package main

import (
	"github.com/simonschneider/pefi/services/accounts/http"
	"github.com/simonschneider/pefi/services/accounts/postgres"
	"github.com/simonschneider/pefi/services/accounts/redis"
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

	service := redis.NewService(&postgres.Service{})

	handler := &http.Handler{*service}

	http.Start(handler)
}
