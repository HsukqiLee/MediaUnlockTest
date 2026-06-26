package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func ABCiView(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://api.iview.abc.net.au/v2/show/abc-kids-live-stream/video/LS1604H001S00?embed=highlightVideo,selectedSeries")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return core.Result{Status: core.StatusFailed}
	}

	if strings.Contains(bodyString, "unavailable outside Australia") {
		return core.Result{Status: core.StatusNo}
	}

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusOK: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
