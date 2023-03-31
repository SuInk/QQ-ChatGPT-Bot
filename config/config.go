package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type Config struct {
	CqHttp struct {
		WebSocket   string `mapstructure:"websocket"`
		AtOnly      bool   `mapstructure:"at_only"`
		UseKeyword  bool   `mapstructure:"use_keyword"`
		KeywordType string `mapstructure:"keyword_type"`
		Keyword     string `mapstructure:"keyword"`
		TimeOut     int    `mapstructure:"timeout"`
	}
	OpenAi struct {
		ApiKey      string `mapstructure:"api_key"`
		Model       string
		Temperature float32
		MaxTokens   int    `mapstructure:"max_tokens"`
		UseProxy    bool   `mapstructure:"use_proxy"`
		ProxyUrl    string `mapstructure:"proxy_url"`
	}
	Identity struct {
		UseIdentity bool     `mapstructure:"use_identity"`
		Prompt      string   `mapstructure:"prompt"`
		Stop        []string `mapstructure:"stop"`
	}
	Context struct {
		PrivateContext bool `mapstructure:"private_context"`
		GroupContext   bool `mapstructure:"group_context"`
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
		// 自动生成配置文件
		_, err = f.Write([]byte("# config.toml 配置文件\n\n" +
			"# cqhttp机器人配置\n[cqhttp]\n" +
			"# go-cqhttp的正向WebSocket地址\n" +
			"websocket = \"ws://127.0.0.1:8080\"\n" +
			"# 群聊是否需要@机器人才能触发\n" +
			"at_only = true\n" +
			"# 是否开启触发关键词\n" +
			"use_keyword = false\n" +
			"# 触发关键词场合 可选值: all, group, private, 开启群聊关键词建议关闭at_only\n" +
			"keyword_type = \"group\"\n" +
			"# 触发关键词\n" +
			"keyword = \"对话\"\n" +
			"# 生成中提醒时间秒数\n" +
			"timeout = 30\n\n" +
			"# openai配置\n[openai]\n" +
			"# 你的 OpenAI API Key, 可以在 https://beta.openai.com/account/api-keys 获取\n" +
			"api_key = \"sk-xxxxxx\"\n" +
			"# 使用的模型，默认是 gpt-3.5-turbo\n" +
			"model = \"gpt-3.5-turbo\"\n" +
			"# 对话温度，越大越随机 参照https://algowriting.medium.com/gpt-3-temperature-setting-101-41200ff0d0be\n" +
			"temperature = 0.3\n" +
			"# 每次对话最大生成字符数\n" +
			"max_tokens = 1000\n" +
			"# openai是否走代理，默认关闭\n" +
			"use_proxy = false\n" +
			"# 代理地址\n" +
			"proxy_url = \"http://127.0.0.1:7890\"\n\n" +
			"# 角色信息配置\n[identity]\n" +
			"# 角色预设功能，默认关闭\n" +
			"use_identity = false\n" +
			"# 角色预设信息(设定可以参考：https://github.com/easydu2002/chat_gpt_oicq/wiki/设定AI人格---以猫娘为案例【chatGPT猫娘】）\n" +
			"prompt = \"（你扮演的角色名称）:你要求AI扮演的角色信息\\n（AI扮演的角色名称）:AI的回应\"\n" +
			"# 扮演的身份名称（前面填对话者，后面填bot要扮演的角色）\n" +
			"stop = [\"（你扮演的角色名称）:\", \"（AI扮演的角色名称）:\"]\n\n" +
			"# 连续对话相关（实际使用中，连续对话似乎会导致更多的token使用，在这里可以设置是否启用这个功能。默认关闭。另注：预设角色不支持连续对话。）\n[context]\n" +
			"# 是否在私聊中启用连续对话\n" +
			"private_context = false\n" +
			"# 是否在群聊中启用连续对话\n" +
			"group_context = false\n"))
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
