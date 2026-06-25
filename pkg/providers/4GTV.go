package mediaunlocktest

import "net/http"

func TW4GTV(c http.Client) Result {
	success, err := PostFormBoolSuccess(c, "https://api2.4gtv.tv//Vod/GetVodUrl3",
		`value=D33jXJ0JVFkBqV%2BZSi1mhPltbejAbPYbDnyI9hmfqjKaQwRQdj7ZKZRAdb16%2FRUrE8vGXLFfNKBLKJv%2BfDSiD%2BZJlUa5Msps2P4IWuTrUP1%2BCnS255YfRadf%2BKLUhIPj`,
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	if success {
		return Result{Status: StatusOK}
	}
	return Result{Status: StatusNo}
}
