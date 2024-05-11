package config

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
		Url         string `mapstructure:"url"` //镜像站点链接，默认使用openai站点
		ApiKey      string `mapstructure:"api_key"`
		Model       string
		Temperature float64
		TopP        float64 `mapstructure:"top_p"`
		MaxTokens   int     `mapstructure:"max_tokens"`
		UseProxy    bool    `mapstructure:"use_proxy"`
		ProxyUrl    string  `mapstructure:"proxy_url"`
	}
	RolePlay struct {
		Role    []string `mapstructure:"role"`    //system,assistant,user（目前支持）
		Content []string `mapstructure:"content"` //角色对应的对话内容
	}
	Context struct {
		PrivateContext bool `mapstructure:"private_context"`
		GroupContext   bool `mapstructure:"group_context"`
	}
}

var Cfg Config

//func init() {
//	log.SetFlags(log.Lshortfile | log.LstdFlags)
//	if _, err := os.Stat("config.cfg"); os.IsNotExist(err) {
//		f, err := os.Create("config.cfg")
//		if err != nil {
//			log.Println(err)
//		}
//		// 自动生成配置文件
//		_, err = f.Write([]byte("# config.toml 配置文件\n\n" +
//			"# cqhttp机器人配置\n[cqhttp]\n" +
//			"# go-cqhttp的正向WebSocket地址\n" +
//			"websocket = \"ws://127.0.0.1:8080\"\n" +
//			"# 群聊是否需要@机器人才能触发\n" +
//			"at_only = true\n" +
//			"# 是否开启触发关键词\n" +
//			"use_keyword = false\n" +
//			"# 触发关键词场合 可选值: all, group, private, 开启群聊关键词建议关闭at_only\n" +
//			"keyword_type = \"group\"\n" +
//			"# 触发关键词\n" +
//			"keyword = \"对话\"\n" +
//			"# 生成中提醒时间秒数\n" +
//			"timeout = 30\n\n" +
//			"# openai配置\n[openai]\n" +
//			"# 镜像站点链接，默认使用openai站点\n" +
//			"url = \"https://api.openai.com/v1/chat/completions\"\n" +
//			"# 你的 OpenAI API Key, 可以在 https://beta.openai.com/account/api-keys 获取\n" +
//			"api_key = \"sk-xxxxxx\"\n" +
//			"# 使用的模型，默认是 gpt-3.5-turbo\n" +
//			"model = \"gpt-3.5-turbo\"\n" +
//			"# 对话温度，越大越随机 参照https://algowriting.medium.com/gpt-3-temperature-setting-101-41200ff0d0be\n" +
//			"temperature = 0.3\n" +
//			"# Top-p，越大越随机\n" +
//			"top_p = 0.9\n" +
//			"# 每次对话最大生成字符数\n" +
//			"max_tokens = 1000\n" +
//			"# openai是否走代理，默认关闭\n" +
//			"use_proxy = false\n" +
//			"# 代理地址\n" +
//			"proxy_url = \"http://127.0.0.1:7890\"\n\n" +
//			"# 角色信息配置。如果关闭角色扮演，请删除role一栏[]中的所有内容\n[roleplay]\n" +
//			"# 角色列表\n" +
//			"role = [\"system\", \"user\", \"assistant\"]\n" +
//			"# 角色对应的对话内容\n" +
//			"content = [\"You are a helpful assistant.\", \"你好\", \"你好，有什么我可以帮助你的吗\"]\n\n" +
//			"# 是否在私聊中启用连续对话\n" +
//			"private_context = false\n" +
//			"# 是否在群聊中启用连续对话\n" +
//			"group_context = false\n"))
//		if err != nil {
//			log.Println(err)
//		}
//		log.Println("配置文件不存在, 已自动生成配置文件, 请修改配置文件后再次运行程序, 5秒后退出程序...")
//		time.Sleep(5 * time.Second)
//		os.Exit(0)
//	}
//	viper.SetConfigName("config")
//	viper.SetConfigType("toml")
//	viper.AddConfigPath(".") // 指定查找配置文件的路径
//	err := viper.ReadInConfig()
//	if err != nil {
//		log.Fatalf("read config failed: %v", err)
//	}
//	err = viper.Unmarshal(&Cfg)
//	if err != nil {
//		log.Fatalf("unmarshal config failed: %v", err)
//	}
//
//	if Cfg.OpenAi.Url == "" {
//		Cfg.OpenAi.Url = "https://api.openai.com/v1/chat/completions"
//	}
//}
