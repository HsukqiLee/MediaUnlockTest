package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	m "MediaUnlockTest"

	pb "github.com/schollz/progressbar/v3"
)

var IPV4 = true
var IPV6 = true
var M, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA bool
var Force bool

type result struct {
	Name    string
	Divider bool
	Value   m.Result
}

var tot int64
var R []*result
var bar *pb.ProgressBar
var wg *sync.WaitGroup

func excute(Name string, F func(client http.Client) m.Result, client http.Client) {
	r := &result{Name: Name}
	R = append(R, r)
	wg.Add(1)
	tot++
	go func() {
		r.Value = F(client)
		bar.Describe(Name + " " + ShowResult(r.Value))
		bar.Add(1)
		wg.Done()
	}()
}

var (
	FontRed     = "\033[31m"
	FontGreen   = "\033[32m"
	FontYellow  = "\033[33m"
	FontBlue    = "\033[34m"
	FontPurple  = "\033[35m"
	FontSkyBlue = "\033[36m"
	FontWhite   = "\033[37m"
	FontSuffix  = "\033[0m"
)

func InitColor() {
	if runtime.GOOS == "windows" {
		FontRed = ""
		FontGreen = ""
		FontYellow = ""
		FontBlue = ""
		FontPurple = ""
		FontSkyBlue = ""
		FontWhite = ""
		FontSuffix = ""
	}
}

func ShowResult(r m.Result) (s string) {
	if r.Status == m.StatusOK {
		s = FontGreen + "YES"
		if r.Region != "" {
			s += " (Region: " + strings.ToUpper(r.Region) + ")"
		}
		s += FontSuffix
		return s
	}
	if r.Status == m.StatusNetworkErr {
		return FontRed + "NO" + FontSuffix + FontYellow + " (Network Err)" + FontSuffix
	}
	if r.Status == m.StatusRestricted {
		if r.Info != "" {
			return FontYellow + "Restricted" + " (" + r.Info + ")" + FontSuffix
		}
		return FontYellow + "Restricted" + FontSuffix
	}
	if r.Status == m.StatusErr {
		s = FontYellow + "Error"
		if r.Err != nil {
			s += ": " + r.Err.Error() + ""
		}
		s += FontSuffix
		return s
	}
	if r.Status == m.StatusNo {
		if r.Info != "" {
			return FontRed + "NO" + FontSuffix + FontYellow + " (" + r.Info + ")" + FontSuffix
		}
		if r.Region != "" {
			return FontRed + "NO (Region: " + strings.ToUpper(r.Region) + ")" + FontSuffix
		}
		return FontRed + "NO" + FontSuffix
	}
	if r.Status == m.StatusBanned {
		if r.Info != "" {
			return FontRed + "Banned" + FontSuffix + FontYellow + " (" + r.Info + ")" + FontSuffix
		}
		return FontRed + "Banned" + FontSuffix
	}
	if r.Status == m.StatusUnexpected {
		return FontPurple + "Unexpected" + FontSuffix
	}
	if r.Status == m.StatusFailed {
		return "Failed"
	}
	return
}

func ShowR() {
	fmt.Println("测试时间: ", FontYellow+time.Now().Format("2006-01-02 15:04:05")+FontSuffix)
	NameLength := 25
	for _, r := range R {
		if len(r.Name) > NameLength {
			NameLength = len(r.Name)
		}
	}
	for _, r := range R {
		if r.Divider {
			s := "[ " + r.Name + " ] "
			for i := NameLength - len(s) + 4; i > 0; i-- {
				s += "="
			}
			if r.Name == "" {
				s = "\n"
			}
			fmt.Println(s)
		} else {
			result := ShowResult(r.Value)
			if r.Value.Status == m.StatusOK && strings.HasSuffix(r.Name, "CDN") {
				result = FontSkyBlue + r.Value.Region + FontSuffix
			}
			fmt.Printf("%-"+strconv.Itoa(NameLength)+"s %s\n", r.Name, result)
		}
	}
}

func NewBar(count int64) *pb.ProgressBar {
	return pb.NewOptions64(
		count,
		pb.OptionSetDescription("testing"),
		pb.OptionSetWriter(os.Stderr),
		pb.OptionSetWidth(20),
		pb.OptionThrottle(100*time.Millisecond),
		pb.OptionShowCount(),
		pb.OptionClearOnFinish(),
		pb.OptionEnableColorCodes(true),
		pb.OptionSpinnerType(14),
	)
}

func Multination(c http.Client) {
	R = append(R, &result{Name: "Multination", Divider: true})
	excute("Dazn", m.Dazn, c)
	excute("Hotstar", m.Hotstar, c)
	excute("Disney+", m.DisneyPlus, c)
	excute("Netflix", m.NetflixRegion, c)
	excute("Netflix CDN", m.NetflixCDN, c)
	excute("Youtube Premium", m.YoutubeRegion, c)
	excute("Youtube CDN", m.YoutubeCDN, c)
	excute("Amazon Prime Video", m.PrimeVideo, c)
	excute("TVBAnywhere+", m.TVBAnywhere, c)
	excute("iQiYi", m.IQiYi, c)
	excute("Viu.com", m.ViuCom, c)
	excute("Spotify", m.Spotify, c)
	excute("Steam", m.Steam, c)
	excute("ChatGPT", m.ChatGPT, c)
	excute("Wikipedia", m.WikipediaEditable, c)
	excute("Reddit", m.Reddit, c)
	excute("TikTok", m.TikTok, c)
	excute("Bing", m.Bing, c)
	excute("Instagram Audio", m.Instagram, c)
	excute("SonyLiv", m.SonyLiv, c)
}

func HongKong(c http.Client) {
	R = append(R, &result{Name: "Hong Kong", Divider: true})
	excute("Now E", m.NowE, c)
	excute("Viu.TV", m.ViuTV, c)
	excute("MyTVSuper", m.MyTvSuper, c)
	excute("HBO GO Aisa", m.HboGoAisa, c)
	excute("BiliBili HongKong/Macau Only", m.BilibiliHKMO, c)
}

func Taiwan(c http.Client) {
	R = append(R, &result{Name: "Taiwan", Divider: true})
	excute("KKTV", m.KKTV, c)
	excute("LiTV", m.LiTV, c)
	excute("MyVideo", m.MyVideo, c)
	excute("4GTV", m.TW4GTV, c)
	excute("LineTV", m.LineTV, c)
	excute("Hami Video", m.HamiVideo, c)
	excute("CatchPlay+", m.Catchplay, c)
	excute("Bahamut Anime", m.BahamutAnime, c)
	excute("HBO GO Aisa", m.HboGoAisa, c)
	excute("Bilibili Taiwan Only", m.BilibiliTW, c)
}

func Japan(c http.Client) {
	R = append(R, &result{Name: "Japan", Divider: true})
	excute("DMM", m.DMM, c)
	excute("DMM TV", m.DMMTV, c)
	excute("Abema", m.Abema, c)
	excute("Niconico", m.Niconico, c)
	excute("music.jp", m.MusicJP, c)
	excute("Telasa", m.Telasa, c)
	excute("Paravi", m.Paravi, c)
	excute("U-NEXT", m.U_NEXT, c)
	excute("Hulu Japan", m.HuluJP, c)
	excute("GYAO!", m.GYAO, c)
	excute("VideoMarket", m.VideoMarket, c)
	excute("FOD(Fuji TV)", m.FOD, c)
	excute("Radiko", m.Radiko, c)
	excute("Karaoke@DAM", m.Karaoke, c)
	excute("J:COM On Demand", m.J_COM_ON_DEMAND, c)
	excute("Kancolle", m.Kancolle, c)
	excute("Pretty Derby Japan", m.PrettyDerbyJP, c)
	excute("Konosuba Fantastic Days", m.KonosubaFD, c)
	excute("Princess Connect Re:Dive Japan", m.PCRJP, c)
	excute("World Flipper Japan", m.WFJP, c)
	excute("Project Sekai: Colorful Stage", m.PJSK, c)
	excute("Rakuten TV JP", m.RakutenTV_JP, c)
	excute("Wowow", m.Wowow, c)
}

func Korea(c http.Client) {
	R = append(R, &result{Name: "Korea", Divider: true})
	excute("Wavve", m.Wavve, c)
	excute("Tving", m.Tving, c)
	excute("Watcha", m.Watcha, c)
	excute("Coupang Play", m.CoupangPlay, c)
	excute("SpotvNow", m.SpotvNow, c)
	excute("NaverTV", m.NaverTV, c)
	excute("Afreeca", m.Afreeca, c)
	excute("KBS", m.KBS, c)
}

func NorthAmerica(c http.Client) {
	R = append(R, &result{Name: "North America", Divider: true})
	R = append(R, &result{Name: "US", Divider: true})
	excute("FOX", m.Fox, c)
	excute("Hulu", m.Hulu, c)
	excute("NFL+", m.NFLPlus, c)
	excute("ESPN+", m.ESPNPlus, c)
	excute("Epix", m.Epix, c)
	excute("Starz", m.Starz, c)
	excute("Philo", m.Philo, c)
	excute("FXNOW", m.FXNOW, c)
	excute("TLC GO", m.TlcGo, c)
	excute("HBO Max", m.HBOMax, c)
	excute("Shudder", m.Shudder, c)
	excute("BritBox", m.BritBox, c)
	excute("CW TV", m.CW_TV, c)
	excute("NBA TV", m.NBA_TV, c)
	excute("Fubo TV", m.FuboTV, c)
	excute("Tubi TV", m.TubiTV, c)
	excute("Sling TV", m.SlingTV, c)
	excute("Pluto TV", m.PlutoTV, c)
	excute("Acorn TV", m.AcornTV, c)
	excute("SHOWTIME", m.SHOWTIME, c)
	excute("encoreTVB", m.EncoreTVB, c)
	//excute("Funimation", m.Funimation, c)
	excute("Discovery+", m.DiscoveryPlus, c)
	excute("Paramount+", m.ParamountPlus, c)
	excute("Peacock TV", m.PeacockTV, c)
	excute("Popcornflix", m.Popcornflix, c)
	excute("Crunchyroll", m.Crunchyroll, c)
	excute("Direct Stream", m.DirectvStream, c)
	R = append(R, &result{Name: "CA", Divider: true})
	excute("CBC Gem", m.CBCGem, c)
	excute("Crave", m.Crave, c)
}

func SouthAmerica(c http.Client) {
    R = append(R, &result{Name: "South America", Divider: true})
    excute("Star Plus", m.StarPlus, c)
    excute("DirecTV GO", m.DirecTVGO, c)
}

func Europe(c http.Client) {
    R = append(R, &result{Name: "Europe", Divider: true})
    excute("Rakuten TV EU", m.RakutenTV_EU, c)
    excute("Setanta Sports", m.SetantaSports, c)
    excute("Sky Show Time", m.SkyShowTime, c)
    R = append(R, &result{Name: "GB", Divider: true})
    excute("BBC iPlayer", m.BBCiPlayer, c)
    excute("Channel 4", m.Channel4, c)
    excute("Channel 5", m.Channel5, c)
    excute("Sky Go", m.SkyGo, c)
    excute("ITVX", m.ITVX, c)
    R = append(R, &result{Name: "IT", Divider: true})
    excute("Rai Play", m.RaiPlay, c)
    R = append(R, &result{Name: "FR/DE", Divider: true})
    excute("Canal Plus", m.CanalPlus, c)
    excute("ZDF", m.ZDF, c)
    excute("Joyn", m.Joyn, c)
    excute("Molotov", m.Molotov, c)
    R = append(R, &result{Name: "NL", Divider: true})
    excute("NPO Start Plus", m.NPOStartPlus, c)
    excute("Video Land", m.VideoLand, c)
    excute("NLZIET", m.NLZIET, c)
    R = append(R, &result{Name: "ES", Divider: true})
    excute("Movistar Plus+", m.MoviStarPlus, c)
    R = append(R, &result{Name: "RU", Divider: true})
    excute("Amediateka", m.Amediateka, c)
}

func Africa(c http.Client) {
    R = append(R, &result{Name: "Africa", Divider: true})
    excute("DSTV", m.DSTV, c)
    excute("Showmax", m.Showmax, c)
}

func Oceania(c http.Client) {
    R = append(R, &result{Name: "Oceania", Divider: true})
    R = append(R, &result{Name: "AU", Divider: true})
    excute("Stan", m.Stan, c)
    excute("Binge", m.Binge, c)
    excute("7Plus", m.SevenPlus, c)
    excute("Channel 9", m.Channel9, c)
    excute("10 Play", m.Channel10, c)
    excute("ABC iView", m.ABCiView, c)
    excute("Optus Sports", m.OptusSports, c)
    excute("SBS on Demand", m.SBSonDemand, c)
    R = append(R, &result{Name: "NZ", Divider: true})
    excute("Neon TV", m.NeonTV, c)
    excute("Three Now", m.ThreeNow, c)
    excute("Maori TV", m.MaoriTV, c)
}

func Ipv6Multination() {
	c := m.Ipv6HttpClient
	R = append(R, &result{Name: "", Divider: true})
	R = append(R, &result{Name: "IPV6 Multination", Divider: true})
	excute("Hotstar", m.Hotstar, c)
	excute("Disney+", m.DisneyPlus, c)
	excute("Netflix", m.NetflixRegion, c)
	excute("Netflix CDN", m.NetflixCDN, c)
	excute("Youtube", m.YoutubeRegion, c)
	excute("Youtube CDN", m.YoutubeCDN, c)
	excute("Wikipedia", m.WikipediaEditable, c)
	excute("Bing", m.Bing, c)
}

func GetIpv4Info() {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.cloudflare.com/cdn-cgi/trace", nil)
	resp, err := m.Ipv4HttpClient.Do(req)
	if err != nil {
		IPV4 = false
		log.Println(err)
		fmt.Println("No IPv4 support")
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		IPV4 = false
		fmt.Println("No IPv4 support")
	}
	s := string(b)
	i := strings.Index(s, "ip=")
	s = s[i+3:]
	i = strings.Index(s, "\n")
	fmt.Println("Your IPV4 address:", FontSkyBlue, s[:i], FontSuffix)
}
func GetIpv6Info() {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.cloudflare.com/cdn-cgi/trace", nil)
	resp, err := m.Ipv6HttpClient.Do(req)
	if err != nil {
		IPV6 = false
		fmt.Println("No IPv6 support")
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("No IPv6 support")
	}
	s := string(b)
	i := strings.Index(s, "ip=")
	s = s[i+3:]
	i = strings.Index(s, "\n")
	fmt.Println("Your IPV6 address:", FontSkyBlue, s[:i], FontSuffix)
}

func ReadSelect() {
	fmt.Println("请选择检测项目,直接按回车将进行全部检测: ")
	fmt.Println("[0]: 跨国平台")
	fmt.Println("[1]: 台湾平台")
	fmt.Println("[2]: 香港平台")
	fmt.Println("[3]: 日本平台")
	fmt.Println("[4]: 韩国平台")
	fmt.Println("[5]: 北美平台")
	fmt.Println("[6]: 南美平台")
	fmt.Println("[7]: 欧洲平台")
	fmt.Println("[8]: 非洲平台")
	fmt.Println("[9]: 大洋洲平台")
	fmt.Print("请输入对应数字,空格分隔(回车确认): ")
	r := bufio.NewReader(os.Stdin)
	l, _, err := r.ReadLine()
	if err != nil {
		M, TW, HK, JP = true, true, true, true
		return
	}
	for _, c := range strings.Split(string(l), " ") {
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
		    OCEA = true
		default:
			M, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA = true, true, true, true, true, true, true, true, true, true
		}
	}
}

type Downloader struct {
	io.Reader
	Total   uint64
	Current uint64
	Pb      *pb.ProgressBar
	done    bool
}

func (d *Downloader) Read(p []byte) (n int, err error) {
	n, err = d.Reader.Read(p)
	d.Current += uint64(n)
	if d.done {
		return
	}
	d.Pb.Add(n)
	// fmt.Printf("\r正在下载，进度：%.2f%% [%s/%s]", float64(d.Current*10000/d.Total)/100, humanize.Bytes(d.Current), humanize.Bytes(d.Total))
	if d.Current == d.Total {
		d.done = true
		d.Pb.Describe("unlock-test下载完成")
		d.Pb.Finish()
	}
	return
}

func checkUpdate() {
	resp, err := http.Get("https://unlock.icmp.ing/test/latest/version")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	version := string(b)
	if version == m.Version {
		fmt.Println("已经是最新版本")
		return
	}
	fmt.Println("检测到新版本", version)
	OS, ARCH := runtime.GOOS, runtime.GOARCH
	fmt.Println("OS:", OS)
	fmt.Println("ARCH:", ARCH)
	
	targetDir := "/usr/bin/"
	switch OS {
	case "darwin":
		targetDir = "/usr/local/bin/"
	case "windows":
		programFiles := os.Getenv("ProgramFiles")
		if programFiles == "" {
			programFiles = "C:\\Program Files"
		}
		targetDir = programFiles + "\\MediaUnlockTest\\"
	}
	
	out, err := os.Create(targetDir+"unlock-test_new")
	if err != nil {
		log.Fatal("[ERR] 创建文件出错:", err)
		return
	}
	defer out.Close()
	log.Println("下载unlock-test中 ...")
	url := "https://unlock.icmp.ing/test/latest/unlock-test_" + runtime.GOOS + "_" + runtime.GOARCH
	resp, err = http.Get(url)
	if err != nil {
		log.Fatal("[ERR] 下载unlock-test时出错:", err)
	}
	defer resp.Body.Close()
	downloader := &Downloader{
		Reader: resp.Body,
		Total:  uint64(resp.ContentLength),
		Pb:     pb.DefaultBytes(resp.ContentLength, "下载进度"),
	}
	if _, err := io.Copy(out, downloader); err != nil {
		log.Fatal("[ERR] 下载unlock-test时出错:", err)
	}
	if OS != "windows" && os.Chmod(targetDir+"unlock-test_new", 0777) != nil {
		log.Fatal("[ERR] 更改unlock-test后端权限出错:", err)
	}
	if _, err := os.Stat(targetDir+"unlock-test"); err == nil {
		if err := os.Remove(targetDir+"unlock-test"); err != nil {
			log.Fatal("[ERR] 删除unlock-test旧版本时出错:", err.Error())
		}
	}
	if os.Rename(targetDir+"unlock-test_new", targetDir+"unlock-test") != nil {
		log.Fatal("[ERR] 更新unlock-test后端时出错:", err)
	}
	if OS == "windows" && os.Rename(targetDir+"unlock-test", targetDir+"unlock-test.exe") != nil {
		log.Fatal("[ERR] 更新unlock-test后端时出错:", err)
	}
	log.Println("[OK] unlock-test后端更新成功")
}


func showCounts() {
	resp, err := http.Get("https://unlock.moe/count.php")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	s := strings.Split(string(b), " ")
	d, m, t := s[0], s[1], s[3]
	fmt.Printf("当天运行共%s次, 本月运行共%s次, 共计运行%s次\n", FontSkyBlue+d+FontSuffix, FontYellow+m+FontSuffix, FontGreen+t+FontSuffix)
}

func showAd() {
	resp, err := http.Get("https://unlock.icmp.ing/ad.txt")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(b))
}

var setSocketOptions = func(network, address string, c syscall.RawConn, interfaceName string) (err error) {
	return
}

func main() {
	client := m.AutoHttpClient
	mode := 0
	showVersion := false
	update := false
	nf := false
	test := false
	Iface := ""
	DnsServers := ""
	httpProxy := ""
	flag.IntVar(&mode, "m", 0, "mode 0(default)/4/6")
	flag.BoolVar(&Force, "f", false, "ipv6 force")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&update, "u", false, "update")
	flag.StringVar(&Iface, "I", "", "source ip / interface")
	flag.StringVar(&DnsServers, "dns-servers", "", "specify dns servers")
	flag.StringVar(&httpProxy, "http-proxy", "", "http proxy")
	flag.BoolVar(&nf, "nf", false, "netflix")
	flag.BoolVar(&test, "test", false, "test")
	flag.Parse()
	InitColor()
	if showVersion {
		fmt.Println(m.Version)
		return
	}
	if update {
		checkUpdate()
		return
	}
	if Iface != "" {
		if IP := net.ParseIP(Iface); IP != nil {
			m.Dialer.LocalAddr = &net.TCPAddr{IP: IP}
		} else {
			m.Dialer.Control = func(network, address string, c syscall.RawConn) error {
				return setSocketOptions(network, address, c, Iface)
			}
		}
	}
	if DnsServers != "" {
		m.Dialer.Resolver = &net.Resolver{
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "udp", DnsServers)
			},
		}
	}
	if httpProxy != "" {
		log.Println(httpProxy)
		// c := httpproxy.Config{HTTPProxy: httpProxy, CGI: true}
		// m.ClientProxy = func(req *http.Request) (*url.URL, error) { return c.ProxyFunc()(req.URL) }
		if u, err := url.Parse(httpProxy); err == nil {
			m.ClientProxy = http.ProxyURL(u)
			m.Ipv4Transport.Proxy = m.ClientProxy
			m.Ipv4HttpClient.Transport = m.Ipv4Transport
			m.Ipv6Transport.Proxy = m.ClientProxy
			m.Ipv6HttpClient.Transport = m.Ipv6Transport
			m.AutoTransport.Proxy = m.ClientProxy
			m.AutoHttpClient.Transport = m.AutoTransport
		}
	}
	if mode == 4 {
		client = m.Ipv4HttpClient
		IPV6 = false
	}
	if mode == 6 {
		client = m.Ipv6HttpClient
		IPV4 = false
		M = true
	}

	if nf {
		fmt.Println("Netflix", ShowResult(m.NetflixRegion(m.AutoHttpClient)))
		return
	}
	
	if test {
	    fmt.Println("wowow", ShowResult(m.Wowow(m.Ipv4HttpClient)))
		return
	}

	fmt.Println("项目地址: " + FontSkyBlue + "https://github.com/HsukqiLee/MediaUnlockTest" + FontSuffix)
	fmt.Println("使用方式: " + FontYellow + "bash <(curl -Ls unlock.icmp.ing/test.sh)" + FontSuffix)
	fmt.Println()

	GetIpv4Info()
	GetIpv6Info()

	if IPV4 || Force {
		ReadSelect()
	}
	wg = &sync.WaitGroup{}
	bar = NewBar(0)
	if IPV4 {
		if M {
			Multination(client)
		}
		if TW {
			Taiwan(client)
		}
		if HK {
			HongKong(client)
		}
		if JP {
			Japan(client)
		}
		if KR {
			Korea(client)
		}
		if NA {
			NorthAmerica(client)
		}
		if SA {
			SouthAmerica(client)
		}
		if EU {
			Europe(client)
		}
		if AFR {
			Africa(client)
		}
		if OCEA {
			Oceania(client)
		}
	}
	if IPV6 {
		if Force {
			if M {
				Multination(m.Ipv6HttpClient)
			}
			if TW {
				Taiwan(m.Ipv6HttpClient)
			}
			if HK {
				HongKong(m.Ipv6HttpClient)
			}
			if JP {
				Japan(m.Ipv6HttpClient)
			}
			if KR {
				Korea(m.Ipv6HttpClient)
			}
			if NA {
				NorthAmerica(m.Ipv6HttpClient)
			}
			if SA {
		    	SouthAmerica(m.Ipv6HttpClient)
	    	}
	    	if EU {
	    		Europe(m.Ipv6HttpClient)
	    	}
	    	if AFR {
	    		Africa(m.Ipv6HttpClient)
	    	}
	    	if OCEA {
			    Oceania(m.Ipv6HttpClient)
		    }
		} else {
			Ipv6Multination()
		}
	}
	bar.ChangeMax64(tot)

	wg.Wait()
	bar.Finish()
	fmt.Println()
	ShowR()
	fmt.Println()
	fmt.Println("检测完毕，感谢您的使用!")
	showCounts()
	fmt.Println()
	showAd()
}
