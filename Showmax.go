package mediaunlocktest

import (
	"net/http"
	"strings"
	"regexp"
	"fmt"
)

func extractShowmaxRegion(body string) string {
    re := regexp.MustCompile(`hterr=([A-Z]{2})`)
    matches := re.FindStringSubmatch(string(body))
    if len(matches) > 1 {
        return matches[1]
    }
    return ""
}

func Showmax(c http.Client) Result {
	resp, err := GET(c, "https://www.showmax.com/",
        H{"Host", "www.showmax.com"},
        H{"Connection", "keep-alive"},
        H{"sec-ch-ua", `"Chromium";v="124", "Microsoft Edge";v="124", "Not-A.Brand";v="99"`},
        H{"sec-ch-ua-mobile", "?0"},
        H{"sec-ch-ua-platform", "Windows"},
        H{"Upgrade-Insecure-Requests", "1"},
        H{"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0"},
        H{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
        H{"sec-fetch-site", "none"},
        H{"sec-fetch-mode", "navigate"},
        H{"sec-fetch-user", "?1"},
        H{"sec-fetch-dest", "document"},
        H{"Accept-Language", "zh-CN,zh;q=0.9"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	

	
	cookie := resp.Header.Get("Set-Cookie")
	
	fmt.Println(cookie)
	if cookie == "" {
		return Result{Status: StatusNo}
	}

    region := extractShowmaxRegion(cookie)
	if region != "" {
		return Result{Status: StatusOK, Region: strings.ToLower(region)}
	}

	return Result{Status: StatusUnexpected}
}
