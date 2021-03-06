package transport

import (
	"context"
	"errors"
	"fmt"
	"github.com/13days/gfly/codec"
	"github.com/13days/gfly/codes"
	"github.com/13days/gfly/flow_compare"
	"github.com/13days/gfly/plugin/consul"
	"net"
	"time"
)

type clientTransport struct {
	opts *ClientTransportOptions
}

var clientTransportMap = make(map[string]ClientTransport)

var RPCTimeoutErr = errors.New("call RPC timeout")

func init() {
	clientTransportMap["default"] = DefaultClientTransport
}

// RegisterClientTransport supports business custom registered ClientTransport
func RegisterClientTransport(name string, clientTransport ClientTransport) {
	if clientTransportMap == nil {
		clientTransportMap = make(map[string]ClientTransport)
	}
	clientTransportMap[name] = clientTransport
}

// Get the ServerTransport
func GetClientTransport(transport string) ClientTransport {

	if v, ok := clientTransportMap[transport]; ok {
		return v
	}

	return DefaultClientTransport
}

// The default ClientTransport
var DefaultClientTransport = New()

// Use the singleton pattern to create a ClientTransport
var New = func() ClientTransport {
	return &clientTransport{
		opts: &ClientTransportOptions{},
	}
}

func (c *clientTransport) Send(ctx context.Context, req []byte, opts ...ClientTransportOption) ([]byte, error) {

	for _, o := range opts {
		o(c.opts)
	}

	var rspBytes []byte
	var err error
	ch := make(chan interface{})
	go func() {
		if c.opts.Network == "tcp" {
			rspBytes, err = c.SendTcpReq(ctx, req)
			ch <- struct{}{}
		}

		if c.opts.Network == "udp" {
			rspBytes, err = c.SendUdpReq(ctx, req)
			ch <- struct{}{}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, RPCTimeoutErr
		case <-ch:
			return rspBytes, err
		}
	}
}

func (c *clientTransport) SendTcpReq(ctx context.Context, req []byte) ([]byte, error) {

	// service discovery
	addr, err := c.opts.Selector.Select(c.opts.ServiceName)
	if err != nil {
		return nil, err
	}

	// defaultSelector returns "", use the target as address
	if addr == "" {
		addr = c.opts.Target
	}

	t1 := time.Now()
	resp, err := c.realSendTcpReq(ctx, req, addr)
	dur1 := time.Now().Sub(t1)
	// flow compare discovery
	flowInfo, _ := c.opts.Selector.Select(consul.FlowCompareTag + "/" + c.opts.ServiceName)
	if flowInfo != "" {
		methodMap, svrAddr := flow_compare.ParseFlowComparePath(flowInfo)
		if rate, ok := methodMap[c.opts.Method]; ok {
			go func() {
				fmt.Println("flow compare req...")
				t2 := time.Now()
				newContext, _ := context.WithTimeout(context.Background(), c.opts.Timeout)
				resp1, err := c.realSendTcpReq(newContext, req, svrAddr)
				if err != nil {
					fmt.Println("flow req, err", err)
					return
				}
				fmt.Println("flow compare resp...")
				dur2 := time.Now().Sub(t2)
				flow_compare.ExecFlowCompare(ctx, resp, resp1, dur1, dur2, codec.GetCodec(c.opts.Protocol), rate)
			}()
		}
	}

	return resp, err
}

func (c *clientTransport) realSendTcpReq(ctx context.Context, req []byte, addr string) ([]byte, error) {
	conn, err := c.opts.Pool.Get(ctx, c.opts.Network, addr)
	//	conn, err := net.DialTimeout("tcp", addr, c.opts.Timeout);
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	sendNum := 0
	num := 0
	for sendNum < len(req) {
		num, err = conn.Write(req[sendNum:])
		if err != nil {
			return nil, err
		}
		sendNum += num

		if err = isDone(ctx); err != nil {
			return nil, err
		}
	}

	// parse frame
	wrapperConn := wrapConn(conn)
	frame, err := wrapperConn.framer.ReadFrame(conn)
	if err != nil {
		return nil, err
	}

	return frame, err
}

func (c *clientTransport) SendUdpReq(ctx context.Context, req []byte) ([]byte, error) {
	// service discovery
	addr, err := c.opts.Selector.Select(c.opts.ServiceName)
	if err != nil {
		return nil, err
	}

	// defaultSelector returns "", use the target as address
	if addr == "" {
		addr = c.opts.Target
	}

	udpAddr, err := net.ResolveUDPAddr(c.opts.Network, addr)
	if err != nil {
		return nil, codes.NewFrameworkError(codes.ClientMsgErrorCode, "addr invalid ...")
	}

	conn, err := net.DialUDP(c.opts.Network, nil, udpAddr)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if n, err := conn.Write(req); n != len(req) || err != nil {
		return nil, err
	}

	recvBuf := make([]byte, 65536)
	n, err := conn.Read(recvBuf)
	if err != nil {
		return nil, err
	}

	rsp := recvBuf[:n]

	return rsp, nil
}

func isDone(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return nil
}
