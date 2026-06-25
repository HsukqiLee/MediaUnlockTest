package providers

import (
	"net/http"

	core "MediaUnlockTest/pkg/core"
)

type TestItem struct {
	Name       string
	Func       func(client http.Client) core.Result
	SupportsV6 bool
}

var GlobeTests = []TestItem{
	{"Amazon Prime Video", PrimeVideo, false},
	{"Apple", Apple, true},
	{"Bing", Bing, true},
	{"Dazn", Dazn, false},
	{"Disney+", DisneyPlus, true},
	{"Google Play Store", GooglePlayStore, true},
	{"iQiYi", IQiYi, false},
	{"Netflix", NetflixRegion, true},
	{"Netflix CDN", NetflixCDN, true},
	{"Reddit", Reddit, false},
	{"Spotify", Spotify, true},
	{"Steam", Steam, false},
	{"TikTok", TikTok, false},
	{"TVBAnywhere+", TVBAnywhere, false},
	{"Viu.com", ViuCom, false},
	{"Wikipedia", WikipediaEditable, true},
	{"Youtube CDN", YoutubeCDN, true},
	{"Youtube Premium", YoutubeRegion, true},
}

var HongKongTests = []TestItem{
	{"Bahamut Anime", BahamutAnime, false},
	{"Bilibili HongKong/Macau Only", BilibiliHKMO, false},
	{"Hoy TV", HoyTV, true},
	{"Max", Max, true},
	{"MyTVSuper", MyTvSuper, false},
	{"NBA TV", NBA_TV, true},
	//{"Now TV", NowTV, false},
	{"SonyLiv", SonyLiv, false},
	{"Viu.TV", ViuTV, false},
}

var TaiwanTests = []TestItem{
	{"4GTV", TW4GTV, false},
	{"Bahamut Anime", BahamutAnime, false},
	{"Bilibili Taiwan Only", BilibiliTW, false},
	{"CatchPlay+", Catchplay, false},
	{"Friday Video", FridayVideo, false},
	{"Hami Video", HamiVideo, false},
	{"KKTV", KKTV, false},
	{"LiTV", LiTV, false},
	{"LineTV", LineTV, false},
	{"Max", Max, true},
	{"MyVideo", MyVideo, false},
	{"Ofiii", Ofiii, false},
}

var JapanTests = []TestItem{
	{"Abema", Abema, false},
	{"AnimeFesta", AnimeFesta, false},
	{"D Anime Store", DAnimeStore, false},
	{"DMM", DMM, false},
	{"DMM TV", DMMTV, true},
	{"EroGameSpace", EroGameSpace, false},
	{"FOD(Fuji TV)", FOD, false},
	{"Hulu Japan", HuluJP, false},
	{"J:COM On Demand", J_COM_ON_DEMAND, false},
	{"Kancolle", Kancolle, false},
	{"Karaoke@DAM", Karaoke, false},
	{"Lemino", Lemino, true},
	{"MGStage", MGStage, false},
	{"Mora", Mora, false},
	{"Music.jp", MusicJP, false},
	{"NHK+", NHKPlus, true},
	{"Niconico", Niconico, false},
	{"Pretty Derby Japan", PrettyDerbyJP, true},
	{"Princess Connect Re:Dive Japan", PCRJP, false},
	{"Project Sekai: Colorful Stage", PJSK, false},
	{"Radiko", Radiko, false},
	{"Rakuten Magazine", RakutenMagazine, false},
	{"Rakuten TV JP", RakutenTV_JP, false},
	{"Telasa", Telasa, true},
	{"TVer", TVer, false},
	{"U-NEXT", U_NEXT, true},
	{"VideoMarket", VideoMarket, false},
	{"Watcha", Watcha, false},
	{"Wowow", Wowow, false},
}

var KoreaTests = []TestItem{
	{"Afreeca", Afreeca, false},
	{"Coupang Play", CoupangPlay, false},
	{"KBS", KBS, false},
	{"Naver TV", NaverTV, false},
	{"Panda TV", PandaTV, false},
	{"Spotv Now", SpotvNow, false},
	{"Tving", Tving, false},
	{"Watcha", Watcha, false},
	{"Wavve", Wavve, false},
}

var NorthAmericaTests = []TestItem{
	{"A&E TV", AETV, false},
	{"Acorn TV", AcornTV, false},
	{"AMC+", AMCPlus, true},
	{"BritBox", BritBox, true},
	{"CBC Gem", CBCGem, false},
	{"Crave", Crave, false},
	{"Crunchyroll", Crunchyroll, false},
	{"CW TV", CW_TV, true},
	{"DirecTV Stream", DirectvStream, true},
	{"Discovery+", DiscoveryPlus, false},
	{"encoreTVB", EncoreTVB, false},
	{"ESPN+", ESPNPlus, true},
	{"FOX", Fox, true},
	{"Fubo TV", FuboTV, false},
	{"FXNOW", FXNOW, false},
	{"Hotstar", Hotstar, true},
	{"Hulu", Hulu, true},
	{"KOCOWA+", KOCOWA, false},
	{"MGM+", MGMPlus, false},
	{"MathsSpot Roblox", MathsSpotRoblox, false},
	{"Max", Max, true},
	{"NBC TV", NBC_TV, true},
	{"NFL+", NFLPlus, false},
	{"NBA TV", NBA_TV, true},
	{"Paramount+", ParamountPlus, true},
	{"Peacock TV", PeacockTV, true},
	{"Philo", Philo, false},
	{"Pluto TV", PlutoTV, false},
	{"SHOWTIME", SHOWTIME, true},
	{"Shudder", Shudder, true},
	{"Sling TV", SlingTV, true},
	{"SonyLiv", SonyLiv, true},
	{"Starz", Starz, false},
	{"TLC GO", TlcGo, true},
	{"Tubi TV", TubiTV, true},
	{"Viaplay", Viaplay, false},
}

var SouthAmericaTests = []TestItem{
	{"DirecTV GO", DirecTVGO, false},
	{"Max", Max, true},
}

var EuropeTests = []TestItem{
	{"Rakuten TV EU", RakutenTV_EU, false},
	{"Sky Show Time", SkyShowTime, true},
	{"Viaplay", Viaplay, true},
	{"TNTSports", TNTSports, false},
	{"Eurosport RO", EurosportRO, false},
	{"Setanta Sports", SetantaSports, true},
	{"KOCOWA+", KOCOWA, false},
	{"MathsSpot Roblox", MathsSpotRoblox, false},
	{"Max", Max, true},
	{"SonyLiv", SonyLiv, true},
	{"GB", nil, true},
	{"BBC iPlayer", BBCiPlayer, false},
	{"BritBox", BritBox, true},
	{"ITVX", ITVX, false},
	{"Channel 4", Channel4, false},
	{"Channel 5", Channel5, false},
	{"Discovery+ UK", DiscoveryPlus_UK, false},
	{"Sky Go", SkyGo, false},
	{"FR", nil, true},
	{"Canal+", CanalPlus, false},
	{"Molotov", Molotov, true},
	{"France TV", FranceTV, true},
	{"DE", nil, false},
	{"Joyn", Joyn, false},
	{"Sky DE", Sky_DE, false},
	{"ZDF", ZDF, false},
	{"NL", nil, true},
	{"NLZIET", NLZIET, false},
	{"Video Land", VideoLand, true},
	{"NPO Start Plus", NPOStartPlus, false},
	{"ES", nil, false},
	{"Movistar Plus+", MoviStarPlus, false},
	{"IT", nil, false},
	{"Rai Play", RaiPlay, false},
	{"CH", nil, false},
	{"Sky CH", Sky_CH, false},
	{"RU", nil, false},
	{"Amediateka", Amediateka, false},
}

var AfricaTests = []TestItem{
	{"DSTV", DSTV, false},
	{"Showmax", Showmax, true},
}

var SouthEastAsiaTests = []TestItem{
	{"Max", Max, true},
	{"Hotstar", Hotstar, true},
	{"NBA TV", NBA_TV, true},
	{"Bilibili SouthEastAsia Only", BilibiliSEA, false},
	{"SG", nil, false},
	{"MeWatch", MeWatch, false},
	{"CatchPlay+", Catchplay, false},
	{"TH", nil, false},
	{"AIS Play", AISPlay, false},
	{"TrueID", TrueID, false},
	{"Bilibili Thailand Only", BilibiliTH, false},
	{"ID", nil, false},
	{"Bilibili Indonesia Only", BilibiliID, false},
	{"VN", nil, false},
	{"Clip TV", ClipTV, false},
	{"Galaxy Play", GalaxyPlay, false},
	{"K+", KPlus, false},
	{"Bilibili Vietnam Only", BilibiliVN, false},
	{"MY", nil, false},
	{"Sooka", Sooka, false},
	{"IN", nil, true},
	{"Tata Play", TataPlay, true},
	{"SonyLiv", SonyLiv, true},
	{"MX Player", MXPlayer, false},
	{"Zee5", Zee5, true},
}

var OceaniaTests = []TestItem{
	{"10 Play", Channel10, false},
	{"7Plus", SevenPlus, true},
	{"ABC iView", ABCiView, false},
	{"Acorn TV", AcornTV, false},
	{"AMC+", AMCPlus, true},
	{"BritBox", BritBox, true},
	{"Channel 9", Channel9, true},
	{"Doc Play", DocPlay, false},
	{"Kayo Sports", KayoSports, false},
	{"KOCOWA+", KOCOWA, false},
	{"Maori TV", MaoriTV, false},
	{"NBA TV", NBA_TV, true},
	{"Neon TV", NeonTV, false},
	{"Optus Sports", OptusSports, true},
	{"Paramount+", ParamountPlus, true},
	{"SBS on Demand", SBSonDemand, false},
	{"Sky Go NZ", SkyGo_NZ, false},
	{"SonyLiv", SonyLiv, true},
	{"Stan", Stan, false},
	{"Three Now", ThreeNow, false},
}

var AITests = []TestItem{
	{"ChatGPT", ChatGPT, true},
	{"Claude", Claude, true},
	{"Copilot", Copilot, true},
	{"Google Gemini", Gemini, true},
	{"Meta AI", MetaAI, true},
	{"Sora", Sora, true},
}

