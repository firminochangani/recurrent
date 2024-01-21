package main

import (
	"context"
	"fmt"
	"time"

	"github.com/flowck/recurrent/recurrent"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	s := recurrent.New()

	s.Every(time.Second * 1).Do(func(ctx context.Context) {
		fmt.Println("--->", time.Now())
	})

	s.Every(time.Second * 2).Do(func(ctx context.Context) {
		fmt.Println("--->", time.Now())
	})

	s.Run(ctx)
}
