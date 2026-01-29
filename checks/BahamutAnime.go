package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func BahamutAnime(c http.Client) Result {
	// 1. 初始化独立的 CookieJar，确保 Cookie 隔离
	c.Jar, _ = cookiejar.New(nil)

	// 2. 定义一个内部安全的 GET 函数
	// 作用：显式设置 User-Agent，绕过 main.go 中 ResetSessionHeaders 导致的全局 Header 竞争问题
	safeGET := func(url string) (*http.Response, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		// 强制硬编码 User-Agent，防止被其他并发协程清空
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		return c.Do(req)
	}

	type apiResponse struct {
		AnimeSn  int    `json:"animeSn"`
		Deviceid string `json:"deviceid"`
	}

	// 步骤 1: 获取 Device ID
	// 使用 safeGET 替代原本的 GET
	resp1, err := safeGET("https://ani.gamer.com.tw/ajax/getdeviceid.php")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	b1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var res1 apiResponse
	if err := json.Unmarshal(b1, &res1); err != nil {
		// 如果解析失败，通常是因为被 WAF 拦截返回了 HTML，说明 IP 不行或 UA 掉了
		if err.Error() == "invalid character '<' looking for beginning of value" {
			return Result{Status: StatusNo}
		}
		return Result{Status: StatusErr, Err: err}
	}

	// 步骤 2: 测试 API (Token)
	resp2, err := safeGET("https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=37783&device=" + res1.Deviceid)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var res2 apiResponse
	if err := json.Unmarshal(b2, &res2); err != nil {
		return Result{Status: StatusErr, Err: err}
	}

	// 步骤 3: 进一步验证
	if res2.AnimeSn != 0 {
		resp3, err := safeGET("https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=38832&device=" + res1.Deviceid)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		defer resp3.Body.Close()
		b3, err := io.ReadAll(resp3.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}

		var res3 apiResponse
		if err := json.Unmarshal(b3, &res3); err != nil {
			return Result{Status: StatusErr, Err: err}
		}

		if res3.AnimeSn != 0 {
			return Result{Status: StatusOK, Region: "tw"}
		}

		// 步骤 4: 通过 Trace 确定地区 (作为兜底)
		resp4, err := safeGET("https://ani.gamer.com.tw/cdn-cgi/trace")
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		defer resp4.Body.Close()
		b4, err := io.ReadAll(resp4.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}

		bodyString := string(b4)
		index := strings.Index(bodyString, "loc=")
		if index == -1 {
			// 如果找不到 loc=，说明返回的不是 trace 信息（可能是 WAF 页面），此时返回 Unexpected 是合理的
			return Result{Status: StatusUnexpected}
		}
		bodyString = bodyString[index+4:]
		index = strings.Index(bodyString, "\n")
		if index == -1 {
			return Result{Status: StatusUnexpected}
		}
		loc := bodyString[:index]
		if len(loc) == 2 {
			return Result{Status: StatusOK, Region: strings.ToLower(loc)}
		}
	}

	// 如果所有尝试都未能确认解锁，且没有报错，则视为 Unexpected
	return Result{Status: StatusUnexpected}
}
