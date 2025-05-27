package deepseek

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	config2 "bark/config"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func removeThinkProcess(content string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	return strings.TrimSpace(re.ReplaceAllString(content, ""))
}

func CallDeepSeek(prompt string) (string, error) {
	config := config2.GetConfig()
	if config.DeepSeek.Token == "" {
		return "", errors.New("DeepSeek token not set")
	}

	reqBody := Request{
		Model: config.DeepSeek.Model,
		Messages: []Message{
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
	var dsResp Response
	if err := json.NewDecoder(resp.Body).Decode(&dsResp); err != nil {
		return "", err
	}
	if len(dsResp.Choices) == 0 {
		return "", errors.New("no response from DeepSeek")
	}
	return removeThinkProcess(dsResp.Choices[0].Message.Content), nil
}
