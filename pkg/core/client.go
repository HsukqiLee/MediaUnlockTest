package core

import (
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
		tls_client.WithTimeoutSeconds(6),
		tls_client.WithClientProfile(profiles.Chrome_120),
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

func doRequest(c HttpClient, method, url string, reqType string, body string, useRealisticHeaders bool, headers ...H) (*http.Response, error) {
	var req *http.Request
	var err error
	if body != "" {
		req, err = http.NewRequest(method, url, strings.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	switch reqType {
	case "json":
		req.Header.Set("content-type", "application/json")
	case "form":
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
	}

	if useRealisticHeaders {
		setRealisticHeaders(req, reqType)
	}

	for _, h := range headers {
		req.Header.Set(h[0], h[1])
	}

	addRandomDelay()
	return DoWithRetry(c, req)
}

// GET performs a GET request with realistic default headers
func GET(c HttpClient, url string, headers ...H) (*http.Response, error) {
	return doRequest(c, "GET", url, "html", "", true, headers...)
}

// GETRaw performs a GET request WITHOUT injecting realistic default headers
func GETRaw(c HttpClient, url string, headers ...H) (*http.Response, error) {
	return doRequest(c, "GET", url, "html", "", false, headers...)
}

// RequestRaw performs a generic request WITHOUT injecting realistic default headers
func RequestRaw(c HttpClient, method, url string, body string, headers ...H) (*http.Response, error) {
	return doRequest(c, method, url, "raw", body, false, headers...)
}

func GET_Dalvik(c HttpClient, url string) (*http.Response, error) {
	return GETRaw(c, url, H{"User-Agent", UA_Dalvik})
}

var ErrNetwork = errors.New("network error")

func DoWithRetry(c HttpClient, req *http.Request) (resp *http.Response, err error) {
	deadline := time.Now().Add(14 * time.Second)
	for i := 0; i < 2; i++ {
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
	return doRequest(c, "POST", url, "json", data, true, headers...)
}

func PostForm(c HttpClient, url string, data string, headers ...H) (*http.Response, error) {
	return doRequest(c, "POST", url, "form", data, true, headers...)
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
