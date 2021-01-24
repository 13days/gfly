package testdata

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Service struct {
}

type HelloRequest struct {
	Msg string
}

type HelloReply struct {
	Msg string
}

func (s *Service) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	rsp := &HelloReply{
		Msg: "world",
	}
	t := rand.Int() % 1000
	time.Sleep(time.Millisecond * time.Duration(t))
	fmt.Println("call....")
	return rsp, nil
}
