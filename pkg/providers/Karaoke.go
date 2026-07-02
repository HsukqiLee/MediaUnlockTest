package providers

import (
	"MediaUnlockTest/pkg/core"
)

func Karaoke(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "http://cds1.clubdam.com/vhls-cds1/site/xbox/sample_1.mp4.m3u8", core.ResultMap{
		403: {Status: core.StatusNo},
		200: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
