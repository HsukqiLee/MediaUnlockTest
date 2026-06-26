package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"regexp"
	"strings"
)

func PrimeVideo(c core.HttpClient) core.Result {
	c.SetFollowRedirect(true)
	resp, err := core.GET(c, "https://www.primevideo.com")
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned, Err: err}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var s string
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		loc := resp.Header.Get("Location")
		if loc != "" {
			if strings.HasPrefix(loc, "/") {
				loc = "https://www.primevideo.com" + loc
			}
			resp2, err := core.GET(c, loc)
			if err == nil {
				b2, _ := io.ReadAll(resp2.Body)
				s = string(b2)
				resp2.Body.Close()
			}
		}
	} else {
		b, err := io.ReadAll(resp.Body)
		if err == nil {
			s = string(b)
		}
	}
	resp.Body.Close()

	// Check if the current page contains a regional Amazon redirection link
	if i := strings.Index(s, `"currentTerritory":`); i == -1 {
		re := regexp.MustCompile(`(https://www\.amazon\.[a-z\.]+/[^"'\s>]+)`)
		matches := re.FindAllStringSubmatch(s, -1)
		for _, match := range matches {
			urlStr := match[1]
			// Only follow storefront links to avoid capturing signin/login redirects
			if strings.Contains(urlStr, "storefront") {
				urlStr = strings.ReplaceAll(urlStr, "&amp;", "&")
				localResp, err := core.GET(c, urlStr)
				if err == nil {
					localBody, _ := io.ReadAll(localResp.Body)
					s = string(localBody) // overwrite s with the regional storefront HTML
					localResp.Body.Close()
					break
				}
			}
		}
	}

	// WAF block or Captcha check
	if strings.Contains(s, "api-services-support@amazon.com") {
		return core.Result{Status: core.StatusNo}
	}

	// Extract the region universally
	if i := strings.Index(s, `"currentTerritory":`); i != -1 {
		return core.Result{
			Status: core.StatusOK,
			Region: strings.ToLower(s[i+20 : i+22]),
		}
	}

	return core.Result{Status: core.StatusNo}
}
