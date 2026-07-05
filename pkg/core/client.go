package core

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type HttpClient = tls_client.HttpClient

var (
	Ipv4HttpClient HttpClient
	Ipv6HttpClient HttpClient
	AutoHttpClient HttpClient
	SocksProxy     string
	HTTPProxy      string
	DNSServers     string
	Dialer         = &net.Dialer{
		Timeout:   15 * time.Second,
		KeepAlive: 30 * time.Second,
	}
)

func buildClientOptions(disableIPv4, disableIPv6 bool) []tls_client.HttpClientOption {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_146),
		tls_client.WithCustomRedirectFunc(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}),
	}
	if disableIPv4 {
		options = append(options, tls_client.WithDisableIPV4())
	}
	if disableIPv6 {
		options = append(options, tls_client.WithDisableIPV6())
	}
	if SocksProxy != "" {
		options = append(options, tls_client.WithProxyUrl(SocksProxy))
	} else if HTTPProxy != "" {
		options = append(options, tls_client.WithProxyUrl(HTTPProxy))
	}
	if DNSServers != "" {
		Dialer.Resolver = &net.Resolver{
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "udp", DNSServers)
			},
		}
	}
	options = append(options, tls_client.WithDialer(*Dialer))
	return options
}

func InitClients() {
	var err error
	Ipv4HttpClient, err = tls_client.NewHttpClient(tls_client.NewNoopLogger(), buildClientOptions(false, true)...)
	if err != nil {
		panic(err)
	}

	Ipv6HttpClient, err = tls_client.NewHttpClient(tls_client.NewNoopLogger(), buildClientOptions(true, false)...)
	if err != nil {
		panic(err)
	}

	AutoHttpClient, err = tls_client.NewHttpClient(tls_client.NewNoopLogger(), buildClientOptions(false, false)...)
	if err != nil {
		panic(err)
	}
}

func NewHttpClient(ipType int) HttpClient {
	var disableIPv4, disableIPv6 bool
	switch ipType {
	case 4:
		disableIPv6 = true
	case 6:
		disableIPv4 = true
	}
	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), buildClientOptions(disableIPv4, disableIPv6)...)
	return client
}

type H [2]string

func doRequest(c HttpClient, method, url string, reqType string, body string, useRealisticHeaders bool, timeout int, headers ...H) (*http.Response, error) {
	var req *http.Request
	var err error

	ctx := context.Background()
	var cancel context.CancelFunc
	if timeout <= 0 {
		timeout = 6 // 默认 6 秒超时
	}
	ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)

	if body != "" {
		req, err = http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}
	if err != nil {
		cancel()
		return nil, err
	}

	switch reqType {
	case "json":
		req.Header.Set("content-type", "application/json")
		req.Header[http.HeaderOrderKey] = append(req.Header[http.HeaderOrderKey], "content-type")
	case "form":
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
		req.Header[http.HeaderOrderKey] = append(req.Header[http.HeaderOrderKey], "content-type")
	}

	if useRealisticHeaders {
		setRealisticHeaders(req, reqType)
	}

	for _, h := range headers {
		req.Header.Set(h[0], h[1])
		if existing, ok := req.Header[http.HeaderOrderKey]; ok {
			req.Header[http.HeaderOrderKey] = append(existing, h[0])
		} else {
			req.Header[http.HeaderOrderKey] = []string{h[0]}
		}
	}

	if GlobalLogLevel >= LevelInfo {
		LogInfo("Request: %s %s", req.Method, req.URL.String())
		LogInfo("Request Headers: %v", req.Header)
		if body != "" {
			LogInfo("Request Body: %s", body)
		}
	}

	addRandomDelay()
	resp, err := DoWithRetry(c, req)

	if GlobalLogLevel >= LevelInfo {
		if err != nil {
			LogError("Response Error for %s: %v", req.URL.String(), err)
		} else if resp != nil {
			LogInfo("Response Status: %s (%d) for %s", resp.Status, resp.StatusCode, req.URL.String())
			LogInfo("Response Headers: %v", resp.Header)

			if resp.Body != nil {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err == nil {
					resp.Body.Close()
					resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					bodyStr := string(bodyBytes)
					if len(bodyStr) > 512 {
						LogInfo("Response Body: %s... (truncated)", bodyStr[:512])
					} else {
						LogInfo("Response Body: %s", bodyStr)
					}
				}
			}
		}
	}

	if resp != nil && resp.Body != nil {
		resp.Body = &cancelTimerBody{ReadCloser: resp.Body, cancel: cancel}
	} else {
		cancel()
	}

	return resp, err
}

// GET performs a GET request with realistic default headers
func GET(c HttpClient, url string, headers ...H) (*http.Response, error) {
	return doRequest(c, "GET", url, "html", "", true, 0, headers...)
}

// GETRaw performs a GET request WITHOUT injecting realistic default headers
func GETRaw(c HttpClient, url string, headers ...H) (*http.Response, error) {
	return doRequest(c, "GET", url, "html", "", false, 0, headers...)
}

// RequestRaw performs a generic request WITHOUT injecting realistic default headers
func RequestRaw(c HttpClient, method, url string, body string, headers ...H) (*http.Response, error) {
	return doRequest(c, method, url, "raw", body, false, 0, headers...)
}

func GET_Dalvik(c HttpClient, url string) (*http.Response, error) {
	return GETRaw(c, url, H{"User-Agent", UA_Dalvik})
}

var ErrNetwork = errors.New("network error")

func DoWithRetry(c HttpClient, req *http.Request) (resp *http.Response, err error) {
	deadline := time.Now().Add(14 * time.Second)
	for range 2 {
		if time.Now().After(deadline) {
			break
		}
		if resp, err = c.Do(req); err == nil {
			return resp, nil
		}
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) && dnsErr.IsNotFound {
			break
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			break
		}
	}
	return nil, err
}

func PostJson(c HttpClient, url string, data string, headers ...H) (*http.Response, error) {
	return doRequest(c, "POST", url, "json", data, true, 0, headers...)
}

func PostForm(c HttpClient, url string, data string, headers ...H) (*http.Response, error) {
	return doRequest(c, "POST", url, "form", data, true, 0, headers...)
}

// IsWAFBlockError checks if the network error is caused by a WAF drop/timeout
func IsWAFBlockError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, context.DeadlineExceeded) || os.IsTimeout(err) ||
		errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNABORTED) {
		return true
	}
	// utls 或 http2 等底层依赖有时只返回字符串错误，且没有导出 Error 变量
	errStr := err.Error()
	if strings.Contains(errStr, "stream error") || strings.Contains(errStr, "handshake failure") || strings.Contains(errStr, "connection reset") {
		return true
	}
	return false
}

func GETWithTimeout(c HttpClient, url string, timeout int, headers ...H) (*http.Response, error) {
	return doRequest(c, "GET", url, "html", "", true, timeout, headers...)
}

func PostJsonWithTimeout(c HttpClient, url string, data string, timeout int, headers ...H) (*http.Response, error) {
	return doRequest(c, "POST", url, "json", data, true, timeout, headers...)
}

type cancelTimerBody struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (b *cancelTimerBody) Close() error {
	b.cancel()
	return b.ReadCloser.Close()
}
