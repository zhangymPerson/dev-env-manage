package constant

import (
	"os"
	"path/filepath"
)

// 基础类型常量
const (
	// 整型常量
	MaxInt   int   = 2147483647
	MinInt   int   = -2147483648
	MaxInt64 int64 = 9223372036854775807
	MinInt64 int64 = -9223372036854775808

	// 浮点型常量
	Pi float64 = 3.141592653589793
	E  float64 = 2.718281828459045

	// 布尔型常量
	True  bool = true
	False bool = false

	// 枚举常量
	StatusOK  int = 200
	StatusErr int = 500

	// 数据文件位置
	DBFileName = "dem_config.db"
)

// 环境类型枚举
const (
	EnvDefault EnvType = iota // 默认环境（未指定时使用）
	EnvDev                    // 开发环境
	EnvTest                   // 测试环境
	EnvPre                    // 预发布环境
	EnvProd                   // 生产环境
	EnvOther                  // 其他环境
)

// EnvType 定义环境类型的底层类型
type EnvType int

// 字符串映射（用于日志/调试）
var envTypeToString = map[EnvType]string{
	EnvDefault: "default",
	EnvDev:     "dev",
	EnvTest:    "test",
	EnvPre:     "pre",
	EnvProd:    "prod",
	EnvOther:   "other",
}

// String 实现Stringer接口，便于打印
func (e EnvType) String() string {
	if str, ok := envTypeToString[e]; ok {
		return str
	}
	return "unknown"
}

// 引用类型模拟（通过函数返回）
var (
	// 模拟切片常量
	DefaultConfig = []string{"config1", "config2", "config3"}

	// 模拟映射常量
	DefaultEnv = map[string]string{
		"GOPATH": "/home/user/go",
		"SHELL":  "/bin/bash",
	}

	// 模拟结构体常量
	DefaultUser = User{
		Name: "admin",
		Role: "superuser",
	}
)

// 结构体定义
type User struct {
	Name string
	Role string
}

// 获取引用类型常量的函数
func GetDefaultConfig() []string {
	return append([]string{}, DefaultConfig...) // 返回副本
}

func GetDefaultEnv() map[string]string {
	env := make(map[string]string)
	for k, v := range DefaultEnv {
		env[k] = v
	}
	return env // 返回副本
}

func GetDefaultUser() User {
	return DefaultUser // 结构体是值类型，可直接返回
}

func GetDBFilePath() string {
	return filepath.Join(GetProjectDir(), DBFileName)
}

func GetProjectDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	demDir := filepath.Join(homeDir, ".dem")
	if _, err := os.Stat(demDir); os.IsNotExist(err) {
		err := os.Mkdir(demDir, 0755)
		if err != nil {
			panic(err)
		}
	}
	return demDir
}
