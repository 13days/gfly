package main

import (
	"context"
	"fmt"
	"github.com/13days/gfly/client"
	"github.com/13days/gfly/testdata"
	"time"
)

func main() {
	opts := []client.Option {
		client.WithTarget("127.0.0.1:5555"),
		client.WithNetwork("tcp"),
		client.WithTimeout(500 * time.Millisecond),
		client.WithSerializationType("msgpack"),
	}
	c := client.DefaultClient

	for o:=0; o<1; o++ {
		go func() {
			for i:=0; i<1; i++ {
				req := &testdata.HelloRequest{
					Msg: "hello",
				}
				rsp := &testdata.HelloReply{}
				err := c.Call(context.Background(), "/gduf.Greeter/SayHello", req, rsp, opts ...)
				fmt.Println(rsp.Msg, err)
			}
		}()
	}
	time.Sleep(time.Second * 10)
}
