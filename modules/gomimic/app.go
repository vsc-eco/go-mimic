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

func (app *App) init() error {
	app.logger = slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: app.cfg.LogFilter},
	))

	fmt.Println("Initialzing app")

	// on a free instance, it could take a while to connect
	const dbConnectTimeout = 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout)
	defer cancel()

	app.logger.Info("connecting to database")

	var err error
	app.mongoClient, err = makeMongoClient(ctx, app.cfg.MongodbUrl)
	if err != nil {
		return fmt.Errorf("failed connect to database: %v", err)
	}

	app.mongoDatabase = app.mongoClient.Database(app.cfg.DatabaseName)
	app.logger.Info("database connected", "db", app.mongoDatabase.Name())

	app.logger.Info("initializing collection")
	if err := initCollections(ctx, app.mongoDatabase); err != nil {
		return fmt.Errorf("failed to initialized collections: %v", err)
	}
	app.logger.Info("collections initialized")

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

	if _, err := routers.Start().Await(ctx); err != nil {
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func makeMongoClient(ctx context.Context, uri string) (*mongo.Client, error) {
	mongoClientOpt := options.Client().ApplyURI(uri)

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

	agg := aggregate.New(
		blockdb.New(db),
		accountdb.New(db),
		transactiondb.New(db),
		// hiveBlocks := blockdb.New(mimicDb)
		// stateDb := state.New(mimicDb)
	)

	if err := agg.Init(); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}

	_, err := agg.Start().Await(ctx)
	return err
}
