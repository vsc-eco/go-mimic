package gomimic

import (
	"context"
	"fmt"
	"log/slog"
	"mimic/modules/admin"
	"mimic/modules/aggregate"
	"mimic/modules/api"
	"mimic/modules/config"
	"mimic/modules/db/mimic/accountdb"
	"mimic/modules/db/mimic/blockdb"
	"mimic/modules/db/mimic/transactiondb"
	"mimic/modules/producers"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	cfg config.AppConfig

	logger        *slog.Logger
	mongoClient   *mongo.Client
	mongoDatabase *mongo.Database
}

func NewApp(conf config.AppConfig) (*App, error) {
	app := &App{cfg: conf}
	if err := app.init(); err != nil {
		return nil, err
	}
	return app, nil
}

// Runs initialization in order of how they are passed in to `Aggregate`
func (app *App) init() error {
	app.logger = slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: app.cfg.LogFilter},
	))

	fmt.Println("Initialzing app")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var err error

	app.mongoClient, err = makeMongoClient(ctx, &app.cfg)
	if err != nil {
		return fmt.Errorf("failed connect to database: %v", err)
	}
	app.mongoDatabase = app.mongoClient.Database(app.cfg.DatabaseName)
	app.logger.Info("connected to database", "db", app.mongoDatabase.Name())

	if err := initCollections(ctx, app.mongoDatabase); err != nil {
		return fmt.Errorf("failed to initialized collections: %v", err)
	}

	return nil
}

func (app *App) Run(ctx context.Context) error {
	app.logger.Info("Starting app")

	routers := aggregate.New(
		api.NewAPIServer(app.cfg.GoMimicPort),
		admin.NewAPIServer(app.cfg.AdminPort, app.cfg.AdminToken),
		producers.New(),
	)

	if err := routers.Init(); err != nil {
		return fmt.Errorf("failed to initialized services: %v", err)
	}

	_, err := routers.Start().Await(ctx)
	return err
}

func makeMongoClient(
	ctx context.Context,
	cfg *config.AppConfig,
) (*mongo.Client, error) {
	mongoClientOpt := options.Client().ApplyURI(cfg.MongodbUrl)
	cx, err := mongo.Connect(ctx, mongoClientOpt)
	if err != nil {
		return nil, err
	}

	if err := cx.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return cx, nil
}

func initCollections(ctx context.Context, db *mongo.Database) error {
	_, err := aggregate.New(
		blockdb.New(db),
		accountdb.New(db),
		transactiondb.New(db),
	// hiveBlocks := blockdb.New(mimicDb)
	// stateDb := state.New(mimicDb)
	).Start().Await(ctx)

	return err
}
