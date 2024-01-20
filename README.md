# schedule

[![main](https://github.com/flowck/schedule/actions/workflows/main.yml/badge.svg)](https://github.com/flowck/schedule/actions/workflows/main.yml)

Golang job scheduling for humans. Run Go functions periodically using a friendly syntax. - Inspired by Python lib [schedule](https://github.com/dbader/schedule)

## Usage

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/flowck/schedule"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := schedule.New()

	s.Every(10).Seconds().Do(func(ctx context.Context) {
		fmt.Println("--->", time.Now())
	})

	s.Every(30).Seconds().Do(func(ctx context.Context) {
		fmt.Println("--->", time.Now())
	})

	s.Run(ctx)
}

```