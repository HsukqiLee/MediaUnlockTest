package core

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

func GenUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func GenBase36(n int) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[secureRandInt(len(letters))]
	}
	return string(b)
}

func MD5Sum(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GenRandomStr(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[secureRandInt(len(charset))]
	}
	return string(b)
}

func secureRandInt(max int) int {
	if max <= 0 {
		return 0
	}
	maxBytes := make([]byte, 4)
	_, err := rand.Read(maxBytes)
	if err != nil {
		return 0
	}
	randUint32 := uint32(maxBytes[0])<<24 | uint32(maxBytes[1])<<16 | uint32(maxBytes[2])<<8 | uint32(maxBytes[3])
	return int(randUint32) % max
}

func secureRandInRange(min, max int) int {
	if min >= max {
		return min
	}
	return secureRandInt(max-min+1) + min
}

func addRandomDelay() {
	if secureRandInt(10) == 0 {
		delay := time.Duration(secureRandInRange(50, 149)) * time.Millisecond
		time.Sleep(delay)
	}
}

func GetCloudflareTraceLoc(c HttpClient, url string, headers ...H) (string, error) {
	resp, err := GET(c, url, headers...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	s := string(b)
	_, after, ok := strings.Cut(s, "loc=")
	if !ok {
		return "", errors.New("loc not found in cloudflare trace")
	}
	loc, _, _ := strings.Cut(after, "\n")
	return strings.TrimSpace(loc), nil
}
