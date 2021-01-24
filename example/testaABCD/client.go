package main

import (
	"context"
	"fmt"
	"github.com/13days/gfly/testdata"
	"github.com/13days/gfly/testdata/nodeA"
	"time"
)

func main() {
	req := &nodeA.CalcABCDRequest{
		A: 1,
		B: 2,
		C: 3,
		D: 4,
	}
	rsp := &nodeA.CalcABCDResponse{}
	err := testdata.Call(context.Background(), testdata.ServerACalcABCD, req, rsp, time.Second*2)
	if err != nil {
		fmt.Printf("call nodeB err:{ %v }\n", err)
	}
	fmt.Println("Answer:",rsp.Answer)
}
