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
	Force   bool
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

func excute(Name string, F func(client http.Client) m.Result, client http.Client) {
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

func Multination(c http.Client) {
	R = append(R, &result{Name: "Multination", Divider: true})
	excute("Dazn", m.Dazn, c)
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
	excute("Google Gemini", m.Gemini, c)
	excute("Google Play Store", m.GooglePlayStore, c)
	excute("Sora", m.Sora, c)
}

func HongKong(c http.Client) {
	R = append(R, &result{Name: "Hong Kong", Divider: true})
	excute("Now E", m.NowE, c)
	excute("Viu.TV", m.ViuTV, c)
	excute("MyTVSuper", m.MyTvSuper, c)
	excute("HBO GO Asia", m.HboGoAsia, c)
	excute("Bilibili HongKong/Macau Only", m.BilibiliHKMO, c)
	excute("SonyLiv", m.SonyLiv, c)
	excute("Bahamut Anime", m.BahamutAnime, c)
	excute("Hoy TV", m.HoyTV, c)
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
	excute("HBO GO Asia", m.HboGoAsia, c)
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
	excute("Project Sekai: Colorful Stage", m.PJSK, c)
	excute("Rakuten TV JP", m.RakutenTV_JP, c)
	excute("Wowow", m.Wowow, c)
	excute("Watcha", m.Watcha, c)
	excute("TVer", m.TVer, c)
	excute("Lemino", m.Lemino, c)
	excute("D Anime Store", m.DAnimeStore, c)
	excute("Mora", m.Mora, c)
	excute("AnimeFesta", m.AnimeFesta, c)
	excute("EroGameSpace", m.EroGameSpace, c)
	excute("NHK+", m.NHKPlus, c)
	excute("Rakuten Magazine", m.RakutenMagazine, c)
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
	excute("Shudder", m.Shudder, c)
	excute("BritBox", m.BritBox, c)
	excute("SonyLiv", m.SonyLiv, c)
	excute("Hotstar", m.Hotstar, c)
	excute("NBA TV", m.NBA_TV, c)
	excute("Fubo TV", m.FuboTV, c)
	excute("Tubi TV", m.TubiTV, c)
	excute("Meta AI", m.MetaAI, c)
	excute("AMC+", m.AMCPlus, c)
	excute("Viaplay", m.Viaplay, c)
	R = append(R, &result{Name: "US", Divider: true})
	excute("FOX", m.Fox, c)
	excute("Hulu", m.Hulu, c)
	excute("NFL+", m.NFLPlus, c)
	excute("ESPN+", m.ESPNPlus, c)
	excute("MGM+", m.MGMPlus, c)
	excute("Starz", m.Starz, c)
	excute("Philo", m.Philo, c)
	excute("FXNOW", m.FXNOW, c)
	excute("TLC GO", m.TlcGo, c)
	excute("HBO Max", m.HBOMax, c)
	excute("CW TV", m.CW_TV, c)
	excute("Sling TV", m.SlingTV, c)
	excute("Pluto TV", m.PlutoTV, c)
	excute("Acorn TV", m.AcornTV, c)
	excute("SHOWTIME", m.SHOWTIME, c)
	excute("encoreTVB", m.EncoreTVB, c)
	excute("Discovery+", m.DiscoveryPlus, c)
	excute("Paramount+", m.ParamountPlus, c)
	excute("Peacock TV", m.PeacockTV, c)
	excute("Crunchyroll", m.Crunchyroll, c)
	excute("DirecTV Stream", m.DirectvStream, c)
	excute("KOCOWA+", m.KOCOWA, c)
	excute("Crackle", m.Crackle, c)
	excute("MathsSpot Roblox", m.MathsSpotRoblox, c)
	R = append(R, &result{Name: "CA", Divider: true})
	excute("CBC Gem", m.CBCGem, c)
	excute("Crave", m.Crave, c)
}

func SouthAmerica(c http.Client) {
	R = append(R, &result{Name: "South America", Divider: true})
	//excute("Star Plus", m.StarPlus, c)
	excute("DirecTV GO", m.DirecTVGO, c)
	excute("HBO Max", m.HBOMax, c)
}

func Europe(c http.Client) {
	R = append(R, &result{Name: "Europe", Divider: true})
	excute("Rakuten TV EU", m.RakutenTV_EU, c)
	excute("Setanta Sports", m.SetantaSports, c)
	excute("Sky Show Time", m.SkyShowTime, c)
	excute("HBO Max", m.HBOMax, c)
	excute("SonyLiv", m.SonyLiv, c)
	excute("KOCOWA+", m.KOCOWA, c)
	excute("Viaplay", m.Viaplay, c)
	R = append(R, &result{Name: "GB", Divider: true})
	excute("BBC iPlayer", m.BBCiPlayer, c)
	excute("Channel 4", m.Channel4, c)
	excute("Channel 5", m.Channel5, c)
	excute("Sky Go", m.SkyGo, c)
	excute("ITVX", m.ITVX, c)
	excute("Hotstar", m.Hotstar, c)
	excute("MathsSpot Roblox", m.MathsSpotRoblox, c)
	R = append(R, &result{Name: "IT", Divider: true})
	excute("Rai Play", m.RaiPlay, c)
	R = append(R, &result{Name: "FR/DE", Divider: true})
	excute("Canal+", m.CanalPlus, c)
	excute("ZDF", m.ZDF, c)
	excute("Joyn", m.Joyn, c)
	excute("Molotov", m.Molotov, c)
	excute("Sky DE", m.Sky_DE, c)
	excute("France TV", m.FranceTV, c)
	R = append(R, &result{Name: "NL", Divider: true})
	excute("NPO Start Plus", m.NPOStartPlus, c)
	excute("Video Land", m.VideoLand, c)
	excute("NLZIET", m.NLZIET, c)
	R = append(R, &result{Name: "ES", Divider: true})
	excute("Movistar Plus+", m.MoviStarPlus, c)
	R = append(R, &result{Name: "RO", Divider: true})
	excute("Eurosport RO", m.EurosportRO, c)
	R = append(R, &result{Name: "CH", Divider: true})
	excute("Sky CH", m.Sky_CH, c)
	R = append(R, &result{Name: "RU", Divider: true})
	excute("Amediateka", m.Amediateka, c)
}

func Africa(c http.Client) {
	R = append(R, &result{Name: "Africa", Divider: true})
	excute("DSTV", m.DSTV, c)
	excute("Showmax", m.Showmax, c)
	excute("Meta AI", m.MetaAI, c)
}

func SouthEastAsia(c http.Client) {
	R = append(R, &result{Name: "South East Asia", Divider: true})
	excute("Bilibili SouthEastAsia Only", m.BilibiliSEA, c)
	excute("SonyLiv", m.SonyLiv, c)
	excute("Hotstar", m.Hotstar, c)
	excute("CatchPlay+", m.Catchplay, c)
	R = append(R, &result{Name: "SG", Divider: true})
	excute("MeWatch", m.MeWatch, c)
	excute("Meta AI", m.MetaAI, c)
	R = append(R, &result{Name: "TH", Divider: true})
	excute("Bilibili Thailand Only", m.BilibiliTH, c)
	excute("AIS Play", m.AISPlay, c)
	excute("TrueID", m.TrueID, c)
	R = append(R, &result{Name: "ID", Divider: true})
	excute("Bilibili Indonesia Only", m.BilibiliID, c)
	R = append(R, &result{Name: "VN", Divider: true})
	excute("Bilibili Vietnam Only", m.BilibiliVN, c)
}

func Oceania(c http.Client) {
	R = append(R, &result{Name: "Oceania", Divider: true})
	excute("NBA TV", m.NBA_TV, c)
	excute("Acorn TV", m.AcornTV, c)
	excute("BritBox", m.BritBox, c)
	excute("Paramount+", m.ParamountPlus, c)
	excute("SonyLiv", m.SonyLiv, c)
	excute("Meta AI", m.MetaAI, c)
	excute("KOCOWA+", m.KOCOWA, c)
	excute("AMC+", m.AMCPlus, c)
	R = append(R, &result{Name: "AU", Divider: true})
	excute("Stan", m.Stan, c)
	excute("Binge", m.Binge, c)
	excute("Doc Play", m.DocPlay, c)
	excute("7Plus", m.SevenPlus, c)
	excute("Channel 9", m.Channel9, c)
	excute("10 Play", m.Channel10, c)
	excute("ABC iView", m.ABCiView, c)
	excute("Optus Sports", m.OptusSports, c)
	excute("SBS on Demand", m.SBSonDemand, c)
	excute("Kayo Sports", m.KayoSports, c)
	R = append(R, &result{Name: "NZ", Divider: true})
	excute("Neon TV", m.NeonTV, c)
	excute("Three Now", m.ThreeNow, c)
	excute("Maori TV", m.MaoriTV, c)
	excute("Sky Go NZ", m.SkyGo_NZ, c)
}

func Ipv6Multination() {
	c := m.Ipv6HttpClient
	R = append(R, &result{Name: "", Divider: true})
	R = append(R, &result{Name: "IPV6 Multination", Divider: true})
	excute("Hotstar", m.Hotstar, c)
	excute("Disney+", m.DisneyPlus, c)
	excute("Netflix", m.NetflixRegion, c)
	excute("Netflix CDN", m.NetflixCDN, c)
	excute("Youtube Premium", m.YoutubeRegion, c)
	excute("Youtube CDN", m.YoutubeCDN, c)
	excute("ChatGPT", m.ChatGPT, c)
	excute("Wikipedia", m.WikipediaEditable, c)
	excute("Bing", m.Bing, c)
	excute("Google Gemini", m.Gemini, c)
	excute("Google Play Store", m.GooglePlayStore, c)
	excute("Sora", m.Sora, c)
}

func GetIpv4Info() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://www.cloudflare.com/cdn-cgi/trace", nil)
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
	fmt.Println("Your IPV4 address:", SkyBlue(s[:i]))
}
func GetIpv6Info() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://www.cloudflare.com/cdn-cgi/trace", nil)
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
	fmt.Println("Your IPV6 address:", SkyBlue(s[:i]))
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
	client := m.AutoHttpClient
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
	flag.BoolVar(&Force, "f", false, "force ipv6")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&update, "u", false, "update")
	flag.BoolVar(&nf, "nf", false, "netflix")
	flag.BoolVar(&test, "test", false, "test")
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
	if mode == 4 {
		client = m.Ipv4HttpClient
		IPV6 = false
	}
	if mode == 6 {
		client = m.Ipv6HttpClient
		IPV4 = false
		M = true
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
		fmt.Println("sora", ShowResult(m.Sora(m.AutoHttpClient)))
		//fmt.Println("DSTV", ShowResult(m.DSTV(m.AutoHttpClient)))
		return
	}

	fmt.Println("项目地址: " + SkyBlue("https://github.com/HsukqiLee/MediaUnlockTest"))
	fmt.Println("使用方式: " + Yellow("bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)"))
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
		if SEA {
			SouthEastAsia(client)
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
			if SEA {
				SouthEastAsia(m.Ipv6HttpClient)
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
	fmt.Println("检测完毕，感谢您的使用！")
	showCounts()
	fmt.Println()
	showAd()
}
