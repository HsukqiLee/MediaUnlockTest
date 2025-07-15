package mediaunlocktest

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	utls "github.com/refraction-networking/utls"

	"golang.org/x/net/http2"
	"golang.org/x/net/proxy"
)

var (
	Version          = "1.8.1"
	StatusOK         = 1
	StatusNetworkErr = -1
	StatusErr        = -2
	StatusRestricted = 2
	StatusNo         = 3
	StatusBanned     = 4
	StatusFailed     = 5
	StatusUnexpected = 6
)

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

var Dialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
	// Resolver:  &net.Resolver{},
}

var ClientProxy = http.ProxyFromEnvironment

func UseLastResponse(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }

type CustomTransport struct {
	Dialer      *net.Dialer
	Resolver    *net.Resolver
	Network     string
	Proxy       func(*http.Request) (*url.URL, error)
	Base        *http.Transport
	SocksDialer proxy.Dialer
}

func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		if t.SocksDialer != nil {
			return t.SocksDialer.Dial(t.Network, addr)
		}
		if t.Resolver != nil {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			ips, err := t.Resolver.LookupIP(ctx, "ip", host)
			if err != nil {
				return nil, err
			}
			var filteredIPs []net.IP
			for _, ip := range ips {
				if (t.Network == "tcp4" && ip.To4() != nil) || (t.Network == "tcp6" && ip.To4() == nil) || t.Network == "tcp" {
					filteredIPs = append(filteredIPs, ip)
				}
			}
			for _, ip := range filteredIPs {
				ipAddr := net.JoinHostPort(ip.String(), port)
				conn, err := t.Dialer.DialContext(ctx, t.Network, ipAddr)
				if err == nil {
					return conn, nil
				}
			}
			return nil, fmt.Errorf("failed to connect to any resolved IP addresses for %s", addr)
		}
		return t.Dialer.DialContext(ctx, t.Network, addr)
	}
	t.Base.DialContext = dialContext
	t.Base.Proxy = t.Proxy
	if err := http2.ConfigureTransport(t.Base); err == nil {
		t.Base.ForceAttemptHTTP2 = true
	}
	return t.Base.RoundTrip(req)
}

var Ipv4Transport = &CustomTransport{
	Dialer:   Dialer,
	Resolver: Dialer.Resolver,
	Network:  "tcp4",
	Proxy:    ClientProxy,
	Base: &http.Transport{
		MaxIdleConns:           100,
		IdleConnTimeout:        90 * time.Second,
		TLSHandshakeTimeout:    30 * time.Second,
		ExpectContinueTimeout:  1 * time.Second,
		TLSClientConfig:        tlsConfig,
		MaxResponseHeaderBytes: 262144,
	},
}

var Ipv4HttpClient = http.Client{
	Timeout:       30 * time.Second,
	CheckRedirect: UseLastResponse,
	Transport:     Ipv4Transport,
}

var Ipv6Transport = &CustomTransport{
	Dialer:   Dialer,
	Resolver: Dialer.Resolver,
	Network:  "tcp6",
	Proxy:    ClientProxy,
	Base: &http.Transport{
		MaxIdleConns:           100,
		IdleConnTimeout:        90 * time.Second,
		TLSHandshakeTimeout:    30 * time.Second,
		ExpectContinueTimeout:  1 * time.Second,
		TLSClientConfig:        tlsConfig,
		MaxResponseHeaderBytes: 262144,
	},
}

var Ipv6HttpClient = http.Client{
	Timeout:       30 * time.Second,
	CheckRedirect: UseLastResponse,
	Transport:     Ipv6Transport,
}

var AutoHttpClient = NewAutoHttpClient()

func AutoTransport() *CustomTransport {
	return &CustomTransport{
		Dialer:   Dialer,
		Resolver: Dialer.Resolver,
		Network:  "tcp",
		Proxy:    ClientProxy,
		Base: &http.Transport{
			MaxIdleConns:           100,
			IdleConnTimeout:        90 * time.Second,
			TLSHandshakeTimeout:    30 * time.Second,
			ExpectContinueTimeout:  1 * time.Second,
			TLSClientConfig:        tlsConfig,
			MaxResponseHeaderBytes: 262144,
		},
	}
}

func NewAutoHttpClient() http.Client {
	return http.Client{
		Timeout:       30 * time.Second,
		CheckRedirect: UseLastResponse,
		Transport:     AutoTransport(),
	}
}

func createEdgeTLSConfig() *tls.Config {
	spec, _ := utls.UTLSIdToSpec(utls.HelloEdge_Auto)

	return &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         spec.TLSVersMin,
		MaxVersion:         spec.TLSVersMax,
		CipherSuites:       spec.CipherSuites,
		ClientSessionCache: tls.NewLRUClientSessionCache(32),
		NextProtos:         []string{"h2", "http/1.1"},
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		PreferServerCipherSuites: false,
		SessionTicketsDisabled:   false,
	}
}

var tlsConfig = createEdgeTLSConfig()

type H [2]string

func GET(c http.Client, url string, headers ...H) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	setRealisticHeaders(req, "html")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	for _, h := range headers {
		req.Header.Set(h[0], h[1])
	}
	addRandomDelay()
	return cdo(c, req)
}

func GET_Dalvik(c http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UA_Dalvik)
	return cdo(c, req)
}

var ErrNetwork = errors.New("network error")

func cdo(c http.Client, req *http.Request) (resp *http.Response, err error) {
	deadline := time.Now().Add(30 * time.Second)
	for i := 0; i < 3; i++ {
		if time.Now().After(deadline) {
			break
		}
		if resp, err = c.Do(req); err == nil {
			return resp, nil
		}
		if strings.Contains(err.Error(), "no such host") {
			break
		}
		if strings.Contains(err.Error(), "timeout") {
			break
		}
	}
	return nil, err
}
func PostJson(c http.Client, url string, data string, headers ...H) (*http.Response, error) {
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
	return cdo(c, req)
}

func PostForm(c http.Client, url string, data string, headers ...H) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	setRealisticHeaders(req, "html")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	for _, h := range headers {
		req.Header.Set(h[0], h[1])
	}
	addRandomDelay()
	return cdo(c, req)
}

func genUUID() string {
	return uuid.New().String()
}

func md5Sum(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func twoToThreeCode(code string) string {
	countryCodes := map[string]string{
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
	return countryCodes[strings.ToUpper(code)]
}

func threeToTwoCode(code string) string {
	countryCodes := map[string]string{
		"AND": "AD", "ARE": "AE", "AFG": "AF", "ATG": "AG", "AIA": "AI", "ALB": "AL", "ARM": "AM", "AGO": "AO", "ATA": "AQ", "ARG": "AR",
		"ASM": "AS", "AUT": "AT", "AUS": "AU", "ABW": "AW", "ALA": "AX", "AZE": "AZ", "BIH": "BA", "BRB": "BB", "BGD": "BD", "BEL": "BE",
		"BFA": "BF", "BGR": "BG", "BHR": "BH", "BDI": "BI", "BEN": "BJ", "BLM": "BL", "BMU": "BM", "BRN": "BN", "BOL": "BO", "BES": "BQ",
		"BRA": "BR", "BHS": "BS", "BTN": "BT", "BVT": "BV", "BWA": "BW", "BLR": "BY", "BLZ": "BZ", "CAN": "CA", "CCK": "CC", "COD": "CD",
		"CAF": "CF", "COG": "CG", "CHE": "CH", "CIV": "CI", "COK": "CK", "CHL": "CL", "CMR": "CM", "CHN": "CN", "COL": "CO", "CRI": "CR",
		"CUB": "CU", "CPV": "CV", "CUW": "CW", "CXR": "CX", "CYP": "CY", "CZE": "CZ", "DEU": "DE", "DJI": "DJ", "DNK": "DK", "DMA": "DM",
		"DOM": "DO", "DZA": "DZ", "ECU": "EC", "EST": "EE", "EGY": "EG", "ESH": "EH", "ERI": "ER", "ESP": "ES", "ETH": "ET", "FIN": "FI",
		"FJI": "FJ", "FLK": "FK", "FSM": "FM", "FRO": "FO", "FRA": "FR", "GAB": "GA", "GBR": "GB", "GRD": "GD", "GEO": "GE", "GUF": "GF",
		"GGY": "GG", "GHA": "GH", "GIB": "GI", "GRL": "GL", "GMB": "GM", "GIN": "GN", "GLP": "GP", "GNQ": "GQ", "GRC": "GR", "SGS": "GS",
		"GTM": "GT", "GUM": "GU", "GNB": "GW", "GUY": "GY", "HKG": "HK", "HMD": "HM", "HND": "HN", "HRV": "HR", "HTI": "HT", "HUN": "HU",
		"IDN": "ID", "IRL": "IE", "ISR": "IL", "IMN": "IM", "IND": "IN", "IOT": "IO", "IRQ": "IQ", "IRN": "IR", "ISL": "IS", "ITA": "IT",
		"JEY": "JE", "JAM": "JM", "JOR": "JO", "JPN": "JP", "KEN": "KE", "KGZ": "KG", "KHM": "KH", "KIR": "KI", "COM": "KM", "KNA": "KN",
		"PRK": "KP", "KOR": "KR", "KWT": "KW", "CYM": "KY", "KAZ": "KZ", "LAO": "LA", "LBN": "LB", "LCA": "LC", "LIE": "LI", "LKA": "LK",
		"LBR": "LR", "LSO": "LS", "LTU": "LT", "LUX": "LU", "LVA": "LV", "LBY": "LY", "MAR": "MA", "MCO": "MC", "MDA": "MD", "MNE": "ME",
		"MAF": "MF", "MDG": "MG", "MHL": "MH", "MKD": "MK", "MLI": "ML", "MMR": "MM", "MNG": "MN", "MAC": "MO", "MNP": "MP", "MTQ": "MQ",
		"MRT": "MR", "MSR": "MS", "MLT": "MT", "MUS": "MU", "MDV": "MV", "MWI": "MW", "MEX": "MX", "MYS": "MY", "MOZ": "MZ", "NAM": "NA",
		"NCL": "NC", "NER": "NE", "NFK": "NF", "NGA": "NG", "NIC": "NI", "NLD": "NL", "NOR": "NO", "NPL": "NP", "NRU": "NR", "NIU": "NU",
		"NZL": "NZ", "OMN": "OM", "PAN": "PA", "PER": "PE", "PYF": "PF", "PNG": "PG", "PHL": "PH", "PAK": "PK", "POL": "PL", "SPM": "PM",
		"PCN": "PN", "PRI": "PR", "PSE": "PS", "PRT": "PT", "PLW": "PW", "PRY": "PY", "QAT": "QA", "REU": "RE", "ROU": "RO", "SRB": "RS",
		"RUS": "RU", "RWA": "RW", "SAU": "SA", "SLB": "SB", "SYC": "SC", "SDN": "SD", "SWE": "SE", "SGP": "SG", "SHN": "SH", "SVN": "SI",
		"SJM": "SJ", "SVK": "SK", "SLE": "SL", "SMR": "SM", "SEN": "SN", "SOM": "SO", "SUR": "SR", "SSD": "SS", "STP": "ST", "SLV": "SV",
		"SXM": "SX", "SYR": "SY", "SWZ": "SZ", "TCA": "TC", "TCD": "TD", "ATF": "TF", "TGO": "TG", "THA": "TH", "TJK": "TJ", "TKL": "TK",
		"TLS": "TL", "TKM": "TM", "TUN": "TN", "TON": "TO", "TUR": "TR", "TTO": "TT", "TUV": "TV", "TWN": "TW", "TZA": "TZ", "UKR": "UA",
		"UGA": "UG", "UMI": "UM", "USA": "US", "URY": "UY", "UZB": "UZ", "VAT": "VA", "VCT": "VC", "VEN": "VE", "VGB": "VG", "VIR": "VI",
		"VNM": "VN", "VUT": "VU", "WLF": "WF", "WSM": "WS", "YEM": "YE", "MYT": "YT", "ZAF": "ZA", "ZMB": "ZM", "ZWE": "ZW",
	}
	return countryCodes[strings.ToUpper(code)]
}

func genRandomStr(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[secureRandInt(len(charset))]
	}
	return string(b)
}

func generateEdgeUserAgent() string {
	edgeVersion := secureRandInRange(136, 140) // secureRandInt(5) + 136
	chromiumVersion := edgeVersion

	return fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.0.0 Safari/537.36 Edg/%d.0.0.0",
		chromiumVersion, edgeVersion)
}

func generateSecChUA() string {
	edgeVersion := secureRandInRange(136, 140) // secureRandInt(5) + 136
	chromiumVersion := edgeVersion
	notBrandVersion := secureRandInRange(20, 29) // secureRandInt(10) + 20

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

func setRealisticHeaders(req *http.Request, requestType string) {
	if ClientSessionHeaders.UserAgent == "" {
		ResetSessionHeaders()
	}
	req.Header.Set("user-agent", ClientSessionHeaders.UserAgent)
	switch requestType {
	case "json":
		req.Header.Set("accept", "application/json, text/plain, */*")
	case "html":
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	default:
		req.Header.Set("accept", "*/*")
	}
	req.Header.Set("sec-ch-ua", ClientSessionHeaders.SecChUA)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("accept-language", ClientSessionHeaders.AcceptLanguage)
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("dnt", ClientSessionHeaders.DNT)
}

type SessionHeaders struct {
	UserAgent      string
	SecChUA        string
	AcceptLanguage string
	DNT            string
}

func ResetSessionHeaders() {
	ClientSessionHeaders.UserAgent = generateEdgeUserAgent()
	ClientSessionHeaders.SecChUA = generateSecChUA()
	ClientSessionHeaders.AcceptLanguage = getRandomAcceptLanguage()
	ClientSessionHeaders.DNT = strconv.Itoa(secureRandInt(2))
}

// secureRandInt generates a cryptographically secure random integer in range [0, max)
func secureRandInt(max int) int {
	if max <= 0 {
		return 0
	}

	// Calculate the number of bytes needed
	maxBytes := make([]byte, 4) // Using 4 bytes for 32-bit int
	_, err := rand.Read(maxBytes)
	if err != nil {
		// Fallback to a simple approach if crypto/rand fails
		return 0
	}

	// Convert bytes to uint32 and then to int
	randUint32 := uint32(maxBytes[0])<<24 | uint32(maxBytes[1])<<16 | uint32(maxBytes[2])<<8 | uint32(maxBytes[3])
	return int(randUint32) % max
}

// secureRandInRange generates a cryptographically secure random integer in range [min, max]
func secureRandInRange(min, max int) int {
	if min >= max {
		return min
	}
	return secureRandInt(max-min+1) + min
}
