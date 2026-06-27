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

// isoAlpha2To3 maps ISO 3166-1 alpha-2 codes to alpha-3 codes.
// ThreeToTwoCode performs the reverse lookup by iterating over this map.
var isoAlpha2To3 = map[string]string{
	"AD": "AND", "AE": "ARE", "AF": "AFG", "AG": "ATG", "AI": "AIA", "AL": "ALB", "AM": "ARM", "AO": "AGO", "AQ": "ATA", "AR": "ARG",
	"AS": "ASM", "AT": "AUT", "AU": "AUS", "AW": "ABW", "AX": "ALA", "AZ": "AZE", "BA": "BIH", "BB": "BRB", "BD": "BGD", "BE": "BEL",
	"BF": "BFA", "BG": "BGR", "BH": "BHR", "BI": "BDI", "BJ": "BEN", "BL": "BLM", "BM": "BMU", "BN": "BRN", "BO": "BOL", "BQ": "BES",
	"BR": "BRA", "BS": "BHS", "BT": "BTN", "BV": "BVT", "BW": "BWA", "BY": "BLR", "BZ": "BLZ", "CA": "CAN", "CC": "CCK", "CD": "COD",
	"CF": "CAF", "CG": "COG", "CH": "CHE", "CI": "CIV", "CK": "COK", "CL": "CHL", "CM": "CMR", "CN": "CHN", "CO": "COL", "CR": "CRI",
	"CU": "CUB", "CV": "CPV", "CW": "CUW", "CX": "CXR", "CY": "CYP", "CZ": "CZE", "DE": "DEU", "DJ": "DJI", "DK": "DNK", "DM": "DMA",
	"DO": "DOM", "DZ": "DZA", "EC": "ECU", "EE": "EST", "EG": "EGY", "EH": "ESH", "ER": "ERI", "ES": "ESP", "ET": "ETH", "FI": "FIN",
	"FJ": "FJI", "FK": "FLK", "FM": "FSM", "FO": "FRO", "FR": "FRA", "GA": "GAB", "GB": "GBR", "GD": "GRD", "GE": "GEO", "GF": "GUF",
	"GG": "GGY", "GH": "GHA", "GI": "GIB", "GL": "GRL", "GM": "GMB", "GN": "GIN", "GP": "GLP", "GQ": "GNQ", "GR": "GRC", "GS": "SGS",
	"GT": "GTM", "GU": "GUM", "GW": "GNB", "GY": "GUY", "HK": "HKG", "HM": "HMD", "HN": "HND", "HR": "HRV", "HT": "HTI", "HU": "HUN",
	"ID": "IDN", "IE": "IRL", "IL": "ISR", "IM": "IMN", "IN": "IND", "IO": "IOT", "IQ": "IRQ", "IR": "IRN", "IS": "ISL", "IT": "ITA",
	"JE": "JEY", "JM": "JAM", "JO": "JOR", "JP": "JPN", "KE": "KEN", "KG": "KGZ", "KH": "KHM", "KI": "KIR", "KM": "COM", "KN": "KNA",
	"KP": "PRK", "KR": "KOR", "KW": "KWT", "KY": "CYM", "KZ": "KAZ", "LA": "LAO", "LB": "LBN", "LC": "LCA", "LI": "LIE", "LK": "LKA",
	"LR": "LBR", "LS": "LSO", "LT": "LTU", "LU": "LUX", "LV": "LVA", "LY": "LBY", "MA": "MAR", "MC": "MCO", "MD": "MDA", "ME": "MNE",
	"MF": "MAF", "MG": "MDG", "MH": "MHL", "MK": "MKD", "ML": "MLI", "MM": "MMR", "MN": "MNG", "MO": "MAC", "MP": "MNP", "MQ": "MTQ",
	"MR": "MRT", "MS": "MSR", "MT": "MLT", "MU": "MUS", "MV": "MDV", "MW": "MWI", "MX": "MEX", "MY": "MYS", "MZ": "MOZ", "NA": "NAM",
	"NC": "NCL", "NE": "NER", "NF": "NFK", "NG": "NGA", "NI": "NIC", "NL": "NLD", "NO": "NOR", "NP": "NPL", "NR": "NRU", "NU": "NIU",
	"NZ": "NZL", "OM": "OMN", "PA": "PAN", "PE": "PER", "PF": "PYF", "PG": "PNG", "PH": "PHL", "PK": "PAK", "PL": "POL", "PM": "SPM",
	"PN": "PCN", "PR": "PRI", "PS": "PSE", "PT": "PRT", "PW": "PLW", "PY": "PRY", "QA": "QAT", "RE": "REU", "RO": "ROU", "RS": "SRB",
	"RU": "RUS", "RW": "RWA", "SA": "SAU", "SB": "SLB", "SC": "SYC", "SD": "SDN", "SE": "SWE", "SG": "SGP", "SH": "SHN", "SI": "SVN",
	"SJ": "SJM", "SK": "SVK", "SL": "SLE", "SM": "SMR", "SN": "SEN", "SO": "SOM", "SR": "SUR", "SS": "SSD", "ST": "STP", "SV": "SLV",
	"SX": "SXM", "SY": "SYR", "SZ": "SWZ", "TC": "TCA", "TD": "TCD", "TF": "ATF", "TG": "TGO", "TH": "THA", "TJ": "TJK", "TK": "TKL",
	"TL": "TLS", "TM": "TKM", "TN": "TUN", "TO": "TON", "TR": "TUR", "TT": "TTO", "TV": "TUV", "TW": "TWN", "TZ": "TZA", "UA": "UKR",
	"UG": "UGA", "UM": "UMI", "US": "USA", "UY": "URY", "UZ": "UZB", "VA": "VAT", "VC": "VCT", "VE": "VEN", "VG": "VGB", "VI": "VIR",
	"VN": "VNM", "VU": "VUT", "WF": "WLF", "WS": "WSM", "YE": "YEM", "YT": "MYT", "ZA": "ZAF", "ZM": "ZMB", "ZW": "ZWE",
}

func TwoToThreeCode(code string) string {
	return isoAlpha2To3[strings.ToUpper(code)]
}

func ThreeToTwoCode(code string) string {
	upper := strings.ToUpper(code)
	for k, v := range isoAlpha2To3 {
		if v == upper {
			return k
		}
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
