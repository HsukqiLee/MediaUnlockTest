package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RunCmd(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InstallService() {
	args := []string{}
	for _, a := range os.Args[1:] {
		if !strings.Contains(a, "-service") {
			args = append(args, a)
		}
	}
	data := []byte(`[Unit]
Description=unlock-monitor
After=network.target
[Service]
Type=simple
LimitCPU=infinity
LimitFSIZE=infinity
LimitDATA=infinity
LimitSTACK=infinity
LimitCORE=infinity
LimitRSS=infinity
LimitNOFILE=infinity
LimitAS=infinity
LimitNPROC=infinity
LimitMEMLOCK=infinity
LimitLOCKS=infinity
LimitSIGPENDING=infinity
LimitMSGQUEUE=infinity
LimitRTPRIO=infinity
LimitRTTIME=infinity
ExecStart=/usr/bin/unlock-monitor ` + strings.Join(args, " ") + `
Restart=always
RestartSec=5
[Install]
WantedBy=multi-user.target`)
	if err := ioutil.WriteFile("/etc/systemd/system/unlock-monitor.service", data, 0644); err != nil {
		log.Fatal("[ERR] 写入systemd守护service时出错:", err)
	}
	log.Println("[OK] systemd守护service成功")
	if err := RunCmd("systemctl", "daemon-reload"); err != nil {
		log.Fatal("[ERR] 重载systemctl时出错:", err)
	}
	if err := RunCmd("systemctl", "restart", "unlock-monitor"); err != nil {
		log.Fatal("[ERR] 启动unlock-monitor服务时出错:", err)
	}
	if err := RunCmd("systemctl", "enable", "unlock-monitor"); err != nil {
		log.Fatal("[ERR] 设置unlock-monitor服务开机自启时出错:", err)
	}
	log.Println("[OK] 初始化服务成功")
}

func AutoUpdateService() {
	var interval string

	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "-interval" && i+1 < len(os.Args) {
			interval = os.Args[i+1]
			break
		}
	}

	if interval == "" {
	    interval = "1d"
	} 
	data1 := []byte(`[Unit]
Description=unlock-monitor auto update service
After=network.target
[Service]
Type=oneshot
LimitCPU=infinity
LimitFSIZE=infinity
LimitDATA=infinity
LimitSTACK=infinity
LimitCORE=infinity
LimitRSS=infinity
LimitNOFILE=infinity
LimitAS=infinity
LimitNPROC=infinity
LimitMEMLOCK=infinity
LimitLOCKS=infinity
LimitSIGPENDING=infinity
LimitMSGQUEUE=infinity
LimitRTPRIO=infinity
LimitRTTIME=infinity
ExecStart=/usr/bin/unlock-monitor -u
[Install]
WantedBy=multi-user.target`)

    data2 := []byte(`[Unit]
Description=Timer for unlock-monitor auto update service

[Timer]
OnActiveSec=`+ interval +`s
Persistent=true

[Install]
WantedBy=timers.target`)

	if err := ioutil.WriteFile("/etc/systemd/system/unlock-monitor-update.service", data1, 0644); err != nil {
		log.Fatal("[ERR] 写入systemd守护service时出错:", err)
	}
	if err := ioutil.WriteFile("/etc/systemd/system/unlock-monitor-update.timer", data2, 0644); err != nil {
		log.Fatal("[ERR] 写入systemd服务timer时出错:", err)
	}
	log.Println("[OK] systemd设置service和timer成功")
	if err := RunCmd("systemctl", "daemon-reload"); err != nil {
		log.Fatal("[ERR] 重载systemctl时出错:", err)
	}
	if err := RunCmd("systemctl", "restart", "unlock-monitor-update.timer"); err != nil {
		log.Fatal("[ERR] 启动unlock-monitor-update服务时出错:", err)
	}
	if err := RunCmd("systemctl", "enable", "unlock-monitor-update.timer"); err != nil {
		log.Fatal("[ERR] 设置unlock-monitor-update服务开机自启时出错:", err)
	}
	log.Println("[OK] 初始化更新服务成功")
}

func UninstallService() {
	if err := RunCmd("systemctl", "stop", "unlock-monitor"); err != nil {
		log.Fatal("[ERR] 停止unlock-monitor服务时出错:", err)
	}
	if err := RunCmd("systemctl", "disable", "unlock-monitor"); err != nil {
		log.Fatal("[ERR] 取消unlock-monitor服务开机自启时出错:", err)
	}
	if err := os.Remove("/etc/systemd/system/unlock-monitor.service"); err != nil {
		log.Fatal("[ERR] 删除unlock-monitor服务时出错:", err)
	}
	if err := RunCmd("systemctl", "stop", "unlock-monitor-update.timer"); err != nil {
		log.Fatal("[ERR] 停止unlock-monitor更新服务时出错:", err)
	}
	if err := RunCmd("systemctl", "disable", "unlock-monitor-update.timer"); err != nil {
		log.Fatal("[ERR] 取消unlock-monitor更新服务开机自启时出错:", err)
	}
	if err := os.Remove("/etc/systemd/system/unlock-monitor-update.service"); err != nil {
		log.Fatal("[ERR] 删除unlock-monitor更新服务时出错:", err)
	}
	if err := os.Remove("/etc/systemd/system/unlock-monitor-update.timer"); err != nil {
		log.Fatal("[ERR] 删除unlock-monitor更新服务定时器时出错:", err)
	}
	log.Println("[OK] 卸载服务成功")
}
