package mediaunlocktest

import (
	"net/http"
)

func Lemino(c http.Client) Result {
	resp, err := PostJson(c, "https://if.lemino.docomo.ne.jp/v1/user/delivery/watch/ready",
	    `{"inflow_flows":[null,"crid://plala.iptvf.jp/group/b100ce3"],"play_type":1,"key_download_only":null,"quality":null,"groupcast":null,"avail_status":"1","terminal_type":3,"test_account":0,"content_list":[{"kind":"main","service_id":null,"cid":"00lm78dz30","lid":"a0lsa6kum1","crid":"crid://plala.iptvf.jp/vod/0000000000_00lm78dymn","preview":0,"trailer":0,"auto_play":0,"stop_position":0}]}`,
	    H{"accept", "application/json, text/plain, */*"},
	    H{"accept-language", "en-US,en;q=0.9"},
	    H{"content-type", "application/json"},
	    H{"origin", "https://lemino.docomo.ne.jp"},
	    H{"referer", "https://lemino.docomo.ne.jp/"},
	    H{"x-service-token", "f365771afd91452fa279863f240c233d" },
	    H{"x-trace-id", "556db33f-d739-4a82-84df-dd509a8aa179"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 500 {
		return Result{Status: StatusErr}
	}
	
	if resp.StatusCode == 200  {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}