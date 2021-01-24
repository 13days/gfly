package main

import (
	"github.com/13days/gfly/testdata"
	"github.com/13days/gfly/testdata/nodeA"
)

func main() {
	testdata.StartServer(testdata.ServerA,  "127.0.0.1:7001", new(nodeA.Service))
}
