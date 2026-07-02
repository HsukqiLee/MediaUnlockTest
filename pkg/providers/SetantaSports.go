package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"strings"
)

func SetantaSports(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://dce-frontoffice.imggaming.com/api/v2/consent-prompt",
		core.H{"Realm", "dce.adjara"},
		core.H{"x-api-key", "857a1e5d-e35e-4fdf-805b-a87b6f8364bf"},
	)
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
