package core

import (
	"encoding/json"
	"io"
)

var (
	StatusOK         = 1
	StatusNetworkErr = -1
	StatusErr        = -2
	StatusRestricted = 2
	StatusNo         = 3
	StatusBanned     = 4
	StatusFailed     = 5
	StatusUnexpected = 6
)

type Result struct {
	Status       int
	Region       string
	Info         string
	Err          error
	CachedResult bool
}

// ResultMap 支持完整 Result 值的映射
type ResultMap map[int]Result

// ResultFromMapping 根据 statusCode 从映射获得 Result，缺省返回 defaultRes[0] 或 StatusUnexpected
func ResultFromMapping(statusCode int, m ResultMap, defaultResult Result) Result {
	if r, ok := m[statusCode]; ok {
		return r
	}
	return defaultResult
}

// CheckGETStatus 使用 GET 请求，并通过 ResultMap 返回对应 Result，支持默认 Result 及可选 headers
func CheckGETStatus(c HttpClient, url string, mapping ResultMap, defaultResult Result, headers ...H) Result {
	resp, err := GET(c, url, headers...)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	return ResultFromMapping(resp.StatusCode, mapping, defaultResult)
}

// CheckDalvikStatus 使用 GET_Dalvik 请求，并通过 ResultMap 返回对应 Result，支持默认 Result
func CheckDalvikStatus(c HttpClient, url string, mapping ResultMap, defaultResult Result) Result {
	resp, err := GET_Dalvik(c, url)
	if err != nil {
		if IsWAFBlockError(err) {
			return Result{Status: StatusBanned}
		}
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	return ResultFromMapping(resp.StatusCode, mapping, defaultResult)
}

func PostFormBoolSuccess(c HttpClient, url, data string, headers ...H) (bool, error) {
	resp, err := PostForm(c, url, data, headers...)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var res struct{ Success bool }
	if err := json.Unmarshal(b, &res); err != nil {
		return false, err
	}
	return res.Success, nil
}

// CheckPostFormStatus 使用 POST 表单请求，并通过 ResultMap 返回对应 Result，支持默认 Result 及可选 headers
func CheckPostFormStatus(c HttpClient, url, data string, mapping ResultMap, defaultResult Result, headers ...H) Result {
	resp, err := PostForm(c, url, data, headers...)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	return ResultFromMapping(resp.StatusCode, mapping, defaultResult)
}

// CheckPostJsonStatus 使用 POST JSON 请求，并通过 ResultMap 返回对应 Result，支持默认 Result 及可选 headers
func CheckPostJsonStatus(c HttpClient, url, data string, mapping ResultMap, defaultResult Result, headers ...H) Result {
	resp, err := PostJson(c, url, data, headers...)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	return ResultFromMapping(resp.StatusCode, mapping, defaultResult)
}

func CheckStatus(c HttpClient, url string, mapping ResultMap, defaultResult Result) Result {
	resp, err := GET(c, url)
	if err != nil {
		if IsWAFBlockError(err) {
			return Result{Status: StatusBanned}
		}
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	return ResultFromMapping(resp.StatusCode, mapping, defaultResult)
}

func CheckStatusWithTimeout(c HttpClient, url string, mapping ResultMap, defaultResult Result, timeout int) Result {
	resp, err := GETWithTimeout(c, url, timeout)
	if err != nil {
		if IsWAFBlockError(err) {
			return Result{Status: StatusBanned, Info: "WAF Timeout"}
		}
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	return ResultFromMapping(resp.StatusCode, mapping, defaultResult)
}
