package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
)

func NBC_TV(c http.Client) core.Result {
	resp, err := core.PostJson(c, "https://geolocation.digitalsvc.apps.nbcuni.com/geolocation/live/nbc", `{"adobeMvpdId": null, "serviceZip": null, "device": "web"}`,
		core.H{"accept", "application/media.geo-v2+json"},
		core.H{"accept-encoding", "gzip, deflate, br, zstd"},
		core.H{"accept-language", "zh-CN,zh;q=0.9"},
		core.H{"app-session-id", "71203A7A-BC80-4452-B8B7-796C564C6EF3"},
		core.H{"authorization", `NBC-Basic key="nbc_live", version="3.0", type="cpc"`},
		core.H{"cache-control", "no-cache"},
		core.H{"client", "oneapp"},
		core.H{"content-type", "application/json"},
		core.H{"origin", "https://www.nbc.com"},
		core.H{"pragma", "no-cache"},
		core.H{"priority", "u=1, i"},
		core.H{"referer", "https://www.nbc.com/"},
		core.H{"sec-ch-ua", `"Microsoft Edge";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`},
		core.H{"sec-ch-ua-mobile", "?0"},
		core.H{"sec-ch-ua-platform", `"Windows"`},
		core.H{"sec-fetch-dest", "empty"},
		core.H{"sec-fetch-mode", "cors"},
		core.H{"sec-fetch-site", "cross-site"},
		core.H{"user-agent"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res struct {
		Restricted  bool `json:"restricted"`
		RequestInfo struct {
			CountryCode string `json:"countryCode"`
		} `json:"requestInfo"`
		RestrictionDetails struct {
			Code string `json:"code"`
		} `json:"restrictionDetails"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	if !res.Restricted {
		return core.Result{Status: core.StatusOK, Region: res.RequestInfo.CountryCode}
	}

	if res.RestrictionDetails.Code == "321" {
		return core.Result{Status: core.StatusNo}
	}

	return core.Result{Status: core.StatusUnexpected}
}

