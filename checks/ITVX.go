package mediaunlocktest

import (
	"fmt"
	"io"
	"net/http"
)

func ITVX(c http.Client) Result {
	resp, err := GET(c, "https://simulcast.itv.com/playlist/itvonline/ITV",
		H{"x-custom-headers", "true"},
		H{"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 Edg/144.0.0.0"},
		H{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		H{"Accept-Encoding", "gzip, deflate, br, zstd"},
		H{"Accept-Language", "zh-HK,zh;q=0.9"},
		H{"Cache-Control", "max-age=0"},
		H{"Sec-Fetch-Dest", "document"},
		H{"Sec-Fetch-Mode", "navigate"},
		H{"Sec-Fetch-Site", "none"},
		H{"Sec-Fetch-User", "?1"},
		H{"Upgrade-Insecure-Requests", "1"},
		H{"sec-ch-ua", "\"Not(A:Brand\";v=\"8\", \"Chromium\";v=\"144\", \"Microsoft Edge\";v=\"144\""},
		H{"sec-ch-ua-mobile", "?0"},
		H{"sec-ch-ua-platform", "\"Windows\""},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	fmt.Printf("ITVX Status: %d\nBody: %s\n", resp.StatusCode, string(b))

	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}
	if resp.StatusCode == 404 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
