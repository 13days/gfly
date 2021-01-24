package main

import (
	"github.com/13days/gfly/testdata"
	"github.com/13days/gfly/testdata/nodeB"
)

func main() {
	testdata.StartServer(testdata.ServerB,  "127.0.0.1:8003", new(nodeB.Service))
}
