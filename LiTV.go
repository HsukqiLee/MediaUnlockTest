package mediaunlocktest

import (
    "io"
	"net/http"
	"strings"
)

func LiTV(c http.Client) Result {
	resp, err := PostJson(c, "https://www.litv.tv/api/get-urls-no-auth",
		`{"AssetId": "vod71211-000001M001_1500K","MediaType": "vod","puid": "d66267c2-9c52-4b32-91b4-3e482943fe7e"}`,
		H{"Cookie", "PUID=34eb9a17-8834-4f83-855c-69382fd656fa; L_PUID=34eb9a17-8834-4f83-855c-69382fd656fa; device-id=f4d7faefc54f476bb2e7e27b7482469a"},
		H{"Origin", "https://www.litv.tv"},
		H{"Referer", "https://www.litv.tv/drama/watch/VOD00331042"},
		H{"Priority", "u=1, i"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	
	bodyBytes, err := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	
	
	if resp.StatusCode == 200 {
	    if strings.Contains(bodyString, "OutsideRegionError") {
		    return Result{Status: StatusNo}
	    }
		return Result{Status: StatusOK}
	}


	return Result{Status: StatusUnexpected}
}
