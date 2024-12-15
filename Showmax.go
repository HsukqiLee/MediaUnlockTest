package mediaunlocktest

import (
	"net/http"
	"regexp"
	"strings"
	//"fmt"
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
		H{"Connection", "keep-alive"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	cookie := resp.Header.Get("Set-Cookie")
	//fmt.Println(cookie)
	if cookie == "" {
		return Result{Status: StatusNo}
	}

	region := extractShowmaxRegion(cookie)
	if region != "" {
		return Result{Status: StatusOK, Region: strings.ToLower(region)}
	}
	return Result{Status: StatusUnexpected}
}
