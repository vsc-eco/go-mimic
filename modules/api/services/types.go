package services

type ServiceHandler interface {
	Expose(mr RegisterMethod)
}

type RegisterMethod func(alias string, methodName string)
