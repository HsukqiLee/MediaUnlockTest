package mediaunlocktest

import (
	"context"
	"net/http"
	"time"
)

func PrettyDerbyJP(c http.Client) Result {
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req, _ := http.NewRequestWithContext(ctx, "GET", "https://api-umamusume.cygames.jp/", nil)
		req.Header.Set("user-agent", UA_Dalvik)
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Set("cache-control", "no-cache")
		req.Header.Set("dnt", "1")
		req.Header.Set("pragma", "no-cache")
		req.Header.Set("sec-ch-ua", `"Not(A:Brand";v="99", "Microsoft Edge";v="133", "Chromium";v="133"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", "macOS")
		req.Header.Set("sec-fetch-dest", "document")
		req.Header.Set("sec-fetch-mode", "navigate")
		req.Header.Set("sec-fetch-site", "none")
		req.Header.Set("sec-fetch-user", "?1")
		req.Header.Set("upgrade-insecure-requests", "1")
		resp, err := c.Do(req)

		if err != nil {
			if err.Error() == `Get "https://api-umamusume.cygames.jp/": context deadline exceeded` {
				return Result{Status: StatusNo}
			}
		}
		defer resp.Body.Close()
		switch resp.StatusCode {
		case 404:
			return Result{Status: StatusOK}
		}
	}
	return Result{Status: StatusNo}
}
