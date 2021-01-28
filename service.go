package gfly

import (
	"context"
	"errors"
	"fmt"
	"github.com/13days/gfly/codec"
	"github.com/13days/gfly/codes"
	"github.com/13days/gfly/interceptor"
	"github.com/13days/gfly/log"
	"github.com/13days/gfly/metadata"
	"github.com/13days/gfly/plugin/consul"
	"github.com/13days/gfly/protocol"
	"github.com/13days/gfly/transport"
	"github.com/13days/gfly/utils"
	"github.com/golang/protobuf/proto"
)


var(
	ServerRPCTimeoutErr = errors.New("RPC handle timeout")
)

// Service defines a generic implementation interface for a specific Service
type Service interface {
	Register(string, Handler)
	Serve(*ServerOptions)
	Close()
	Name() string
}

func NewService(opts *ServerOptions) Service {
	return &service{
		opts: opts,
	}
}

type service struct {
	svr         interface{}        // server
	ctx         context.Context    // Each service is managed in one context
	cancel      context.CancelFunc // controller of context
	serviceName string             // service name
	handlers    map[string]Handler
	opts        *ServerOptions // parameter options

	closing bool // whether the service is closing
}

// ServiceDesc is a detailed description of a service
type ServiceDesc struct {
	Svr         interface{}
	ServiceName string
	Methods     []*MethodDesc
	HandlerType interface{}
}

// MethodDesc is a detailed description of a method
type MethodDesc struct {
	MethodName string
	Handler    Handler
}

// Handler is the handler of a method
type Handler func(context.Context, interface{}, func(interface{}) error, []interceptor.ServerInterceptor) (interface{}, error)

func (s *service) Register(handlerName string, handler Handler) {
	if s.handlers == nil {
		s.handlers = make(map[string]Handler)
	}
	s.handlers[handlerName] = handler
}

func (s *service) Serve(opts *ServerOptions) {

	s.opts = opts

	transportOpts := []transport.ServerTransportOption{
		transport.WithServerAddress(s.opts.address),
		transport.WithServerNetwork(s.opts.network),
		transport.WithHandler(s),
		transport.WithServerTimeout(s.opts.timeout),
		transport.WithSerializationType(s.opts.serializationType),
		transport.WithProtocol(s.opts.protocol),
	}

	serverTransport := transport.GetServerTransport(s.opts.protocol)

	s.ctx, s.cancel = context.WithCancel(context.Background())

	if err := serverTransport.ListenAndServe(s.ctx, transportOpts...); err != nil {
		log.Errorf("%s serve error, %v", s.opts.network, err)
		return
	}

	fmt.Printf("%s service serving at %s ... \n", s.opts.protocol, s.opts.address)

	<-s.ctx.Done()
}

func (s *service) Close() {
	s.closing = true
	if s.cancel != nil {
		s.cancel()
	}
	consul.Delete()
	fmt.Println("service closing ...")
}

func (s *service) Name() string {
	return s.serviceName
}

func (s *service) Handle(ctx context.Context, reqbuf []byte) ([]byte, error) {

	// parse protocol header
	request := &protocol.Request{}
	if err := proto.Unmarshal(reqbuf, request); err != nil {
		return nil, err
	}

	// 全链路处理
	ctx = metadata.WithServerMetadata(ctx, request.Metadata)

	// 更新md
	request.Metadata = metadata.WithMetadataTimeout(request.Metadata, s.opts.timeout)

	// 根据新md做动作
	var cancel context.CancelFunc
	ctx, cancel = metadata.WithMetadataTimeoutContext(ctx, request.Metadata)
	defer cancel()

	serverSerialization := codec.GetSerialization(s.opts.serializationType)

	dec := func(req interface{}) error {

		if err := serverSerialization.Unmarshal(request.Payload, req); err != nil {
			return err
		}
		return nil
	}

	// 服务器超时, 处理结束取消所有后代
	if s.opts.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.opts.timeout)
		defer cancel()
	}

	_, method, err := utils.ParseServicePath(string(request.ServicePath))
	if err != nil {
		return nil, codes.New(codes.ClientMsgErrorCode, "method is invalid")
	}

	handler := s.handlers[method]
	if handler == nil {
		return nil, errors.New("handlers is nil")
	}

	var rsp interface{}
	ch := make(chan interface{})

	// 异步处理请求
	go func() {
		rsp, err = handler(ctx, s.svr, dec, s.opts.interceptors)
		ch <- struct {}{}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ServerRPCTimeoutErr
		case <-ch:
			if err != nil {
				return nil, err
			}

			rspBuf, err := serverSerialization.Marshal(rsp)
			if err != nil {
				return nil, err
			}
			return rspBuf, err
		}
	}
}
