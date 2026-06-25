package mediaunlocktest

import (
	"net/http"
	"encoding/json"
	"io"
)

func FranceTV(c http.Client) Result {
	resp, err := GET(c, "https://geo-info.ftven.fr/ws/edgescape.json")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
    
    b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	
    var res struct {
		Response struct{
		    GeoInfo struct {
		        CountryCode string `json:"country_code"`
		    } `json:"geo_info"`
		} `json:"reponse"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}
	
	if res.Response.GeoInfo.CountryCode == "FR" {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusNo}
}