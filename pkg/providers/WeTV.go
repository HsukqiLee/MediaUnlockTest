package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"regexp"
	"strings"
)

func WeTV(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://wetv.vip/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}

		body := string(b)
		if strings.Contains(body, "static.wetvinfo.com/static/vg-cn/tar") {
			return core.Result{Status: core.StatusNo, Region: "cn"}
		}

		re := regexp.MustCompile(`(?i)"areaPhoneId"\s*:\s*"\+([0-9]+)"`)
		match := re.FindStringSubmatch(body)
		if len(match) > 1 {
			region := core.CountryCodeToAlpha2(match[1])
			if region != "" {
				return core.Result{Status: core.StatusOK, Region: strings.ToLower(region)}
			}
		}

		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
