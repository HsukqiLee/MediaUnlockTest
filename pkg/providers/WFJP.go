package mediaunlocktest

import "net/http"

// World Flipper Japan
func WFJP(c http.Client) Result {
	return CheckDalvikStatus(c, "https://api.worldflipper.jp/", ResultMap{
		http.StatusOK:        {Status: StatusOK},
		http.StatusForbidden: {Status: StatusNo},
	}, Result{Status: StatusUnexpected})
}
