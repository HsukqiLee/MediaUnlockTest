package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	
	selfUpdate "github.com/inconshreveable/go-update"
	pb "github.com/schollz/progressbar/v3"
)

var (
	Version   = "1.4.2"
	buildTime string
)

func serviceExists(serviceName string) bool {
	cmd := exec.Command("systemctl", "is-enabled", serviceName)
	err := cmd.Run()
	return err == nil
}

func serviceIsActive(serviceName string) bool {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	err := cmd.Run()
	return err == nil
}

func restartService(serviceName string) error {
	cmd := exec.Command("systemctl", "restart", serviceName)
	return cmd.Run()
}

func checkUpdate() {
	resp, err := http.Get("https://unlock.icmp.ing/monitor/latest/version")
	if err != nil {
		log.Println("[ERR] 获取版本信息时出错:", err)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERR] 读取版本信息时出错:", err)
		return
	}
	version := string(b)
	if version == Version {
		fmt.Println("已经是最新版本")
		return
	}
	fmt.Println("检测到新版本", version)

	OS, ARCH := runtime.GOOS, runtime.GOARCH
	fmt.Println("OS:", OS)
	fmt.Println("ARCH:", ARCH)
	
	if OS == "linux" && strings.Contains(os.Getenv("PREFIX"), "com.termux") {
	    OS = "android"
	}

	url := "https://unlock.icmp.ing/monitor/latest/unlock-monitor_" + OS + "_" + ARCH
	if OS == "windows" {
	    url += ".exe"
	}

	resp, err = http.Get(url)
	if err != nil {
		log.Fatal("[ERR] 下载unlock-monitor时出错:", err)
		return
	}
	defer resp.Body.Close()
	
	bar := pb.DefaultBytes(
		resp.ContentLength,
		"下载进度",
	)

	body := io.TeeReader(resp.Body, bar)

	if resp.StatusCode != http.StatusOK {
		log.Fatal("[ERR] 下载unlock-monitor时出错: 非预期的状态码", resp.StatusCode)
		return
	}

	err = selfUpdate.Apply(body, selfUpdate.Options{})
	if err != nil {
		log.Fatal("[ERR] 更新unlock-monitor时出错:", err)
		return
	}

	fmt.Println("[OK] unlock-monitor后端更新成功")
	
	if serviceExists("unlock-monitor.service") && serviceIsActive("unlock-monitor.service") {
		err = restartService("unlock-monitor.service")
		if err != nil {
			log.Fatal("[ERR] 重启unlock-monitor服务时出错:", err)
			return
		}
		fmt.Println("[OK] unlock-monitor服务已重启")
	} else {
		fmt.Println("[INFO] unlock-monitor服务不存在或未运行，无需重启")
	}
}
