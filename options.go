package gfly

import (
	"github.com/13days/gfly/interceptor"
	"time"
)

// ServerOptions defines the server serve parameters
type ServerOptions struct {
	address           string        // listening address, e.g. :( ip://127.0.0.1:8080、 dns://www.google.com)
	network           string        // network type, e.g. : tcp、udp
	protocol          string        // protocol type, e.g. : proto、json
	timeout           time.Duration // timeout
	serializationType string        // serialization type, default: proto

	selectorSvrAddr    string   // service discovery server address, required when using the third-party service discovery plugin
	tracingSvrAddr     string   // tracing plugin server address, required when using the third-party tracing plugin
	tracingSpanName    string   // tracing span name, required when using the third-party tracing plugin
	flowCompareMethods []string // flow compare method, when then method of the service need compare, add to it
	flowCompareRate    []int    // flow compare method percentage rate, the order as same as flowCompareMethods, e.g. : 1=1%,10=10%
	pluginNames        []string // plugin name
	interceptors       []interceptor.ServerInterceptor
}

type ServerOption func(*ServerOptions)

func WithAddress(address string) ServerOption {
	return func(options *ServerOptions) {
		options.address = address
	}
}

func WithNetwork(network string) ServerOption {
	return func(o *ServerOptions) {
		o.network = network
	}
}

func WithProtocol(protocol string) ServerOption {
	return func(o *ServerOptions) {
		o.protocol = protocol
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.timeout = timeout
	}
}

func WithSerializationType(serializationType string) ServerOption {
	return func(o *ServerOptions) {
		o.serializationType = serializationType
	}
}

func WithSelectorSvrAddr(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.selectorSvrAddr = addr
	}
}

func WithPlugin(pluginName ...string) ServerOption {
	return func(o *ServerOptions) {
		o.pluginNames = append(o.pluginNames, pluginName...)
	}
}

func WithInterceptor(interceptors ...interceptor.ServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.interceptors = append(o.interceptors, interceptors...)
	}
}

func WithTracingSvrAddr(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.tracingSvrAddr = addr
	}
}

func WithTracingSpanName(name string) ServerOption {
	return func(o *ServerOptions) {
		o.tracingSpanName = name
	}
}

func WithFlowCompareMethod(method string, rate int) ServerOption {
	return func(o *ServerOptions) {
		o.flowCompareMethods = append(o.flowCompareMethods, method)
		o.flowCompareRate = append(o.flowCompareRate, rate)
	}
}
