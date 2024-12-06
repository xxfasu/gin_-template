// main.go

package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// system 配置结构体
type system struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// mysql 配置结构体
type mysql struct {
	Source string `mapstructure:"source"`
}

// redis 配置结构体
type redis struct {
	Addr     string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// oss 配置结构体
type oss struct {
	Provider string `mapstructure:"provider"`
}

// aliyunOSS 配置结构体
type aliyunOSS struct {
	Endpoint  string `mapstructure:"endpoint"`
	KeyID     string `mapstructure:"key_id"`
	KeySecret string `mapstructure:"key_secret"`
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	RoleArn   string `mapstructure:"role_arn"`
	Domain    string `mapstructure:"domain"`
}

// qiniuyunOSS 配置结构体
type qiniuyunOSS struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Zone      string `mapstructure:"zone"`
}

type zapLog struct {
	LogLevel    string `mapstructure:"log_level"`
	Encoding    string `mapstructure:"encoding"`
	LogFileName string `mapstructure:"log_file_name"`
	MaxBackups  int    `mapstructure:"max_backups"`
	MaxAge      int    `mapstructure:"max_age"`
	MaxSize     int    `mapstructure:"max_size"`
	Compress    bool   `mapstructure:"compress"`
}

// Config 总配置结构体
type config struct {
	System      system      `mapstructure:"system"`
	Mysql       mysql       `mapstructure:"mysql"`
	Redis       redis       `mapstructure:"redis"`
	OSS         oss         `mapstructure:"oss"`
	AliyunOSS   aliyunOSS   `mapstructure:"aliyun_oss"`
	Log         zapLog      `mapstructure:"zap_log"`
	QiniuyunOSS qiniuyunOSS `mapstructure:"qiniuyun_oss"`
}

var Config *config
var Env *env

func InitConfig() error {
	err := loadEnv()
	if err != nil {
		return err
	}
	switch Env.Environment {
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
	if err := viper.Unmarshal(&Env); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
		return err
	}

	// 打印解析后的配置内容
	fmt.Printf("Env配置: %+v\n", Env.Environment)
	fmt.Printf("Nacos配置: %+v\n", Env.Nacos)
	return nil
}
