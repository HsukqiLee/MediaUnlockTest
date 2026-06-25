package mediaunlocktest

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func Mora(c http.Client) Result {
	c.Jar, _ = cookiejar.New(nil)
	resp1, err := PostForm(c, "https://mora.jp/env/regist",
		`returnUrl=`+url.QueryEscape(`/buy?__requestToken=1713764407153&amp;returnUrl=https%3A%2F%2Fmora.jp%2Fpackage%2F43000087%2FTFDS01006B00Z%2F%3Ffmid%3DTOPRNKS%26trackMaterialNo%3D31168909&amp;fromMoraUx=false&amp;deleteMaterial=`)+`&userAgent=`+UA_Browser+`&onTouchend=true`,
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()

	resp2, err := GET(c, "https://mora.jp/buy?__requestToken=1713764407153&returnUrl=https%3A%2F%2Fmora.jp%2Fpackage%2F43000087%2FTFDS01006B00Z%2F%3Ffmid%3DTOPRNKS%26trackMaterialNo%3D31168909&fromMoraUx=false&deleteMaterial=")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	if resp2.StatusCode == 403 || resp2.StatusCode == 500 {
		return Result{Status: StatusNo}
	}

	if resp2.StatusCode == 302 {
		if strings.Contains(resp2.Header.Get("Location"), "error") {
			return Result{Status: StatusNo}
		}
		if strings.Contains(resp2.Header.Get("Location"), "signin") {
			return Result{Status: StatusOK}
		}
	}

	return Result{Status: StatusUnexpected}
}
