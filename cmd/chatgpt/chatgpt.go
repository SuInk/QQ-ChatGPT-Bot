package chatgpt

import (
	"QQ-ChatGPT-Bot/config"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const Openaiapiurl1 = "https://api.openai.com/v1/chat/completions" //对话使用的url
const Openaiapiurl2 = "https://api.openai.com/v1/completions"      //角色扮演使用的url

// 对话使用的Request body
type postData struct {
	Model       string        `json:"model"`
	Messages    []interface{} `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

// 角色扮演使用的Request body
type postDataWithIdentity struct {
	Model       string        `json:"model"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
	Prompt      []interface{} `json:"prompt"`
	Stop        []string      `json:"stop"`
}

// OpenAiRcv 对话使用的Response
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
func ChooseGenerateWay(text string) string {
	log.Println("正在调用OpenAI API生成文本...", text)
	if config.Cfg.Identity.Prompt == "" {
		return GenerateText(text)
	} else {
		return GenerateTextWithIdentity(text)
	}
}

// GenerateText 调用openai的API生成文本
func GenerateText(text string) string {
	message := struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		Role:    "user",
		Content: text,
	}
	postDataTemp := postData{
		Model:       config.Cfg.OpenAi.Model,
		Messages:    []interface{}{message},
		MaxTokens:   config.Cfg.OpenAi.MaxTokens,
		Temperature: float64(config.Cfg.OpenAi.Temperature),
	}
	postDataBytes, err := json.Marshal(postDataTemp)
	if err != nil {
		log.Println(err)
		return ""
	}
	req, _ := http.NewRequest("POST", Openaiapiurl1, bytes.NewBuffer(postDataBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Cfg.OpenAi.ApiKey)
	client, err := Client()
	if err != nil {
		log.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	if resp == nil {
		log.Println("response is nil")
		return ""
	}
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

// GenerateTextWithIdentity 使用身份的时候，使用这个生成文本
func GenerateTextWithIdentity(text string) string {
	postDataTemp := postDataWithIdentity{
		Model:       config.Cfg.OpenAi.Model,
		MaxTokens:   config.Cfg.OpenAi.MaxTokens,
		Temperature: float64(config.Cfg.OpenAi.Temperature),
		Prompt:      []interface{}{config.Cfg.Identity.Prompt + "\n" + config.Cfg.Identity.Stop[0] + ":" + text + "\n" + config.Cfg.Identity.Stop[1] + ":"},
		Stop:        config.Cfg.Identity.Stop,
	}
	postDataBytes, err := json.Marshal(postDataTemp)
	if err != nil {
		log.Println(err)
		return ""
	}
	req, _ := http.NewRequest("POST", Openaiapiurl2, bytes.NewBuffer(postDataBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Cfg.OpenAi.ApiKey)
	client, err := Client()
	if err != nil {
		log.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	if resp == nil {
		log.Println("response is nil")
		return ""
	}
	body, _ := io.ReadAll(resp.Body)
	var openAiRcvWithIdentity OpenAiRcvWithIdentity
	err = json.Unmarshal(body, &openAiRcvWithIdentity)
	if err != nil {
		log.Println(err)
	}
	if len(openAiRcvWithIdentity.Choices) == 0 {
		log.Println("OpenAI API调用失败，返回内容：", string(body))
		return string(body)
	}
	openAiRcvWithIdentity.Choices[0].Text = strings.Replace(openAiRcvWithIdentity.Choices[0].Text, "\n\n", "\n", 1)
	log.Printf("Model: %s TotalTokens: %d+%d=%d", openAiRcvWithIdentity.Model, openAiRcvWithIdentity.Usage.PromptTokens, openAiRcvWithIdentity.Usage.CompletionTokes, openAiRcvWithIdentity.Usage.TotalTokens)
	return openAiRcvWithIdentity.Choices[0].Text
}
