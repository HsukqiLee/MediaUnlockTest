package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func DAnimeStore(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://animestore.docomo.ne.jp/animestore/reg_pc")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	if strings.Contains(bodyString, "海外") {
		return core.Result{Status: core.StatusNo}
	}

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusFound: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
