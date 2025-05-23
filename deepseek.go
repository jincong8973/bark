package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekRequest struct {
	Model    string            `json:"model"`
	Messages []DeepSeekMessage `json:"messages"`
	Stream   bool              `json:"stream"`
}

type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func callDeepSeek(diff string) (string, error) {
	config := GetConfig()
	if config.DeepSeek.Token == "" {
		return "", errors.New("DeepSeek token not set")
	}

	reqBody := DeepSeekRequest{
		Model: "deepseek-reasoner",
		Messages: []DeepSeekMessage{
			{Role: "system", Content: "你是一个专业的代码审查助手，请根据 diff 内容给出详细的代码 review 建议, 你再给出review建议的同时请给出相应的代码文件名称和行号。"},
			{Role: "user", Content: "我们的服务主要使用Golang来编写,也有K8S Yaml,因为我们在实现K8S相关的配套服务.diff 内容如下：\n" + diff},
		},
		Stream: false,
	}
	data, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", config.DeepSeek.URL, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.DeepSeek.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var dsResp DeepSeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&dsResp); err != nil {
		return "", err
	}
	if len(dsResp.Choices) == 0 {
		return "", errors.New("no response from DeepSeek")
	}
	return dsResp.Choices[0].Message.Content, nil
}
