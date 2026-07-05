package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"regexp"
	"strings"
)

var kimiRegionRegex = regexp.MustCompile(`"useRegion":"REGION_([^"]+)"`)

func Kimi(c core.HttpClient) core.Result {
	resp, err := core.RequestRaw(c, "POST", "https://www.kimi.com/apiv2/kimi.gateway.order.v1.GoodsService/ListGoods", "{}", core.H{"content-type", "application/json"})
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
			matches := kimiRegionRegex.FindStringSubmatch(string(b))
			if len(matches) > 1 {
				return core.Result{Status: core.StatusOK, Region: strings.ToLower(matches[1])}
			}
		}
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
