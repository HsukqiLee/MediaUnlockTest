package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func KBS(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://vod.kbs.co.kr/index.html?source=episode&sname=vod&stype=vod&program_code=T2022-0690&program_id=PS-2022164275-01-000&broadcast_complete_yn=N&local_station_code=00&section_code=03")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	if resp.StatusCode == 200 && strings.Contains(bodyString, `\"Domestic\": true`) {
		return core.Result{Status: core.StatusOK}
	}

	if strings.Contains(bodyString, ">새로고침<") {
		return core.Result{Status: core.StatusNo}
	}

	return core.Result{Status: core.StatusUnexpected}
}
