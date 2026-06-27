package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"strconv"
)

func bilibili(c core.HttpClient, url string) core.Result {
	resp, err := core.GET(c, url)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		Code int
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if resp.StatusCode == 412 {
		return core.Result{Status: core.StatusFailed}
	}
	switch res.Code {
	case -10403, 10004001, 10003003:
		return core.Result{Status: core.StatusNo}
	case 0:
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusUnexpected}
}

func BilibiliAnime(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://api.bilibili.com/x/web-interface/zone")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		Code int `json:"code"`
		Data struct {
			CountryCode int `json:"country_code"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.Code != 0 {
		return core.Result{Status: core.StatusNo}
	}

	region := core.CountryCodeToAlpha2(strconv.Itoa(res.Data.CountryCode))
	if region == "" {
		region = strconv.Itoa(res.Data.CountryCode)
	}

	var testUrl string
	switch region {
	case "HK", "MO":
		testUrl = "https://api.bilibili.com/pgc/player/web/playurl?avid=473502608&cid=845838026&qn=0&type=&otype=json&ep_id=678506&fourk=1&fnver=0&fnval=16&module=bangumi"
	case "TW":
		testUrl = "https://api.bilibili.com/pgc/player/web/playurl?avid=50762638&cid=100279344&qn=0&type=&otype=json&ep_id=268176&fourk=1&fnver=0&fnval=16&module=bangumi"
	case "TH":
		testUrl = "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=10077726"
	case "ID":
		testUrl = "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=11130043"
	case "VN":
		testUrl = "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=11405745"
	case "MY", "SG", "PH", "BN", "KH", "LA", "MM", "TL":
		testUrl = "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=347666"
	default:
		return core.Result{Status: core.StatusOK, Region: region}
	}

	testRes := bilibili(c, testUrl)
	switch testRes.Status {
	case core.StatusOK:
		testRes.Region = region
		return testRes
	case core.StatusNo, core.StatusFailed:
		return core.Result{Status: core.StatusRestricted, Region: region}
	}

	testRes.Region = region
	return testRes
}
