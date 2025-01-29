package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	stability "distributed.patterns/Stability"
)

var count int

func EmulateTransientError(ctx context.Context) (string, error) {
	count++

	if count <= 3 {
		return "intentional fail", errors.New("error")
	} else {
		return "success", nil
	}
}

func main() {
	r := stability.Retry(EmulateTransientError, 5, 2*time.Second)
	res, err := r(context.Background())
	fmt.Println(res, err)
}
