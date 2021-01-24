package rpc

import (
	"context"
	"fmt"
	"github.com/13days/gfly/testdata"
	"github.com/13days/gfly/testdata/nodeC"
	"time"
)

func CalcCD(ctx context.Context, c, d int64) int64 {
	req := &nodeC.CalcCDRequest{
		C: c,
		D: d,
	}
	rsp := &nodeC.CalcCDResponse{}
	err := testdata.Call(ctx, testdata.ServerCCalcCD, req, rsp, time.Second*2)
	if err != nil {
		fmt.Printf("call nodeB err:{ %v }\n", err)
	}
	return rsp.Answer
}
