package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	CqHttp struct {
		WebSocket string `mapstructure:"websocket"`
		AtOnly    bool   `mapstructure:"at_only"`
		TimeOut   int    `mapstructure:"timeout"`
	}
	OpenAi struct {
		ApiKey      string `mapstructure:"api_key"`
		Model       string
		Temperature float32
		MaxTokens   int `mapstructure:"max_tokens"`
	}
}

var Cfg Config

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".") // 指定查找配置文件的路径
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read config failed: %v", err)
	}
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		log.Fatalf("unmarshal config failed: %v", err)
	}
}
