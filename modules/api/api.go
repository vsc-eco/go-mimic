package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mitchellh/mapstructure"

	"github.com/vsc-eco/go-mimic/modules/api/services"
	// ← v1 import path
	// ← v1 JSON codec
)

type APIServer struct {
	r chi.Router

	rpcRoutes map[string]*ServiceMethod
	services  map[string]reflect.Value
}

func (s *APIServer) RegisterMethod(alias, methodName string, servc any) ServiceMethod {
	servType := reflect.TypeOf(servc)
	method, success := servType.MethodByName(methodName)

	for i := 0; i < servType.NumMethod(); i++ {
		fmt.Println("method", servType.Method(i).Name)
	}

	fmt.Println("method", methodName, method, servType)
	if success != true {
		panic("method not found")
	}

	mtype := method.Type

	fmt.Println("mtype.NumIn()", success)
	return ServiceMethod{
		method:    method,
		argType:   mtype.In(1).Elem(),
		replyType: mtype.In(2).Elem(),
	}
}

func (s *APIServer) RegisterService(service services.ServiceHandler, name string) {
	service.Expose(func(alias string, methodName string) {
		serv := s.RegisterMethod(alias, methodName, service)
		s.rpcRoutes[name+"."+alias] = &serv
		s.services[name] = reflect.ValueOf(service)
	})
}

func (s *APIServer) Init() {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("go-mimic v1.0.0; Hive blockchain end to end simulation. To learn more, visit https://github.com/vsc-eco/go-mimic"))
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte{})
	})

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		fmt.Println("req", req)

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
		err := mapstructure.Decode(req["params"], args.Interface())

		fmt.Println("args", args, err)

		reply := reflect.New(s.rpcRoutes[method].replyType)

		strs := strings.Split(method, ".")
		methodSpec.method.Func.Call([]reflect.Value{
			s.services[strs[0]],
			args,
			reply,
		})

		fmt.Println("req[\"Method\"]", req["method"], reply)

		res := map[string]any{
			"jsonrpc": "2.0",
			"result":  reply.Interface(),
			"id":      req["id"],
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(res)
	})

	s.r = router
}

func (s *APIServer) Start() {

	service := &services.Condenser{}

	s.RegisterService(service, "condenser_api")
	go http.ListenAndServe(":3000", s.r)
}

func NewAPIServer() *APIServer {
	return &APIServer{
		rpcRoutes: make(map[string]*ServiceMethod),
		services:  make(map[string]reflect.Value),
	}
}
