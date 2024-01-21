# recurrent

[![main](https://github.com/flowck/schedule/actions/workflows/main.yml/badge.svg)](https://github.com/flowck/schedule/actions/workflows/main.yml)

A Go package to run tasks recurrently. - Inspired by the Python lib [schedule](https://github.com/dbader/schedule) 

- [x] Parallel execution of jobs
- [x] Cancellation of running jobs via context 

## Usage

```go
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

```