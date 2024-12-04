// main.go

package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ServerConfig 配置结构体
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// Database 配置结构体
type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

// Features 配置结构体
type FeaturesConfig struct {
	EnableLogging bool `mapstructure:"enable_logging"`
	EnableCache   bool `mapstructure:"enable_cache"`
}

// Config 总配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Features FeaturesConfig `mapstructure:"features"`
}

func main() {
	var config Config

	// 设置默认值
	viper.SetDefault("server.port", 8000)
	viper.SetDefault("server.host", "127.0.0.1")

	// 设置配置文件的名称（不带扩展名）
	viper.SetConfigName("config")
	// 设置配置文件的类型
	viper.SetConfigType("toml")
	// 添加配置文件所在的路径
	viper.AddConfigPath(".") // 当前目录

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 将配置文件内容反序列化到 Config 结构体
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 打印解析后的配置内容
	fmt.Printf("服务器配置: %+v\n", config.Server)
	fmt.Printf("数据库配置: %+v\n", config.Database)
	fmt.Printf("功能配置: %+v\n", config.Features)

	// 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已修改:", e.Name)
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("重新解析配置文件失败: %v", err)
		}
		fmt.Printf("更新后的服务器配置: %+v\n", config.Server)
		fmt.Printf("更新后的数据库配置: %+v\n", config.Database)
		fmt.Printf("更新后的功能配置: %+v\n", config.Features)
	})

	// 阻塞主 goroutine
	select {}
}
