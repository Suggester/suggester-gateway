package main

import (
	"context"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Shard struct {
	ID     int
	config *Config

	Session *discordgo.Session
	ctx     context.Context
	wg      *sync.WaitGroup
	cancel  context.CancelFunc
}

func NewManagedShard(ctx context.Context, wg *sync.WaitGroup, cfg *Config, id int) Shard {
	ctx, cancel := context.WithCancel(ctx)
	return Shard{
		config: cfg,
		ID:     id,
		ctx:    ctx,
		wg:     wg,
		cancel: cancel,
	}
}

func (sh *Shard) Up() {
	if wg := sh.wg; wg != nil {
		defer wg.Done()
	}

	s, err := discordgo.New("Bot " + sh.config.Token)
	if err != nil {
		log.Fatalf("[Shard=%v] failed to create session: %v\n", sh.ID, err)
	}

	s.State = nil
	s.StateEnabled = false
	s.Identify.Shard = &[2]int{sh.ID, sh.config.Shards}
	s.Identify.Intents = discordgo.IntentsDirectMessages

	s.AddHandler(func(session *discordgo.Session, r *discordgo.Ready) {
		log.Printf("[Shard=%v] ready\n", sh.ID)
	})

	sh.Session = s

	err = s.Open()
	if err != nil {
		log.Fatalf("[Shard=%v] failed to open websocket connection: %v\n", sh.ID, err)
	}

	defer sh.Down()
	<-sh.ctx.Done()
}

func (sh *Shard) Down() {
	sh.cancel()
}
