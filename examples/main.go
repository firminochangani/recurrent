package main

import (
	"context"
	"fmt"
	"time"

	"github.com/flowck/recurrently/schedule"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	s := schedule.New()

	s.Every(1).Seconds().Do(func(ctx context.Context) {
		fmt.Println("--->", time.Now())
	})

	s.Every(2).Seconds().Do(func(ctx context.Context) {
		fmt.Println("--->", time.Now())
	})

	s.Run(ctx)
}
