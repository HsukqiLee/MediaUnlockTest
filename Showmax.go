package mediaunlocktest

import (
    "io"
	"net/http"
	"strings"
)

func Showmax(c http.Client) Result {
	resp, err := GET(c, "https://www.showmax.com/",
	    H{"host", "www.showmax.com"},
        H{"connection", "keep-alive"},
        H{"upgrade-insecure-requests", "1"},
        H{"accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

    regionStart := strings.Index(bodyString, "activeTerritory")
    if regionStart == -1 {
        return Result{Status: StatusNo}
    }

    regionEnd := strings.Index(bodyString[regionStart:], "\n")
    region := strings.TrimSpace(bodyString[regionStart+len("activeTerritory")+1 : regionStart+regionEnd])
    

	if region != "" {
		return Result{Status: StatusOK, Region: strings.ToLower(region)}
	}

	return Result{Status: StatusUnexpected}
}

