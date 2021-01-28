package flow_compare

import (
	"context"
	"errors"
	"fmt"
	"github.com/13days/gfly/codec"
	"github.com/13days/gfly/protocol"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const FlowCompareTag = "flowCompare"

const methodRateSplit = "@"
const methodSplit = "$"

func RegisterFlowComparor(name string, c FlowComparor) {
	if flowComparorMap == nil {
		flowComparorMap = make(map[string]FlowComparor)
	}
	flowComparorMap[name] = c
}

func ExecFlowCompare(ctx context.Context, bs1 []byte, bs2 []byte, dur1, dur2 time.Duration, codec codec.Codec, rate int64) {
	rInt := rand.Int() % 100
	if rInt > int(rate) {
		return
	}
	rsp1 := getResp(bs1, codec)
	rsp2 := getResp(bs2, codec)
	flowReq := &FlowCompareReq{
		p1: rsp1,
		p2: rsp2,
		t1: dur1,
		t2: dur2,
	}

	for _, c := range flowComparorMap {
		c.Compare(ctx, flowReq)
	}
}

func getResp(bs []byte, codec codec.Codec) *protocol.Response {
	if len(bs) == 0 {
		return &protocol.Response{RetMsg: "len(bs) == 0"}
	}
	respBuf, err := codec.Decode(bs)
	if err != nil {
		fmt.Println("flow compare Decode err:", err)
		return nil
	}

	// parse protocol header
	resp := &protocol.Response{}
	if err = proto.Unmarshal(respBuf, resp); err != nil {
		fmt.Println("flow compare Unmarshal err:", err)
		return nil
	}
	return resp
}

func GetFlowComparor(name string) FlowComparor {
	if flowComparorMap, ok := flowComparorMap[name]; ok {
		return flowComparorMap
	}
	return &FlowComparorDefault{}
}

func GenFlowComparePath(serviceName, svrAddr string, method []string, rate []int) (string, error) {
	if len(method) != len(rate) {
		return "", errors.New("len(method) != len(rate)")
	}

	newPath := FlowCompareTag + "/" + serviceName + "/" + svrAddr
	for idx, name := range method {
		newPath += methodSplit + name + methodRateSplit + strconv.FormatInt(int64(rate[idx]), 10)
	}

	return newPath, nil
}

func ParseFlowComparePath(path string) (map[string]int64, string) {
	methodInfos := make(map[string]int64)
	methodRates := strings.Split(path, methodSplit)
	svrAddr := methodRates[0]
	for _, methodRate := range methodRates[1:] {
		methodInfo := strings.Split(methodRate, methodRateSplit)
		rate, err := strconv.ParseInt(methodInfo[1], 10, 64)
		if err != nil {
			return nil, ""
		}
		methodInfos[methodInfo[0]] = rate
	}
	return methodInfos, svrAddr
}
