package main

import (
	"github.com/13days/gfly"
	"github.com/13days/gfly/testdata"
	"time"
)

func main() {
	opts := []gfly.ServerOption{
		gfly.WithAddress("127.0.0.1:5555"),
		gfly.WithNetwork("tcp"),
		gfly.WithSerializationType("msgpack"),
		gfly.WithTimeout(time.Millisecond * 2000),
	}
	s := gfly.NewServer(opts...)
	if err := s.RegisterService("/gduf.Greeter", new(testdata.Service)); err != nil {
		panic(err)
	}
	s.Serve()
}
