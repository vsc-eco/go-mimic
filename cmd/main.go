package main

import (
	"context"
	"log"
	"mimic/lib/utils"
	"mimic/modules/config"
	"mimic/modules/gomimic"
	"os"
)

var (
	cfg = config.AppConfig{}
)

const (
	mimicServerPort uint16 = 3000
	adminServerPort uint16 = 3001
)

func main() {
	cfg = config.AppConfig{
		GoMimicPort:  mimicServerPort,
		AdminPort:    adminServerPort,
		AdminToken:   os.Getenv("ADMIN_TOKEN"),
		MongodbUrl:   utils.EnvOrPanic("MONGODB_URL"),
		DatabaseName: utils.EnvOrPanic("MONGODB_DB_NAME"),
		LogFilter:    config.DefaultLogLevel(),
	}

	app, err := gomimic.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(app.Run(context.Background()))
}
