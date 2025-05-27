package main

import (
	"fmt"
	"os"

	"bark/api"
	"bark/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bark",
	Short: "Bark - GitLab MR Review Bot",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		if err := config.InitConfig(configPath); err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		r := gin.Default()
		r.POST("/webhook", api.HandleWebhook)
		r.POST("/precommit", api.HandlePreCommit)
		addr := fmt.Sprintf("0.0.0.0:%d", config.GetConfig().Server.Port)
		if err := r.Run(addr); err != nil {
			fmt.Printf("服务启动失败: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "config.yaml", "配置文件路径")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
