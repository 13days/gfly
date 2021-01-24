package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	REQUEST_ID    = "request_id"
	TIME_TERMINAL = "time_terminal"
)

type clientMD struct {}
type serverMD struct {}

type clientMetadata map[string][]byte

type serverMetadata map[string][]byte

// ClientMetadata creates a new context with key-value pairs attached.
func ClientMetadata(ctx context.Context) clientMetadata {
	if md, ok := ctx.Value(clientMD{}).(clientMetadata); ok {
		return md
	}
	md := make(map[string][]byte)
	WithClientMetadata(ctx, md)
	return md
}

// WithClientMetadata creates a new context with the specified metadata
func WithClientMetadata(ctx context.Context, metadata map[string][]byte) context.Context{
	return context.WithValue(ctx, clientMD{}, clientMetadata(metadata))
}

// ServerMetadata creates a new context with key-value pairs attached.
func ServerMetadata(ctx context.Context) serverMetadata {
	if md, ok := ctx.Value(serverMD{}).(serverMetadata); ok {
		return md
	}
	md := make(map[string][]byte)
	WithServerMetadata(ctx, md)
	return md
}

// WithServerMetadata creates a new context with the specified metadata
func WithServerMetadata(ctx context.Context, metadata map[string][]byte) context.Context{
	return context.WithValue(ctx, serverMD{}, serverMetadata(metadata))
}

func WithMetadataTimeoutContext(ctx context.Context, metadata map[string][]byte) (context.Context, context.CancelFunc) {
	timeData, ok := metadata[TIME_TERMINAL]
	if !ok {
		return ctx, func() {}
	}
	timeDeadline := &time.Time{}
	err := json.Unmarshal(timeData, timeDeadline)
	if err != nil {
		return ctx, func() {}
	}
	if dur := timeDeadline.Sub(time.Now()); dur > 0 {
		fmt.Println("dur:", dur)
		// 服务端根据client传来的md更新context
		ctx = context.WithValue(ctx, TIME_TERMINAL, time.Now().Add(dur))
		return context.WithTimeout(ctx, dur)
	}
	return ctx, func() {}
}

func WithMetadataTimeout(md clientMetadata, dur time.Duration) clientMetadata {
	fmt.Println("WithMetadataTimeout: ", dur)
	fmt.Printf("md:%v\n", md)
	timeData, ok := md[TIME_TERMINAL]
	if !ok {
		fmt.Println("not ok")
		return addTimeout(md, dur)
	}
	timeDeadline := &time.Time{}
	err := json.Unmarshal(timeData, timeDeadline)
	if err != nil {
		fmt.Printf("json.Unmarshal err %v\n", err)
		return addTimeout(md, dur)
	}
	fmt.Printf("ok, time:%v\n", timeDeadline)
	// oldDur来自上游, dur来自设置
	if oldDur := timeDeadline.Sub(time.Now()); oldDur > dur {
		fmt.Println("new dur:", dur)
		return addTimeout(md, dur)
	} else {
		fmt.Printf("old dur:%v, new dur:%v\n", oldDur, dur)
	}
	return md
}

func ContextToMetaData(ctx context.Context, md clientMetadata) clientMetadata {
	fmt.Println("ctx1:", ctx)
	val := ctx.Value(TIME_TERMINAL)
	fmt.Println("time val", val)
	if val == nil {
		return md
	}
	bytes, err := json.Marshal(val)
	if err != nil {
		return md
	}
	md[TIME_TERMINAL] = bytes
	return md
}

func addTimeout(md clientMetadata, dur time.Duration) clientMetadata {
	timeDeadline := time.Now().Add(dur)
	if bytes, err := json.Marshal(timeDeadline); err == nil {
		md[TIME_TERMINAL] = bytes
	}
	return md
}