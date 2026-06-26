package providers

import (
	"MediaUnlockTest/pkg/core"
)

func Crackle(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://prod-api.crackle.com/appconfig",
		core.H{"Origin", "https://www.crackle.com"},
		core.H{"Referer", "https://www.crackle.com/"},
		core.H{"x-crackle-apiversion", "v2.0.0"},
		core.H{"x-crackle-brand", "crackle"},
		core.H{"x-crackle-platform", "5FE67CCA-069A-42C6-A20F-4B47A8054D46"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch region := resp.Header.Get("x-crackle-region"); {
	case region == "US":
		return core.Result{Status: core.StatusOK}
	case region != "":
		return core.Result{Status: core.StatusNo}
	default:
		return core.Result{Status: core.StatusUnexpected}
	}
}
