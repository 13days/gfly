package log

import (
	"context"
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	Trace("test....")
	Tracef("testTracef...")
	Debug("test....")
	Debugf("testDebugf...")
	Info("test....")
	Infof("testInfof...")
	Warning("test....")
	Warningf("testWarningf...")
	Error("test....")
	Errorf("testErrorf...")
	Fatal("test....")
	Fatalf("testFatalf...")
}

func f() func() {
	return func() {
		fmt.Println(1)
	}
}
func TestFun(t *testing.T) {
	fmt.Println(1)
	f()
}

func TestLogs(t *testing.T) {
	context.Background()
	Logs("测试")
}