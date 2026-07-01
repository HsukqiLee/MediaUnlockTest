package core

import (
	"context"
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"io"
	"net"
	"strings"
	"time"
)

type IPInfo struct {
	IP           string `json:"ip"`
	Country      string `json:"country"`
	Region       string `json:"region"`
	City         string `json:"city"`
	Timezone     string `json:"timezone"`
	ASN          int    `json:"asn"`
	Organization string `json:"organization"`
}

func GetDetailedIPInfo(url string, ipType int) (*IPInfo, error) {
	timeout := 6
	if ipType == 6 {
		timeout = 3
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var client HttpClient
	switch ipType {
	case 6:
		client = Ipv6HttpClient
	case 4:
		client = Ipv4HttpClient
	case 0:
		client = AutoHttpClient
	default:
		return nil, fmt.Errorf("IP type %d is invalid", ipType)
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info IPInfo
	if err := json.Unmarshal(b, &info); err != nil {
		return nil, err
	}
	if net.ParseIP(info.IP) == nil {
		return nil, fmt.Errorf("invalid IP returned: %s", info.IP)
	}
	return &info, nil
}

func GetIPInfo(url string, ipType int, formatType string) (string, error) {
	timeout := 6
	if ipType == 6 {
		timeout = 3
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var client HttpClient
	switch ipType {
	case 6:
		client = Ipv6HttpClient
	case 4:
		client = Ipv4HttpClient
	case 0:
		client = AutoHttpClient
	default:
		return "", fmt.Errorf("IP type %d is invalid", ipType)
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("dnt", "1")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("sec-ch-ua", `"Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	switch formatType {
	case "cloudflare":
		s := string(b)
		i := strings.Index(s, "ip=")
		if i == -1 {
			return "", fmt.Errorf("ip not found in cloudflare trace")
		}
		s = s[i+3:]
		i = strings.Index(s, "\n")
		if i != -1 {
			s = s[:i]
		}
		ip := strings.TrimSpace(s)
		if net.ParseIP(ip) == nil {
			return "", fmt.Errorf("invalid IP from cloudflare: %s", ip)
		}
		return ip, nil

	default:
		ip := strings.TrimSpace(string(b))
		if net.ParseIP(ip) == nil {
			return "", fmt.Errorf("invalid IP format: %s", ip)
		}
		return ip, nil
	}
}
