package main

import (
	"github.com/13days/gfly"
	"github.com/13days/gfly/plugin/consul"
	"github.com/13days/gfly/testdata"
	"time"

)

// consul agent -dev
func main() {
	opts := []gfly.ServerOption{
		gfly.WithAddress("127.0.0.1:8000"),
		gfly.WithNetwork("tcp"),
		gfly.WithSerializationType("msgpack"),
		gfly.WithTimeout(time.Millisecond * 2000),
		gfly.WithSelectorSvrAddr("localhost:8500"),
		gfly.WithPlugin(consul.Name),
	}
	s := gfly.NewServer(opts ...)
	if err := s.RegisterService("helloworld.Greeter", new(testdata.Service)); err != nil {
		panic(err)
	}
	s.Serve()
}