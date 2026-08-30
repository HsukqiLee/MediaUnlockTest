package main

import (
	core "MediaUnlockTest/pkg/core"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
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

// ShowTableResult prints one six-column service/result table for each tested IP
// version. Region and subgroup dividers are intentionally omitted.
func ShowTableResult() {
	for _, ipVersion := range []int{4, 6, 0} {
		headings, values := tableRows(ResultLines, ipVersion)
		if len(headings) == 0 {
			continue
		}

		label := fmt.Sprintf("IPv%d", ipVersion)
		if ipVersion == 0 {
			label = "Auto"
		}
		fmt.Println(label + ":")
		rows := serviceResultTable(headings, values)
		widths := serviceResultWidths(rows)
		fmt.Println(tableBorder(widths, "┌", "┬", "┐"))
		for i, row := range rows {
			fmt.Println(tableRow(row, widths))
			if i == len(rows)-1 {
				fmt.Println(tableBorder(widths, "└", "┴", "┘"))
			} else {
				fmt.Println(tableBorder(widths, "├", "┼", "┤"))
			}
		}
	}
}

func tableRows(results []*result, ipVersion int) ([]string, []string) {
	var headings []string
	var values []string
	for _, r := range results {
		if r.Divider || r.IPVersion != ipVersion || r.Value.Status == 0 {
			continue
		}
		headings = append(headings, tableCell(compactServiceName(r.Name)))
		values = append(values, tableCell(tableResultValue(r.Value)))
	}
	return headings, values
}

func compactServiceName(name string) string {
	switch {
	case strings.EqualFold(name, "Amazon Prime Video"):
		return "Amazon"
	case strings.EqualFold(name, "Google Play Store"):
		return "Google Play"
	case strings.EqualFold(name, "Spotify Registration"):
		return "Spotify"
	case strings.EqualFold(name, "Wikipedia Editability"):
		return "Wikipedia"
	case strings.EqualFold(name, "Youtube CDN"):
		return "YouTube CDN"
	case strings.EqualFold(name, "Youtube Premium"):
		return "YouTube Premium"
	}
	return name
}

func serviceResultTable(services, results []string) [][]string {
	const pairsPerRow = 3
	rows := [][]string{{"Service", "Result", "Service", "Result", "Service", "Result"}}
	for start := 0; start < len(services); start += pairsPerRow {
		row := make([]string, pairsPerRow*2)
		for pair := 0; pair < pairsPerRow; pair++ {
			index := start + pair
			if index >= len(services) {
				break
			}
			row[pair*2] = services[index]
			row[pair*2+1] = results[index]
		}
		rows = append(rows, row)
	}
	return rows
}

func serviceResultWidths(rows [][]string) []int {
	const columnCount = 6
	serviceWidth := tableCellWidth("Service")
	resultWidth := tableCellWidth("Result")
	for _, row := range rows {
		for column, cell := range row {
			if column%2 == 0 {
				serviceWidth = max(serviceWidth, tableCellWidth(cell))
			} else {
				resultWidth = max(resultWidth, tableCellWidth(cell))
			}
		}
	}
	widths := make([]int, columnCount)
	for column := range widths {
		if column%2 == 0 {
			widths[column] = serviceWidth
		} else {
			widths[column] = resultWidth
		}
	}
	return widths
}

func tableRow(cells []string, widths []int) string {
	padded := make([]string, len(cells))
	for i, cell := range cells {
		padding := widths[i] - tableCellWidth(cell)
		leftPadding := padding / 2
		rightPadding := padding - leftPadding
		padded[i] = strings.Repeat(" ", leftPadding) + cell + strings.Repeat(" ", rightPadding)
	}
	return "│ " + strings.Join(padded, " │ ") + " │"
}

func tableBorder(widths []int, left, middle, right string) string {
	sections := make([]string, len(widths))
	for i, width := range widths {
		sections[i] = strings.Repeat("─", width+2)
	}
	return left + strings.Join(sections, middle) + right
}

var ansiColorPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func tableCellWidth(value string) int {
	return runewidth.StringWidth(ansiColorPattern.ReplaceAllString(value, ""))
}

func tableCell(value string) string {
	value = strings.ReplaceAll(value, "\r", " ")
	return strings.ReplaceAll(value, "\n", " ")
}

func tableResultValue(r core.Result) string {
	switch r.Status {
	case core.StatusOK:
		if r.Region != "" {
			return core.Green(strings.ToUpper(r.Region))
		}
		return core.Green("YES")
	case core.StatusRestricted:
		return core.Yellow("RESTRICTED")
	case core.StatusNo:
		return core.Red("NO")
	case core.StatusBanned:
		return core.Red("BANNED")
	case core.StatusNetworkErr, core.StatusErr, core.StatusUnexpected:
		return core.Red("ERR")
	case core.StatusFailed:
		return core.Blue("FAILED")
	default:
		return core.Purple("UNKNOWN")
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
		newDesc = fmt.Sprintf("等待测试开始 (%s)...", progressIPLabel)
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
			newDesc = fmt.Sprintf("正在测试 (%s): %s 等 %d 个测试", progressIPLabel, strings.Join(displayNames, ", "), len(activeList))
		} else {
			newDesc = fmt.Sprintf("正在测试 (%s): %s", progressIPLabel, strings.Join(displayNames, ", "))
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
