package mediaunlocktest

import (
    "io"
	"net/http"
	"strings"
)

func DAnimeStore(c http.Client) Result {
	resp, err := GET(c, "https://animestore.docomo.ne.jp/animestore/reg_pc")
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    
    if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	
	if resp.StatusCode == 403 && strings.Contains(bodyString, "海外") {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 200 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}