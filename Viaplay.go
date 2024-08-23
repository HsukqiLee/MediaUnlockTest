package mediaunlocktest

import (
	"net/http"
    "strings"
    "regexp"
)

func extractViaplayRegion1(url string) string {
    re := regexp.MustCompile(`/([a-z]{2})/`)
    matches := re.FindStringSubmatch(string(url))
    if len(matches) > 1 {
        return matches[1]
    }
    return ""
}

func extractViaplayRegion2(url string) string {
    re := regexp.MustCompile(`viaplay.([a-z]{2})`)
    matches := re.FindStringSubmatch(string(url))
    if len(matches) > 1 {
        return matches[1]
    }
    return ""
}

func Viaplay(c http.Client) Result {
	resp1, err := GET(c, "https://checkout.viaplay.pl/?recommended=viaplay")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
    
    if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	if resp1.StatusCode == 403 {
	    return Result{Status: StatusBanned}
	}

	if resp1.StatusCode == 302 && resp1.Header.Get("Location") == "/region-blocked" {
		return Result{Status: StatusNo}
	}
	
	if resp1.StatusCode == 200 {
		resp2, err := GET(c, "https://viaplay.com/")
    	if err != nil {
    		return Result{Status: StatusNetworkErr, Err: err}
    	}
    	defer resp2.Body.Close()
        
        if err != nil {
    		return Result{Status: StatusNetworkErr, Err: err}
    	}
        if resp2.StatusCode == 404 {
            return Result{Status: StatusNo}
        }
    	if resp2.StatusCode == 302 {
    	    if region := extractViaplayRegion1(resp2.Header.Get("Location")); region != "" {
    	        return Result{Status: StatusOK, Region: strings.ToLower(region)}
    	    }
    	    if region := extractViaplayRegion2(resp2.Header.Get("Location")); region != "" {
    	        return Result{Status: StatusOK, Region: region}
    	    }
    	}
	}

	return Result{Status: StatusUnexpected}
}