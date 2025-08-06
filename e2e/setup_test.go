package e2e_test

import (
	"context"
	"log/slog"
	"mimic/lib/utils"
	"mimic/modules/app"
	"mimic/modules/config"
	"time"
)

func setupTest(a **app.App) error {
	appConf := config.AppConfig{
		GoMimic: config.GoMimicConfig{
			Port: goMimicPort,
		},
		Admin: config.AdminConfig{
			Port:  3001,
			Token: utils.EnvOrPanic("ADMIN_TOKEN"),
		},
		LogFilter:    slog.LevelError,
		MongodbUrl:   utils.EnvOrPanic("MONGODB_URL"),
		DatabaseName: "go-mimic-test-db",
	}

	var err error

	*a, err = app.NewApp(appConf)
	if err != nil {
		return err
	}

	return nil
}

func teardownTest(a *app.App) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.Database().Drop(ctx); err != nil {
		slog.Error("failed to drop test database",
			"db", a.Database().Name(),
			"err", err,
		)
	}
}
