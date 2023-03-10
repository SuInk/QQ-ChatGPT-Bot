package chatgpt

import (
	"QQ-ChatGPT-Bot/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const OpenaiApiUrl = "https://api.openai.com/v1/chat/completions"

type OpenAiRcv struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens    int `json:"prompt_tokens"`
		CompletionTokes int `json:"completion_tokens"`
		TotalTokens     int `json:"total_tokens"`
	}
}

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

// GenerateText 调用openai的API生成文本
func GenerateText(text string) string {
	log.Println("正在调用OpenAI API生成文本...", text)
	postData := []byte(fmt.Sprintf(`{
	  "model": "%s",
	  "messages": %s,
	  "max_tokens": %d,
	  "temperature": %.1f
	}`, config.Cfg.OpenAi.Model, "[{\"role\": \"user\", \"content\": \""+text+"\"}]", config.Cfg.OpenAi.MaxTokens, config.Cfg.OpenAi.Temperature))
	req, _ := http.NewRequest("POST", OpenaiApiUrl, bytes.NewBuffer(postData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Cfg.OpenAi.ApiKey)
	client, err := Client()
	if err != nil {
		log.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var openAiRcv OpenAiRcv
	err = json.Unmarshal(body, &openAiRcv)
	if err != nil {
		log.Println(err)
	}
	if len(openAiRcv.Choices) == 0 {
		log.Println("OpenAI API调用失败，返回内容：", string(body))
		return string(body)
	}
	openAiRcv.Choices[0].Message.Content = strings.Replace(openAiRcv.Choices[0].Message.Content, "\n\n", "\n", 1)
	log.Printf("Model: %s TotalTokens: %d+%d=%d", openAiRcv.Model, openAiRcv.Usage.PromptTokens, openAiRcv.Usage.CompletionTokes, openAiRcv.Usage.TotalTokens)
	return openAiRcv.Choices[0].Message.Content
}
