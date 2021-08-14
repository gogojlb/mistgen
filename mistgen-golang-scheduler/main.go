package main

import (
	"golang-mistgen/app"
	"golang-mistgen/utils"

	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// type sensorReq struct {
// 	Props []string `json:"props"`
// }

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "debug"
	}
	// parse string, this is built-in feature of logrus
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.DebugLevel
	}
	// set global log level
	log.SetLevel(ll)
}

func main() {
	config := utils.NewConfig()
	utils.MigratePostgresUp(config.Postgres)
	pgConnection := utils.NewPostgresConnection(config.Postgres)
	// requests := utils.NewRequests(config.MiApi)
	ticker := time.NewTicker(config.Interval)
	jobResult := &app.JobResult{}
	for range ticker.C {
		currentJobResult := app.StartJob(config.Mist, pgConnection, jobResult, config.Grafana)
		if currentJobResult.Executed {
			jobResult = currentJobResult
		}
	}

}
