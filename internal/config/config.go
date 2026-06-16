package config

import (
	"os"
	"path/filepath"
)

var (
	// DataDir 用户指定的数据目录
	DataDir string
)

// Load 加载配置
// 优先级：命令行参数 > 环境变量 > 默认值 (data/)
func Load() string {
	// 1. 优先使用命令行参数
	if DataDir != "" {
		return DataDir
	}

	// 2. 环境变量
	envDataDir := os.Getenv("CHRONICLE_DATA_DIR")
	if envDataDir != "" {
		return envDataDir
	}

	// 3. 默认值，保持向后兼容
	return "data"
}

// GetDBPath 获取数据库文件路径
func GetDBPath() string {
	dir := Load()
	return filepath.Join(dir, "app.db")
}

// EnsureDataDir 确保数据目录存在
func EnsureDataDir() error {
	dir := Load()
	return os.MkdirAll(dir, 0755)
}
