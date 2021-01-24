package nodeA

import (
	"context"
	"fmt"
	"github.com/13days/gfly/testdata/nodeA/rpc"
)

type Service struct {

}

type CalcABCDRequest struct {
	A int64
	B int64
	C int64
	D int64
}

type CalcABCDResponse struct {
	Answer int64
}

func (s *Service) CalcABCD(ctx context.Context, req *CalcABCDRequest) (*CalcABCDResponse, error) {
	fmt.Println("CalcABCD....")
	bcd := rpc.CalcBCD(ctx, req.B, req.C, req.D)
	answer := req.A * bcd
	rsp := &CalcABCDResponse{
		Answer: answer,
	}
	return rsp, nil
}

