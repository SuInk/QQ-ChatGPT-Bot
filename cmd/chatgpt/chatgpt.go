package chatgpt

import (
	"QQ-ChatGPT-Bot/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const OpenaiApiUrl = "https://api.openai.com/v1/completions"

type OpenAiRcv struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens    int `json:"prompt_tokens"`
		CompletionTokes int `json:"completion_tokens"`
		TotalTokens     int `json:"total_tokens"`
	}
}

// GenerateText 调用openai的API生成文本
func GenerateText(text string) string {
	log.Println("正在调用OpenAI API生成文本...", text)
	postData := []byte(fmt.Sprintf(`{
	  "model": "%s",
	  "prompt": "%s",
	  "max_tokens": %d,
	  "temperature": %.1f
	}`, config.Cfg.OpenAi.Model, text, config.Cfg.OpenAi.MaxTokens, config.Cfg.OpenAi.Temperature))
	req, _ := http.NewRequest("POST", OpenaiApiUrl, bytes.NewBuffer(postData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Cfg.OpenAi.ApiKey)
	resp, err := http.DefaultClient.Do(req)
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
	openAiRcv.Choices[0].Text = strings.ReplaceAll(openAiRcv.Choices[0].Text, "\n\n", "")
	log.Printf("Model: %s TotalTokens: %d+%d=%d", openAiRcv.Model, openAiRcv.Usage.PromptTokens, openAiRcv.Usage.CompletionTokes, openAiRcv.Usage.TotalTokens)
	return openAiRcv.Choices[0].Text
}
