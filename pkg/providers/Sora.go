package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func Sora(c http.Client) core.Result {
	c.Jar, _ = cookiejar.New(nil)
	resp, err := core.GET(c, "https://sora.com/cdn-cgi/trace")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	i := strings.Index(s, "loc=")
	if i == -1 {
		return core.Result{Status: core.StatusUnexpected}
	}
	s = s[i+4:]
	i = strings.Index(s, "\n")
	if i == -1 {
		return core.Result{Status: core.StatusUnexpected}
	}
	loc := s[:i]

	resp, err = core.GET(c, "https://sora.com/backend/authenticate")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	if strings.Contains(string(b), "Attention Required") {
		return core.Result{Status: core.StatusBanned, Info: "VPN Blocked"}
	}
	if resp.StatusCode == 429 {
		return core.Result{Status: core.StatusRestricted, Region: strings.ToLower(loc), Info: "429 Rate limit"}
	}
	if loc == "T1" {
		return core.Result{Status: core.StatusOK, Region: "tor"}
	}
	if SupportGPT(loc) {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(loc)}
	}
	return core.Result{Status: core.StatusNo, Region: strings.ToLower(loc)}
}

