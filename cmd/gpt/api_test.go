package gpt_test

import (
	"QQ-ChatGPT-Bot/cmd/gpt"
	config2 "QQ-ChatGPT-Bot/config"
	"testing"
)

// 测试API能否正常使用
func TestAPI(t *testing.T) {
	var config config2.Config

	config.OpenAi.Url = ""
	config.OpenAi.ApiKey = ""
	config.OpenAi.Model = "gpt-3.5-turbo-0613"
	config.OpenAi.MaxTokens = 2048

	config.RolePlay.Role = []string{"user", "assistant"}
	config.RolePlay.Content = []string{"你的名字是李天所", "我的名字叫李天所"}

	config2.Cfg = config

	//FIXME:角色扮演似乎仍然有问题。
	result, _ := gpt.GenerateText("test", "你的名字是什么？", false)
	t.Log(result)
}
