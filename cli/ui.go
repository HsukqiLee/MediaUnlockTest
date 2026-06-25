package main

import (
	core "MediaUnlockTest/pkg/core"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/term"
)

func ShowSingleResult(r core.Result) (s string) {
	switch r.Status {
	case core.StatusOK:
		s = core.Green("YES")
		if r.Region != "" {
			s += core.Green(" (Region: " + strings.ToUpper(r.Region) + ")")
		}
		if Debug && r.CachedResult {
			s += " (Cached)"
		}
		return s

	case core.StatusNetworkErr:
		s = core.Red("ERR")
		if Debug {
			s += core.Yellow(" (Network Err: " + r.Err.Error() + ")")
		} else {
			s += core.Yellow(" (Network Err)")
		}
		if Debug && r.CachedResult {
			s += " (Cached)"
		}
		return s

	case core.StatusRestricted:
		s = core.Yellow("Restricted")
		if r.Info != "" {
			s = core.Yellow("Restricted (" + r.Info + ")")
		}
		if Debug && r.CachedResult {
			s += " (Cached)"
		}
		return s

	case core.StatusErr:
		s = core.Red("ERR")
		if r.Err != nil && Debug {
			s += core.Yellow(" (Err: " + r.Err.Error() + ")")
		}
		if Debug && r.CachedResult {
			s += " (Cached)"
		}
		return s

	case core.StatusNo:
		if r.Info != "" {
			return core.Red("NO ") + core.Yellow(" (Info: "+r.Info+")")
		}
		if r.Region != "" {
			return core.Red("NO  (Region: " + strings.ToUpper(r.Region) + ")")
		}
		return core.Red("NO")

	case core.StatusBanned:
		if r.Info != "" {
			return core.Red("Banned") + core.Yellow(" ("+r.Info+")")
		}
		return core.Red("Banned")

	case core.StatusUnexpected:
		return core.Purple("Unexpected")

	case core.StatusFailed:
		return core.Blue("Failed")

	default:
		return
	}
}

func ShowFinalResult() {
	fmt.Println("测试时间: ", core.Yellow(time.Now().Format("2006-01-02 15:04:05")))

	NameLength := 25
	for _, r := range ResultLines {
		if len(r.Name) > NameLength {
			NameLength = len(r.Name)
		}
	}

	sortedResultLines := ResultLines

	for i, r := range sortedResultLines {
		if r.Divider {

			isRegionGroup := strings.Contains(r.Name, "IPv4") || strings.Contains(r.Name, "IPv6") || strings.Contains(r.Name, "Auto")

			if i > 0 && isRegionGroup {
				fmt.Println()
			}
			s := "[ " + r.Name + " ] "
			for i := NameLength - len(s) + 4; i > 0; i-- {
				s += "="
			}
			if r.Name == "" {
				s = "\n"
			}
			fmt.Println(s)
		} else {
			result := ShowSingleResult(r.Value)
			if r.Value.Status == core.StatusOK && strings.HasSuffix(r.Name, "CDN") {
				result = core.SkyBlue(r.Value.Region)
			}
			fmt.Printf("%-"+strconv.Itoa(NameLength)+"s %s\n", r.Name, result)
		}
	}
}

func newProgressBar(count int64, desc string) *progressbar.ProgressBar {
	width := 30
	if w, _, err := term.GetSize(int(os.Stderr.Fd())); err == nil {

		if w := w - 65; w > 10 {
			width = w
		}
		if width > 80 {
			width = 80
		}
	}
	return progressbar.NewOptions64(count,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(width),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprintln(os.Stderr)
		}),
	)
}

// 新增：更新进度条描述，显示正在进行的测试
func updateProgressBarDescription() {
	if bar == nil {
		return
	}

	if !ShowActive {
		return
	}

	activeTestsMutex.RLock()
	var activeList []string
	for testName, isActive := range activeTests {
		if isActive {
			activeList = append(activeList, testName)
		}
	}
	activeTestsMutex.RUnlock()

	var newDesc string
	if len(activeList) == 0 {
		newDesc = "等待测试开始..."
	} else {

		maxLen := 40
		currentLen := 0
		var displayNames []string

		for _, name := range activeList {

			if currentLen+len(name)+2 > maxLen {
				break
			}
			displayNames = append(displayNames, name)
			currentLen += len(name) + 2
		}

		if len(displayNames) < len(activeList) {
			newDesc = fmt.Sprintf("正在测试: %s 等 %d 个测试", strings.Join(displayNames, ", "), len(activeList))
		} else {
			newDesc = "正在测试: " + strings.Join(displayNames, ", ")
		}
	}

	targetWidth := 45
	descWidth := runewidth.StringWidth(newDesc)

	if descWidth > targetWidth {

		runes := []rune(newDesc)
		width := 0
		for i, r := range runes {
			w := runewidth.RuneWidth(r)
			if width+w > targetWidth-3 {
				newDesc = string(runes[:i]) + "..."
				break
			}
			width += w
		}
	}

	currentWidth := runewidth.StringWidth(newDesc)
	if currentWidth < targetWidth {
		newDesc += strings.Repeat(" ", targetWidth-currentWidth)
	}

	progressDescMu.Lock()
	if progressDescriptionCache != newDesc {
		bar.Describe(newDesc)
		progressDescriptionCache = newDesc
	}
	progressDescMu.Unlock()
}

// 新增：启动进度条更新协程
func startProgressUpdater() {

	if !ShowActive {
		return
	}

	stopProgressUpdater()

	updaterMutex.Lock()
	updaterStopChan = make(chan struct{})
	ch := updaterStopChan
	updaterMutex.Unlock()

	go func() {
		ticker := time.NewTicker(150 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				updateProgressBarDescription()
			case <-ch:
				return
			}
		}
	}()
}

func stopProgressUpdater() {
	updaterMutex.Lock()
	defer updaterMutex.Unlock()
	if updaterStopChan != nil {
		close(updaterStopChan)
		updaterStopChan = nil
	}
}

func ShowCounts() {
	resp, err := http.Get("https://unlock.moe/count.php")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	s := strings.Split(string(b), " ")
	d, m, t := s[0], s[1], s[3]
	fmt.Printf("当天运行共%s次, 本月运行共%s次, 共计运行%s次\n", core.SkyBlue(d), core.Yellow(m), core.Green(t))
}

func ShowAD() {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://unlock.icmp.ing/ad.txt", nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(b))
}
