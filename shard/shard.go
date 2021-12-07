package shard

import (
	"context"
	"log"
	"sync"

	"github.com/suggester/suggester-gateway/config"

	"github.com/bwmarrin/discordgo"
)

type Shard struct {
	sync.Mutex

	ID     int
	config *config.Config

	Session *discordgo.Session
	ctx     context.Context
	wg      *sync.WaitGroup
	done    chan struct{}
}

func NewManaged(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, id int) Shard {
	return Shard{
		config: cfg,
		ID:     id,
		ctx:    ctx,
		wg:     wg,
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

	s.AddHandler(func(session *discordgo.Session, r *discordgo.Ready) {
		log.Printf("[Shard=%v] ready\n", sh.ID)
	})

	sh.Session = s

	err = s.Open()
	if err != nil {
		log.Fatalf("[Shard=%v] failed to open websocket connection: %v\n", sh.ID, err)
	}

	select {
	case <-sh.ctx.Done():
	case <-sh.done:
		sh.Down()
	}
	log.Printf("[Shard=%v] stopping shard", sh.ID)
}

func (sh *Shard) Down() {
	sh.done <- struct{}{}
}
