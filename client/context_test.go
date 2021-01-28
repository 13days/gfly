package client

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T)  {
	messages := make(chan int, 10)

	// producer
	for i := 0; i < 10; i++ {
		messages <- i
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// consumer
	go func(ctx context.Context) {
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <-ctx.Done():
				fmt.Println("child process interrupt...")
				return
			default:
				fmt.Printf("send message: %d\n", <-messages)
			}
		}
	}(ctx)

	defer close(messages)
	defer cancel()

	select {
	case <-ctx.Done():
		time.Sleep(1 * time.Second)
		fmt.Println("main process exit!")
	}
}

func TestParentTimeout(t *testing.T)  {
	ctx := context.Background()
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("first parent exit!")
		}
	}()

	ctx = context.WithValue(ctx, "key", "val")
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("second parent exit!")
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, time.Hour)
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("three parent exit!")
		}
	}()

	messages := make(chan int, 10)
	// producer
	for i := 0; i < 10; i++ {
		messages <- i
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	// consumer
	go func(ctx context.Context) {
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <-ctx.Done():
				fmt.Println("child process interrupt...")
				return
			default:
				fmt.Printf("send message: %d\n", <-messages)
			}
		}
	}(ctx)

	defer close(messages)
	defer cancel()


	// 子节点
	ctx = context.WithValue(ctx, "key", "val")
	ctx = context.WithValue(ctx, "key", "val")

	select {
	case <-ctx.Done():
		time.Sleep(2 * time.Second)
		fmt.Println("main process exit!")
	}
}


func TestContext1(t *testing.T){
	a, b, c := 1, 2, 3
	if a == 1 && b != 2 || c ==3 {
		fmt.Println("ok")
	}
}