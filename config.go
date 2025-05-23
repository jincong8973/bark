package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	GitLab struct {
		Token string `mapstructure:"token"`
		URL   string `mapstructure:"url"`
	} `mapstructure:"gitlab"`

	DeepSeek struct {
		Token string `mapstructure:"token"`
		URL   string `mapstructure:"url"`
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

	// 绑定环境变量
	viper.BindEnv("gitlab.token", "BARK_GITLAB_TOKEN")
	viper.BindEnv("gitlab.url", "BARK_GITLAB_URL")
	viper.BindEnv("deepseek.token", "BARK_DEEPSEEK_TOKEN")
	viper.BindEnv("deepseek.url", "BARK_DEEPSEEK_URL")
	viper.BindEnv("server.port", "BARK_SERVER_PORT")

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
