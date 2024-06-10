package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
	"os"
	"strings"
    
    mt "MediaUnlockTest"
	"github.com/kardianos/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	AutoUpdate     bool
	UpdateInterval uint64
	Version        string = mt.Version
	buildTime      string
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.scheduleUpdate()
	go p.run()
	return nil
}

func (p *program) scheduleUpdate() {
	if AutoUpdate {
		if UpdateInterval == 0 {
			UpdateInterval = 86400
		}
		ticker := time.NewTicker(time.Duration(UpdateInterval) * time.Second)
		for {
			select {
		    case <-ticker.C:
				checkUpdate()
			}
		}
	}
}

func (p *program) run() {
	go recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	log.Println("serve on " + Listen)
	http.ListenAndServe(Listen, nil)
}

func (p *program) Stop(s service.Service) error {
	//log.Println("Service is stopping...")
	return nil
}

func main() {
	var install bool
	var uninstall bool
	var update bool
	var version bool

	flag.Uint64Var(&Interval, "interval", 60, "check interval (s)")
	flag.Uint64Var(&UpdateInterval, "update-interval", 0, "update check interval (s)")
	flag.StringVar(&Listen, "listen", ":9101", "listen address")
	flag.StringVar(&Node, "node", "", "node")
	flag.BoolVar(&MUL, "mul", true, "Mutation")
	flag.BoolVar(&HK, "hk", false, "Hong Kong")
	flag.BoolVar(&TW, "tw", false, "Taiwan")
	flag.BoolVar(&JP, "jp", false, "Japan")
	flag.BoolVar(&KR, "kr", false, "Korea")
	flag.BoolVar(&NA, "na", false, "North America")
	flag.BoolVar(&SA, "sa", false, "South America")
	flag.BoolVar(&EU, "eu", false, "Europe")
	flag.BoolVar(&AFR, "afr", false, "Africa")
	flag.BoolVar(&OCEA, "ocea", false, "Oceania")
	flag.BoolVar(&update, "u", false, "check update")
	flag.BoolVar(&version, "v", false, "show version")
	flag.BoolVar(&AutoUpdate, "auto-update", false, "set auto update")
	flag.BoolVar(&install, "install", false, "install service")
	flag.BoolVar(&uninstall, "uninstall", false, "uninstall service")

	flag.Parse()

	if version {
		fmt.Println("unlock-monitor " + Version + " (" + runtime.GOOS + "_" + runtime.GOARCH + " " + runtime.Version() + " build-time: " + buildTime + ")")
		return
	}
	
	if update {
	    checkUpdate()
		return
	}
	
    args := []string{}
	for _, a := range os.Args[1:] {
		if !strings.Contains(a, "-install") && !strings.Contains(a, "-uninstall") {
			args = append(args, a)
		}
	}
	
	svcConfig := &service.Config{
		Name:        "unlock-monitor",
		DisplayName: "unlock-monitor",
		Description: "Service to monitor media unlock status.",
		Arguments:   args,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if install {
	    installService(s)
	    return
	}

	if uninstall {
	    uninstallService(s)
	    return
	}

	err = s.Run()
	if err != nil {
		log.Fatal("[ERR] 运行服务时出现错误", err)
	}
}
