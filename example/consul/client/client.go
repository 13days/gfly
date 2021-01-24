package main

import (
	"context"
	"fmt"
	"github.com/13days/gfly/client"
	"github.com/13days/gfly/plugin/consul"
	"github.com/13days/gfly/testdata"
	"time"

)

func main() {
	opts := []client.Option {
		client.WithNetwork("tcp"),
		client.WithTimeout(2000 * time.Millisecond),
		client.WithSelectorName(consul.Name),
	}
	c := client.DefaultClient
	req := &testdata.HelloRequest{
		Msg: "hello",
	}
	rsp := &testdata.HelloReply{}

	consul.Init("localhost:8500")
	for i:=0; i<1; i++ {
		err := c.Call(context.Background(), "/helloworld.Greeter/SayHello", req, rsp, opts ...)
		fmt.Println(rsp.Msg, err)
	}
}
