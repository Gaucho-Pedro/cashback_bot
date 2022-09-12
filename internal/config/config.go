package config

import (
	"encoding/json"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel    string
	BotDebug    bool
	BotToken    string
	AdminChatID int64
}

var config Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		file, err := os.Open("configs/config.json")
		if err != nil {
			log.Fatal("f", err)
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&config); err != nil {
			log.Fatal(err)
		}
	})
	return &config
}
