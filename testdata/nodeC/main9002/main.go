package main

import (
	"github.com/13days/gfly/testdata"
	"github.com/13days/gfly/testdata/nodeC"
)

func main() {
	testdata.StartServer(testdata.ServerC,  "127.0.0.1:9002", new(nodeC.Service))
}
