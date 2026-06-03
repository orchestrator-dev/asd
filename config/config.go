package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Theme string `toml:"theme"`
}

func Load() (Config, error) {
	var cfg Config
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, err
	}

	path := filepath.Join(home, ".config", "asd", "config.toml")
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	_, err = toml.Decode(string(b), &cfg)
	return cfg, err
}
