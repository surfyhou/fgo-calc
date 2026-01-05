package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port        string `json:"port"`
	DataDir     string `json:"data_dir"`
	MaxCost     int    `json:"max_cost"`
	MaxSvtLimit int    `json:"max_svt_limit"`
	MaxCeLimit  int    `json:"max_ce_limit"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
