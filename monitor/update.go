package main

import (
	core "MediaUnlockTest/pkg/core"
	"log"

	"github.com/kardianos/service"
)

// checkUpdate is called both by the foreground CLI (`monitor -u`) and the background AutoUpdate goroutine.
func checkUpdate(silent bool) {
	cfg := core.UpdateConfig{
		AppName:         "unlock-monitor",
		VersionURL:      "https://unlock.icmp.ing/monitor/latest/version",
		BinaryURLPrefix: "https://unlock.icmp.ing/monitor/latest/unlock-monitor",
		Silent:          silent,
	}

	updated := core.CheckUpdate(cfg)
	if updated {
		serviceConfig := &service.Config{
			Name:        "unlock-monitor",
			DisplayName: "Media Unlock Monitor Service",
			Description: "Service to monitor media unlock status.",
		}

		prg := &program{}
		s, err := service.New(prg, serviceConfig)
		if err != nil {
			log.Fatal(err)
		}

		if serviceExists(s) && serviceIsActive(s) {
			if !silent {
				log.Println("[INFO] 正在重启 unlock-monitor 服务...")
			}
			restartService(s)
		} else {
			if !silent {
				log.Println("[INFO] unlock-monitor 服务未运行")
			}
		}
	}
}
