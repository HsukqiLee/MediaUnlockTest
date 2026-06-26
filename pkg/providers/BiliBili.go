package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
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

func BilibiliHKMO(c core.HttpClient) core.Result {
	return bilibili(c, "https://api.bilibili.com/pgc/player/web/playurl?avid=473502608&cid=845838026&qn=0&type=&otype=json&ep_id=678506&fourk=1&fnver=0&fnval=16&module=bangumi")
}

func BilibiliTW(c core.HttpClient) core.Result {
	return bilibili(c, "https://api.bilibili.com/pgc/player/web/playurl?avid=50762638&cid=100279344&qn=0&type=&otype=json&ep_id=268176&fourk=1&fnver=0&fnval=16&module=bangumi")
}

func BilibiliSEA(c core.HttpClient) core.Result {
	return bilibili(c, "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=347666")
}

func BilibiliTH(c core.HttpClient) core.Result {
	return bilibili(c, "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=10077726")
}

func BilibiliID(c core.HttpClient) core.Result {
	return bilibili(c, "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=11130043")
}

func BilibiliVN(c core.HttpClient) core.Result {
	return bilibili(c, "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=11405745")
}
