package providers

import (
	"MediaUnlockTest/pkg/core"
	"slices"
	"strings"
)

func SupportPoe(loc string) bool {
	var POE_SUPPORT_COUNTRY = []string{
		"AF", "AL", "DZ", "AS", "AD", "AO", "AI", "AG", "AR", "AM", "AW", "AU", "AT", "AZ", "BS",
		"BH", "BD", "BB", "BY", "BE", "BZ", "BJ", "BM", "BT", "BO", "BA", "BW", "BR", "VG", "BN",
		"BG", "BF", "BI", "KH", "CM", "CA", "CV", "KY", "CF", "TD", "IO", "CL", "CN", "CX", "CC",
		"CO", "KM", "CG", "CD", "CK", "CR", "CI", "HR", "CU", "CY", "CZ", "DK", "DJ", "DM", "DO",
		"EC", "EG", "SV", "GQ", "ER", "EE", "SZ", "ET", "FK", "FO", "FJ", "FI", "FR", "GF", "PF",
		"GA", "GM", "GE", "DE", "GH", "GI", "GR", "GL", "GD", "GP", "GU", "GT", "GG", "GN", "GW",
		"GY", "HT", "HN", "HK", "HU", "IS", "IN", "ID", "IR", "IQ", "IE", "IM", "IL", "IT", "JM",
		"JP", "JE", "JO", "KZ", "KE", "KI", "KW", "KG", "LA", "LV", "LB", "LS", "LR", "LY", "LI",
		"LT", "LU", "MO", "MG", "MW", "MY", "MV", "ML", "MT", "MH", "MQ", "MR", "MU", "YT", "MX",
		"FM", "MD", "MC", "MN", "ME", "MS", "MA", "MZ", "MM", "NA", "NR", "NP", "NL", "NC", "NZ",
		"NI", "NE", "NG", "NU", "NF", "KP", "MK", "MP", "NO", "OM", "PK", "PW", "PS", "PA", "PG",
		"PY", "PE", "PH", "PN", "PL", "PT", "PR", "QA", "RE", "RO", "RU", "RW", "WS", "SM", "ST",
		"SA", "SN", "RS", "SC", "SL", "SG", "SK", "SI", "GS", "SB", "SO", "ZA", "KR", "ES", "LK",
		"BL", "SH", "KN", "LC", "MF", "PM", "VC", "SD", "SR", "SJ", "SE", "CH", "SY", "TW", "TJ",
		"TZ", "TH", "TL", "TG", "TK", "TO", "TT", "TN", "TR", "TM", "TC", "TV", "VI", "UG", "UA",
		"AE", "GB", "US", "UY", "UZ", "VU", "VA", "VE", "VN", "WF", "YE", "ZM", "ZW",
	}
	return slices.Contains(POE_SUPPORT_COUNTRY, loc)
}

func Poe(c core.HttpClient) core.Result {
	loc, err := core.GetCloudflareTraceLoc(c, "https://poe.com/cdn-cgi/trace")
	if err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	resp, err := core.GET(c, "https://poe.com/")
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusNo, Region: strings.ToLower(loc)}
	}

	if SupportPoe(loc) {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(loc)}
	}

	return core.Result{Status: core.StatusNo, Region: strings.ToLower(loc)}
}
