package nodeB

import (
	"context"
	"fmt"
	"github.com/13days/gfly/testdata/nodeB/rpc"
)

type Service struct {

}

type CalcBCDRequest struct {
	B int64
	C int64
	D int64
}

type CalcBCDResponse struct {
	Answer int64
}

func (s *Service) CalcBCD(ctx context.Context, req *CalcBCDRequest) (*CalcBCDResponse, error) {
	fmt.Println("CalcBCD....")

	bd := rpc.CalcCD(ctx, req.C, req.D)
	answer := req.B * bd
	rsp := &CalcBCDResponse{
		Answer:answer,
	}
	return rsp, nil
}

