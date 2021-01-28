package transport

import (
	"github.com/13days/gfly/pool/connpool"
	"github.com/13days/gfly/selector"
	"time"
)

// ClientTransportOptions includes all ClientTransport parameter options
type ClientTransportOptions struct {
	Target      string
	ServiceName string
	Method      string
	Network     string
	Protocol	string
	Pool        connpool.Pool
	Selector    selector.Selector
	Timeout     time.Duration
}

// Use the Options mode to wrap the ClientTransportOptions
type ClientTransportOption func(*ClientTransportOptions)

// WithServiceName returns a ClientTransportOption which sets the value for serviceName
func WithServiceName(serviceName string) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.ServiceName = serviceName
	}
}

func WithMethod(method string) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Method = method
	}
}

// WithClientTarget returns a ClientTransportOption which sets the value for target
func WithClientTarget(target string) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Target = target
	}
}

// WithClientNetwork returns a ClientTransportOption which sets the value for network
func WithClientNetwork(network string) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Network = network
	}
}

// WithClientPool returns a ClientTransportOption which sets the value for pool
func WithClientPool(pool connpool.Pool) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Pool = pool
	}
}

// WithSelector returns a ClientTransportOption which sets the value for selector
func WithSelector(selector selector.Selector) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Selector = selector
	}
}

// WithTimeout returns a ClientTransportOption which sets the value for timeout
func WithTimeout(timeout time.Duration) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Timeout = timeout
	}
}
