package mediaunlocktest

import (
	"net/http"
	"encoding/json"
	"io"
)

func Molotov(c http.Client) Result {
	resp, err := GET(c, "https://fapi.molotov.tv/v1/open-europe/is-france")
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()
    
    b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	
    var res struct {
		isFrance bool `json:"is_france"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		// log.Println(err)
		return Result{Status: StatusFailed, Err: err}
	}
	
	if res.isFrance {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusNo}
}