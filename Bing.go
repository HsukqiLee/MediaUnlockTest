package mediaunlocktest

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
)

func extractBingRegion(responseBody string) string {
	re := regexp.MustCompile(`Region:"([^"]*)"`)
	match := re.FindStringSubmatch(responseBody)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func Bing(c http.Client) Result {
	c.Jar, _ = cookiejar.New(nil)
	resp, err := GET(c, "https://www.bing.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.Header.Get("Location") == "https://www.bing.com/?brdr=1" {
		resp, err = GET(c, "https://www.bing.com/")
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		defer resp.Body.Close()
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return Result{Status: StatusFailed}
	}

	if resp.StatusCode == 200 {
		region := extractBingRegion(bodyString)
		if region == "CN" {
			return Result{Status: StatusNo, Region: "cn"}
		}
		if region != "" {
			return Result{Status: StatusOK, Region: strings.ToLower(region)}
		}
	}

	if strings.Contains(bodyString, "cn.bing.com") {
		return Result{Status: StatusNo, Region: "cn"}
	}

	if resp.StatusCode == 403 {
		return Result{Status: StatusBanned}
	}

	fmt.Println(resp.Header.Get("Location"))
	return Result{Status: StatusUnexpected}
}
