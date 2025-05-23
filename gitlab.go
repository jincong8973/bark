package main

import "github.com/xanzy/go-gitlab"

var client *gitlab.Client

func GetGitlabClient() *gitlab.Client {
	if client == nil {
		config := GetConfig()
		newClient, err := gitlab.NewClient(config.GitLab.Token, gitlab.WithBaseURL(config.GitLab.URL))
		if err != nil {
			panic("GitLab client initialization failed")
		}
		client = newClient
	}
	return client
}
