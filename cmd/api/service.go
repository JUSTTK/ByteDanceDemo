// Package api @Author: youngalone [2023/8/7]
package api

import (
	config2 "bytedancedemo/config"
	"bytedancedemo/dao"
	"bytedancedemo/database/mysql"
	redis2 "bytedancedemo/database/redis"
	"bytedancedemo/service"
	redis3 "bytedancedemo/utils"
	"bytedancedemo/utils/sensetive"

	"bytedancedemo/middleware/rabbitmq"
	"bytedancedemo/middleware/redis"
	"bytedancedemo/router"

	"bytedancedemo/utils/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

var (
	config   string
	mode     string
	StartCmd = &cobra.Command{
		Use:   "server",
		Short: "服务入口",
		Long:  "抖音极简版APP服务器",
		PreRun: func(cmd *cobra.Command, args []string) {
			config2.Init(config)
			log.InitLogger(mode) //日志重复
			mysql.Init()
			redis.InitRedis()
			redis2.Init()
			redis3.Init()
			rabbitmq.InitRabbitMQ()
			rabbitmq.InitCommentRabbitMQ()
			rabbitmq.InitFollowRabbitMQ()
			dao.SetDefault(mysql.DB)
			service.StartFavoriteWorkerPool() // 初始化点赞worker池
			service.InitUserWorkerPool(10, 100) // 初始化用户查询worker池（最小10，最大100）
			sensetive.InitFilter()           // 初始化敏感词过滤器
		},
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/settings.yml", "配置文件路径")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "debug", "运行模式")
}

func run() {
	go router.Setup()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	zap.L().Info("监听中断中...")
	<-quit
	zap.L().Sync()
	zap.L().Info("关闭服务器...")
}
