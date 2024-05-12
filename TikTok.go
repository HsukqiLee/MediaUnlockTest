package mediaunlocktest

import (
    "io/ioutil"
    "net/http"
    "strings"
    "regexp"
)

func extractRegion(body string) string {
    // Use regular expression to extract region information
    re := regexp.MustCompile(`"region":\s*"([^"]+)"`)
    match := re.FindStringSubmatch(body)
    if len(match) > 0 {
        return match[1]
    }

    // If direct match fails, try to extract from a different pattern
    re = regexp.MustCompile(`"region"\s*:\s*"([^"]+)"`)
    match = re.FindStringSubmatch(body)
    if len(match) > 0 {
        return match[1]
    }

    return ""
}

func TikTok(c http.Client) Result {
    resp, err := GET(c, "https://www.tiktok.com/")
    if err != nil {
        return Result{Status: StatusNetworkErr}
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return Result{Status: StatusFailed}
    }

    // Extract region information from response
    if region := extractRegion(string(body)); region != "" {
        return Result{Status: StatusOK, Region: strings.ToLower(region)}
    }

    return Result{Status: StatusNo}
}


