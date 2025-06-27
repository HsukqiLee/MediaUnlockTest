package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
)

func NBC_TV(c http.Client) Result {
	resp, err := PostJson(c, "https://geolocation.digitalsvc.apps.nbcuni.com/geolocation/live/nbc", `{"adobeMvpdId": null, "serviceZip": null, "device": "web"}`,
		H{"accept", "application/media.geo-v2+json"},
		H{"accept-encoding", "gzip, deflate, br, zstd"},
		H{"accept-language", "zh-CN,zh;q=0.9"},
		H{"app-session-id", "71203A7A-BC80-4452-B8B7-796C564C6EF3"},
		H{"authorization", `NBC-Basic key="nbc_live", version="3.0", type="cpc"`},
		H{"cache-control", "no-cache"},
		H{"client", "oneapp"},
		H{"content-type", "application/json"},
		H{"origin", "https://www.nbc.com"},
		H{"pragma", "no-cache"},
		H{"priority", "u=1, i"},
		H{"referer", "https://www.nbc.com/"},
		H{"sec-ch-ua", `"Microsoft Edge";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`},
		H{"sec-ch-ua-mobile", "?0"},
		H{"sec-ch-ua-platform", `"Windows"`},
		H{"sec-fetch-dest", "empty"},
		H{"sec-fetch-mode", "cors"},
		H{"sec-fetch-site", "cross-site"},
		H{"user-agent"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
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
		return Result{Status: StatusErr, Err: err}
	}

	if !res.Restricted {
		return Result{Status: StatusOK, Region: res.RequestInfo.CountryCode}
	}

	if res.RestrictionDetails.Code == "321" {
		return Result{Status: StatusNo}
	}

	return Result{Status: StatusUnexpected}
}
