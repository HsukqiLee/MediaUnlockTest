package mediaunlocktest

import (
	"net/http"
)

func AnimeFesta(c http.Client) Result {
	resp, err := GET(c, "https://api-animefesta.iowl.jp/v1/titles/1560",
		H{"accept", "application/json"},
		H{"accept-language", "en-US,en;q=0.9"},
		H{"anime-user-tracking-id", "yEZr4P_U7JEdBucZOkv1Y"},
		H{"authorization", ""},
		H{"origin", "https://animefesta.iowl.jp"},
		H{"referer", "https://animefesta.iowl.jp/"},
		H{"sec-gpc", "1"},
		H{"x-requested-with", "XMLHttpRequest"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}

	if resp.StatusCode == 200 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
