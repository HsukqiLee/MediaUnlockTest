package core

import (
	"strconv"
	"sync"

	http "github.com/bogdanfinn/fhttp"
)

var (
	UA_Browser      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36"
	UA_Dalvik       = "Dalvik/2.1.0 (Linux; U; Android 14; M2006J10C Build/RP1A.200720.011)"
	SecChUA_Browser = `"Not A(Brand";v="99", "Google Chrome";v="146", "Chromium";v="146"`

	ClientSessionHeaders = &SessionHeaders{
		UserAgent:      "",
		SecChUA:        "",
		AcceptLanguage: "",
		DNT:            "0",
	}
)

var sessionMutex sync.RWMutex

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

func GetRealisticHeaders(requestType string) []H {
	headers := make([]H, 0)
	acceptLanguage := getRandomAcceptLanguage()
	dnt := strconv.Itoa(secureRandInt(2))
	headers = append(headers, H{"user-agent", UA_Browser})
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
	headers = append(headers, H{"sec-ch-ua", SecChUA_Browser})
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
	var order []string
	for _, header := range headers {
		req.Header.Set(header[0], header[1])
		order = append(order, header[0])
	}
	// append to existing HeaderOrderKey if any
	if existing, ok := req.Header[http.HeaderOrderKey]; ok {
		req.Header[http.HeaderOrderKey] = append(existing, order...)
	} else {
		req.Header[http.HeaderOrderKey] = order
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
		UserAgent:      UA_Browser,
		SecChUA:        SecChUA_Browser,
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
