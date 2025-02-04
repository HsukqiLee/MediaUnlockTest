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

	"github.com/fatih/color"
	selfUpdate "github.com/inconshreveable/go-update"
	pb "github.com/schollz/progressbar/v3"
	"golang.org/x/net/proxy"
)

var (
	IPV4    bool = true
	IPV6    bool = true
	M       bool
	HK      bool
	TW      bool
	JP      bool
	KR      bool
	NA      bool
	SA      bool
	EU      bool
	AFR     bool
	SEA     bool
	OCEA    bool
	Debug   bool   = false
	Conc    uint64 = 0
	sem     chan struct{}
	tot     int64
	R       []*result
	bar     *pb.ProgressBar
	wg      *sync.WaitGroup
	Red     = color.New(color.FgRed).SprintFunc()
	Green   = color.New(color.FgGreen).SprintFunc()
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Blue    = color.New(color.FgBlue).SprintFunc()
	Purple  = color.New(color.FgMagenta).SprintFunc()
	SkyBlue = color.New(color.FgCyan).SprintFunc()
	White   = color.New(color.FgWhite).SprintFunc()
)

type result struct {
	Name    string
	Divider bool
	Value   m.Result
}

func execute(Name string, F func(client http.Client) m.Result, client http.Client) {
	r := &result{Name: Name}
	R = append(R, r)
	wg.Add(1)
	tot++
	go func() {
		if Conc > 0 {
			sem <- struct{}{}
			defer func() {
				<-sem
				wg.Done()
			}()
		} else {
			defer wg.Done()
		}
		r.Value = F(client)
		bar.Describe(Name + " " + ShowResult(r.Value))
		bar.Add(1)
	}()
}

func ShowResult(r m.Result) (s string) {
	switch r.Status {
	case m.StatusOK:
		s = Green("YES")
		if r.Region != "" {
			s += Green(" (Region: " + strings.ToUpper(r.Region) + ")")
		}
		return s

	case m.StatusNetworkErr:
		if Debug {
			return Red("ERR") + Yellow(" (Network Err: "+r.Err.Error()+")")
		}
		return Red("ERR") + Yellow(" (Network Err)")

	case m.StatusRestricted:
		if r.Info != "" {
			return Yellow("Restricted (" + r.Info + ")")
		}
		return Yellow("Restricted")

	case m.StatusErr:
		s = Red("ERR")
		if r.Err != nil && Debug {
			s += Yellow(" (Err: " + r.Err.Error() + ")")
		}
		return s

	case m.StatusNo:
		if r.Info != "" {
			return Red("NO ") + Yellow(" (Info: "+r.Info+")")
		}
		if r.Region != "" {
			return Red("NO  (Region: " + strings.ToUpper(r.Region) + ")")
		}
		return Red("NO")

	case m.StatusBanned:
		if r.Info != "" {
			return Red("Banned") + Yellow(" ("+r.Info+")")
		}
		return Red("Banned")

	case m.StatusUnexpected:
		return Purple("Unexpected")

	case m.StatusFailed:
		return Blue("Failed")

	default:
		return
	}
}

func ShowR() {
	fmt.Println("测试时间: ", Yellow(time.Now().Format("2006-01-02 15:04:05")))
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
				result = SkyBlue(r.Value.Region)
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

func Globe(c http.Client, ipType int) {
	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}
	R = append(R, &result{Name: fmt.Sprintf("%s Globe", ipTypeStr), Divider: true})
	executeTests(MultinationTests, c, ipType)

}

func HongKong(c http.Client, ipType int) {
	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}
	R = append(R, &result{Name: fmt.Sprintf("%s Hong Kong", ipTypeStr), Divider: true})
	executeTests(HongKongTests, c, ipType)

}

func Taiwan(c http.Client, ipType int) {
	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}
	R = append(R, &result{Name: fmt.Sprintf("%s Taiwan", ipTypeStr), Divider: true})
	executeTests(TaiwanTests, c, ipType)

}

func Japan(c http.Client, ipType int) {
	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}
	R = append(R, &result{Name: fmt.Sprintf("%s Japan", ipTypeStr), Divider: true})
	executeTests(JapanTests, c, ipType)
}

func Korea(c http.Client, ipType int) {
	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}
	R = append(R, &result{Name: fmt.Sprintf("%s Korea", ipTypeStr), Divider: true})
	if ipType == 6 {
		R = append(R, &result{Name: "No Korean platform supports IPv6", Divider: false})
	}
	executeTests(KoreaTests, c, ipType)
}

func NorthAmerica(c http.Client, ipType int) {
	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}
	R = append(R, &result{Name: fmt.Sprintf("%s North America", ipTypeStr), Divider: true})
	executeTests(NorthAmericaTests, c, ipType)
}

func SouthAmerica(c http.Client, ipType int) {
	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}
	R = append(R, &result{Name: fmt.Sprintf("%s South America", ipTypeStr), Divider: true})
	executeTests(SouthAmericaTests, c, ipType)
}

func Europe(c http.Client, ipType int) {
	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}
	R = append(R, &result{Name: fmt.Sprintf("%s Europe", ipTypeStr), Divider: true})
	executeTests(EuropeTests, c, ipType)
}

func Africa(c http.Client, ipType int) {
	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}
	R = append(R, &result{Name: fmt.Sprintf("%s Africa", ipTypeStr), Divider: true})
	executeTests(AfricaTests, c, ipType)
}

func SouthEastAsia(c http.Client, ipType int) {
	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}
	R = append(R, &result{Name: fmt.Sprintf("%s South East Asia", ipTypeStr), Divider: true})
	executeTests(SouthEastAsiaTests, c, ipType)
}

func Oceania(c http.Client, ipType int) {
	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}
	R = append(R, &result{Name: fmt.Sprintf("%s Oceania", ipTypeStr), Divider: true})
	executeTests(OceaniaTests, c, ipType)
}

type TestItem struct {
	Name       string
	Func       func(client http.Client) m.Result
	SupportsV6 bool
}

var MultinationTests = []TestItem{
	{"Dazn", m.Dazn, false},
	{"Disney+", m.DisneyPlus, true},
	{"Netflix", m.NetflixRegion, true},
	{"Netflix CDN", m.NetflixCDN, true},
	{"Youtube Premium", m.YoutubeRegion, true},
	{"Youtube CDN", m.YoutubeCDN, true},
	{"Amazon Prime Video", m.PrimeVideo, false},
	{"TVBAnywhere+", m.TVBAnywhere, false},
	{"iQiYi", m.IQiYi, false},
	{"Viu.com", m.ViuCom, false},
	{"Spotify", m.Spotify, true},
	{"Steam", m.Steam, false},
	{"ChatGPT", m.ChatGPT, true},
	{"Wikipedia", m.WikipediaEditable, true},
	{"Reddit", m.Reddit, false},
	{"TikTok", m.TikTok, false},
	{"Bing", m.Bing, true},
	{"Instagram Audio", m.Instagram, true},
	{"Google Gemini", m.Gemini, true},
	{"Google Play Store", m.GooglePlayStore, true},
	{"Sora", m.Sora, true},
	{"Claude", m.Claude, true},
}

var HongKongTests = []TestItem{
	{"Now E", m.NowE, false},
	{"Viu.TV", m.ViuTV, false},
	{"MyTVSuper", m.MyTvSuper, false},
	{"Max", m.HBOMax, true},
	{"Bilibili HongKong/Macau Only", m.BilibiliHKMO, false},
	{"SonyLiv", m.SonyLiv, false},
	{"Bahamut Anime", m.BahamutAnime, false},
	{"Hoy TV", m.HoyTV, true},
	{"NBA TV", m.NBA_TV, true},
}

var TaiwanTests = []TestItem{
	{"KKTV", m.KKTV, false},
	{"LiTV", m.LiTV, false},
	{"MyVideo", m.MyVideo, false},
	{"4GTV", m.TW4GTV, false},
	{"LineTV", m.LineTV, false},
	{"Hami Video", m.HamiVideo, false},
	{"CatchPlay+", m.Catchplay, false},
	{"Bahamut Anime", m.BahamutAnime, false},
	{"Max", m.HBOMax, true},
	{"Bilibili Taiwan Only", m.BilibiliTW, false},
	{"Ofiii", m.Ofiii, false},
	{"Friday Video", m.FridayVideo, false},
}

var JapanTests = []TestItem{
	{"DMM", m.DMM, false},
	{"DMM TV", m.DMMTV, true},
	{"Abema", m.Abema, false},
	{"Niconico", m.Niconico, false},
	{"music.jp", m.MusicJP, false},
	{"Telasa", m.Telasa, true},
	{"U-NEXT", m.U_NEXT, true},
	{"Hulu Japan", m.HuluJP, false},
	{"GYAO!", m.GYAO, false},
	{"VideoMarket", m.VideoMarket, false},
	{"FOD(Fuji TV)", m.FOD, false},
	{"Radiko", m.Radiko, false},
	{"Karaoke@DAM", m.Karaoke, false},
	{"J:COM On Demand", m.J_COM_ON_DEMAND, false},
	{"Kancolle", m.Kancolle, false},
	{"Pretty Derby Japan", m.PrettyDerbyJP, true},
	{"Konosuba Fantastic Days", m.KonosubaFD, false},
	{"Princess Connect Re:Dive Japan", m.PCRJP, false},
	{"Project Sekai: Colorful Stage", m.PJSK, false},
	{"Rakuten TV JP", m.RakutenTV_JP, false},
	{"Wowow", m.Wowow, false},
	{"Watcha", m.Watcha, false},
	{"TVer", m.TVer, false},
	{"Lemino", m.Lemino, true},
	{"D Anime Store", m.DAnimeStore, false},
	{"Mora", m.Mora, false},
	{"AnimeFesta", m.AnimeFesta, false},
	{"EroGameSpace", m.EroGameSpace, false},
	{"NHK+", m.NHKPlus, true},
	{"Rakuten Magazine", m.RakutenMagazine, false},
	{"MGStage", m.MGStage, false},
}

var KoreaTests = []TestItem{
	{"Wavve", m.Wavve, false},
	{"Tving", m.Tving, false},
	{"Watcha", m.Watcha, false},
	{"Coupang Play", m.CoupangPlay, false},
	{"Spotv Now", m.SpotvNow, false},
	{"Naver TV", m.NaverTV, false},
	{"Afreeca", m.Afreeca, false},
	{"KBS", m.KBS, false},
	{"Panda TV", m.PandaTV, false},
}

var NorthAmericaTests = []TestItem{
	{"Shudder", m.Shudder, true},
	{"BritBox", m.BritBox, true},
	{"SonyLiv", m.SonyLiv, true},
	{"Hotstar", m.Hotstar, true},
	{"NBA TV", m.NBA_TV, true},
	{"Fubo TV", m.FuboTV, false},
	{"Tubi TV", m.TubiTV, true},
	{"Meta AI", m.MetaAI, true},
	{"AMC+", m.AMCPlus, true},
	{"Viaplay", m.Viaplay, false},
	{"A&E TV", m.AETV, false},
	{"FOX", m.Fox, true},
	{"Hulu", m.Hulu, true},
	{"NFL+", m.NFLPlus, false},
	{"ESPN+", m.ESPNPlus, true},
	{"MGM+", m.MGMPlus, false},
	{"Starz", m.Starz, false},
	{"Philo", m.Philo, false},
	{"FXNOW", m.FXNOW, false},
	{"TLC GO", m.TlcGo, true},
	{"Max", m.HBOMax, true},
	{"NBC TV", m.NBC_TV, true},
	{"CW TV", m.CW_TV, true},
	{"Sling TV", m.SlingTV, true},
	{"Pluto TV", m.PlutoTV, false},
	{"Acorn TV", m.AcornTV, false},
	{"SHOWTIME", m.SHOWTIME, true},
	{"encoreTVB", m.EncoreTVB, false},
	{"Discovery+", m.DiscoveryPlus, false},
	{"Paramount+", m.ParamountPlus, true},
	{"Peacock TV", m.PeacockTV, true},
	{"Crunchyroll", m.Crunchyroll, false},
	{"DirecTV Stream", m.DirectvStream, true},
	{"KOCOWA+", m.KOCOWA, false},
	{"Crackle", m.Crackle, true},
	{"MathsSpot Roblox", m.MathsSpotRoblox, false},
	{"CBC Gem", m.CBCGem, false},
	{"Crave", m.Crave, false},
}

var SouthAmericaTests = []TestItem{
	{"DirecTV GO", m.DirecTVGO, false},
	{"Max", m.HBOMax, true},
}

var EuropeTests = []TestItem{
	{"Rakuten TV EU", m.RakutenTV_EU, false},
	{"Setanta Sports", m.SetantaSports, true},
	{"Sky Show Time", m.SkyShowTime, true},
	{"Max", m.HBOMax, true},
	{"SonyLiv", m.SonyLiv, true},
	{"KOCOWA+", m.KOCOWA, false},
	{"Viaplay", m.Viaplay, true},
	{"BBC iPlayer", m.BBCiPlayer, false},
	{"Channel 4", m.Channel4, false},
	{"Channel 5", m.Channel5, false},
	{"Sky Go", m.SkyGo, false},
	{"ITVX", m.ITVX, false},
	{"Hotstar", m.Hotstar, true},
	{"MathsSpot Roblox", m.MathsSpotRoblox, false},
	{"Rai Play", m.RaiPlay, false},
	{"Canal+", m.CanalPlus, false},
	{"ZDF", m.ZDF, false},
	{"Joyn", m.Joyn, false},
	{"Molotov", m.Molotov, true},
	{"Sky DE", m.Sky_DE, false},
	{"France TV", m.FranceTV, true},
	{"NPO Start Plus", m.NPOStartPlus, false},
	{"Video Land", m.VideoLand, true},
	{"NLZIET", m.NLZIET, false},
	{"Movistar Plus+", m.MoviStarPlus, false},
	{"Eurosport RO", m.EurosportRO, false},
	{"Sky CH", m.Sky_CH, false},
	{"Amediateka", m.Amediateka, false},
}

var AfricaTests = []TestItem{
	{"DSTV", m.DSTV, false},
	{"Showmax", m.Showmax, true},
	{"Meta AI", m.MetaAI, true},
}

var SouthEastAsiaTests = []TestItem{
	{"Bilibili SouthEastAsia Only", m.BilibiliSEA, false},
	{"SonyLiv", m.SonyLiv, true},
	{"Hotstar", m.Hotstar, true},
	{"CatchPlay+", m.Catchplay, false},
	{"MeWatch", m.MeWatch, false},
	{"Meta AI", m.MetaAI, true},
	{"Bilibili Thailand Only", m.BilibiliTH, false},
	{"AIS Play", m.AISPlay, false},
	{"TrueID", m.TrueID, false},
	{"Bilibili Indonesia Only", m.BilibiliID, false},
	{"Bilibili Vietnam Only", m.BilibiliVN, false},
}

var OceaniaTests = []TestItem{
	{"NBA TV", m.NBA_TV, true},
	{"Acorn TV", m.AcornTV, false},
	{"BritBox", m.BritBox, true},
	{"Paramount+", m.ParamountPlus, true},
	{"SonyLiv", m.SonyLiv, true},
	{"Meta AI", m.MetaAI, true},
	{"KOCOWA+", m.KOCOWA, false},
	{"AMC+", m.AMCPlus, true},
	{"Stan", m.Stan, false},
	{"Binge", m.Binge, true},
	{"Doc Play", m.DocPlay, false},
	{"7Plus", m.SevenPlus, true},
	{"Channel 9", m.Channel9, true},
	{"10 Play", m.Channel10, false},
	{"ABC iView", m.ABCiView, false},
	{"Optus Sports", m.OptusSports, true},
	{"SBS on Demand", m.SBSonDemand, false},
	{"Kayo Sports", m.KayoSports, false},
	{"Neon TV", m.NeonTV, false},
	{"Three Now", m.ThreeNow, false},
	{"Maori TV", m.MaoriTV, false},
	{"Sky Go NZ", m.SkyGo_NZ, false},
}

func executeTests(tests []TestItem, client http.Client, ipType int) {
	for _, test := range tests {
		if ipType == 6 && !test.SupportsV6 {
			continue
		}
		execute(test.Name, test.Func, client)
	}
	R = append(R, &result{Name: "", Divider: false})
}

func GetIPInfo(url string, ipType int, isCloudflare bool) (string, error) {
	timeout := 6
	if ipType == 6 {
		timeout = 3
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var client http.Client
	if ipType == 6 {
		client = m.Ipv6HttpClient
	} else if ipType == 4 {
		client = m.Ipv4HttpClient
	} else if ipType == 0 {
		client = m.AutoHttpClient
	} else {
		return "", fmt.Errorf("IP type %d is invalid", ipType)
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("dnt", "1")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("sec-ch-ua", `"Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if isCloudflare {
		s := string(b)
		i := strings.Index(s, "ip=")
		s = s[i+3:]
		i = strings.Index(s, "\n")
		return s[:i], nil
	} else {
		return strings.TrimSpace(string(b)), nil
	}
}

func ReadSelect() {
	fmt.Println("请选择检测项目,直接按回车将进行全部检测: ")
	fmt.Println("[0] :   跨国平台")
	fmt.Println("[1] :   台湾平台")
	fmt.Println("[2] :   香港平台")
	fmt.Println("[3] :   日本平台")
	fmt.Println("[4] :   韩国平台")
	fmt.Println("[5] :   北美平台")
	fmt.Println("[6] :   南美平台")
	fmt.Println("[7] :   欧洲平台")
	fmt.Println("[8] :   非洲平台")
	fmt.Println("[9] : 东南亚平台")
	fmt.Println("[10]: 大洋洲平台")
	fmt.Print("请输入对应数字，空格分隔（回车确认）: ")
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
			SEA = true
		case "10":
			OCEA = true
		default:
			M, TW, HK, JP, KR, NA, SA, EU, AFR, SEA, OCEA = true, true, true, true, true, true, true, true, true, true, true
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
		log.Println("[ERR] 获取版本信息时出错:", err)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERR] 读取版本信息时出错:", err)
		return
	}

	parts := strings.Split(string(b), "-")
	if len(parts) != 2 {
		log.Println("[ERR] 版本号格式错误:", err)
		return
	}
	version := parts[0]

	if version == m.Version {
		fmt.Println("已经是最新版本")
		return
	}

	timestampInt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Println("[ERR] 版本号时间戳错误:", err)
		return
	}
	timestamp := time.Unix(timestampInt, 0)

	fmt.Println("最新版本：", version)
	fmt.Println("发布时间：", timestamp.Format("2006-01-02 15:04:05"))

	OS, ARCH := runtime.GOOS, runtime.GOARCH
	fmt.Println("运行系统：", OS)
	fmt.Println("运行架构：", ARCH)

	if OS == "android" && strings.Contains(os.Getenv("PREFIX"), "com.termux") {
		target_path := os.Getenv("PREFIX") + "/bin"
		out, err := os.Create(target_path + "/unlock-test_new")
		if err != nil {
			log.Fatal("[ERR] 创建文件出错:", err)
			return
		}
		defer out.Close()
		log.Println("下载unlock-test中 ...")
		url := "https://unlock.icmp.ing/test/latest/unlock-test_" + OS + "_" + ARCH
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
		if err := os.Chmod(target_path+"/unlock-test_new", 0777); err != nil {
			log.Fatal("[ERR] 更改unlock-test后端权限出错:", err)
		}
		if _, err := os.Stat(target_path + "/unlock-test"); err == nil {
			if err := os.Remove(target_path + "/unlock-test"); err != nil {
				log.Fatal("[ERR] 删除unlock-test旧版本时出错:", err.Error())
			}
		}
		if err := os.Rename(target_path+"/unlock-test_new", target_path+"/unlock-test"); err != nil {
			log.Fatal("[ERR] 更新unlock-test后端时出错:", err)
		}
	} else {
		url := "https://unlock.icmp.ing/test/latest/unlock-test_" + OS + "_" + ARCH
		if OS == "windows" {
			url += ".exe"
		}

		resp, err = http.Get(url)
		if err != nil {
			log.Fatal("[ERR] 下载unlock-test时出错:", err)
			return
		}
		defer resp.Body.Close()

		bar := pb.DefaultBytes(
			resp.ContentLength,
			"下载进度",
		)

		body := io.TeeReader(resp.Body, bar)

		if resp.StatusCode != http.StatusOK {
			log.Fatal("[ERR] 下载unlock-test时出错: 非预期的状态码", resp.StatusCode)
			return
		}

		err = selfUpdate.Apply(body, selfUpdate.Options{})
		if err != nil {
			log.Fatal("[ERR] 更新unlock-test时出错:", err)
			return
		}
	}

	fmt.Println("[OK] unlock-test后端更新成功")
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
	fmt.Printf("当天运行共%s次, 本月运行共%s次, 共计运行%s次\n", SkyBlue(d), Yellow(m), Green(t))
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
	Iface := ""
	DnsServers := ""
	httpProxy := ""
	socksProxy := ""
	showVersion := false
	update := false
	nf := false
	test := false
	mode := 0
	flag.StringVar(&Iface, "I", "", "source ip / interface")
	flag.StringVar(&DnsServers, "dns-servers", "", "specify dns servers")
	flag.StringVar(&httpProxy, "http-proxy", "", "http proxy")
	flag.StringVar(&socksProxy, "socks-proxy", "", "socks5 proxy")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&update, "u", false, "update")
	flag.BoolVar(&nf, "nf", false, "only test netflix")
	flag.BoolVar(&test, "test", false, "test mode")
	flag.BoolVar(&Debug, "debug", false, "debug mode")
	flag.IntVar(&mode, "m", 0, "mode 0(default)/4/6")
	flag.Uint64Var(&Conc, "conc", 0, "concurrency of tests")
	flag.Parse()
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
		m.Ipv4Transport.Resolver = m.Dialer.Resolver
		m.Ipv6Transport.Resolver = m.Dialer.Resolver
		m.AutoHttpClient.Transport.(*m.CustomTransport).Resolver = m.Dialer.Resolver
	}
	if httpProxy != "" {
		if u, err := url.Parse(httpProxy); err == nil {
			m.ClientProxy = http.ProxyURL(u)
			m.Ipv4Transport.Proxy = m.ClientProxy
			m.Ipv6Transport.Proxy = m.ClientProxy
			m.AutoHttpClient.Transport.(*m.CustomTransport).Proxy = m.ClientProxy
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

		m.Ipv4Transport.SocksDialer = dialer
		m.Ipv6Transport.SocksDialer = dialer
		m.AutoHttpClient.Transport.(*m.CustomTransport).SocksDialer = dialer
	}
	if Conc > 0 {
		sem = make(chan struct{}, Conc)
	}

	if nf {
		fmt.Println("Netflix", ShowResult(m.NetflixRegion(m.AutoHttpClient)))
		return
	}

	if test {
		//GetIpv4Info()
		//GetIpv6Info()
		fmt.Println("watcha", ShowResult(m.Watcha(m.AutoHttpClient)))
		//fmt.Println("DSTV", ShowResult(m.DSTV(m.AutoHttpClient)))
		return
	}

	fmt.Println("项目地址: " + SkyBlue("https://github.com/HsukqiLee/MediaUnlockTest"))
	fmt.Println("使用方式: " + Yellow("bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)"))
	fmt.Println()

	fmt.Println("正在从 ipw.cn 获取 IP...")
	var IP4_1, IP6_1, IP4_2, IP6_2 string
	var err error
	if mode == 0 || mode == 4 {
		IP4_1, err = GetIPInfo("https://4.ipw.cn", mode, false)
		if err != nil {
			if Debug {
				fmt.Println(Red("No IPv4 address (") + Yellow(err.Error()) + Red(")"))
			} else {
				fmt.Println(Red("No IPv4 address"))
			}
		} else {
			fmt.Println(SkyBlue("IPv4 address: ") + Green(IP4_1))
		}
	}
	if mode == 0 || mode == 6 {
		IP6_1, err = GetIPInfo("https://6.ipw.cn", mode, false)
		if err != nil {
			if Debug {
				fmt.Println(Red("No IPv6 address (") + Yellow(err.Error()) + Red(")"))
			} else {
				fmt.Println(Red("No IPv6 address"))
			}
		} else {
			fmt.Println(SkyBlue("IPv6 address: ") + Green(IP6_1))
		}
	}

	fmt.Println("正在从 ip.sb 获取 IP...")
	if mode == 0 || mode == 4 {
		IP4_2, err = GetIPInfo("https://api-ipv4.ip.sb/ip", mode, false)
		if err != nil {
			if Debug {
				fmt.Println(Red("No IPv4 address (") + Yellow(err.Error()) + Red(")"))
			} else {
				fmt.Println(Red("No IPv4 address"))
			}
		} else {
			fmt.Println(SkyBlue("IPv4 address: ") + Green(IP4_2))
		}
	}
	if mode == 0 || mode == 6 {
		IP6_2, err = GetIPInfo("https://api-ipv6.ip.sb/ip", mode, false)
		if err != nil {
			if Debug {
				fmt.Println(Red("No IPv6 address (") + Yellow(err.Error()) + Red(")"))
			} else {
				fmt.Println(Red("No IPv6 address"))
			}
		} else {
			fmt.Println(SkyBlue("IPv6 address: ") + Green(IP6_2))
		}
	}

	fmt.Println("正在检测系统代理...")
	isProxy := false
	if mode == 0 || mode == 4 {
		IP4, err := GetIPInfo("https://www.cloudflare.com/cdn-cgi/trace", 4, true)
		if err != nil {
			if IP4_1 != "" || IP4_2 != "" {
				isProxy = true
				fmt.Println(Yellow("正在使用系统代理，且无法通过 IPv4 连接代理"))
			} else {
				IPV4 = false
				fmt.Println(Red("未使用 IPv4 代理，无 IPv4 网络"))
			}
		} else {
			IPV4 = true
			if IP4_1 != IP4_2 || IP4_1 != IP4 {
				isProxy = true
				fmt.Println(Yellow("正在使用监听地址为 IPv4 的代理，出口 IP：") + Red(IP4))
			} else if IP4 == IP4_1 {
				fmt.Println(Green("未使用 IPv4 代理，有 IPv4 网络"))
			} else {
				fmt.Println(Red("无法强制使用 IPv4 网络测试，可能使用 IPv4 代理"))
				IPV4 = false
				if mode == 4 {
					IPV6 = false
				}
			}
		}
	}
	if mode == 0 || mode == 6 {
		IP6, err := GetIPInfo("https://www.cloudflare.com/cdn-cgi/trace", 6, true)
		if err != nil {
			if IP6_1 != "" && IP6_2 != "" {
				isProxy = true
				fmt.Println(Yellow("正在使用系统代理，且无法通过 IPv6 连接代理"))
			} else {
				IPV6 = false
				fmt.Println(Red("未使用 IPv6 代理，无 IPv6 网络"))
			}
		} else {
			IPV6 = true
			if IP6_1 != IP6_2 && IP6_1 != IP6 {
				isProxy = true
				fmt.Println(Yellow("正在使用监听地址为 IPv6 的代理，出口 IP：") + Red(IP6))
			} else if IP6 == IP6_1 {
				fmt.Println(Green("未使用 IPv6 代理，有 IPv6 网络"))
			} else {
				fmt.Println(Red("无法强制使用 IPv6 网络测试，可能使用 IPv6 代理"))
				IPV6 = false
				if mode == 6 {
					IPV4 = false
				}
			}
		}
	}

	if isProxy {
		fmt.Println(Red("提示：") + Yellow("正在使用系统代理，此时连接行为全部受代理控制"))
	}
	if mode != 0 {
		if mode == 4 {
			IPV6 = false
		} else if mode == 6 {
			IPV4 = false
		}
	}
	fmt.Println()

	if IPV4 || IPV6 {
		ReadSelect()
	}
	regions := []struct {
		enabled bool
		name    string
		fn      func(http.Client, int)
	}{
		{M, "Globe", Globe},
		{TW, "Taiwan", Taiwan},
		{HK, "HongKong", HongKong},
		{JP, "Japan", Japan},
		{KR, "Korea", Korea},
		{NA, "NorthAmerica", NorthAmerica},
		{SA, "SouthAmerica", SouthAmerica},
		{EU, "Europe", Europe},
		{AFR, "Africa", Africa},
		{SEA, "SouthEastAsia", SouthEastAsia},
		{OCEA, "Oceania", Oceania},
	}
	wg = &sync.WaitGroup{}
	bar = NewBar(0)
	if isProxy {
		for _, region := range regions {
			if region.enabled {
				region.fn(m.AutoHttpClient, 0)
			}
		}
	} else {
		if IPV4 {
			for _, region := range regions {
				if region.enabled {
					region.fn(m.Ipv4HttpClient, 4)
				}
			}
		}
		if IPV6 {
			for _, region := range regions {
				if region.enabled {
					region.fn(m.Ipv6HttpClient, 6)
				}
			}
		}
	}

	bar.ChangeMax64(tot)

	wg.Wait()
	bar.Finish()
	fmt.Println()
	ShowR()
	fmt.Println()
	fmt.Println("检测完毕，感谢您的使用！")
	showCounts()
	fmt.Println()
	showAd()
}
