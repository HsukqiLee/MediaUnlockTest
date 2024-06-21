package mediaunlocktest

import (
	"net/http"
)

func Watcha(c http.Client) Result {
	resp, err := GET(c, "https://watcha.com/",
	    H{"User-Agent", UA_Browser},
	    H{"Host", "watcha.com"},
	    H{"Connection", "keep-alive"},
	    H{"Upgrade-Insecure-Requests", "1"},
	    H{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	
	if resp.StatusCode == 451 {
		return Result{Status: StatusNo}
	}
	if resp.StatusCode == 403 {
		return Result{Status: StatusBanned}
	}
	
	if resp.StatusCode == 200  {
		return Result{Status: StatusOK, Region: "kr"}
	}
	
	if resp.StatusCode == 302 && resp.Header.Get("Location") == "/ja-JP/" {
		return Result{Status: StatusOK, Region: "jp"}
	}

	return Result{Status: StatusUnexpected}
}