package providers

import (
	"MediaUnlockTest/pkg/core"
	"crypto/tls"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var deepseekRegionRegex = regexp.MustCompile(`<meta\s+name="region"\s+content="([^"]+)"`)

func DeepSeek(c core.HttpClient) core.Result {
	// Temporarily bypass DNS to hit CN IP directly
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 10 * time.Second,
	}

	req, _ := http.NewRequest("GET", "https://116.205.40.114/sign_in", nil)
	req.Header.Set("Host", "chat.deepseek.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	resp, err := client.Do(req)
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 403:
		return core.Result{Status: core.StatusNo}
	case 200:
		b, err := io.ReadAll(resp.Body)
		if err == nil {
			matches := deepseekRegionRegex.FindStringSubmatch(string(b))
			if len(matches) > 1 {
				return core.Result{Status: core.StatusOK, Region: strings.ToLower(matches[1])}
			}
		}
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
