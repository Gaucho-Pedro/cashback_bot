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

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		file, err := os.Open("configs/config.json")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(instance); err != nil {
			log.Fatal(err)
		}
	})
	return instance
}
