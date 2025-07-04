package mediaunlocktest

import (
	"io"
	"net/http"
	"regexp"
)

func extractAETVCountryCode(html string) string {
	re := regexp.MustCompile(`<meta\s+name=["']aetn:countryCode["']\s+content=["']([A-Z]{2})["']\s*/?>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

func AETV(c http.Client) Result {
	resp, err := GET(c, "https://www.aetv.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	region := extractAETVCountryCode(string(b))
	switch region {
	case "US":
		return Result{Status: StatusOK}
	case "":
		return Result{Status: StatusUnexpected}
	default:
		return Result{Status: StatusNo}
	}
}
