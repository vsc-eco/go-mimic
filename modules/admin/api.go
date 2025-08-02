package admin

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"mimic/lib/httputil"
	"mimic/lib/utils"
	"mimic/modules/db/mimic/accountdb"
	"net/http"

	"github.com/chebyrash/promise"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// At run time, AdminAPI checks for the env variable `ADMIN_TOKEN`, if
// it isn't set, the server won't run.
//
// All requests issued to this server must have a matching value for the
// header `X-ADMIN-TOKEN`, otherwise an http status 401|403 is returned.
//
// AdminAPI implements aggregate.Plugins
type AdminAPI struct {
	adminToken string

	// if the env token is not set, this is nil
	mux      *chi.Mux
	httpAddr string

	handler *serverHandler
}

type serverHandler struct {
	logger *slog.Logger
	db     accountdb.AccountQuery
}

func NewAPIServer(httpPort uint16, token string) *AdminAPI {
	srv := &AdminAPI{
		adminToken: token,
		mux:        nil,
		httpAddr:   fmt.Sprintf("0.0.0.0:%d", httpPort),
		handler:    nil,
	}

	return srv
}

// Runs initialization in order of how they are passed in to `Aggregate`
func (a *AdminAPI) Init() error {
	a.handler.logger = slog.Default().WithGroup("admin-api")

	// load admin token
	adminDisabled := len(a.adminToken) == 0
	if adminDisabled {
		a.handler.logger.Info("admin server disabled.")
		a.mux = nil
		return nil
	}

	adminToken, err := hex.DecodeString(a.adminToken)
	if err != nil {
		return fmt.Errorf("invalid admin token hex: %v", err)
	}

	// db
	a.handler.db = accountdb.Collection()

	// initialize mux
	a.mux = chi.NewRouter()

	requestLogger := slog.Default().WithGroup("admin-trace")
	a.mux.Use(middleware.Logger)
	a.mux.Use(httputil.AuthMiddleware(adminToken, requestLogger))

	a.mux.Post("/user", a.handler.newUser)
	a.mux.Put("/user", a.handler.updateUser)

	return nil
}

// Runs startup and should be non blocking
func (a *AdminAPI) Start() *promise.Promise[any] {
	// `ADMIN_TOKEN` isn't set
	if a.mux == nil {
		return utils.PromiseResolve[any](nil)
	}

	a.handler.logger.Info("starting admin API server.", "addr", a.httpAddr)
	go func(mux *chi.Mux) {
		err := http.ListenAndServe(a.httpAddr, mux)
		if err != nil {
			a.handler.logger.Error("failed to start server.", "err", err)
		}
	}(a.mux)

	return utils.PromiseResolve[any](nil)
}

// Runs cleanup once the `Aggregate` is finished
func (a *AdminAPI) Stop() error {
	return nil
}
