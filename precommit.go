package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type PreCommitRequest struct {
	Diff string `json:"diff"`
}

type PreCommitResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// handlePreCommit 处理 pre-commit 检查请求.
func handlePreCommit(c *gin.Context) {
	var req PreCommitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, PreCommitResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}
	review, err := reviewCode(req.Diff)
	if err != nil {
		c.JSON(http.StatusOK, PreCommitResponse{
			Success: false,
			Message: fmt.Sprintf("Error reviewing code: %v", err),
		})
	}
	fmt.Println(review)

	if strings.Contains(review, "FIXIT!") {
		c.JSON(http.StatusOK, PreCommitResponse{
			Success: false,
			Message: "Code review issues,message: " + review,
		})
		return
	}

	c.JSON(http.StatusOK, PreCommitResponse{
		Success: true,
		Message: "All files passed code review,message: " + review,
	})
}

// reviewCode 使用 DeepSeek 进行代码审查.
func reviewCode(content string) (string, error) {
	prompt := fmt.Sprintf(config.Prompt.Precommit, content)
	response, err := callDeepSeek(prompt)
	if err != nil {
		return "", err
	}
	return response, nil
}
