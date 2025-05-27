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

func callDeepSeek(prompt string) (string, error) {
	config := GetConfig()
	if config.DeepSeek.Token == "" {
		return "", errors.New("DeepSeek token not set")
	}

	reqBody := DeepSeekRequest{
		Model: config.DeepSeek.Model,
		Messages: []DeepSeekMessage{
			{Role: "system", Content: config.DeepSeek.Messages.System},
			{Role: "user", Content: prompt},
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
