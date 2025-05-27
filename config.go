package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	GitLab struct {
		Token string `mapstructure:"token"`
		URL   string `mapstructure:"url"`
	} `mapstructure:"gitlab"`
	Prompt struct {
		MergeRequest string `mapstructure:"merge_request"`
		Precommit    string `mapstructure:"precommit"`
	} `mapstructure:"prompt"`
	DeepSeek struct {
		Token    string `mapstructure:"token"`
		URL      string `mapstructure:"url"`
		Model    string `mapstructure:"model"`
		Messages struct {
			System string `mapstructure:"system"`
			User   string `mapstructure:"user"`
		} `mapstructure:"messages"`
	} `mapstructure:"deepseek"`

	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

var config Config

func initConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	// 设置环境变量前缀
	viper.SetEnvPrefix("BARK")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return viper.Unmarshal(&config)
}

func GetConfig() *Config {
	return &config
}
