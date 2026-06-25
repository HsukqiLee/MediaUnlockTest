package mediaunlocktest

import "net/http"

func Karaoke(c http.Client) Result {
	return CheckGETStatus(c, "http://cds1.clubdam.com/vhls-cds1/site/xbox/sample_1.mp4.m3u8", ResultMap{
		http.StatusForbidden: {Status: StatusNo},
		http.StatusOK:        {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
