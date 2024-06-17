package mediaunlocktest

import (
	"net/http"
	"regexp"
)

func extractDirecTVGORegion (url string) string {
    re := regexp.MustCompile(`https?://www\.directvgo\.com/([^/]+)/`)

    matches := re.FindStringSubmatch(url)
    if len(matches) > 1 {
        return matches[1]
    }
    return ""
}

func DirecTVGO(c http.Client) Result {
	resp, err := GET(c, "https://www.directvgo.com/registrarse")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
    
    if err != nil {
		return Result{Status: StatusFailed}
	}
	
	if resp.StatusCode == 403 {
	    return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 301 {
	    if region := extractDirecTVGORegion(resp.Header.Get("Location")); region != ""{
	        return Result{Status: StatusOK, Region: region}
	    }
		return Result{Status: StatusUnexpected}
	}
	
	return Result{Status: StatusUnexpected}
}