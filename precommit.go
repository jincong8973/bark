package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type PreCommitRequest struct {
	Files []string `json:"files"`
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

	// 检查文件列表
	if len(req.Files) == 0 {
		c.JSON(http.StatusOK, PreCommitResponse{
			Success: false,
			Message: "No files provided",
		})
		return
	}

	// 对每个文件进行代码审查
	for _, file := range req.Files {
		// 只检查特定类型的文件
		if !shouldCheckFile(file) {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			c.JSON(http.StatusOK, PreCommitResponse{
				Success: false,
				Message: fmt.Sprintf("Error reading file %s: %v", file, err),
			})
			return
		}

		// 调用 DeepSeek 进行代码审查
		review, err := reviewCode(string(content))
		if err != nil {
			c.JSON(http.StatusOK, PreCommitResponse{
				Success: false,
				Message: fmt.Sprintf("Error reviewing file %s: %v", file, err),
			})
			return
		}

		// 如果审查发现问题，返回错误
		if review != "" {
			c.JSON(http.StatusOK, PreCommitResponse{
				Success: false,
				Message: fmt.Sprintf("Code review issues in %s:\n%s", file, review),
			})
			return
		}
	}

	c.JSON(http.StatusOK, PreCommitResponse{
		Success: true,
		Message: "All files passed code review",
	})
}

// shouldCheckFile 判断是否需要检查该文件.
func shouldCheckFile(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	// 只检查代码文件
	return ext == ".go" || ext == ".py" || ext == ".js" || ext == ".ts" || ext == ".java"
}

// reviewCode 使用 DeepSeek 进行代码审查.
func reviewCode(content string) (string, error) {
	prompt := fmt.Sprintf("请对以下代码进行审查，指出潜在的问题和改进建议：\n\n%s", content)
	response, err := callDeepSeek(prompt)
	if err != nil {
		return "", err
	}
	return response, nil
}
