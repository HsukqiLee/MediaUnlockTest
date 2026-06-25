package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
	"encoding/json"
	"io"
)

func FranceTV(c http.Client) core.Result {
	resp, err := core.GET(c, "https://geo-info.ftven.fr/ws/edgescape.json")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
    
    b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	
    var res struct {
		Response struct{
		    GeoInfo struct {
		        CountryCode string `json:"country_code"`
		    } `json:"geo_info"`
		} `json:"reponse"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}
	
	if res.Response.GeoInfo.CountryCode == "FR" {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusNo}
}
