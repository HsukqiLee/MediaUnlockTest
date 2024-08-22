package mediaunlocktest

import (
    "io"
    "regexp"
	"net/http"
	"strings"
)

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

    bodyBytes, err := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    
    if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	regionThreeCode := extractGeminiRegion(bodyString)
	regionTwoCode := threeToTwoCode(regionThreeCode)
	
	if regionThreeCode != "" && regionTwoCode != "" {
	    if !strings.Contains(bodyString, "45617354,null,true") {
		    return Result{Status: StatusNo, Region: strings.ToLower(regionTwoCode)}
	    }
	    return Result{Status: StatusOK, Region: strings.ToLower(regionTwoCode)}
	}

	return Result{Status: StatusUnexpected}
}
