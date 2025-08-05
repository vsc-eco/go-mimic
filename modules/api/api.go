package api

import (
	"fmt"
	"log"
	"log/slog"
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
	"github.com/go-chi/chi/v5/middleware"
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

	mtype := method.Type

	if err := validateMethod(mtype); err != nil {
		msg := fmt.Sprintf(
			`invalid method function signature %s.%s:
expected: func(*self, any) (any, *jsonrpc2.Error)
error: %v`,
			reflect.TypeOf(servc), methodName, err,
		)
		panic(msg)
	}

	return apijsonrpc.ServiceMethod{
		Method:  method,
		ArgType: mtype.In(1).Elem(),
	}
}

func validateMethod(m reflect.Type) error {
	if m.NumIn() != 2 {
		return errors.New("invalid argument count")
	}

	if m.NumOut() != 2 {
		return errors.New("invalid return value count")
	}

	errOut := m.Out(1).Elem()
	if errOut.String() != "jsonrpc2.Error" {
		return errors.New("invalid error type")
	}

	return nil
}

func (s *GoMimicAPI) RegisterService(
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
	rcService := &services.RcApi{}
	blockApi := &services.BlockAPI{}
	accountHistoryApi := &services.AccountHistoryApi{}
	condenser := &condenser.Condenser{
		Logger:    s.rpc.Logger,
		BlockDB:   blockdb.Collection(),
		AccountDB: accountdb.Collection(),
	}

	s.RegisterService(condenser, "condenser_api")
	s.RegisterService(rcService, "rc_api")
	s.RegisterService(blockApi, "block_api")
	s.RegisterService(accountHistoryApi, "account_history_api")

	// intialize router
	s.mux = chi.NewRouter()
	s.mux.Use(middleware.DefaultLogger)
	s.mux.Get("/", s.http.root)
	s.mux.Get("/health", s.http.health)
	s.mux.Post("/", s.rpc.Handle)

	return nil
}

func (s *APIServer) Start() *promise.Promise[any] {
	s.rpc.Logger.Info("starting go-mimic API server.", "addr", s.addr)
	go func(addr string, mux *chi.Mux) {
		log.Fatal(http.ListenAndServe(addr, mux))
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
