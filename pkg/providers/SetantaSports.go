package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"net/http"
	"strings"
)

func SetantaSports(c http.Client) core.Result {
	req, err := http.NewRequest("GET", "https://dce-frontoffice.imggaming.com/api/v2/consent-prompt", nil)
	if err != nil {
		return core.Result{Status: core.StatusFailed}
	}
	req.Header.Set("Realm", "dce.adjara")
	req.Header.Set("x-api-key", "857a1e5d-e35e-4fdf-805b-a87b6f8364bf")

	resp, err := core.Cdo(c, req)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return core.Result{Status: core.StatusUnexpected}
	}

	result, ok := data["outsideAllowedTerritories"].(bool)
	if !ok {
		return core.Result{Status: core.StatusUnexpected}
	}

	if strings.HasPrefix(resp.Status, "200") {
		if result {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}


