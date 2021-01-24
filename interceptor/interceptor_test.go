package interceptor

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntercept(t *testing.T) {
	ivk := func (ctx context.Context, req, rsp interface{}) error {
		fmt.Println("invoker...")
		return nil
	}

	inter1 := func(ctx context.Context, req, rsp interface{}, ivk Invoker) error {
		fmt.Println("interceptor1...")
		return ivk(ctx, req,rsp)
	}

	inter2 := func(ctx context.Context, req,rsp interface{},  ivk Invoker) error {
		fmt.Println("interceptor2...")
		return ivk(ctx, req,rsp)
	}
	ceps := []ClientInterceptor{inter1, inter2}

	err := ClientIntercept(context.Background(), nil ,nil, ceps , ivk)
	assert.Nil(t, err)
}

func TestServerIntercept(t *testing.T) {
	han := func(ctx context.Context, req interface{}) (interface{}, error){
		fmt.Println("handle...")
		return nil, nil
	}
	inter1 := func(ctx context.Context, req interface{}, handler Handler) (interface{}, error){
		fmt.Println("interceptor1...")
		return handler(ctx, req)
	}

	inter2 := func(ctx context.Context, req interface{}, handler Handler) (interface{}, error){
		fmt.Println("interceptor2...")
		return handler(ctx, req)
	}

	ceps := []ServerInterceptor{inter1, inter2}

	err, _ := ServerIntercept(context.Background(), nil, ceps, han)
	assert.Nil(t, err)
}