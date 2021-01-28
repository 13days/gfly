package consul

import (
	"errors"
	"fmt"
	"github.com/13days/gfly/flow_compare"
	"github.com/13days/gfly/plugin"
	"github.com/13days/gfly/selector"
	"github.com/hashicorp/consul/api"
	"net/http"
	"strings"
	"time"
)

// Consul implements the server discovery specification
type Consul struct {
	opts         *plugin.Options
	client       *api.Client
	config       *api.Config
	balancerName string // load balancing mode, including random, polling, weighted polling, consistent hash, etc
	writeOptions *api.WriteOptions
	queryOptions *api.QueryOptions
	closeChan    chan struct{}
	doneChan     chan struct{}
	keys     []string
}

const Name = "consul"
const FlowCompareTag = "flowCompare"

func init() {
	plugin.Register(Name, ConsulSvr)
	selector.RegisterSelector(Name, ConsulSvr)
}

// global consul objects for framework
var ConsulSvr = &Consul{
	opts:      &plugin.Options{},
	closeChan: make(chan struct{}),
	doneChan:  make(chan struct{}),
}

func (c *Consul) InitConfig() error {

	config := api.DefaultConfig()
	c.config = config

	config.HttpClient = http.DefaultClient
	config.Address = c.opts.SelectorSvrAddr
	config.Scheme = "http"

	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	c.client = client

	return nil
}

func (c *Consul) Resolve(serviceName string) ([]*selector.Node, error) {
	pairs, _, err := c.client.KV().List(serviceName, nil)
	if err != nil {
		return nil, err
	}

	if len(pairs) == 0 {
		return nil, fmt.Errorf("no services find in path : %s", serviceName)
	}
	var nodes []*selector.Node
	for _, pair := range pairs {
		nodes = append(nodes, &selector.Node{
			Key:   pair.Key,
			Value: pair.Value,
		})
		//fmt.Printf("key:%v, val:%v\n", pair.Key, pair.Value)
	}
	return nodes, nil
}

// implements selector Select method
func (c *Consul) Select(serviceName string) (string, error) {

	fmt.Println("serviceName:", serviceName)
	nodes, err := c.Resolve(serviceName)

	if nodes == nil || len(nodes) == 0 || err != nil {
		return "", err
	}

	balancer := selector.GetBalancer(c.balancerName)
	node := balancer.Balance(serviceName, nodes)

	if node == nil {
		return "", fmt.Errorf("no services find in %s", serviceName)
	}

	return parseAddrFromNode(node)
}

func parseAddrFromNode(node *selector.Node) (string, error) {
	if node.Key == "" {
		return "", errors.New("addr is empty")
	}

	strs := strings.Split(node.Key, "/")

	return strs[len(strs)-1], nil
}

func (c *Consul) Init(opts ...plugin.Option) error {

	for _, o := range opts {
		o(c.opts)
	}

	if len(c.opts.Services) == 0 || c.opts.SvrAddr == "" || c.opts.SelectorSvrAddr == "" {
		return fmt.Errorf("consul init error, len(services) : %d, svrAddr : %s, selectorSvrAddr : %s",
			len(c.opts.Services), c.opts.SvrAddr, c.opts.SelectorSvrAddr)
	}

	if err := c.InitConfig(); err != nil {
		return err
	}

	defer c.asyncListenServicesSingOut()

	var err error
	for _, serviceName := range c.opts.Services {
		nodeName := fmt.Sprintf("%s/%s", serviceName, c.opts.SvrAddr)

		// 流量对比注册
		if len(c.opts.FlowCompareMethods) != 0 {
			nodeName, err = flow_compare.GenFlowComparePath(serviceName, c.opts.SvrAddr, c.opts.FlowCompareMethods, c.opts.FlowCompareRate)
			fmt.Println("nodeName:", nodeName)
			if err != nil {
				return err
			}
		}

		kvPair := &api.KVPair{
			Key:   nodeName,
			Value: []byte(c.opts.SvrAddr),
			Flags: api.LockFlagValue,
		}

		fmt.Println(kvPair)
		if _, err = c.client.KV().Put(kvPair, c.writeOptions); err != nil {
			return err
		}
		c.keys = append(c.keys, nodeName)
	}

	return nil
}

func (c *Consul) asyncListenServicesSingOut()  {
	// 监听每个service下线
	go func() {
		select {
		case <-c.closeChan:
			for _, key := range c.keys {
				c.client.KV().Delete(key, c.writeOptions)
			}
			c.doneChan <- struct{}{}
		}
	}()
}

// Init implements the initialization of the consul configuration when the framework is loaded
func Init(consulSvrAddr string, opts ...plugin.Option) error {
	for _, o := range opts {
		o(ConsulSvr.opts)
	}

	ConsulSvr.opts.SelectorSvrAddr = consulSvrAddr
	err := ConsulSvr.InitConfig()
	return err
}

func Delete() {
	ConsulSvr.closeChan <- struct{}{}
	select {
	case <-ConsulSvr.doneChan:
		// 延时一段时间退出
		time.Sleep(time.Second * 3)
		return
	}
}
