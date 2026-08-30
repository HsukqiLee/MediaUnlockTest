package main

import (
	core "MediaUnlockTest/pkg/core"
	"reflect"
	"testing"
	"time"
)

func TestTableRowsSeparatesIPVersionsAndOmitsDividers(t *testing.T) {
	results := []*result{
		{Name: "Globe (IPv4)", Divider: true},
		{Name: "Netflix", IPVersion: 4, Value: core.Result{Status: core.StatusOK, Region: "us"}},
		{Name: "AI", Divider: true},
		{Name: "Gemini", IPVersion: 4, Value: core.Result{Status: core.StatusNo}},
		{Name: "Globe (IPv6)", Divider: true},
		{Name: "Netflix", IPVersion: 6, Value: core.Result{Status: core.StatusNetworkErr}},
		{Name: "No Korean platform supports IPv6"},
	}

	headings, values := tableRows(results, 4)
	if want := []string{"Netflix", "Gemini"}; !reflect.DeepEqual(headings, want) {
		t.Fatalf("headings = %#v, want %#v", headings, want)
	}
	if want := []string{"US", "NO"}; !reflect.DeepEqual(values, want) {
		t.Fatalf("values = %#v, want %#v", values, want)
	}

	headings, values = tableRows(results, 6)
	if want := []string{"Netflix"}; !reflect.DeepEqual(headings, want) {
		t.Fatalf("IPv6 headings = %#v, want %#v", headings, want)
	}
	if want := []string{"ERR"}; !reflect.DeepEqual(values, want) {
		t.Fatalf("IPv6 values = %#v, want %#v", values, want)
	}
}

func TestTableRowsSanitizesCells(t *testing.T) {
	results := []*result{{
		Name:      "Service Name Extra",
		IPVersion: 4,
		Value:     core.Result{Status: core.StatusRestricted, Info: "plan|only\nupgrade"},
	}}

	headings, values := tableRows(results, 4)
	if want := []string{"Service Name Extra"}; !reflect.DeepEqual(headings, want) {
		t.Fatalf("headings = %#v, want %#v", headings, want)
	}
	if want := []string{"RESTRICTED"}; !reflect.DeepEqual(values, want) {
		t.Fatalf("values = %#v, want %#v", values, want)
	}
}

func TestTableResultValuePreservesStatuses(t *testing.T) {
	tests := map[int]string{
		core.StatusOK:         "YES",
		core.StatusRestricted: "RESTRICTED",
		core.StatusNo:         "NO",
		core.StatusBanned:     "BANNED",
		core.StatusNetworkErr: "ERR",
		core.StatusErr:        "ERR",
		core.StatusUnexpected: "ERR",
		core.StatusFailed:     "FAILED",
	}
	for status, want := range tests {
		if got := tableResultValue(core.Result{Status: status}); got != want {
			t.Errorf("tableResultValue(%d) = %q, want %q", status, got, want)
		}
	}
}

func TestServiceResultTableHasSixAlignedColumns(t *testing.T) {
	services := []string{"Amazon", "Apple", "Bilibili", "Netflix"}
	results := []string{"SG", "US", "NO", "JP"}
	rows := serviceResultTable(services, results)
	widths := serviceResultWidths(rows)

	if got, want := tableBorder(widths, "┌", "┬", "┐"), "┌──────────┬────────┬──────────┬────────┬──────────┬────────┐"; got != want {
		t.Fatalf("top border = %q, want %q", got, want)
	}
	if got, want := tableRow(rows[0], widths), "│ Service  │ Result │ Service  │ Result │ Service  │ Result │"; got != want {
		t.Fatalf("header row = %q, want %q", got, want)
	}
	if got, want := tableRow(rows[1], widths), "│  Amazon  │   SG   │  Apple   │   US   │ Bilibili │   NO   │"; got != want {
		t.Fatalf("first data row = %q, want %q", got, want)
	}
	if got, want := tableRow(rows[2], widths), "│ Netflix  │   JP   │          │        │          │        │"; got != want {
		t.Fatalf("padded final row = %q, want %q", got, want)
	}
}

func TestPerformanceStatsUsesOneLine(t *testing.T) {
	got := performanceStats(13, 2780*time.Millisecond)
	want := "性能统计: 总测试数量: 13 | 总耗时: 2.78s | 平均每个测试耗时: 213.85ms | 测试速度: 4.68 测试/秒"
	if got != want {
		t.Fatalf("performanceStats() = %q, want %q", got, want)
	}
}

func TestExecutionStatsCombinesIPRuns(t *testing.T) {
	stats := executionStats{}
	stats.add(executionStats{TotalTests: 21, TotalDuration: 2420 * time.Millisecond})
	stats.add(executionStats{TotalTests: 13, TotalDuration: 4370 * time.Millisecond})

	if stats.TotalTests != 34 {
		t.Fatalf("combined tests = %d, want 34", stats.TotalTests)
	}
	if stats.TotalDuration != 6790*time.Millisecond {
		t.Fatalf("combined duration = %v, want 6.79s", stats.TotalDuration)
	}
}

func TestCompactServiceName(t *testing.T) {
	tests := map[string]string{
		"Amazon Prime Video": "Amazon",
		"Google Play Store":  "Google Play",
		"Netflix":            "Netflix",
		"Netflix CDN":        "Netflix CDN",
		"Youtube CDN":        "YouTube CDN",
		"Youtube Premium":    "YouTube Premium",
	}
	for input, want := range tests {
		if got := compactServiceName(input); got != want {
			t.Errorf("compactServiceName(%q) = %q, want %q", input, got, want)
		}
	}
}
