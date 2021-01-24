package testdata

import (
	"context"
	"github.com/13days/gfly"
	"github.com/13days/gfly/client"
	"github.com/13days/gfly/plugin/consul"
	"time"
)

const (
	ServerA         = "gfly.GDUF/ServerA"
	ServerACalcABCD = "/gfly.GDUF/ServerA/CalcABCD"
	ServerB         = "gfly.GDUF/ServerB"
	ServerBCalcBCD  = "/gfly.GDUF/ServerB/CalcBCD"
	ServerC         = "gfly.GDUF/ServerC"
	ServerCCalcCD   = "/gfly.GDUF/ServerC/CalcCD"
)

func getOptions(dur time.Duration) []client.Option {
	opts := []client.Option{
		client.WithNetwork("tcp"),
		client.WithTimeout(dur),
		client.WithSelectorName(consul.Name),
	}

	return opts
}

func Call(ctx context.Context, serverName string, req, rsp interface{}, dur time.Duration) error {
	consul.Init("localhost:8500")
	return client.DefaultClient.Call(ctx, serverName, req, rsp, getOptions(dur)...)
}

func StartServer(serverName, addr string, svr interface{}) {
	opts := []gfly.ServerOption{
		gfly.WithAddress(addr),
		gfly.WithNetwork("tcp"),
		gfly.WithSerializationType("msgpack"),
		gfly.WithTimeout(time.Millisecond * 2000),
		gfly.WithSelectorSvrAddr("localhost:8500"),
		gfly.WithPlugin(consul.Name),
	}
	s := gfly.NewServer(opts...)
	if err := s.RegisterService(serverName, svr); err != nil {
		panic(err)
	}
	s.Serve()
}
