package api

import (
	"fmt"
	"log/slog"
	"mimic/lib/httputil"
	"mimic/lib/utils"
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
	mux     *chi.Mux
	addr    string
	handler requestHandler
}

type requestHandler struct {
	logger    *slog.Logger
	rpcRoutes map[string]*ServiceMethod
	services  map[string]reflect.Value
}

func (s *APIServer) RegisterMethod(
	alias, methodName string,
	servc any,
) ServiceMethod {
	servType := reflect.TypeOf(servc)

	method, success := servType.MethodByName(methodName)
	if !success {
		panic("method not found")
	}

	slog.Debug("Method registered.",
		"methodName", methodName,
		"methodNum", servType.NumMethod())

	mtype := method.Type

	return ServiceMethod{
		method:    method,
		argType:   mtype.In(1).Elem(),
		replyType: mtype.In(2).Elem(),
	}
}

func (s *APIServer) RegisterService(
	service services.ServiceHandler,
	name string,
) {
	service.Expose(func(alias string, methodName string) {
		serv := s.RegisterMethod(alias, methodName, service)
		s.handler.rpcRoutes[name+"."+alias] = &serv
		s.handler.services[name] = reflect.ValueOf(service)
	})
}

func (s *APIServer) Init() error {
	// initialize jsonrpc methods
	rcService := &services.RcApi{}
	blockApi := &services.BlockAPI{}
	accountHistoryApi := &services.AccountHistoryApi{}
	condenser := &condenser.Condenser{
		BlockDB:   blockdb.Collection(),
		AccountDB: accountdb.Collection(),
	}

	s.RegisterService(condenser, "condenser_api")
	s.RegisterService(rcService, "rc_api")
	s.RegisterService(blockApi, "block_api")
	s.RegisterService(accountHistoryApi, "account_history_api")

	// intialize router
	s.mux = chi.NewRouter()
	s.mux.Use(httputil.RequestTrace(slog.Default().WithGroup("mimic-trace")))
	s.mux.Get("/", s.handler.root)
	s.mux.Get("/health", s.handler.health)
	s.mux.Post("/", s.handler.jsonrpc)

	return nil
}

func (s *APIServer) Start() *promise.Promise[any] {
	s.handler.logger.Info("APIServer accepting requests.", "addr", s.addr)
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
		handler: requestHandler{
			rpcRoutes: make(map[string]*ServiceMethod),
			services:  make(map[string]reflect.Value),
			logger:    slog.Default().WithGroup("mimic"),
		},
	}
}
