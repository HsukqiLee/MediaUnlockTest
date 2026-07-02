package providers

import (
	"MediaUnlockTest/pkg/core"
)

// Project Sekai: Colorful Stage
func PJSK(c core.HttpClient) core.Result {
	resp, err := core.GETRaw(c, "https://game-version.sekai.colorfulpalette.org/1.8.1/3ed70b6a-8352-4532-b819-108837926ff5", core.H{"User-Agent", "pjsekai/48 CFNetwork/1240.0.4 Darwin/20.6.0"})
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return core.Result{Status: core.StatusOK}
	case 403:
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusUnexpected}
}
