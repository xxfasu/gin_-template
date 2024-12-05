// main.go

package conf

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// System 配置结构体
type System struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// Mysql 配置结构体
type Mysql struct {
	Source string `mapstructure:"source"`
}

// Redis 配置结构体
type Redis struct {
	Addr     string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// OSS 配置结构体
type OSS struct {
	Provider string `mapstructure:"provider"`
}

// AliyunOSS 配置结构体
type AliyunOSS struct {
	Endpoint  string `mapstructure:"endpoint"`
	KeyID     string `mapstructure:"key_id"`
	KeySecret string `mapstructure:"key_secret"`
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	RoleArn   string `mapstructure:"role_arn"`
	Domain    string `mapstructure:"domain"`
}

// Config 总配置结构体
type config struct {
	System    System    `mapstructure:"system"`
	Mysql     Mysql     `mapstructure:"mysql"`
	Redis     Redis     `mapstructure:"redis"`
	OSS       OSS       `mapstructure:"oss"`
	AliyunOSS AliyunOSS `mapstructure:"aliyun_oss"`
}

var Config *config
var env *Env

func InitConfig() error {
	err := loadEnv()
	if err != nil {
		return err
	}
	switch env.Environment {
	case "local":
		err = loadLocal()
	case "prod":
		err = loadNacos()
	}
	if err != nil {
		return err
	}
	return nil
}

func loadEnv() error {
	// 设置默认值
	viper.SetDefault("environment", "local")

	// 设置配置文件的名称（不带扩展名）
	viper.SetConfigName("env")
	// 设置配置文件的类型
	viper.SetConfigType("toml")
	// 添加配置文件所在的路径
	viper.AddConfigPath("./config") // 当前目录

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
		return err
	}

	// 将配置文件内容反序列化到 Config 结构体
	if err := viper.Unmarshal(&env); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
		return err
	}

	// 打印解析后的配置内容
	fmt.Printf("Env配置: %+v\n", env.Environment)
	fmt.Printf("Nacos配置: %+v\n", env.Nacos)
	return nil
}
