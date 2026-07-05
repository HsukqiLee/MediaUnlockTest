package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"strings"
)

func SetantaSports(c core.HttpClient) core.Result {
	resp, err := core.GETWithTimeout(c, "https://dce-frontoffice.imggaming.com/api/v2/consent-prompt", 30,
		core.H{"Realm", "dce.adjara"},
		core.H{"x-api-key", "857a1e5d-e35e-4fdf-805b-a87b6f8364bf"},
	)
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	var data map[string]any
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return core.Result{Status: core.StatusUnexpected}
	}

	result, ok := data["outsideAllowedTerritories"].(bool)
	if !ok {
		return core.Result{Status: core.StatusUnexpected}
	}

	region := ""
	resp2, err := core.GETWithTimeout(c, "https://dce-frontoffice.imggaming.com/api/v3/i18n/country-codes", 30,
		core.H{"Realm", "dce.adjara"},
		core.H{"x-api-key", "857a1e5d-e35e-4fdf-805b-a87b6f8364bf"},
	)
	if err == nil {
		defer resp2.Body.Close()
		var data2 map[string]any
		if json.NewDecoder(resp2.Body).Decode(&data2) == nil {
			if cc, ok := data2["callerCountryCode"].(string); ok {
				region = strings.ToLower(cc)
			}
		}
	}

	if strings.HasPrefix(resp.Status, "200") {
		if result {
			return core.Result{Status: core.StatusNo, Region: region}
		}
		return core.Result{Status: core.StatusOK, Region: region}
	}

	return core.Result{Status: core.StatusUnexpected}
}
