package flow_compare

import (
	"context"
	"fmt"
	"github.com/13days/gfly/protocol"
	"time"
)

type FlowCompareReq struct {
	p1, p2 *protocol.Response
	t1, t2 time.Duration
}

type FlowCompareResp struct {
	diff string
}

type FlowComparor interface {
	Compare(context.Context, *FlowCompareReq) error
}

var flowComparorMap = make(map[string]FlowComparor)

func init() {
	flowComparorMap["default"] = &FlowComparorDefault{}
}

type FlowComparorDefault struct{}

func (FlowComparorDefault) Compare(ctx context.Context, req *FlowCompareReq) error {
	compareName := "default"
	fmt.Println(compareName, ":", req)
	return nil
}
