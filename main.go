package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"time"
)

var cfg Config

func init() {
	cfgPath := flag.String("c", "config.toml", "path to config.toml")
	cfg = ParseConfig(*cfgPath)
}

func main() {
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, os.Interrupt)
		<-sc

		cancel()
	}()

CreateShard:
	for i := 0; i < cfg.Shards; i++ {
		select {
		case <-ctx.Done():
			break CreateShard
		default:
			wg.Add(1)

			sh := NewManagedShard(ctx, &cfg, i)
			go sh.Up(wg)
			<-time.NewTimer(time.Millisecond * 5_500).C
		}
	}

	wg.Wait()
}
