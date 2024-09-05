package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	GitRepoPath string `json:"gitRepoPath"`
}

func Load() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
