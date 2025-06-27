package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func MetaAI(c http.Client) Result {
	resp, err := GET(c, "https://www.meta.ai/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	if strings.Contains(bodyString, "GeoBlockedErrorRoot") {
		return Result{Status: StatusNo}
	}

	if strings.Contains(bodyString, "HomeRootQuery") {
		return Result{Status: StatusOK}
	}
	return Result{Status: StatusUnexpected}
}
