package gfly

import (
	"os"
	"os/signal"
	"plugin"
	"syscall"
)

// gfly Server, a Server can have one or more Services
type Server struct {
	opts     *ServerOptions
	services map[string]Service
	plugins  []plugin.Plugin

	closing bool // whether the server is closing
}

func NerServer(opts ...ServerOption) *Server {
	server := &Server{
		opts: &ServerOptions{},
		services: make(map[string]Service),
	}

	for _, opt := range opts{
		opt(server.opts)
	}
	return server
}

func (s *Server)Serve()  {
	for _, service := range s.services{
		go service.Serve(s.opts)
	}
	
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV)
	<-ch
	
	s.Close()
}

func (s *Server)Close()  {
	
}