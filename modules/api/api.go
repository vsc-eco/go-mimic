package api

import (
	"fmt"
	"log/slog"
	"mimic/lib/httputil"
	"mimic/lib/utils"
	apijsonrpc "mimic/modules/api/jsonrpc"
	"mimic/modules/api/services"
	"mimic/modules/api/services/condenser"
	"mimic/modules/db/mimic/accountdb"
	"mimic/modules/db/mimic/blockdb"
	"net/http"
	"reflect"

	"github.com/chebyrash/promise"
	"github.com/go-chi/chi/v5"
	// ← v1 import path
	// ← v1 JSON codec
)

type APIServer struct {
	mux  *chi.Mux
	addr string
	rpc  apijsonrpc.Handler
	http httpHandler
}

func (s *APIServer) RegisterMethod(
	alias, methodName string,
	servc any,
) apijsonrpc.ServiceMethod {
	servType := reflect.TypeOf(servc)

	method, success := servType.MethodByName(methodName)
	if !success {
		panic(fmt.Sprintf("method not found: %s", methodName))
	}

	slog.Debug("Method registered.",
		"methodName", methodName,
		"methodNum", servType.NumMethod())

	mtype := method.Type

	return apijsonrpc.ServiceMethod{
		Method:  method,
		ArgType: mtype.In(1).Elem(),
	}
}

func (s *APIServer) RegisterService(
	service services.ServiceHandler,
	name string,
) {
	service.Expose(func(alias string, methodName string) {
		serv := s.RegisterMethod(alias, methodName, service)
		s.rpc.Routes[name+"."+alias] = &serv
		s.rpc.Services[name] = reflect.ValueOf(service)
	})
}

func (s *APIServer) Init() error {
	s.rpc.Logger = slog.Default().WithGroup("api")
	// initialize jsonrpc methods
	// rcService := &services.RcApi{}
	// blockApi := &services.BlockAPI{}
	// accountHistoryApi := &services.AccountHistoryApi{}
	condenser := &condenser.Condenser{
		// Logger:    s.rpc.logger,
		BlockDB:   blockdb.Collection(),
		AccountDB: accountdb.Collection(),
	}

	s.RegisterService(condenser, "condenser_api")
	// s.RegisterService(rcService, "rc_api")
	// s.RegisterService(blockApi, "block_api")
	// s.RegisterService(accountHistoryApi, "account_history_api")

	// intialize router
	s.mux = chi.NewRouter()
	s.mux.Use(httputil.RequestTrace(slog.Default().WithGroup("mimic-trace")))
	s.mux.Get("/", s.http.root)
	s.mux.Get("/health", s.http.health)
	s.mux.Post("/", s.rpc.Handle)

	return nil
}

func (s *APIServer) Start() *promise.Promise[any] {
	s.rpc.Logger.Info("APIServer accepting requests.", "addr", s.addr)
	go func(addr string, mux *chi.Mux) {
		http.ListenAndServe(addr, mux)
	}(s.addr, s.mux)

	return utils.PromiseResolve[any](nil)
}

func (a *APIServer) Stop() error {
	return nil
}

func NewAPIServer(httpPort uint16) *APIServer {
	return &APIServer{
		addr: fmt.Sprintf("0.0.0.0:%d", httpPort),
		rpc: apijsonrpc.Handler{
			Routes:   make(map[string]*apijsonrpc.ServiceMethod),
			Services: make(map[string]reflect.Value),
			Logger:   slog.Default().WithGroup("api-rpc"),
		},
		http: httpHandler{
			logger: slog.Default().WithGroup("api-http"),
		},
	}
}
