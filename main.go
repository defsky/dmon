package main

import (
	//"log"
	//"net/http"
	"fmt"
	"os"
	"pmon/app"

	"github.com/jander/golog/logger"
	"github.com/kardianos/service"
)

/**
* MAIN函数，程序入口
 */

func main() {
	svcConfig := &service.Config{
		Name:        "pmonsvc",         //服务名称
		DisplayName: "Process Monitor", //服务显示名称
		Description: "此服务用作监控指定进程是否存活", //服务描述
	}

	prg := &app.Program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Fatal(err)
	}

	if err != nil {
		logger.Fatal(err)
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			s.Install()
			logger.Println("服务安装成功")
			return
		}

		if os.Args[1] == "remove" {
			s.Uninstall()
			logger.Println("服务卸载成功")
			return
		}

		if os.Args[1] == "-h" {
			fmt.Println(os.Args[0], "[install | remove | -h]")
			fmt.Println("\tinstall\tInstall program as a service")
			fmt.Println("\tremove\tRemove service of this program")
			fmt.Println("\t-h\tDisplay this help message")
			return
		}
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
