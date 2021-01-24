package rpc

import (
	"context"
	"fmt"
	"github.com/13days/gfly/testdata"
	"github.com/13days/gfly/testdata/nodeB"
	"time"
)

func CalcBCD(ctx context.Context, b, c, d int64) int64 {
	req := &nodeB.CalcBCDRequest{
		B: b,
		C: c,
		D: d,
	}
	rsp := &nodeB.CalcBCDResponse{}
	err := testdata.Call(ctx, testdata.ServerBCalcBCD, req, rsp, time.Second*3)
	if err != nil {
		fmt.Printf("call nodeB err:{ %v }\n", err)
	}
	return rsp.Answer
}