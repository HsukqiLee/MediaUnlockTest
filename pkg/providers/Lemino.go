package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func Lemino(c http.Client) core.Result {
	resp, err := core.PostJson(c, "https://if.lemino.docomo.ne.jp/v1/user/delivery/watch/ready",
		`{"inflow_flows":[null,"crid://plala.iptvf.jp/group/b100ce3"],"play_type":1,"key_download_only":null,"quality":null,"groupcast":null,"avail_status":"1","terminal_type":3,"test_account":0,"content_list":[{"kind":"main","service_id":null,"cid":"00lm78dz30","lid":"a0lsa6kum1","crid":"crid://plala.iptvf.jp/vod/0000000000_00lm78dymn","preview":0,"trailer":0,"auto_play":0,"stop_position":0}]}`,
		core.H{"accept", "application/json, text/plain, */*"},
		core.H{"accept-language", "en-US,en;q=0.9"},
		core.H{"content-type", "application/json"},
		core.H{"origin", "https://lemino.docomo.ne.jp"},
		core.H{"referer", "https://lemino.docomo.ne.jp/"},
		core.H{"x-service-token", "f365771afd91452fa279863f240c233d"},
		core.H{"x-trace-id", "556db33f-d739-4a82-84df-dd509a8aa179"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusOK:                  {Status: core.StatusOK},
		http.StatusForbidden:           {Status: core.StatusNo},
		http.StatusInternalServerError: {Status: core.StatusErr},
	}, core.Result{Status: core.StatusUnexpected})
}

