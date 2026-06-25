package mediaunlocktest

import (
	"net/http"
)

func HoyTV(c http.Client) Result {
	return CheckGETStatus(c, "https://hoytv-live-stream.hoy.tv/ch78/index-fhd.m3u8", ResultMap{
		http.StatusForbidden: {Status: StatusNo},
		http.StatusOK:        {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
