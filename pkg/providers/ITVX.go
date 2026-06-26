package providers

import (
	"MediaUnlockTest/pkg/core"
)

func ITVX(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://simulcast.itv.com/playlist/itvonline/ITV",
		core.H{"x-custom-headers", "true"},
		core.H{"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 Edg/144.0.0.0"},
		core.H{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		core.H{"Accept-Encoding", "gzip, deflate, br, zstd"},
		core.H{"Accept-Language", "zh-HK,zh;q=0.9"},
		core.H{"Cache-Control", "max-age=0"},
		core.H{"Sec-Fetch-Dest", "document"},
		core.H{"Sec-Fetch-Mode", "navigate"},
		core.H{"Sec-Fetch-Site", "none"},
		core.H{"Sec-Fetch-User", "?1"},
		core.H{"Upgrade-Insecure-Requests", "1"},
		core.H{"sec-ch-ua", "\"Not(A:Brand\";v=\"8\", \"Chromium\";v=\"144\", \"Microsoft Edge\";v=\"144\""},
		core.H{"sec-ch-ua-mobile", "?0"},
		core.H{"sec-ch-ua-platform", "\"Windows\""},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusNo}
	}
	if resp.StatusCode == 404 {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
