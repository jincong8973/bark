package api

import (
	"fmt"
	"io"
	"net/http"

	"bark/config"
	"bark/llm/deepseek"
	"bark/thirdparty"

	"github.com/gin-gonic/gin"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func HandleWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	event, err := gitlab.ParseWebhook(gitlab.HookEventType(c.Request), body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook"})
		return
	}
	mrEvent, ok := event.(*gitlab.MergeEvent)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"msg": "not a merge request event"})
		return
	}

	// 只在 MR 创建和更新时触发 review
	state := mrEvent.ObjectAttributes.State
	// doc : https://docs.gitlab.com/user/project/integrations/webhook_events/#merge-request-events
	action := mrEvent.ObjectAttributes.Action

	// 必须是打开状态的 MR 才进行审查.
	if state != "opened" {
		c.JSON(http.StatusOK, gin.H{"msg": "ignored state: " + state})
		return
	}

	if action != "update" && action != "open" {
		c.JSON(http.StatusOK, gin.H{"msg": "ignored action: merged"})
		return
	}

	if action == "update" && mrEvent.ObjectAttributes.OldRev == "" {
		c.JSON(http.StatusOK, gin.H{"msg": "update event without old revision, no code changes, ignored"})
		return
	}

	projectID := mrEvent.Project.ID
	mrIID := mrEvent.ObjectAttributes.IID

	// 获取 MR diff
	changes, _, err := thirdparty.GetGitlabClient().MergeRequests.ListMergeRequestDiffs(projectID, mrIID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get diff failed"})
		return
	}

	// 拼接 diff 内容
	var diffText string
	for _, change := range changes {
		status := "修改"
		if change.NewFile {
			status = "新建"
		} else if change.DeletedFile {
			status = "删除"
		} else if change.RenamedFile {
			status = "重命名"
		}
		diffText += fmt.Sprintf("文件: %s\n状态: %s\n差异: %s\n", change.NewPath, status, change.Diff)
	}

	cfg := config.GetConfig()
	review, err := deepseek.CallDeepSeek(fmt.Sprintf(cfg.Prompt.MergeRequest, diffText))
	if err != nil {
		fmt.Println("call deepseek failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "call deepseek failed"})
		return
	}

	fmt.Println("Review length:", len(review))

	// 评论 MR
	noteOpt := &gitlab.CreateMergeRequestNoteOptions{
		Body: &review,
	}
	_, _, err = thirdparty.GetGitlabClient().Notes.CreateMergeRequestNote(projectID, mrIID, noteOpt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "comment failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "reviewed"})
}
