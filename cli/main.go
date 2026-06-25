package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	core "MediaUnlockTest/pkg/core"
	m "MediaUnlockTest/pkg/providers"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/proxy"
)

var (
	IPV4        bool = true
	IPV6        bool = true
	M           bool
	HK          bool
	TW          bool
	JP          bool
	KR          bool
	NA          bool
	SA          bool
	EU          bool
	AFR         bool
	SEA         bool
	OCEA        bool
	AI          bool
	Debug       bool   = false
	Conc        uint64 = 0
	Cache       bool   = false
	sem         chan struct{}
	ResultLines []*result
	bar         *progressbar.ProgressBar
	Red         = color.New(color.FgRed).SprintFunc()
	Green       = color.New(color.FgGreen).SprintFunc()
	Yellow      = color.New(color.FgYellow).SprintFunc()
	Blue        = color.New(color.FgBlue).SprintFunc()
	Purple      = color.New(color.FgMagenta).SprintFunc()
	SkyBlue     = color.New(color.FgCyan).SprintFunc()
	White       = color.New(color.FgWhite).SprintFunc()
	resultCache = make(map[string]core.Result)
	cacheMutex  sync.RWMutex

	// 全局超时控制 - 优化超时时间
	testTimeout   = 15 * time.Second // 单个测试超时时间从30秒减少到15秒
	regionTimeout = 3 * time.Minute  // 整个地区测试超时时间从5分钟减少到3分钟

	// 新增：正在进行的测试状态管理
	activeTestsMutex sync.RWMutex
	activeTests      = make(map[string]bool)

	// 新增：全局进度条显示控制
	ShowActive bool = true // 是否显示正在进行的测试

	// 新增：进度条描述缓存，避免重复更新
	progressDescriptionCache string
	progressDescMu           sync.Mutex

	// 进度条更新协程控制
	updaterStopChan chan struct{}
	updaterMutex    sync.Mutex
)

type TestItem struct {
	Name       string
	Func       func(client http.Client) core.Result
	SupportsV6 bool
}

type result struct {
	Name    string
	Divider bool
	Value   core.Result
}

type regionItem struct {
	Enabled bool
	Name    string
	Tests   []m.TestItem
}

func ReadSelect() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signalChan)

	fmt.Println("请选择检测项目：")
	fmt.Println(core.Green("直接按回车进行全部检测") + "，" + core.Yellow("按 Ctrl+C 取消检测") + "。")
	fmt.Println("")
	fmt.Println("[0]  : 　跨国平台")
	fmt.Println("[1]  : 　台湾平台")
	fmt.Println("[2]  : 　香港平台")
	fmt.Println("[3]  : 　日本平台")
	fmt.Println("[4]  : 　韩国平台")
	fmt.Println("[5]  : 　北美平台")
	fmt.Println("[6]  : 　南美平台")
	fmt.Println("[7]  : 　欧洲平台")
	fmt.Println("[8]  : 　非洲平台")
	fmt.Println("[9]  : 东南亚平台")
	fmt.Println("[10] : 大洋洲平台")
	fmt.Println("[11] : 　ＡＩ平台")
	fmt.Println("")
	fmt.Print("请输入对应数字，空格分隔，回车确认: ")

	inputChan := make(chan string, 1)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-signalChan:
				fmt.Println("")
				fmt.Println(core.Yellow("输入中止，检测已取消。"))
				os.Exit(0)
			default:
				if input, err := reader.ReadString('\n'); err == nil {
					inputChan <- strings.TrimSpace(input)
					return
				}
			}
		}
	}()

	select {
	case <-signalChan:
		fmt.Println("")
		fmt.Println(core.Yellow("输入中止，检测已取消。"))
		os.Exit(0)
	case input := <-inputChan:
		for _, c := range strings.Split(input, " ") {
			switch c {
			case "0":
				M = true
			case "1":
				TW = true
			case "2":
				HK = true
			case "3":
				JP = true
			case "4":
				KR = true
			case "5":
				NA = true
			case "6":
				SA = true
			case "7":
				EU = true
			case "8":
				AFR = true
			case "9":
				SEA = true
			case "10":
				OCEA = true
			case "11":
				AI = true
			default:
				M, TW, HK, JP, KR, NA, SA, EU, AFR, SEA, OCEA, AI = true, true, true, true, true, true, true, true, true, true, true, true
			}
		}
	}
}


func main() {
	var (
		Interface   string
		DNSServers  string
		HTTPProxy   string
		SocksProxy  string
		ShowVersion bool
		CheckUpdate bool
		NF          bool
		TestMode    bool
		IPMode      int
		IP4_1       string
		IP6_1       string
		IP4_2       string
		IP6_2       string
		err         error
		IsProxy     bool
	)
	flag.StringVar(&Interface, "I", "", "Source IP or network interface to use for connections")
	flag.StringVar(&DNSServers, "dns-servers", "", "Custom DNS servers (format: ip:port)")
	flag.StringVar(&HTTPProxy, "http-proxy", "", "HTTP proxy URL (format: http://user:pass@host:port)")
	flag.StringVar(&SocksProxy, "socks-proxy", "", "SOCKS5 proxy URL (format: socks5://user:pass@host:port)")
	flag.BoolVar(&ShowVersion, "v", false, "Show version information and exit")
	flag.BoolVar(&CheckUpdate, "u", false, "Update to latest version")
	flag.BoolVar(&NF, "nf", false, "Only test Netflix availability")
	flag.BoolVar(&TestMode, "test", false, "Run in test mode (only checks Viaplay)")
	flag.BoolVar(&Debug, "debug", false, "Enable debug mode for verbose output")
	flag.IntVar(&IPMode, "m", 0, "Connection mode: 0=auto (default), 4=IPv4 only, 6=IPv6 only")
	flag.Uint64Var(&Conc, "conc", 0, "Max concurrent tests (0=unlimited)")
	flag.BoolVar(&ShowActive, "show-active", true, "Show active tests in progress bar (default: true)")
	flag.BoolVar(&Cache, "cache", false, "Enable caching and sequential region execution (default: false)")
	flag.Parse()
	if ShowVersion {
		fmt.Println(core.Version)
		return
	}
	if CheckUpdate {
		checkUpdate()
		return
	}
	if Interface != "" {
		if IP := net.ParseIP(Interface); IP != nil {
			core.Dialer.LocalAddr = &net.TCPAddr{IP: IP}
		} else {
			core.Dialer.Control = func(network, address string, c syscall.RawConn) error {
				return core.SetSocketOptions(network, address, c, Interface)
			}
		}
	}
	if DNSServers != "" {
		core.Dialer.Resolver = &net.Resolver{
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "udp", DNSServers)
			},
		}
		core.Ipv4Transport.Resolver = core.Dialer.Resolver
		core.Ipv6Transport.Resolver = core.Dialer.Resolver
		core.AutoHttpClient.Transport.(*core.CustomTransport).Resolver = core.Dialer.Resolver
	}
	if HTTPProxy != "" {
		if u, err := url.Parse(HTTPProxy); err == nil {
			core.ClientProxy = http.ProxyURL(u)
			core.Ipv4Transport.Proxy = core.ClientProxy
			core.Ipv6Transport.Proxy = core.ClientProxy
			core.AutoHttpClient.Transport.(*core.CustomTransport).Proxy = core.ClientProxy
		}
	}
	if SocksProxy != "" {
		proxyURL, err := url.Parse(SocksProxy)
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

		core.Ipv4Transport.SocksDialer = dialer
		core.Ipv6Transport.SocksDialer = dialer
		core.AutoHttpClient.Transport.(*core.CustomTransport).SocksDialer = dialer
	}
	if Conc > 0 {
		sem = make(chan struct{}, Conc)
	}

	if NF {
		fmt.Println("Netflix", ShowSingleResult(m.NetflixRegion(core.AutoHttpClient)))
		return
	}

	if TestMode {

		fmt.Println("Amediateka", ShowSingleResult(m.Amediateka(core.AutoHttpClient)))

		return
	}

	fmt.Println("")
	fmt.Println("[ 项目地址: " + core.SkyBlue("https://github.com/HsukqiLee/MediaUnlockTest") + " ]")
	fmt.Println("[ 使用方式: " + core.Yellow("bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)") + " ]")
	fmt.Println()

	if !Debug {
		info4, err := core.GetDetailedIPInfo("https://unlock.icmp.ing/api/ip-info", 4)
		if err != nil {
			fmt.Println(core.Red("无法获取 IPv4 地址"))
			IPV4 = false
		} else {
			IP4_2 = info4.IP
			fmt.Println(core.SkyBlue("IPv4 地址：") + core.Green(info4.IP))
			fmt.Println(core.SkyBlue("地区：") + core.Yellow(info4.Country) + core.SkyBlue(" / ") + core.Yellow(info4.Region) + core.SkyBlue(" / ") + core.Yellow(info4.City))
			fmt.Println(core.SkyBlue("ISP：") + core.Green(info4.Organization) + core.Purple(" (AS"+strconv.Itoa(info4.ASN)+")"))
			IPV4 = true
		}
		info6, err := core.GetDetailedIPInfo("https://unlock.icmp.ing/api/ip-info", 6)
		if err != nil {
			fmt.Println(core.Red("无法获取 IPv6 地址"))
			IPV6 = false
		} else {
			IP6_2 = info6.IP
			fmt.Println(core.SkyBlue("IPv6 地址：") + core.Green(info6.IP))
			fmt.Println(core.SkyBlue("地区：") + core.Yellow(info6.Country) + core.SkyBlue(" / ") + core.Yellow(info6.Region) + core.SkyBlue(" / ") + core.Yellow(info6.City))
			fmt.Println(core.SkyBlue("ISP：") + core.Green(info6.Organization) + core.Purple(" (AS"+strconv.Itoa(info6.ASN)+")"))
			IPV6 = true
		}
	} else {
		fmt.Println("[ 正在获取国内分流 IP... ]")
		if IPMode == 0 || IPMode == 4 {
			IP4_1, err = core.GetIPInfo("http://4.itdog.cn/", 4, "plain")
			if err != nil {
				if Debug {
					fmt.Println(core.Red("无法获取国内分流 IPv4 地址 (") + core.Yellow(err.Error()) + core.Red(")"))
				} else {
					fmt.Println(core.Red("无法获取国内分流 IPv4 地址"))
				}
			} else {
				fmt.Println(core.SkyBlue("IPv4 地址： ") + core.Green(IP4_1))
			}
		}
		if IPMode == 0 || IPMode == 6 {
			IP6_1, err = core.GetIPInfo("http://6.itdog.cn/", 6, "plain")
			if err != nil {
				if Debug {
					fmt.Println(core.Red("无法获取国内分流 IPv6 地址 (") + core.Yellow(err.Error()) + core.Red(")"))
				} else {
					fmt.Println(core.Red("无法获取国内分流 IPv6 地址"))
				}
			} else {
				fmt.Println(core.SkyBlue("IPv6 地址： ") + core.Green(IP6_1))
			}
		}
		fmt.Println("")
		fmt.Println("[ 正在获取国外分流 IP... ]")
		if IPMode == 0 || IPMode == 4 {
			info4, err := core.GetDetailedIPInfo("https://unlock.icmp.ing/api/ip-info", 4)
			if err != nil {
				if Debug {
					fmt.Println(core.Red("无法获取国外 IPv4 地址 (") + core.Yellow(err.Error()) + core.Red(")"))
				} else {
					fmt.Println(core.Red("无法获取国外 IPv4 地址"))
				}
			} else {
				IP4_2 = info4.IP
				fmt.Println(core.SkyBlue("IPv4 地址：") + core.Green(info4.IP))
				fmt.Println(core.SkyBlue("地区：") + core.Yellow(info4.Country) + core.SkyBlue("/") + core.Yellow(info4.Region) + core.SkyBlue("/") + core.Yellow(info4.City))
				fmt.Println(core.SkyBlue("ISP：") + core.Green(info4.Organization) + core.Purple(" (AS"+strconv.Itoa(info4.ASN)+")"))
			}
		}
		if IPMode == 0 || IPMode == 6 {
			info6, err := core.GetDetailedIPInfo("https://unlock.icmp.ing/api/ip-info", 6)
			if err != nil {
				if Debug {
					fmt.Println(core.Red("无法获取国外 IPv6 地址 (") + core.Yellow(err.Error()) + core.Red(")"))
				} else {
					fmt.Println(core.Red("无法获取国外 IPv6 地址"))
				}
			} else {
				IP6_2 = info6.IP
				fmt.Println(core.SkyBlue("IPv6 地址：") + core.Green(info6.IP))
				fmt.Println(core.SkyBlue("地区：") + core.Yellow(info6.Country) + core.SkyBlue("/") + core.Yellow(info6.Region) + core.SkyBlue("/") + core.Yellow(info6.City))
				fmt.Println(core.SkyBlue("ISP：") + core.Green(info6.Organization) + core.Purple(" (AS"+strconv.Itoa(info6.ASN)+")"))
			}
		}
		fmt.Println("")
		fmt.Println("[ 正在检测系统代理... ]")

		if IPMode == 0 || IPMode == 4 {
			IP4, err := core.GetIPInfo("https://www.cloudflare.com/cdn-cgi/trace", 4, "cloudflare")
			if err != nil {
				if IP4_1 != "" || IP4_2 != "" {
					IsProxy = true
					fmt.Println(core.Yellow("正在使用系统代理，且无法通过 IPv4 连接代理"))
				} else {
					IPV4 = false
					fmt.Println(core.Red("未使用 IPv4 代理，无 IPv4 网络"))
				}
			} else {
				IPV4 = true
				if IP4_1 != IP4_2 || IP4_1 != IP4 {
					IsProxy = true
					fmt.Println(core.Yellow("正在使用监听地址为 IPv4 的代理，出口 IP：") + core.Red(IP4))
				} else if IP4 == IP4_1 {
					fmt.Println(core.Green("未使用 IPv4 代理，有 IPv4 网络"))
				} else {
					fmt.Println(core.Red("无法强制使用 IPv4 网络测试，可能使用 IPv4 代理"))
					IPV4 = false
					if IPMode == 4 {
						IPV6 = false
					}
				}
			}
		}
		if IPMode == 0 || IPMode == 6 {
			IP6, err := core.GetIPInfo("https://www.cloudflare.com/cdn-cgi/trace", 6, "cloudflare")
			if err != nil {
				if IP6_1 != "" && IP6_2 != "" {
					IsProxy = true
					fmt.Println(core.Yellow("正在使用系统代理，且无法通过 IPv6 连接代理"))
				} else {
					IPV6 = false
					fmt.Println(core.Red("未使用 IPv6 代理，无 IPv6 网络"))
				}
			} else {
				IPV6 = true
				if IP6_1 != IP6_2 && IP6_1 != IP6 {
					IsProxy = true
					fmt.Println(core.Yellow("正在使用监听地址为 IPv6 的代理，出口 IP：") + core.Red(IP6))
				} else if IP6 == IP6_1 {
					fmt.Println(core.Green("未使用 IPv6 代理，有 IPv6 网络"))
				} else {
					fmt.Println(core.Red("无法强制使用 IPv6 网络测试，可能使用 IPv6 代理"))
					IPV6 = false
					if IPMode == 6 {
						IPV4 = false
					}
				}
			}
		}
	}

	if IsProxy {
		fmt.Println(core.Yellow("提示：正在使用系统代理，此时连接行为全部受代理控制"))
	}
	if IPMode != 0 {
		switch IPMode {
		case 4:
			IPV6 = false
		case 6:
			IPV4 = false
		}
	}
	fmt.Println()

	if IPV4 || IPV6 {
		ReadSelect()
	}
	regions := []regionItem{
		{Enabled: M, Name: "Globe", Tests: m.GlobeTests},
		{Enabled: TW, Name: "Taiwan", Tests: m.TaiwanTests},
		{Enabled: HK, Name: "HongKong", Tests: m.HongKongTests},
		{Enabled: JP, Name: "Japan", Tests: m.JapanTests},
		{Enabled: KR, Name: "Korea", Tests: m.KoreaTests},
		{Enabled: NA, Name: "NorthAmerica", Tests: m.NorthAmericaTests},
		{Enabled: SA, Name: "SouthAmerica", Tests: m.SouthAmericaTests},
		{Enabled: EU, Name: "Europe", Tests: m.EuropeTests},
		{Enabled: AFR, Name: "Africa", Tests: m.AfricaTests},
		{Enabled: SEA, Name: "SouthEastAsia", Tests: m.SouthEastAsiaTests},
		{Enabled: OCEA, Name: "Oceania", Tests: m.OceaniaTests},
		{Enabled: AI, Name: "AI", Tests: m.AITests},
	}
	if IsProxy {
		if Cache {
			ExecuteTests(regions, core.AutoHttpClient, 0)
		} else {
			ExecuteTestsParallel(regions, core.AutoHttpClient, 0)
		}
	} else {
		if IPV4 {
			if Cache {
				ExecuteTests(regions, core.Ipv4HttpClient, 4)
			} else {
				ExecuteTestsParallel(regions, core.Ipv4HttpClient, 4)
			}
		}
		if IPV6 {
			if Cache {
				ExecuteTests(regions, core.Ipv6HttpClient, 6)
			} else {
				ExecuteTestsParallel(regions, core.Ipv6HttpClient, 6)
			}
		}
	}
	fmt.Println()
	ShowFinalResult()
	fmt.Println()
	fmt.Println("检测完毕，感谢您的使用！")
	ShowCounts()
	fmt.Println()
	ShowAD()
}
