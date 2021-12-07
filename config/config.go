package config

import (
	"io/ioutil"
	"log"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Token  string `toml:"token"`
	Shards int    `toml:"shards"`
}

func Parse(file string) Config {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("error reading config file: %v\n", err)
	}

	var cfg Config
	err = toml.Unmarshal(f, &cfg)
	if err != nil {
		log.Fatalf("error parsing config: %v\n", err)
	}
	return cfg
}
