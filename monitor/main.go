package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	mt "MediaUnlockTest/checks"

	"github.com/kardianos/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/proxy"
)

var (
	AutoUpdate     bool
	UpdateInterval uint64
	Version        string = mt.Version
	buildTime      string
	Iface          string = ""
	DnsServers     string
	httpProxy      string
	socksProxy     string
)

type program struct {
	authToken   string
	metricsPath string
}

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
		for range ticker.C {
			checkUpdate()
		}
	}
}

func (p *program) run() {
	go recordMetrics()
	handler := promhttp.Handler()
	http.HandleFunc(p.metricsPath, func(w http.ResponseWriter, r *http.Request) {
		if p.authToken != "" {
			token := r.URL.Query().Get("token")
			if token == "" {
				token = r.Header.Get("token")
			}
			if token != p.authToken {
				w.Write([]byte("wrong token"))
				return
			}
		}
		handler.ServeHTTP(w, r)
	})
	log.Println("serve on " + Listen)
	http.ListenAndServe(Listen, nil)
}

func (p *program) Stop(s service.Service) error {
	return nil
}

var setSocketOptions = func(network, address string, c syscall.RawConn, interfaceName string) (err error) {
	return
}

func main() {
	var (
		install     bool
		uninstall   bool
		start       bool
		stop        bool
		restart     bool
		update      bool
		authToken   string
		metricsPath string
		version     bool
	)

	flag.Uint64Var(&Interval, "interval", 60, "check interval (s)")
	flag.Uint64Var(&UpdateInterval, "update-interval", 0, "update check interval (s)")
	flag.Uint64Var(&Conc, "conc", 0, "concurrency of tests")
	flag.StringVar(&Listen, "listen", ":9101", "listen address")
	flag.StringVar(&Node, "node", "", "Prometheus node field")
	flag.StringVar(&Iface, "I", "", "source ip / interface")
	flag.StringVar(&DnsServers, "dns-servers", "", "specify dns servers")
	flag.StringVar(&httpProxy, "http-proxy", "", "http proxy")
	flag.StringVar(&socksProxy, "socks-proxy", "", "socks5 proxy")
	flag.StringVar(&authToken, "token", "", "check token in http headers or queries")
	flag.StringVar(&metricsPath, "metrics-path", "/metrics", "custom metrics path")
	flag.BoolVar(&MUL, "mul", true, "Multination")
	flag.BoolVar(&HK, "hk", false, "Hong Kong")
	flag.BoolVar(&TW, "tw", false, "Taiwan")
	flag.BoolVar(&JP, "jp", false, "Japan")
	flag.BoolVar(&KR, "kr", false, "Korea")
	flag.BoolVar(&NA, "na", false, "North America")
	flag.BoolVar(&SA, "sa", false, "South America")
	flag.BoolVar(&EU, "eu", false, "Europe")
	flag.BoolVar(&AFR, "afr", false, "Africa")
	flag.BoolVar(&SEA, "sea", false, "South East Asia")
	flag.BoolVar(&OCEA, "ocea", false, "Oceania")
	flag.BoolVar(&AI, "ai", true, "AI Platforms")
	flag.BoolVar(&update, "u", false, "check update")
	flag.BoolVar(&version, "v", false, "show version")
	flag.BoolVar(&AutoUpdate, "auto-update", false, "set auto update")
	flag.BoolVar(&install, "install", false, "install service")
	flag.BoolVar(&uninstall, "uninstall", false, "uninstall service")
	flag.BoolVar(&start, "start", false, "start service")
	flag.BoolVar(&stop, "stop", false, "stop service")
	flag.BoolVar(&restart, "restart", false, "restart service")

	flag.Parse()

	if version {
		fmt.Println("unlock-monitor " + Version + " (" + runtime.GOOS + "_" + runtime.GOARCH + " " + runtime.Version() + " build-time: " + buildTime + ")")
		return
	}

	if update {
		checkUpdate()
		return
	}

	if Iface != "" {
		if IP := net.ParseIP(Iface); IP != nil {
			mt.Dialer.LocalAddr = &net.TCPAddr{IP: IP}
		} else {
			mt.Dialer.Control = func(network, address string, c syscall.RawConn) error {
				return setSocketOptions(network, address, c, Iface)
			}
		}
	}
	if DnsServers != "" {
		mt.Dialer.Resolver = &net.Resolver{
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "udp", DnsServers)
			},
		}

		mt.AutoHttpClient.Transport.(*mt.CustomTransport).Resolver = mt.Dialer.Resolver
	}
	if httpProxy != "" {
		if u, err := url.Parse(httpProxy); err == nil {
			mt.ClientProxy = http.ProxyURL(u)
			mt.AutoHttpClient.Transport.(*mt.CustomTransport).Proxy = mt.ClientProxy
		}
	}
	if socksProxy != "" {
		proxyURL, err := url.Parse(socksProxy)
		if err != nil {
			log.Fatal("SOCKS5 地址不合法：", err)
		}
		var auth *proxy.Auth
		if proxyURL.User != nil {
			username := proxyURL.User.Username()
			password, _ := proxyURL.User.Password()
			auth = &proxy.Auth{
				User:     username,
				Password: password,
			}
		}
		dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, auth, proxy.Direct)
		if err != nil {
			log.Fatal("创建 SOCKS5 连接失败：", err)
		}
		mt.AutoHttpClient.Transport.(*mt.CustomTransport).SocksDialer = dialer
	}

	args := []string{}
	for _, a := range os.Args[1:] {
		if !strings.Contains(a, "-install") && !strings.Contains(a, "-uninstall") {
			args = append(args, a)
		}
	}

	svcConfig := &service.Config{
		Name:        "unlock-monitor",
		DisplayName: "Media Unlock Monitor Service",
		Description: "Service to monitor media unlock status.",
		Arguments:   args,
	}

	prg := &program{
		authToken:   authToken,
		metricsPath: metricsPath,
	}
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

	if start {
		startService(s)
		return
	}

	if stop {
		stopService(s)
		return
	}

	if restart {
		restartService(s)
		return
	}

	err = s.Run()
	if err != nil {
		log.Fatal("[ERR] 运行服务时出现错误", err)
	}
}
