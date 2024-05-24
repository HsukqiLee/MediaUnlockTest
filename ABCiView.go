package mediaunlocktest

import (
    "io"
	"net/http"
	"strings"
)

func ABCiView(c http.Client) Result {
	resp, err := GET(c, "https://api.iview.abc.net.au/v2/show/abc-kids-live-stream/video/LS1604H001S00?embed=highlightVideo,selectedSeries")
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    
    if err != nil {
		return Result{Status: StatusFailed}
	}
	
	if strings.Contains(bodyString, "unavailable outside Australia") {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 200 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}