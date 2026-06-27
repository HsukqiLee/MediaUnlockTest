package core

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

var (
	Version          = "1.9.2"
	StatusOK         = 1
	StatusNetworkErr = -1
	StatusErr        = -2
	StatusRestricted = 2
	StatusNo         = 3
	StatusBanned     = 4
	StatusFailed     = 5
	StatusUnexpected = 6
)

type HttpClient = tls_client.HttpClient

type Result struct {
	Status       int
	Region       string
	Info         string
	Err          error
	CachedResult bool
}

var (
	UA_Browser = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36 Edg/137.0.0.0"
	UA_Dalvik  = "Dalvik/2.1.0 (Linux; U; Android 11; M2006J10C Build/RP1A.200720.011)"

	ClientSessionHeaders = &SessionHeaders{
		UserAgent:      "",
		SecChUA:        "",
		AcceptLanguage: "",
		DNT:            "0",
	}
)

var (
	Ipv4HttpClient HttpClient
	Ipv6HttpClient HttpClient
	AutoHttpClient HttpClient
	SocksProxy     string
	HTTPProxy      string
	DNSServers     string
	Dialer         = &net.Dialer{
		Timeout:   15 * time.Second,
		KeepAlive: 30 * time.Second,
	}
)

func buildClientOptions(disableIPv4, disableIPv6 bool) []tls_client.HttpClientOption {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(6),
		tls_client.WithClientProfile(profiles.Chrome_120),
		tls_client.WithCustomRedirectFunc(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}),
	}
	if disableIPv4 {
		options = append(options, tls_client.WithDisableIPV4())
	}
	if disableIPv6 {
		options = append(options, tls_client.WithDisableIPV6())
	}
	if SocksProxy != "" {
		options = append(options, tls_client.WithProxyUrl(SocksProxy))
	} else if HTTPProxy != "" {
		options = append(options, tls_client.WithProxyUrl(HTTPProxy))
	}
	if DNSServers != "" {
		Dialer.Resolver = &net.Resolver{
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "udp", DNSServers)
			},
		}
	}
	options = append(options, tls_client.WithDialer(*Dialer))
	return options
}

func InitClients() {
	var err error
	Ipv4HttpClient, err = tls_client.NewHttpClient(tls_client.NewNoopLogger(), buildClientOptions(false, true)...)
	if err != nil {
		panic(err)
	}

	Ipv6HttpClient, err = tls_client.NewHttpClient(tls_client.NewNoopLogger(), buildClientOptions(true, false)...)
	if err != nil {
		panic(err)
	}

	AutoHttpClient, err = tls_client.NewHttpClient(tls_client.NewNoopLogger(), buildClientOptions(false, false)...)
	if err != nil {
		panic(err)
	}
}

type H [2]string

func GET(c HttpClient, url string, headers ...H) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	hasCustomHeaders := false
	for i := 0; i < len(headers); i++ {
		if headers[i][0] == "x-custom-headers" && headers[i][1] == "true" {
			hasCustomHeaders = true
			headers = append(headers[:i], headers[i+1:]...)
			break
		}
	}
	if !hasCustomHeaders {
		setRealisticHeaders(req, "html")
	}

	for _, h := range headers {
		req.Header.Set(h[0], h[1])
	}
	addRandomDelay()
	return Cdo(c, req)
}

func GET_Dalvik(c HttpClient, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UA_Dalvik)
	return Cdo(c, req)
}

var ErrNetwork = errors.New("network error")

func Cdo(c HttpClient, req *http.Request) (resp *http.Response, err error) {
	deadline := time.Now().Add(14 * time.Second)
	for i := 0; i < 2; i++ {
		if time.Now().After(deadline) {
			break
		}
		if resp, err = c.Do(req); err == nil {
			return resp, nil
		}
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) && dnsErr.IsNotFound {
			break
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			break
		}
	}
	return nil, err
}
func PostJson(c HttpClient, url string, data string, headers ...H) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	setRealisticHeaders(req, "json")
	for _, h := range headers {
		req.Header.Set(h[0], h[1])
	}
	addRandomDelay()
	return Cdo(c, req)
}

func PostForm(c HttpClient, url string, data string, headers ...H) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	setRealisticHeaders(req, "html")

	for _, h := range headers {
		req.Header.Set(h[0], h[1])
	}
	addRandomDelay()
	return Cdo(c, req)
}

func GenUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func GenBase36(n int) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[secureRandInt(len(letters))]
	}
	return string(b)
}

func MD5Sum(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

type Country struct {
	Alpha2      string
	Alpha3      string
	CountryCode string // Calling code
}

var (
	alpha2Map = make(map[string]*Country)
	alpha3Map = make(map[string]*Country)
	codeMap   = make(map[string]*Country)
)

var countryList = []Country{
	{"AD", "AND", "376"},
	{"AE", "ARE", "971"},
	{"AF", "AFG", "93"},
	{"AG", "ATG", "1268"},
	{"AI", "AIA", "1264"},
	{"AL", "ALB", "355"},
	{"AM", "ARM", "374"},
	{"AO", "AGO", "244"},
	{"AQ", "ATA", "672"},
	{"AR", "ARG", "54"},
	{"AS", "ASM", "1684"},
	{"AT", "AUT", "43"},
	{"AU", "AUS", "61"},
	{"AW", "ABW", "297"},
	{"AX", "ALA", "35818"},
	{"AZ", "AZE", "994"},
	{"BA", "BIH", "387"},
	{"BB", "BRB", "1246"},
	{"BD", "BGD", "880"},
	{"BE", "BEL", "32"},
	{"BF", "BFA", "226"},
	{"BG", "BGR", "359"},
	{"BH", "BHR", "973"},
	{"BI", "BDI", "257"},
	{"BJ", "BEN", "229"},
	{"BL", "BLM", "590"},
	{"BM", "BMU", "1441"},
	{"BN", "BRN", "673"},
	{"BO", "BOL", "591"},
	{"BQ", "BES", "599"},
	{"BR", "BRA", "55"},
	{"BS", "BHS", "1242"},
	{"BT", "BTN", "975"},
	{"BV", "BVT", "47"},
	{"BW", "BWA", "267"},
	{"BY", "BLR", "375"},
	{"BZ", "BLZ", "501"},
	{"CA", "CAN", "1"},
	{"CC", "CCK", "61"},
	{"CD", "COD", "243"},
	{"CF", "CAF", "236"},
	{"CG", "COG", "242"},
	{"CH", "CHE", "41"},
	{"CI", "CIV", "225"},
	{"CK", "COK", "682"},
	{"CL", "CHL", "56"},
	{"CM", "CMR", "237"},
	{"CN", "CHN", "86"},
	{"CO", "COL", "57"},
	{"CR", "CRI", "506"},
	{"CU", "CUB", "53"},
	{"CV", "CPV", "238"},
	{"CW", "CUW", "599"},
	{"CX", "CXR", "61"},
	{"CY", "CYP", "357"},
	{"CZ", "CZE", "420"},
	{"DE", "DEU", "49"},
	{"DJ", "DJI", "253"},
	{"DK", "DNK", "45"},
	{"DM", "DMA", "1767"},
	{"DO", "DOM", "1"},
	{"DZ", "DZA", "213"},
	{"EC", "ECU", "593"},
	{"EE", "EST", "372"},
	{"EG", "EGY", "20"},
	{"EH", "ESH", "212"},
	{"ER", "ERI", "291"},
	{"ES", "ESP", "34"},
	{"ET", "ETH", "251"},
	{"FI", "FIN", "358"},
	{"FJ", "FJI", "679"},
	{"FK", "FLK", "500"},
	{"FM", "FSM", "691"},
	{"FO", "FRO", "298"},
	{"FR", "FRA", "33"},
	{"GA", "GAB", "241"},
	{"GB", "GBR", "44"},
	{"GD", "GRD", "1473"},
	{"GE", "GEO", "995"},
	{"GF", "GUF", "594"},
	{"GG", "GGY", "44"},
	{"GH", "GHA", "233"},
	{"GI", "GIB", "350"},
	{"GL", "GRL", "299"},
	{"GM", "GMB", "220"},
	{"GN", "GIN", "224"},
	{"GP", "GLP", "590"},
	{"GQ", "GNQ", "240"},
	{"GR", "GRC", "30"},
	{"GS", "SGS", "500"},
	{"GT", "GTM", "502"},
	{"GU", "GUM", "1671"},
	{"GW", "GNB", "245"},
	{"GY", "GUY", "592"},
	{"HK", "HKG", "852"},
	{"HM", "HMD", "672"},
	{"HN", "HND", "504"},
	{"HR", "HRV", "385"},
	{"HT", "HTI", "509"},
	{"HU", "HUN", "36"},
	{"ID", "IDN", "62"},
	{"IE", "IRL", "353"},
	{"IL", "ISR", "972"},
	{"IM", "IMN", "44"},
	{"IN", "IND", "91"},
	{"IO", "IOT", "246"},
	{"IQ", "IRQ", "964"},
	{"IR", "IRN", "98"},
	{"IS", "ISL", "354"},
	{"IT", "ITA", "39"},
	{"JE", "JEY", "44"},
	{"JM", "JAM", "1876"},
	{"JO", "JOR", "962"},
	{"JP", "JPN", "81"},
	{"KE", "KEN", "254"},
	{"KG", "KGZ", "996"},
	{"KH", "KHM", "855"},
	{"KI", "KIR", "686"},
	{"KM", "COM", "269"},
	{"KN", "KNA", "1869"},
	{"KP", "PRK", "850"},
	{"KR", "KOR", "82"},
	{"KW", "KWT", "965"},
	{"KY", "CYM", "1345"},
	{"KZ", "KAZ", "7"},
	{"LA", "LAO", "856"},
	{"LB", "LBN", "961"},
	{"LC", "LCA", "1758"},
	{"LI", "LIE", "423"},
	{"LK", "LKA", "94"},
	{"LR", "LBR", "231"},
	{"LS", "LSO", "266"},
	{"LT", "LTU", "370"},
	{"LU", "LUX", "352"},
	{"LV", "LVA", "371"},
	{"LY", "LBY", "218"},
	{"MA", "MAR", "212"},
	{"MC", "MCO", "377"},
	{"MD", "MDA", "373"},
	{"ME", "MNE", "382"},
	{"MF", "MAF", "590"},
	{"MG", "MDG", "261"},
	{"MH", "MHL", "692"},
	{"MK", "MKD", "389"},
	{"ML", "MLI", "223"},
	{"MM", "MMR", "95"},
	{"MN", "MNG", "976"},
	{"MO", "MAC", "853"},
	{"MP", "MNP", "1670"},
	{"MQ", "MTQ", "596"},
	{"MR", "MRT", "222"},
	{"MS", "MSR", "1664"},
	{"MT", "MLT", "356"},
	{"MU", "MUS", "230"},
	{"MV", "MDV", "960"},
	{"MW", "MWI", "265"},
	{"MX", "MEX", "52"},
	{"MY", "MYS", "60"},
	{"MZ", "MOZ", "258"},
	{"NA", "NAM", "264"},
	{"NC", "NCL", "687"},
	{"NE", "NER", "227"},
	{"NF", "NFK", "672"},
	{"NG", "NGA", "234"},
	{"NI", "NIC", "505"},
	{"NL", "NLD", "31"},
	{"NO", "NOR", "47"},
	{"NP", "NPL", "977"},
	{"NR", "NRU", "674"},
	{"NU", "NIU", "683"},
	{"NZ", "NZL", "64"},
	{"OM", "OMN", "968"},
	{"PA", "PAN", "507"},
	{"PE", "PER", "51"},
	{"PF", "PYF", "689"},
	{"PG", "PNG", "675"},
	{"PH", "PHL", "63"},
	{"PK", "PAK", "92"},
	{"PL", "POL", "48"},
	{"PM", "SPM", "508"},
	{"PN", "PCN", "64"},
	{"PR", "PRI", "1"},
	{"PS", "PSE", "970"},
	{"PT", "PRT", "351"},
	{"PW", "PLW", "680"},
	{"PY", "PRY", "595"},
	{"QA", "QAT", "974"},
	{"RE", "REU", "262"},
	{"RO", "ROU", "40"},
	{"RS", "SRB", "381"},
	{"RU", "RUS", "7"},
	{"RW", "RWA", "250"},
	{"SA", "SAU", "966"},
	{"SB", "SLB", "677"},
	{"SC", "SYC", "248"},
	{"SD", "SDN", "249"},
	{"SE", "SWE", "46"},
	{"SG", "SGP", "65"},
	{"SH", "SHN", "290"},
	{"SI", "SVN", "386"},
	{"SJ", "SJM", "47"},
	{"SK", "SVK", "421"},
	{"SL", "SLE", "232"},
	{"SM", "SMR", "378"},
	{"SN", "SEN", "221"},
	{"SO", "SOM", "252"},
	{"SR", "SUR", "597"},
	{"SS", "SSD", "211"},
	{"ST", "STP", "239"},
	{"SV", "SLV", "503"},
	{"SX", "SXM", "1721"},
	{"SY", "SYR", "963"},
	{"SZ", "SWZ", "268"},
	{"TC", "TCA", "1649"},
	{"TD", "TCD", "235"},
	{"TF", "ATF", "262"},
	{"TG", "TGO", "228"},
	{"TH", "THA", "66"},
	{"TJ", "TJK", "992"},
	{"TK", "TKL", "690"},
	{"TL", "TLS", "670"},
	{"TM", "TKM", "993"},
	{"TN", "TUN", "216"},
	{"TO", "TON", "676"},
	{"TR", "TUR", "90"},
	{"TT", "TTO", "1868"},
	{"TV", "TUV", "688"},
	{"TW", "TWN", "886"},
	{"TZ", "TZA", "255"},
	{"UA", "UKR", "380"},
	{"UG", "UGA", "256"},
	{"UM", "UMI", "1"},
	{"US", "USA", "1"},
	{"UY", "URY", "598"},
	{"UZ", "UZB", "998"},
	{"VA", "VAT", "3906698"},
	{"VC", "VCT", "1784"},
	{"VE", "VEN", "58"},
	{"VG", "VGB", "1284"},
	{"VI", "VIR", "1340"},
	{"VN", "VNM", "84"},
	{"VU", "VUT", "678"},
	{"WF", "WLF", "681"},
	{"WS", "WSM", "685"},
	{"YE", "YEM", "967"},
	{"YT", "MYT", "262"},
	{"ZA", "ZAF", "27"},
	{"ZM", "ZMB", "260"},
	{"ZW", "ZWE", "263"},
}

func init() {
	for i := range countryList {
		c := &countryList[i]
		alpha2Map[c.Alpha2] = c
		alpha3Map[c.Alpha3] = c
		if c.CountryCode != "" {
			codeMap[c.CountryCode] = c
		}
	}
}

func TwoToThreeCode(code string) string {
	if c, ok := alpha2Map[strings.ToUpper(code)]; ok {
		return c.Alpha3
	}
	return ""
}

func ThreeToTwoCode(code string) string {
	if c, ok := alpha3Map[strings.ToUpper(code)]; ok {
		return c.Alpha2
	}
	return ""
}

func CountryCodeToAlpha2(code string) string {
	if c, ok := codeMap[strings.ToUpper(code)]; ok {
		return c.Alpha2
	}
	return ""
}

func CountryCodeToAlpha3(code string) string {
	if c, ok := codeMap[strings.ToUpper(code)]; ok {
		return c.Alpha3
	}
	return ""
}

func GenRandomStr(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[secureRandInt(len(charset))]
	}
	return string(b)
}

func generateEdgeUserAgent() string {
	edgeVersion := secureRandInRange(136, 140)
	chromiumVersion := edgeVersion

	return fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.0.0 Safari/537.36 Edg/%d.0.0.0",
		chromiumVersion, edgeVersion)
}

func generateSecChUA() string {
	edgeVersion := secureRandInRange(136, 140)
	chromiumVersion := edgeVersion
	notBrandVersion := secureRandInRange(20, 29)

	return fmt.Sprintf(`"Microsoft Edge";v="%d", "Chromium";v="%d", "Not/A)Brand";v="%d"`,
		edgeVersion, chromiumVersion, notBrandVersion)
}

func getRandomAcceptLanguage() string {
	languages := []string{
		"en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
		"en-US,en;q=0.9",
		"zh-CN,zh;q=0.9,en;q=0.8",
		"zh-CN,zh;q=0.9",
		"en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,ja;q=0.6",
	}
	return languages[secureRandInt(len(languages))]
}

func addRandomDelay() {
	if secureRandInt(10) == 0 {
		delay := time.Duration(secureRandInRange(50, 149)) * time.Millisecond // secureRandInt(100)+50
		time.Sleep(delay)
	}
}

func GetRealisticHeaders(requestType string) []H {
	headers := make([]H, 0)
	ua := generateEdgeUserAgent()
	secChUa := generateSecChUA()
	acceptLanguage := getRandomAcceptLanguage()
	dnt := strconv.Itoa(secureRandInt(2))
	headers = append(headers, H{"user-agent", ua})
	secFetchMode := "cors"
	secFetchDest := "empty"
	switch requestType {
	case "json":
		headers = append(headers, H{"accept", "application/json, text/plain, */*"})
	case "html":
		headers = append(headers, H{"accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"})
		secFetchMode = "navigate"
		secFetchDest = "document"
		headers = append(headers, H{"sec-fetch-user", "?1"})
		headers = append(headers, H{"upgrade-insecure-requests", "1"})
	default:
		headers = append(headers, H{"accept", "*/*"})
	}
	headers = append(headers, H{"sec-ch-ua", secChUa})
	headers = append(headers, H{"sec-ch-ua-mobile", "?0"})
	headers = append(headers, H{"sec-ch-ua-platform", `"Windows"`})
	headers = append(headers, H{"accept-language", acceptLanguage})
	headers = append(headers, H{"cache-control", "no-cache"})
	headers = append(headers, H{"pragma", "no-cache"})
	headers = append(headers, H{"sec-fetch-site", "cross-site"})
	headers = append(headers, H{"sec-fetch-mode", secFetchMode})
	headers = append(headers, H{"sec-fetch-dest", secFetchDest})
	headers = append(headers, H{"dnt", dnt})
	return headers
}

func setRealisticHeaders(req *http.Request, requestType string) {
	// Generate fresh headers for each request (default behavior)
	headers := GetRealisticHeaders(requestType)
	for _, header := range headers {
		req.Header.Set(header[0], header[1])
	}
}

type SessionHeaders struct {
	UserAgent      string
	SecChUA        string
	AcceptLanguage string
	DNT            string
}

func NewSessionHeaders() *SessionHeaders {
	return &SessionHeaders{
		UserAgent:      generateEdgeUserAgent(),
		SecChUA:        generateSecChUA(),
		AcceptLanguage: getRandomAcceptLanguage(),
		DNT:            strconv.Itoa(secureRandInt(2)),
	}
}

func (s *SessionHeaders) Headers() []H {
	return []H{
		{"user-agent", s.UserAgent},
		{"sec-ch-ua", s.SecChUA},
		{"accept-language", s.AcceptLanguage},
		{"dnt", s.DNT},
	}
}

func secureRandInt(max int) int {
	if max <= 0 {
		return 0
	}
	maxBytes := make([]byte, 4)
	_, err := rand.Read(maxBytes)
	if err != nil {
		return 0
	}
	randUint32 := uint32(maxBytes[0])<<24 | uint32(maxBytes[1])<<16 | uint32(maxBytes[2])<<8 | uint32(maxBytes[3])
	return int(randUint32) % max
}

func secureRandInRange(min, max int) int {
	if min >= max {
		return min
	}
	return secureRandInt(max-min+1) + min
}

// ResultMap 支持完整 Result 值的映射
type ResultMap map[int]Result

// ResultFromMapping 根据 statusCode 从映射获得 Result，缺省返回 defaultRes[0] 或 StatusUnexpected
func ResultFromMapping(statusCode int, m ResultMap, defaultResult Result) Result {
	if r, ok := m[statusCode]; ok {
		return r
	}
	return defaultResult
}

// CheckStatus 使用 GET 请求，并通过 ResultMap 返回对应 Result，支持默认 Result 及可选 headers
func CheckGETStatus(c HttpClient, url string, mapping ResultMap, defaultResult Result, headers ...H) Result {
	resp, err := GET(c, url, headers...)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	return ResultFromMapping(resp.StatusCode, mapping, defaultResult)
}

// CheckDalvikStatus 使用 GET_Dalvik 请求，并通过 ResultMap 返回对应 Result，支持默认 Result
func CheckDalvikStatus(c HttpClient, url string, mapping ResultMap, defaultResult Result) Result {
	resp, err := GET_Dalvik(c, url)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	return ResultFromMapping(resp.StatusCode, mapping, defaultResult)
}

func PostFormBoolSuccess(c HttpClient, url, data string, headers ...H) (bool, error) {
	resp, err := PostForm(c, url, data, headers...)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var res struct{ Success bool }
	if err := json.Unmarshal(b, &res); err != nil {
		return false, err
	}
	return res.Success, nil
}

// CheckPostFormStatus 使用 POST 表单请求，并通过 ResultMap 返回对应 Result，支持默认 Result 及可选 headers
func CheckPostFormStatus(c HttpClient, url, data string, mapping ResultMap, defaultResult Result, headers ...H) Result {
	resp, err := PostForm(c, url, data, headers...)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	return ResultFromMapping(resp.StatusCode, mapping, defaultResult)
}

// CheckPostJsonStatus 使用 POST JSON 请求，并通过 ResultMap 返回对应 Result，支持默认 Result 及可选 headers
func CheckPostJsonStatus(c HttpClient, url, data string, mapping ResultMap, defaultResult Result, headers ...H) Result {
	resp, err := PostJson(c, url, data, headers...)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	return ResultFromMapping(resp.StatusCode, mapping, defaultResult)
}

var sessionMutex sync.RWMutex

func NewHttpClient(ipType int) HttpClient {
	var disableIPv4, disableIPv6 bool
	switch ipType {
	case 4:
		disableIPv6 = true
	case 6:
		disableIPv4 = true
	}
	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), buildClientOptions(disableIPv4, disableIPv6)...)
	return client
}

func ResetSessionHeaders() {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	ClientSessionHeaders.UserAgent = ""
	ClientSessionHeaders.SecChUA = ""
	ClientSessionHeaders.AcceptLanguage = ""
	ClientSessionHeaders.DNT = "0"
}

func SetSessionHeaders(h *SessionHeaders) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	ClientSessionHeaders = h
}

// IsWAFBlockError checks if the network error is caused by a WAF drop/timeout
func IsWAFBlockError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, context.DeadlineExceeded) || os.IsTimeout(err) ||
		errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNABORTED) {
		return true
	}
	// utls 或 http2 等底层依赖有时只返回字符串错误，且没有导出 Error 变量
	errStr := err.Error()
	if strings.Contains(errStr, "stream error") || strings.Contains(errStr, "handshake failure") || strings.Contains(errStr, "connection reset") {
		return true
	}
	return false
}
