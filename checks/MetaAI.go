package mediaunlocktest

import (
	"net/http"
	"regexp"
	"strings"
)

func MetaAI(c http.Client) Result {
	res := CheckGETStatus(c, "https://www.meta.ai/ajax", ResultMap{
		200: {Status: StatusNo},
		400: {Status: StatusOK},
		404: {Status: StatusOK},
	}, Result{Status: StatusUnexpected})

	if res.Status != StatusOK {
		return res
	}

	resp, err := GET(c, "https://www.meta.com/legal/")
	if err != nil {
		return res
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		loc, err := resp.Location()
		if err == nil {
			path := loc.Path
			// match /xx/legal
			if len(path) >= 9 && path[3] == '/' && path[4:] == "legal/" {
				res.Region = strings.ToUpper(path[1:3])
				return res
			}
			// regex fallback or stricter check
			re := regexp.MustCompile(`^/([a-z]{2})/legal`)
			matches := re.FindStringSubmatch(path)
			if len(matches) > 1 {
				res.Region = strings.ToUpper(matches[1])
				return res
			}
		}
	}

	return res
}
