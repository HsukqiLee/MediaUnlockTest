package mediaunlocktest

import (
    "io"
    //"regexp"
	"net/http"
	"strings"
)

/*func extractMetaAIRegion(body string) string {
    re := regexp.MustCompile(`"code"\s*:\s*"([^"]+)"`)
    matches := re.FindStringSubmatch(body)
    if len(matches) > 1 {
        return matches[1]
    }
    return ""
}*/

func MetaAI(c http.Client) Result {
	resp, err := GET(c, "https://www.meta.ai/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    
    if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	
	if strings.Contains(bodyString, "AbraGeoBlockedErrorRoot") {
		return Result{Status: StatusNo}
	}
	
	if strings.Contains(bodyString, "AbraHomeRootConversationQuery") {
	    /*region := extractMetaAIRegion(bodyString)
	    if region != "" {
		    return Result{Status: StatusOK, Region: strings.ToLower(region)}
	    }*/
	    return Result{Status: StatusOK}
	}
	
	return Result{Status: StatusUnexpected}
}