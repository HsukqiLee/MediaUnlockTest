package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
)

func KPlus(c http.Client) Result {
	resp, err := PostJson(c, "https://tvapi-sgn.solocoo.tv/v1/provision", `{"osVersion":"Windows 10","deviceModel":"Edge","deviceType":"PC","deviceSerial":"w7ab83550-c0aa-11ee-bf07-531681e47537","deviceOem":"Edge","devicePrettyName":"Edge 121.0.0.0","appVersion":"11.0","language":"en_US","brand":"vstv","featureLevel":5}`)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var result struct {
		Session struct {
			GeoCountryCode string `json:"geoCountryCode"`
		} `json:"session"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return Result{Status: StatusErr, Err: err}
	}

	region := result.Session.GeoCountryCode
	switch region {
	case "VN":
		return Result{Status: StatusOK, Region: "vn"}
	case "":
		return Result{Status: StatusUnexpected}
	default:
		return Result{Status: StatusNo}
	}
}
