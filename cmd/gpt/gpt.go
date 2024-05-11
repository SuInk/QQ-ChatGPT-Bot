package gpt

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"QQ-ChatGPT-Bot/config"
)

var Cache CacheInterface

func init() {
	Cache = GetSessionCache()
}

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request body
type postData struct {
	Model       string     `json:"model"`
	Messages    []Messages `json:"messages"`
	MaxTokens   int        `json:"max_tokens"`
	Temperature float64    `json:"temperature"`
	TopP        float64    `json:"top_p"`
}

// OpenAiRcv 对话使用的Response
type OpenAiRcv struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      Messages `json:"message"`
		FinishReason string   `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens    int `json:"prompt_tokens"`
		CompletionTokes int `json:"completion_tokens"`
		TotalTokens     int `json:"total_tokens"`
	}
}

// OpenAiRcvWithIdentity 角色扮演使用的Response
type OpenAiRcvWithIdentity struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		Logprobs     int    `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens    int `json:"prompt_tokens"`
		CompletionTokes int `json:"completion_tokens"`
		TotalTokens     int `json:"total_tokens"`
	}
}

// Client 返回代理客户端
func Client() (http.Client, error) {
	if config.Cfg.OpenAi.UseProxy == false {
		return http.Client{}, nil
	}
	// 设置clash代理
	uri, err := url.Parse(config.Cfg.OpenAi.ProxyUrl)
	if err != nil {
		log.Fatal(err)
		return http.Client{}, nil
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}
	return client, nil
}

// ChooseGenerateWay 选择生成方式
func ChooseGenerateWay(session string, text string, useContext bool) (string, error) {
	log.Println("正在调用OpenAI API生成文本...", text)
	return GenerateText(session, text, useContext)
}

// GenerateText 调用openai的API生成文本
func GenerateText(session string, text string, useContext bool) (string, error) {
	var ms []Messages
	//检测是否有角色信息
	if config.Cfg.RolePlay.Role != nil {
		//遍历角色信息
		for index, role := range config.Cfg.RolePlay.Role {
			var message Messages
			message.Role = role
			message.Content = config.Cfg.RolePlay.Content[index]

			ms = append(ms, message)
		}
	}

	// 读取上下文
	if useContext {
		ms = append(ms, Cache.GetMsg(session)...)
	}

	// 填入本次对话内容
	message := &Messages{
		Role:    "user",
		Content: text,
	}
	ms = append(ms, *message)

	// 构造请求体
	postDataTemp := postData{
		Model:       config.Cfg.OpenAi.Model,
		Messages:    ms,
		MaxTokens:   config.Cfg.OpenAi.MaxTokens,
		Temperature: config.Cfg.OpenAi.Temperature,
	}

	//请求体序列化
	postDataBytes, err := json.Marshal(postDataTemp)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// 创建请求
	req, _ := http.NewRequest("POST", config.Cfg.OpenAi.Url, bytes.NewBuffer(postDataBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Cfg.OpenAi.ApiKey)
	req.Header.Set("api-key", config.Cfg.OpenAi.ApiKey) //兼容azure openai api
	client, err := Client()
	if err != nil {
		log.Println(err)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err.Error(), err
	}

	defer resp.Body.Close()
	if resp == nil {
		log.Println("response is nil")
		return "", err
	}
	body, _ := io.ReadAll(resp.Body)
	var openAiRcv OpenAiRcv
	err = json.Unmarshal(body, &openAiRcv)
	if err != nil {
		log.Println(err)
		return err.Error(), err
	}
	if len(openAiRcv.Choices) == 0 {
		log.Println("OpenAI API调用失败，返回内容：", string(body))
		return string(body), err
	}

	// 保存上下文
	if useContext {
		ms = append(ms, openAiRcv.Choices[0].Message)
		Cache.SetMsg(session, ms)
	}
	openAiRcv.Choices[0].Message.Content = strings.Replace(openAiRcv.Choices[0].Message.Content, "\n\n", "\n", 1)
	log.Printf("Model: %s TotalTokens: %d+%d=%d", openAiRcv.Model, openAiRcv.Usage.PromptTokens, openAiRcv.Usage.CompletionTokes, openAiRcv.Usage.TotalTokens)
	return openAiRcv.Choices[0].Message.Content, err
}
