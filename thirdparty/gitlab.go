package thirdparty

import (
	config2 "bark/config"
	"gitlab.com/gitlab-org/api/client-go"
)

var client *gitlab.Client

func GetGitlabClient() *gitlab.Client {
	if client == nil {
		config := config2.GetConfig()
		newClient, err := gitlab.NewClient(config.GitLab.Token, gitlab.WithBaseURL(config.GitLab.URL))
		if err != nil {
			panic("GitLab client initialization failed")
		}
		client = newClient
	}
	return client
}
