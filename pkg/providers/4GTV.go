package providers

import (
	"MediaUnlockTest/pkg/core"
)

func TW4GTV(c core.HttpClient) core.Result {
	success, err := core.PostFormBoolSuccess(c, "https://api2.4gtv.tv//Vod/GetVodUrl3",
		`value=D33jXJ0JVFkBqV%2BZSi1mhPltbejAbPYbDnyI9hmfqjKaQwRQdj7ZKZRAdb16%2FRUrE8vGXLFfNKBLKJv%2BfDSiD%2BZJlUa5Msps2P4IWuTrUP1%2BCnS255YfRadf%2BKLUhIPj`,
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	if success {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}
