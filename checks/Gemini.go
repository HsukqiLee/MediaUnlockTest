package mediaunlocktest

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

func SupportGemini(loc string) bool {
	var GEMINI_SUPPORT_COUNTRY = []string{
		"AX", "AL", "DZ", "AS", "AD", "AO", "AI", "AQ", "AG", "AR",
		"AM", "AW", "AU", "AT", "AZ", "BH", "BD", "BB", "BE", "BZ",
		"BJ", "BM", "BT", "BO", "BA", "BW", "BR", "IO", "VG", "BN",
		"BG", "BF", "BI", "CV", "KH", "CM", "CA", "BQ", "KY", "CF",
		"TD", "CL", "CX", "CC", "CO", "KM", "CK", "CR", "CI", "HR",
		"CW", "CZ", "CD", "DK", "DJ", "DM", "DO", "EC", "EG", "SV",
		"GQ", "ER", "EE", "SZ", "ET", "FK", "FO", "FJ", "FI", "FR",
		"GF", "PF", "TF", "GA", "GE", "DE", "GH", "GI", "GR", "GL",
		"GD", "GP", "GU", "GT", "GG", "GN", "GW", "GY", "HT", "HM",
		"HN", "HU", "IS", "IN", "ID", "IQ", "IE", "IM", "IL", "IT",
		"JM", "JP", "JE", "JO", "KZ", "KE", "KI", "XK", "KW", "KG",
		"LA", "LV", "LB", "LS", "LR", "LY", "LI", "LT", "LU", "MG",
		"MW", "MY", "MV", "ML", "MT", "MH", "MQ", "MR", "MU", "YT",
		"MX", "FM", "MD", "MC", "MN", "ME", "MS", "MA", "MZ", "MM",
		"NA", "NR", "NP", "NL", "NC", "NZ", "NI", "NE", "NG", "NU",
		"NF", "MK", "MP", "NO", "OM", "PK", "PW", "PS", "PA", "PG",
		"PY", "PE", "PH", "PN", "PL", "PT", "PR", "QA", "CY", "CG",
		"RE", "RO", "RW", "BL", "SH", "KN", "LC", "MF", "PM", "VC",
		"WS", "SM", "ST", "SA", "SN", "RS", "SC", "SL", "SG", "SX",
		"SK", "SI", "SB", "SO", "ZA", "GS", "KR", "SS", "ES", "LK",
		"SD", "SR", "SJ", "SE", "CH", "TW", "TJ", "TZ", "TH", "BS",
		"GM", "TL", "TG", "TK", "TO", "TT", "TN", "TR", "TM", "TC",
		"TV", "VI", "UG", "UA", "AE", "GB", "US", "UM", "UY", "UZ",
		"VU", "VA", "VE", "VN", "WF", "EH", "YE", "ZM", "ZW",
	}
	for _, s := range GEMINI_SUPPORT_COUNTRY {
		if loc == s {
			return true
		}
	}
	return false
}

func extractGeminiRegion(body string) string {
	regex := regexp.MustCompile(`,2,1,200,"([A-Z]{3})"`)
	matches := regex.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func Gemini(c http.Client) Result {
	resp, err := GET(c, "https://gemini.google.com/?hl=en",
		H{"Cookie", "SOCS=CAISOAgeEitib3FfaWRlbnRpdHlmcm9udGVuZHVpc2VydmVyXzIwMjQwODExLjA4X3AwGgV6aC1DTiACGgYIgOfvtQY; NID=516=I_1jljreYQcLJzkCvUL8k6eIiBqaC8_mkGLq8XaPbFdY2MMgL9re9rZw7NC9scKoZJyI3pzxTfTqhklea0KVsKRKeRMyl-TOLZgzpb_rqoL2vH7J2_5Wlef7KTXVuLfArRqcGeFWF4OZ6HBfQqu7BQc_YiFfiXshUK1bAp19DZQOQ_nmzgacv-HSMnOG6wPOJsBD7qNXmf4IuQ;__Secure-ENID=delete"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return Result{Status: StatusBanned}
	}

	if resp.StatusCode == 302 {
		return Result{Status: StatusFailed}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	regionThreeCode := extractGeminiRegion(bodyString)
	regionTwoCode := threeToTwoCode(regionThreeCode)

	if regionThreeCode != "" && regionTwoCode != "" {
		if !SupportGemini(regionTwoCode) {
			return Result{Status: StatusNo, Region: strings.ToLower(regionTwoCode)}
		}
		return Result{Status: StatusOK, Region: strings.ToLower(regionTwoCode)}
	}

	return Result{Status: StatusUnexpected}
}
