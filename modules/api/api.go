package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"

	"mimic/modules/api/services"
	"mimic/modules/db/mimic/blockdb"
	// ← v1 import path
	// ← v1 JSON codec
)

type APIServer struct {
	r chi.Router

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
		s.rpcRoutes[name+"."+alias] = &serv
		s.services[name] = reflect.ValueOf(service)
	})
}

func (s *APIServer) Init() {
	router := chi.NewRouter()

	loggerOpts := &httplog.Options{
		// Level defines the verbosity of the request logs:
		// slog.LevelDebug - log all responses (incl. OPTIONS)
		// slog.LevelInfo  - log responses (excl. OPTIONS)
		// slog.LevelWarn  - log 4xx and 5xx responses only (except for 429)
		// slog.LevelError - log 5xx responses only
		Level: slog.LevelInfo,

		// Set log output to Elastic Common Schema (ECS) format.
		Schema: httplog.SchemaECS,

		// RecoverPanics recovers from panics occurring in the underlying HTTP handlers
		// and middlewares. It returns HTTP 500 unless response status was already set.
		//
		// NOTE: Panics are logged as errors automatically, regardless of this setting.
		RecoverPanics: true,
	}

	router.Use(httplog.RequestLogger(slog.Default(), loggerOpts))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(
			[]byte(
				"go-mimic v1.0.0; Hive blockchain end to end simulation. To learn more, visit https://github.com/vsc-eco/go-mimic",
			),
		)
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte{})
	})

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Error("failed to decode incoming requests.", "err", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		method, valid := req["method"].(string)

		if !valid {
			http.Error(w, "invalid method", http.StatusBadRequest)
			return
		}

		if s.rpcRoutes[method] == nil {
			http.Error(w, "method not found", http.StatusNotFound)
			return
		}

		methodSpec := s.rpcRoutes[method]

		args := reflect.New(methodSpec.argType)
		paramsJSON, err := json.Marshal(req["params"])
		if err != nil {
			http.Error(w, "invalid params", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(paramsJSON, args.Interface()); err != nil {
			slog.Error("Failed to decode params",
				"raw", paramsJSON, "err", err)
			http.Error(w, "failed to decode params", http.StatusBadRequest)
			return
		}
		reply := reflect.New(s.rpcRoutes[method].replyType)

		strs := strings.Split(method, ".")
		methodSpec.method.Func.Call([]reflect.Value{
			s.services[strs[0]],
			args,
			reply,
		})

		res := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result":  reply.Interface(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	})

	s.r = router
}

func (s *APIServer) Start() {

	condenser := services.NewCondenser(blockdb.Collection())
	rcService := &services.RcApi{}
	blockApi := &services.BlockAPI{}
	accountHistoryApi := &services.AccountHistoryApi{}
	broadcastOps := &services.BroadcastOps{}

	s.RegisterService(condenser, "condenser_api")
	s.RegisterService(rcService, "rc_api")
	s.RegisterService(blockApi, "block_api")
	s.RegisterService(accountHistoryApi, "account_history_api")
	s.RegisterService(broadcastOps, "broadcast_ops")

	port := "3000"
	slog.Info("APIServer accepting requests.", "port", port)
	http.ListenAndServe(":"+port, s.r)
}

func NewAPIServer() *APIServer {
	return &APIServer{
		rpcRoutes: make(map[string]*ServiceMethod),
		services:  make(map[string]reflect.Value),
	}
}
