package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/maczh/logs"
	"github.com/maczh/mgconfig"
	"github.com/maczh/mgi18n"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const config_file = "test.yml"

//@title	OSS集中管理服务模块
//@version 	0.0.1(oss-server)
//@description	OSS集中管理服务模块

func main() {
	//初始化配置，自动连接数据库和Nacos服务注册
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	mgconfig.InitConfig(path + "/" + config_file)

	//GIN的模式，生产环境可以设置成release
	gin.SetMode("debug")

	//初始化国际化错误代码
	mgi18n.Init()

	engine := setupRouter()

	server := &http.Server{
		Addr:    ":" + mgconfig.GetConfigString("go.application.port"),
		Handler: engine,
	}
	serverSsl := &http.Server{
		Addr:    ":" + mgconfig.GetConfigString("go.application.port_ssl"),
		Handler: engine,
	}

	logs.Info("|-----------------------------------|")
	logs.Info("|       OSS集中管理服务模块  0.0.1     |")
	logs.Info("|-----------------------------------|")
	logs.Info("|  Go Http Server Start Successful  |")
	logs.Info("|    Port:" + mgconfig.GetConfigString("go.application.port") + "     Pid:" + fmt.Sprintf("%d", os.Getpid()) + "        |")
	logs.Info("|-----------------------------------|")

	logs.Debug("====================================")
	logs.Debug("| {}启动成功!   侦听端口:{}     |", mgconfig.GetConfigString("go.application.name"), mgconfig.GetConfigString("go.application.port"))
	logs.Debug("====================================")

	//http端口侦听
	if mgconfig.GetConfigString("go.application.port") != "" {
		go func() {
			var err error
			err = server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				logs.Error("HTTP server listen: {}", err.Error())
			}
		}()
	}
	//https端口侦听
	if mgconfig.GetConfigString("go.application.cert") != "" {
		go func() {
			var err error
			err = serverSsl.ListenAndServeTLS(path+"/"+mgconfig.GetConfigString("go.application.cert"), path+"/"+mgconfig.GetConfigString("go.application.key"))
			if err != nil && err != http.ErrServerClosed {
				logs.Error("HTTPS server listen: {}", err.Error())
			}
		}()
	}

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalChan
	logs.Error("Get Signal:" + sig.String())
	logs.Error("Shutdown Server ...")
	mgconfig.SafeExit()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logs.Error("Server Shutdown:" + err.Error())
	}
	logs.Error("Server exiting")

}
