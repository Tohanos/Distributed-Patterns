package stability

import (
	"context"
	"fmt"
	"time"
)

type SlowFunction func(string) (string, error)
type WithContext func(context.Context, string) (string, error)

func Timeout(f SlowFunction) WithContext {
	return func(ctx context.Context, arg string) (string, error) {
		chres := make(chan string)
		cherr := make(chan error)

		go func() {
			res, err := f(arg)
			chres <- res
			cherr <- err
		}()

		select {
		case res := <-chres:
			return res, <-cherr
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

var Slow SlowFunction = func(s string) (string, error) {
	time.Sleep(5 * time.Second)
	return "Slow function started: " + s, nil
}

func UsageExample() {
	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	timeout := Timeout(Slow)
	res, err := timeout(ctxt, "Go go go")

	fmt.Println(res, err)
}
