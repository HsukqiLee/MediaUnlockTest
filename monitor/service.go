package main

import (
	"log"
	
	"github.com/kardianos/service"
)

func serviceExists(s service.Service) bool {
	status, err := s.Status()
	return err == nil && status != service.StatusUnknown
}

func serviceIsActive(s service.Service) bool {
	status, err := s.Status()
	return err == nil && status == service.StatusRunning
}

func restartService(s service.Service) error {
	return s.Restart()
}

func installService(s service.Service) {
    if serviceExists(s) {
	    uninstallService(s)
	}
	if err := service.Control(s, "install"); err != nil {
		log.Fatal("[ERR] 设置unlock-monitor服务时出错:", err)
	}
	log.Println("[OK] 设置unlock-monitor服务成功")
	if err := service.Control(s, "start"); err != nil {
		log.Fatal("[ERR] 启动unlock-monitor服务失败", err)
	} else {
		log.Println("[OK] 启动unlock-monitor服务成功")
	}
	return
}

func uninstallService(s service.Service) {
    if err := service.Control(s, "stop"); err != nil {
		log.Println("[OK] 停止unlock-monitor服务失败:", err)
	}
	if err := service.Control(s, "uninstall"); err != nil {
		log.Fatal("[ERR] 卸载unlock-monitor服务失败", err)
	} else {
		log.Println("[OK] 卸载unlock-monitor服务成功")
	}
	return
}

func startService(s service.Service) {
    if !serviceExists(s) {
        log.Println("[ERR] unlock-monitor服务不存在")
    } else if err := service.Control(s, "start"); err != nil {
		log.Fatal("[ERR] 启动unlock-monitor服务失败", err)
	} else {
		log.Println("[OK] 启动unlock-monitor服务成功")
	}
	return
}

func stopService(s service.Service) {
    if !serviceExists(s) {
        log.Println("[ERR] unlock-monitor服务不存在")
    } else if err := service.Control(s, "stop"); err != nil {
		log.Println("[OK] 停止unlock-monitor服务失败:", err)
	} else {
	    log.Println("[OK] 停止unlock-monitor服务成功:", err)
	}
	return
}

