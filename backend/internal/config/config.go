package config

import "os"

type Config struct {
	DataDir string
	Port    string
}

func LoadConfig() *Config {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "../data" // 默认假设在 backend 目录下运行，数据在上一层
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":30006"
	}
	return &Config{
		DataDir: dataDir,
		Port:    port,
	}
}

