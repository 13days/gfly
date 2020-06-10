package gfly

import "context"

// Service defines a generic implementation interface for a specific Service
type Service interface {
	Register(string, Handler)
	Serve(*ServerOptions)
	Close()
}

type service struct {
	svr         interface{}        // server
	ctx         context.Context    // Each service is managed in one context
	cancel      context.CancelFunc // controller of context
	serviceName string             // service name
	handlers    map[string]Handler // much handler for a service,one handler for one method
	opts        *ServerOptions     // parameter options
	closing     bool               // whether the service is closing
}

// Handler is the handler of a method
type Handler func(interface{}, context.Context, func(interface{}) error, []interceptor.ServerInterceptor) (interface{}, error)
