package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
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
		MaxTokens   int    `mapstructure:"max_tokens"`
		UseProxy    bool   `mapstructure:"use_proxy"`
		ProxyUrl    string `mapstructure:"proxy_url"`
	}
}

var Cfg Config

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	if _, err := os.Stat("config.cfg"); os.IsNotExist(err) {
		f, err := os.Create("config.cfg")
		if err != nil {
			log.Println(err)
		}
		_, err = f.Write([]byte("# config.cfg 配置文件\n\n# cqhttp机器人配置\n[cqhttp]\n# go-cqhttp的正向WebSocket地址\nwebsocket = \"ws://127.0.0.1:8080\"\n# 是否需要@机器人才能触发\nat_only = true\n# 生成中提醒时间秒数\ntimeout = 30\n\n# openai配置\n[openai]\n# 你的 OpenAI API Key, 可以在 https://beta.openai.com/account/api-keys 获取\napi_key = \"sk-xxxxxx\"\n# 使用的模型，默认是 gpt-3.5-turbo\nmodel = \"gpt-3.5-turbo\"\n# 对话温度，越大越随机 参照https://algowriting.medium.com/gpt-3-temperature-setting-101-41200ff0d0be\ntemperature = 0.3\n# 每次对话最大生成字符数\nmax_tokens = 1000\n# openai是否走代理，默认关闭\nuse_proxy = false\n# 代理地址\nproxy_url = \"http://127.0.0.1:7890\"\n\n"))
		if err != nil {
			log.Println(err)
		}
		log.Println("配置文件不存在, 已自动生成配置文件, 请修改配置文件后再次运行程序, 5秒后退出程序...")
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
	viper.SetConfigName("config.cfg")
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
