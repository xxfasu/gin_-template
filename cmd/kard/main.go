package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
)

// 定义程序结构体
type program struct {
	server *gin.Engine
	svc    service.Service
}

// Start 方法在服务启动时调用
func (p *program) Start(s service.Service) error {
	// 在一个新的 goroutine 中启动服务
	go p.run()
	return nil
}

// Stop 方法在服务停止时调用
func (p *program) Stop(s service.Service) error {
	// 这里可以添加清理资源的代码
	return nil
}

// 运行 Gin 服务器
func (p *program) run() {
	// 配置 Gin
	p.server = gin.Default()

	// 注册路由
	p.server.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// 启动服务器
	if err := p.server.Run(":8080"); err != nil {
		log.Fatalf("Gin server failed to run: %v", err)
	}
}

func main() {
	// 定义服务配置
	svcConfig := &service.Config{
		Name:        "MyGinService",
		DisplayName: "My Gin Web Service",
		Description: "This is a Gin web application managed by kardianos/service.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	prg.svc = s

	// 设置日志
	errs := make(chan error)
	logger, err := s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	// 启动服务
	if err := s.Run(); err != nil {
		logger.Error(err)
	}

	// 处理错误
	go func() {
		for {
			select {
			case err := <-errs:
				log.Println("Error:", err)
			}
		}
	}()

	// 等待中断信号以优雅关闭
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutting down...")
}
