package context

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestGetRequestId(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestId(ctx, "1")
	fmt.Println(GetRequestId(ctx))
}

func TestGetTerminalTime(t *testing.T) {
	ctx := context.Background()
	ctx = WithTimeout(ctx, 100 * time.Second)
	fmt.Println(GetTerminalTime(ctx))
}

func TestGetRpcParamsMeta(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestId(ctx, "1324242342")
	oldT := 100 * time.Second
	ctx = WithTimeout(ctx, oldT)
	fmt.Println(ctx.Value(TIME_TERMINAL))
	meta := GetRpcParamsMeta(ctx)
	fmt.Println(meta)
	data,_ := json.Marshal(meta)
	tt := make(map[string][]byte)
	json.Unmarshal(data, &tt)
	newT := &time.Time{}
	json.Unmarshal(tt[TIME_TERMINAL], newT)
	fmt.Println(tt)
	fmt.Println(newT)
}