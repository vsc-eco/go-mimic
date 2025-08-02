package main

import (
	"context"
	"log"
	"mimic/lib/utils"
	"mimic/modules/config"
	"mimic/modules/gomimic"
	"os"
	"time"
)

var (
	cfg = config.AppConfig{}
)

const (
	mimicServerPort uint16 = 3000
	adminServerPort uint16 = 3001
)

func init() {
	cfg = config.AppConfig{
		GoMimicPort:  mimicServerPort,
		AdminPort:    adminServerPort,
		AdminToken:   os.Getenv("ADMIN_TOKEN"),
		MongodbUrl:   utils.EnvOrPanic("MONGODB_URL"),
		DatabaseName: utils.EnvOrPanic("MONGODB_DB_NAME"),
		LogFilter:    config.DefaultLogLevel(),
	}
}

func main() {
	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()

	app, err := gomimic.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(app.Run(ctx))
}
