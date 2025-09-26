package main

import (
	"fmt"
	"os"
	"path/filepath"

	dem "github.com/zhangymPerson/dev-env-manage/src"
	"github.com/zhangymPerson/dev-env-manage/src/constant"
	"github.com/zhangymPerson/dev-env-manage/src/db"
	"github.com/zhangymPerson/dev-env-manage/src/log"
)

// 全局变量用于存储构建信息
// 这些变量将在构建时通过 ldflags 设置
var (
	GitCommit = "main" // 默认值
	GitBranch = "main" // 默认值

)

func main() {
	// 最先配置日志
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	logFile := filepath.Join(homeDir, ".dem", "dem.log")
	if err := log.Configure(false, logFile); err != nil {
		panic(fmt.Sprintf("Failed to configure logger: %v", err))
	}

	// 确保在程序结束时关闭日志文件
	defer func() {
		if err := log.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close log file: %v\\n", err)
		}
	}()

	// 然后执行其他初始化
	dem.Init()
	db.InitDB(constant.GetDBFilePath())
	dem.Options(GitBranch, GitCommit)
}
