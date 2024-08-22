package main

import (
	mt "MediaUnlockTest"
	"net/http"
	"sync"
	"time"
)

var (
	MUL  bool
	HK   bool
	TW   bool
	JP   bool
	KR   bool
	NA   bool
	SA   bool
	EU   bool
	AFR  bool
	SEA  bool
	OCEA bool
)

type TEST struct {
	Client  http.Client
	Results []*result
	Wg      *sync.WaitGroup
}

func NewTest() *TEST {
	t := &TEST{
		Client:  mt.NewAutoHttpClient(),
		Results: make([]*result, 0),
		Wg:      &sync.WaitGroup{},
	}
	return t
}

func (T *TEST) Check() bool {
	if MUL {
		T.Multination()
	}
	if HK {
		T.HongKong()
	}
	if TW {
		T.Taiwan()
	}
	if JP {
		T.Japan()
	}
	if KR {
		T.Korea()
	}
	if NA {
		T.NorthAmerica()
	}
	if SA {
        T.SouthAmerica()
	}
	if EU {
		T.Europe()
	}
	if AFR {
		T.Africa()
	}
	if SEA {
		T.SouthEastAsia()
	}
	if OCEA {
		T.Oceania()
	}

	ch := make(chan struct{})
	go func() {
		defer close(ch)
		T.Wg.Wait()
	}()
	select {
	case <-ch:
		return false
	case <-time.After(30 * time.Second):
		return true
	}
}

type result struct {
	Type  string
	Name  string
	Value mt.Result
}

func (T *TEST) excute(Name string, F func(client http.Client) mt.Result) {
	r := &result{Name: Name}
	T.Results = append(T.Results, r)
	T.Wg.Add(1)
	go func() {
		res := F(T.Client)
		r.Value = res
		T.Wg.Done()
	}()
}

func (T *TEST) Multination() {
	// R = append(R, &result{Name: "Multination", Divider: true})
	T.excute("Dazn", mt.Dazn)
	T.excute("Disney+", mt.DisneyPlus)
	T.excute("Netflix", mt.NetflixRegion)
	T.excute("Netflix CDN", mt.NetflixCDN)
	T.excute("Youtube Premium", mt.YoutubeRegion)
	T.excute("Youtube CDN", mt.YoutubeCDN)
	T.excute("Amazon Prime Video", mt.PrimeVideo)
	T.excute("TVBAnywhere+", mt.TVBAnywhere)
	T.excute("iQiYi", mt.IQiYi)
	T.excute("Viu.com", mt.ViuCom)
	T.excute("Spotify", mt.Spotify)
	T.excute("Steam", mt.Steam)
	T.excute("ChatGPT", mt.ChatGPT)
	T.excute("Wikipedia", mt.WikipediaEditable)
	T.excute("Reddit", mt.Reddit)
	T.excute("TikTok", mt.TikTok)
	T.excute("Bing", mt.Bing)
	T.excute("Instagram Audio", mt.Instagram)
	T.excute("Google Gemini", mt.Gemini)
}

func (T *TEST) HongKong() {
	// R = append(R, &result{Name: "Hong Kong", Divider: true})
	T.excute("Now E", mt.NowE)
	T.excute("Viu.TV", mt.ViuTV)
	T.excute("MyTVSuper", mt.MyTvSuper)
	T.excute("HBO GO Asia", mt.HboGoAsia)
	T.excute("BiliBili HK/MO", mt.BilibiliHKMO)
	T.excute("SonyLiv", mt.SonyLiv)
	T.excute("Bahamut Anime", mt.BahamutAnime)
	T.excute("Hoy TV", mt.HoyTV)
}

func (T *TEST) Taiwan() {
	// R = append(R, &result{Name: "Taiwan", Divider: true})
	T.excute("KKTV", mt.KKTV)
	T.excute("LiTV", mt.LiTV)
	T.excute("MyVideo", mt.MyVideo)
	T.excute("4GTV", mt.TW4GTV)
	T.excute("LineTV", mt.LineTV)
	T.excute("Hami Video", mt.HamiVideo)
	T.excute("CatchPlay+", mt.Catchplay)
	T.excute("Bahamut Anime", mt.BahamutAnime)
	T.excute("HBO GO Asia", mt.HboGoAsia)
	T.excute("Bilibili TW", mt.BilibiliTW)
}

func (T *TEST) Japan() {
	// R = append(R, &result{Name: "Japan", Divider: true})
	T.excute("DMM", mt.DMM)
	T.excute("DMM TV", mt.DMMTV)
	T.excute("Abema", mt.Abema)
	T.excute("Niconico", mt.Niconico)
	T.excute("music.jp", mt.MusicJP)
	T.excute("Telasa", mt.Telasa)
	T.excute("Paravi", mt.Paravi)
	T.excute("U-NEXT", mt.U_NEXT)
	T.excute("Hulu Japan", mt.HuluJP)
	T.excute("GYAO!", mt.GYAO)
	T.excute("VideoMarket", mt.VideoMarket)
	T.excute("FOD(Fuji TV)", mt.FOD)
	T.excute("Radiko", mt.Radiko)
	T.excute("Karaoke@DAM", mt.Karaoke)
	T.excute("J:COM On Demand", mt.J_COM_ON_DEMAND)
	T.excute("Kancolle", mt.Kancolle)
	T.excute("Pretty Derby Japan", mt.PrettyDerbyJP)
	T.excute("Konosuba Fantastic Days", mt.KonosubaFD)
	T.excute("Princess Connect Re:Dive Japan", mt.PCRJP)
	T.excute("Project Sekai: Colorful Stage", mt.PJSK)
	T.excute("Rakuten TV", mt.RakutenTV_JP)
	T.excute("Wowow", mt.Wowow)
	T.excute("Watcha", mt.Watcha)
	T.excute("TVer", mt.TVer)
	T.excute("Lemino", mt.Lemino)
	T.excute("D Anime Store", mt.DAnimeStore)
	T.excute("Mora", mt.Mora)
	T.excute("AnimeFesta", mt.AnimeFesta)
	T.excute("EroGameSpace", mt.EroGameSpace)
	T.excute("NHK+", mt.NHKPlus)
	T.excute("Rakuten Magazine", mt.RakutenMagazine)
}

func (T *TEST) Korea() {
	// R = append(R, &result{Name: "Korea", Divider: true})
	T.excute("Wavve", mt.Wavve)
	T.excute("Tving", mt.Tving)
	T.excute("Watcha", mt.Watcha)
	T.excute("Coupang Play", mt.CoupangPlay)
	T.excute("SpotvNow", mt.SpotvNow)
	T.excute("NaverTV", mt.NaverTV)
	T.excute("Afreeca", mt.Afreeca)
	T.excute("KBS", mt.KBS)
}

func (T *TEST) NorthAmerica() {
	// R = append(R, &result{Name: "North America", Divider: true})
	T.excute("FOX", mt.Fox)
	T.excute("Hulu", mt.Hulu)
	T.excute("NFL+", mt.NFLPlus)
	T.excute("ESPN+", mt.ESPNPlus)
	T.excute("MGM+", mt.MGMPlus)
	T.excute("Starz", mt.Starz)
	T.excute("Philo", mt.Philo)
	T.excute("FXNOW", mt.FXNOW)
	T.excute("TLC GO", mt.TlcGo)
	T.excute("HBO Max", mt.HBOMax)
	T.excute("Shudder", mt.Shudder)
	T.excute("BritBox", mt.BritBox)
	T.excute("CW TV", mt.CW_TV)
	T.excute("NBA TV", mt.NBA_TV)
	T.excute("Fubo TV", mt.FuboTV)
	T.excute("Tubi TV", mt.TubiTV)
	T.excute("Sling TV", mt.SlingTV)
	T.excute("Pluto TV", mt.PlutoTV)
	T.excute("Acorn TV", mt.AcornTV)
	T.excute("SHOWTIME", mt.SHOWTIME)
	T.excute("encoreTVB", mt.EncoreTVB)
	T.excute("Discovery+", mt.DiscoveryPlus)
	T.excute("Paramount+", mt.ParamountPlus)
	T.excute("Peacock TV", mt.PeacockTV)
	T.excute("Crunchyroll", mt.Crunchyroll)
	T.excute("DirecTV Stream", mt.DirectvStream)
	T.excute("SonyLiv", mt.SonyLiv)
	T.excute("Hotstar", mt.Hotstar)
	T.excute("Meta AI", mt.MetaAI)
	T.excute("AMC+", mt.AMCPlus)
	T.excute("Crackle", mt.Crackle)
	T.excute("MathsSpot Roblox", mt.MathsSpotRoblox)
	T.excute("KOCOWA+", mt.KOCOWA)
	T.excute("Viaplay", mt.Viaplay)
	// R = append(R, &result{Name: "CA", Divider: true})
	T.excute("CBC Gem", mt.CBCGem)
	T.excute("Crave", mt.Crave)
}

func (T *TEST) SouthAmerica() {
    //R = append(R, &result{Name: "South America", Divider: true})
    //T.excute("Star Plus", mt.StarPlus)
    T.excute("DirecTV GO", mt.DirecTVGO)
}

func (T *TEST) Europe() {
    //R = append(R, &result{Name: "Europe", Divider: true})
    T.excute("Rakuten TV", mt.RakutenTV_EU)
    T.excute("Setanta Sports", mt.SetantaSports)
    T.excute("Sky Show Time", mt.SkyShowTime)
    T.excute("HBO Max", mt.HBOMax)
    T.excute("SonyLiv", mt.SonyLiv)
    T.excute("BBC iPlayer", mt.BBCiPlayer)
    T.excute("Channel 4", mt.Channel4)
    T.excute("Channel 5", mt.Channel5)
    T.excute("Sky Go", mt.SkyGo)
    T.excute("ITVX", mt.ITVX)
    T.excute("Rai Play", mt.RaiPlay)
    T.excute("Canal+", mt.CanalPlus)
    T.excute("ZDF", mt.ZDF)
    T.excute("Joyn", mt.Joyn)
    T.excute("Sky DE", mt.Sky_DE)
    T.excute("Molotov", mt.Molotov)
    T.excute("NPO Start Plus", mt.NPOStartPlus)
    T.excute("Video Land", mt.VideoLand)
    T.excute("NLZIET", mt.NLZIET)
    T.excute("Movistar Plus+", mt.MoviStarPlus)
    T.excute("Eurosport RO", mt.EurosportRO)
    T.excute("Sky CH", mt.Sky_CH)
    T.excute("Amediateka", mt.Amediateka)
    T.excute("Hotstar", mt.Hotstar)
    T.excute("MathsSpot Roblox", mt.MathsSpotRoblox)
    T.excute("KOCOWA+", mt.KOCOWA)
    T.excute("Meta AI", mt.MetaAI)
    T.excute("France TV", mt.FranceTV)
    T.excute("Viaplay", mt.Viaplay)
}

func (T *TEST) Africa() {
    //R = append(R, &result{Name: "Europe", Divider: true})
    T.excute("DSTV", mt.DSTV)
    T.excute("Showmax", mt.Showmax)
    T.excute("Meta AI", mt.MetaAI)
}

func (T *TEST) SouthEastAsia() {
    //R = append(R, &result{Name: "South East Asia", Divider: true})
    T.excute("Bilibili SouthEastAsia Only", mt.BilibiliSEA)
    T.excute("SonyLiv", mt.SonyLiv)
    T.excute("MeWatch", mt.MeWatch)
    T.excute("Bilibili Thailand Only", mt.BilibiliTW)
    T.excute("AIS Play", mt.AISPlay)
    T.excute("TrueID", mt.TrueID)
    T.excute("Bilibili Indonesia Only", mt.BilibiliID)
    T.excute("Bilibili Vietnam Only", mt.BilibiliVN)
    T.excute("Hotstar", mt.Hotstar)
    T.excute("Meta AI", mt.MetaAI)
    T.excute("CatchPlay+", mt.Catchplay)
}

func (T *TEST) Oceania() {
    //R = append(R, &result{Name: "Oceania", Divider: true})
    T.excute("NBA TV", mt.NBA_TV)
    T.excute("Acorn TV", mt.AcornTV)
    T.excute("BritBox", mt.BritBox)
    T.excute("Paramount+", mt.ParamountPlus)
    T.excute("SonyLiv", mt.SonyLiv)
    T.excute("Stan", mt.Stan)
    T.excute("Binge", mt.Binge)
    T.excute("Doc Play", mt.DocPlay)
    T.excute("7Plus", mt.SevenPlus)
    T.excute("Channel 9", mt.Channel9)
    T.excute("10 Play", mt.Channel10)
    T.excute("ABC iView", mt.ABCiView)
    T.excute("Optus Sports", mt.OptusSports)
    T.excute("SBS on Demand", mt.SBSonDemand)
    T.excute("Kayo Sports", mt.KayoSports)
    T.excute("Neon TV", mt.NeonTV)
    T.excute("Three Now", mt.ThreeNow)
    T.excute("Maori TV", mt.MaoriTV)
    T.excute("Sky Go NZ", mt.SkyGo_NZ)
    T.excute("AMC+", mt.AMCPlus)
    T.excute("KOCOWA+", mt.KOCOWA)
    T.excute("Meta AI", mt.MetaAI)
}