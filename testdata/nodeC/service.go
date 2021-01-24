package nodeC

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
}

type CalcCDRequest struct {
	C int64
	D int64
}

type CalcCDResponse struct {
	Answer int64
}

func (s *Service) CalcCD(ctx context.Context, req *CalcCDRequest) (*CalcCDResponse, error) {
	fmt.Println("CalcCD...")
	time.Sleep(time.Second * 3)
	rsp := &CalcCDResponse{
		Answer: req.C * req.D,
	}
	return rsp, nil
}
