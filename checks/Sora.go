package mediaunlocktest

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func Sora(c http.Client) Result {
	c.Jar, _ = cookiejar.New(nil)
	resp, err := GET(c, "https://sora.com/cdn-cgi/trace")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	s := string(b)
	i := strings.Index(s, "loc=")
	if i == -1 {
		return Result{Status: StatusUnexpected}
	}
	s = s[i+4:]
	i = strings.Index(s, "\n")
	if i == -1 {
		return Result{Status: StatusUnexpected}
	}
	loc := s[:i]

	resp, err = GET(c, "https://sora.com/backend/authenticate")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	if strings.Contains(string(b), "Attention Required") {
		return Result{Status: StatusBanned, Info: "VPN Blocked"}
	}
	if resp.StatusCode == 429 {
		return Result{Status: StatusRestricted, Region: strings.ToLower(loc), Info: "429 Rate limit"}
	}
	if loc == "T1" {
		return Result{Status: StatusOK, Region: "tor"}
	}
	if SupportGPT(loc) {
		return Result{Status: StatusOK, Region: strings.ToLower(loc)}
	}
	return Result{Status: StatusNo, Region: strings.ToLower(loc)}
}
