package mediaunlocktest

import (
	"encoding/json"
	"net/http"
	"strings"
)

func SetantaSports(c http.Client) Result {
	req, err := http.NewRequest("GET", "https://dce-frontoffice.imggaming.com/api/v2/consent-prompt", nil)
	if err != nil {
		return Result{Status: StatusFailed}
	}
	req.Header.Set("Realm", "dce.adjara")
	req.Header.Set("x-api-key", "857a1e5d-e35e-4fdf-805b-a87b6f8364bf")

	resp, err := cdo(c, req)
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()
	
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return Result{Status: StatusUnexpected}
	}

	result, ok := data["outsideAllowedTerritories"].(bool)
	if !ok {
		return Result{Status: StatusUnexpected}
	}

	if strings.HasPrefix(resp.Status, "200") {
		if result {
			return Result{Status: StatusNo}
		}
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}

