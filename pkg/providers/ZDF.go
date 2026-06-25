package mediaunlocktest

import (
	"net/http"
)

func ZDF(c http.Client) Result {
	return CheckDalvikStatus(c, "https://ssl.zdf.de/geo/de/geo.txt", ResultMap{
		http.StatusForbidden: {Status: StatusNo},
		http.StatusOK:        {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
