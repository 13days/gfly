package context

import (
	"context"
	"encoding/json"
	"time"
)

const (
	REQUEST_ID    = "request_id"
	TIME_TERMINAL = "time_terminal"
)

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, REQUEST_ID, requestId)
}

func WithTimeout(ctx context.Context, dur time.Duration) context.Context {
	val := time.Now().Add(dur)
	return context.WithValue(ctx, TIME_TERMINAL, &val)
}

func GetRequestId(ctx context.Context) ([]byte, error) {
	v := ctx.Value(REQUEST_ID)
	if v == nil {
		return nil, nil
	}
	return []byte(v.(string)), nil
}

func GetTerminalTime(ctx context.Context) ([]byte, error) {
	v := ctx.Value(TIME_TERMINAL)
	if v == nil {
		return nil, nil
	}
	return json.Marshal(v)
}

func GetRpcParamsMeta(ctx context.Context) map[string][]byte {
	meta := make(map[string][]byte)
	if data, err := GetRequestId(ctx); err == nil {
		meta[REQUEST_ID] = data
	}
	if data, err := GetTerminalTime(ctx); err == nil {
		meta[TIME_TERMINAL] = data
	}
	return meta
}
