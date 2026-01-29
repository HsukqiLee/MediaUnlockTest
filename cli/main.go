package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	m "MediaUnlockTest/checks"

	"sort"

	"github.com/fatih/color"
	selfUpdate "github.com/inconshreveable/go-update"
	pb "github.com/schollz/progressbar/v3"
	"golang.org/x/net/proxy"
)

var (
	IPV4        bool = true
	IPV6        bool = true
	M           bool
	HK          bool
	TW          bool
	JP          bool
	KR          bool
	NA          bool
	SA          bool
	EU          bool
	AFR         bool
	SEA         bool
	OCEA        bool
	AI          bool
	Debug       bool   = false
	Conc        uint64 = 0
	Cache       bool   = false // 新增：是否启用缓存和串行地区执行
	sem         chan struct{}
	ResultLines []*result
	bar         *pb.ProgressBar
	Red         = color.New(color.FgRed).SprintFunc()
	Green       = color.New(color.FgGreen).SprintFunc()
	Yellow      = color.New(color.FgYellow).SprintFunc()
	Blue        = color.New(color.FgBlue).SprintFunc()
	Purple      = color.New(color.FgMagenta).SprintFunc()
	SkyBlue     = color.New(color.FgCyan).SprintFunc()
	White       = color.New(color.FgWhite).SprintFunc()
	resultCache = make(map[string]m.Result)
	cacheMutex  sync.RWMutex

	// 全局超时控制 - 优化超时时间
	testTimeout   = 15 * time.Second // 单个测试超时时间从30秒减少到15秒
	regionTimeout = 3 * time.Minute  // 整个地区测试超时时间从5分钟减少到3分钟

	// 新增：正在进行的测试状态管理
	activeTestsMutex sync.RWMutex
	activeTests      = make(map[string]bool)

	// 新增：全局进度条显示控制
	ShowActive bool = true // 是否显示正在进行的测试

	// 新增：进度条描述缓存，避免重复更新
	progressDescriptionCache string
)

type testItem struct {
	Name       string
	Func       func(client http.Client) m.Result
	SupportsV6 bool
}

type result struct {
	Name    string
	Divider bool
	Value   m.Result
}

type regionItem struct {
	Enabled bool
	Name    string
	Tests   []testItem
}

type Downloader struct {
	io.Reader
	Total   uint64
	Current uint64
	Pb      *pb.ProgressBar
	done    bool
}

func ShowSingleResult(r m.Result) (s string) {
	switch r.Status {
	case m.StatusOK:
		s = Green("YES")
		if r.Region != "" {
			s += Green(" (Region: " + strings.ToUpper(r.Region) + ")")
		}
		if Debug && r.CachedResult {
			s += " (Cached)"
		}
		return s

	case m.StatusNetworkErr:
		s = Red("ERR")
		if Debug {
			s += Yellow(" (Network Err: " + r.Err.Error() + ")")
		} else {
			s += Yellow(" (Network Err)")
		}
		if Debug && r.CachedResult {
			s += " (Cached)"
		}
		return s

	case m.StatusRestricted:
		s = Yellow("Restricted")
		if r.Info != "" {
			s = Yellow("Restricted (" + r.Info + ")")
		}
		if Debug && r.CachedResult {
			s += " (Cached)"
		}
		return s

	case m.StatusErr:
		s = Red("ERR")
		if r.Err != nil && Debug {
			s += Yellow(" (Err: " + r.Err.Error() + ")")
		}
		if Debug && r.CachedResult {
			s += " (Cached)"
		}
		return s

	case m.StatusNo:
		if r.Info != "" {
			return Red("NO ") + Yellow(" (Info: "+r.Info+")")
		}
		if r.Region != "" {
			return Red("NO  (Region: " + strings.ToUpper(r.Region) + ")")
		}
		return Red("NO")

	case m.StatusBanned:
		if r.Info != "" {
			return Red("Banned") + Yellow(" ("+r.Info+")")
		}
		return Red("Banned")

	case m.StatusUnexpected:
		return Purple("Unexpected")

	case m.StatusFailed:
		return Blue("Failed")

	default:
		return
	}
}

func ShowFinalResult() {
	fmt.Println("测试时间: ", Yellow(time.Now().Format("2006-01-02 15:04:05")))

	// 计算最大名称长度用于对齐
	NameLength := 25
	for _, r := range ResultLines {
		if len(r.Name) > NameLength {
			NameLength = len(r.Name)
		}
	}

	// 智能排序：保持分隔行的逻辑分组，对每个分组内的测试项目排序
	sortedResultLines := sortResultLinesWithGrouping(ResultLines)

	// 输出排序后的结果
	for i, r := range sortedResultLines {
		if r.Divider {
			// 判断是否为地区分组（包含IPv4/IPv6/Auto字样）还是子分组
			isRegionGroup := strings.Contains(r.Name, "IPv4") || strings.Contains(r.Name, "IPv6") || strings.Contains(r.Name, "Auto")

			// 只有地区分组才在标题前添加换行（除了第一个标题）
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
			if r.Value.Status == m.StatusOK && strings.HasSuffix(r.Name, "CDN") {
				result = SkyBlue(r.Value.Region)
			}
			fmt.Printf("%-"+strconv.Itoa(NameLength)+"s %s\n", r.Name, result)
		}
	}
}

// sortResultLinesWithGrouping 智能排序结果行，保持分隔行的逻辑分组
func sortResultLinesWithGrouping(lines []*result) []*result {
	if len(lines) == 0 {
		return lines
	}

	var sortedLines []*result
	var currentGroup []*result
	var currentDivider *result

	for _, r := range lines {
		if r.Divider {
			// 如果遇到分隔行，先处理当前分组
			if len(currentGroup) > 0 && currentDivider != nil {
				// 先添加分隔行
				sortedLines = append(sortedLines, currentDivider)
				// 再添加排序后的测试项目
				sortedGroup := sortTestItemsInGroup(currentGroup)
				sortedLines = append(sortedLines, sortedGroup...)
				currentGroup = nil
			}

			// 保存当前分隔行，等待处理其后的测试项目
			currentDivider = r
		} else {
			// 将测试项目添加到当前分组
			currentGroup = append(currentGroup, r)
		}
	}

	// 处理最后一个分组
	if len(currentGroup) > 0 && currentDivider != nil {
		sortedLines = append(sortedLines, currentDivider)
		sortedGroup := sortTestItemsInGroup(currentGroup)
		sortedLines = append(sortedLines, sortedGroup...)
	}

	return sortedLines
}

// sortTestItemsInGroup 对分组内的测试项目进行排序
func sortTestItemsInGroup(group []*result) []*result {
	if len(group) <= 1 {
		return group
	}

	// 创建副本进行排序，避免修改原数组
	sortedGroup := make([]*result, len(group))
	copy(sortedGroup, group)

	// 按名称进行排序
	sort.Slice(sortedGroup, func(i, j int) bool {
		return sortedGroup[i].Name < sortedGroup[j].Name
	})

	return sortedGroup
}

func NewBar(count int64) *pb.ProgressBar {
	return pb.NewOptions64(
		count,
		pb.OptionSetDescription("正在测试..."),
		pb.OptionSetWriter(os.Stderr),
		pb.OptionSetWidth(60),                   // 增加宽度以显示更多信息
		pb.OptionThrottle(200*time.Millisecond), // 减少刷新频率以提升性能
		pb.OptionShowCount(),
		pb.OptionClearOnFinish(),
		pb.OptionEnableColorCodes(true),
		pb.OptionSpinnerType(14),
		pb.OptionSetRenderBlankState(false), // 防止空白状态渲染
	)
}

// 新增：更新进度条描述，显示正在进行的测试
func updateProgressBarDescription() {
	if bar == nil {
		return
	}

	// 检查是否启用显示正在进行的测试
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

	// 缓存上一次的描述，避免重复更新
	var newDesc string
	if len(activeList) == 0 {
		newDesc = "等待测试开始..."
	} else if len(activeList) <= 3 {
		newDesc = "正在测试: " + strings.Join(activeList, ", ")
	} else {
		newDesc = fmt.Sprintf("正在测试: %s 等 %d 个测试", strings.Join(activeList[:3], ", "), len(activeList))
	}

	// 只有当描述发生变化时才更新，减少不必要的刷新
	if progressDescriptionCache != newDesc {
		bar.Describe(newDesc)
		progressDescriptionCache = newDesc
	}
}

// 新增：启动进度条更新协程
func startProgressUpdater() {
	// 如果用户选择不显示正在进行的测试，则不启动更新协程
	if !ShowActive {
		return
	}

	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond) // 减少更新频率到1秒，减少多行输出
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				updateProgressBarDescription()
			}
		}
	}()
}

// ExecuteTestsParallel 并行执行所有测试（不按地区分组）
func ExecuteTestsParallel(regions []regionItem, client http.Client, ipType int) {
	startTime := time.Now() // 记录开始时间

	// 收集所有启用的测试，并记录地区信息
	type testWithRegion struct {
		test   testItem
		region string
	}
	var allTests []testWithRegion
	var regionOrder []string                     // 记录地区顺序
	var regionSubGroups map[string][]interface{} // 记录每个地区的测试和子分组

	for _, region := range regions {
		if !region.Enabled {
			continue
		}
		regionOrder = append(regionOrder, region.Name)
		var regionItems []interface{}

		for _, test := range region.Tests {
			if test.Func == nil {
				// 这是子分组分隔行
				regionItems = append(regionItems, test)
			} else if !(ipType == 6 && !test.SupportsV6) {
				// 这是实际的测试
				allTests = append(allTests, testWithRegion{test: test, region: region.Name})
				regionItems = append(regionItems, test)
			}
		}

		if regionSubGroups == nil {
			regionSubGroups = make(map[string][]interface{})
		}
		regionSubGroups[region.Name] = regionItems
	}

	if len(allTests) == 0 {
		return
	}

	// 创建进度条
	bar = NewBar(int64(len(allTests)))
	startProgressUpdater()

	// 优化并发数量
	maxWorkers := 30
	if Conc > 0 {
		maxWorkers = int(Conc)
	} else if len(allTests) > 50 {
		maxWorkers = 40
	} else if len(allTests) < 20 {
		maxWorkers = 25
	}

	// 创建工作池和结果通道
	workerPool := make(chan struct{}, maxWorkers)
	resultChan := make(chan *result, len(allTests))
	completedTests := 0

	var wg sync.WaitGroup

	// 启动所有测试
	for _, testWithRegion := range allTests {
		workerPool <- struct{}{}
		wg.Add(1)

		// 创建局部变量来避免闭包问题
		test := testWithRegion.test

		go func() {
			defer func() {
				<-workerPool
				wg.Done()
			}()

			// 标记测试开始
			activeTestsMutex.Lock()
			activeTests[test.Name] = true
			activeTestsMutex.Unlock()

			// 执行测试（不使用缓存）
			done := make(chan m.Result, 1)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						select {
						case done <- m.Result{
							Status: m.StatusFailed,
							Err:    fmt.Errorf("测试panic: %v", r),
						}:
						default:
						}
					}
				}()

				m.ResetSessionHeaders()
				result := test.Func(client)
				result.CachedResult = false

				select {
				case done <- result:
				default:
				}
			}()

			// 等待测试完成或超时
			select {
			case testResult := <-done:
				r := &result{Name: test.Name, Value: testResult}
				select {
				case resultChan <- r:
				case <-time.After(1 * time.Second):
				}
			case <-time.After(testTimeout):
				timeoutResult := m.Result{
					Status: m.StatusFailed,
					Err:    fmt.Errorf("测试超时 (%v)", testTimeout),
				}
				r := &result{Name: test.Name, Value: timeoutResult}
				select {
				case resultChan <- r:
				default:
				}
			}

			// 标记测试完成
			activeTestsMutex.Lock()
			delete(activeTests, test.Name)
			activeTestsMutex.Unlock()
		}()
	}

	// 等待所有测试完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	resultsByRegion := make(map[string][]*result)
	for completedTests < len(allTests) {
		select {
		case r, ok := <-resultChan:
			if !ok {
				break
			}
			// 根据测试名称找到对应的地区
			var testRegion string
			for _, testWithRegion := range allTests {
				if testWithRegion.test.Name == r.Name {
					testRegion = testWithRegion.region
					break
				}
			}

			if testRegion != "" {
				resultsByRegion[testRegion] = append(resultsByRegion[testRegion], r)
			}
			completedTests++

			if bar != nil {
				bar.Add(1)
			}
		}
	}

	// 按地区顺序重新组织结果，添加分隔行和子分组
	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}

	for _, regionName := range regionOrder {
		if regionItems, exists := regionSubGroups[regionName]; exists && len(regionItems) > 0 {
			// 添加地区分隔行
			ResultLines = append(ResultLines, &result{Name: fmt.Sprintf("%s (%s)", regionName, ipTypeStr), Divider: true})

			// 按原有顺序添加该地区的所有项目（包括子分组和测试结果）
			for _, item := range regionItems {
				if test, ok := item.(testItem); ok {
					if test.Func == nil {
						// 这是子分组分隔行
						ResultLines = append(ResultLines, &result{Name: test.Name, Divider: true})
					} else {
						// 这是测试结果，需要从resultsByRegion中找到对应的结果
						for _, result := range resultsByRegion[regionName] {
							if result.Name == test.Name {
								ResultLines = append(ResultLines, result)
								break
							}
						}
					}
				}
			}
		}
	}

	// 清理
	if bar != nil {
		bar.Finish()
		fmt.Fprint(os.Stderr, "\r\033[K")
	}

	activeTestsMutex.Lock()
	for testName := range activeTests {
		delete(activeTests, testName)
	}
	activeTestsMutex.Unlock()

	// 显示性能统计
	totalDuration := time.Since(startTime)
	if len(allTests) > 0 {
		avgTime := totalDuration / time.Duration(len(allTests))
		fmt.Printf("\n性能统计:\n")
		fmt.Printf("总测试数量: %d\n", len(allTests))
		fmt.Printf("总耗时: %v\n", totalDuration)
		fmt.Printf("平均每个测试耗时: %v\n", avgTime)
		fmt.Printf("测试速度: %.2f 测试/秒\n", float64(len(allTests))/totalDuration.Seconds())
	}

}

func ExecuteTests(regions []regionItem, client http.Client, ipType int) {
	startTime := time.Now() // 记录开始时间

	for _, region := range regions {
		if !region.Enabled {
			continue
		}
		regionStartTime := time.Now() // 记录每个地区的开始时间

		ipTypeStr := fmt.Sprintf("IPv%d", ipType)
		if ipType == 0 {
			ipTypeStr = "Auto"
		}
		fmt.Printf("\n正在检测 %s (%s) ...\n", region.Name, ipTypeStr)

		ResultLines = append(ResultLines, &result{Name: fmt.Sprintf("%s (%s)", region.Name, ipTypeStr), Divider: true})
		if region.Name == "Korea" && ipType == 6 {
			ResultLines = append(ResultLines, &result{Name: "No Korean platform supports IPv6", Divider: false})
		}

		// 创建带超时的context
		ctx, cancel := context.WithTimeout(context.Background(), regionTimeout)
		defer cancel()

		// 计算实际需要执行的测试数量
		totalTests := 0
		var validTests []testItem
		for _, test := range region.Tests {
			if test.Func != nil && !(ipType == 6 && !test.SupportsV6) {
				totalTests++
				validTests = append(validTests, test)
			}
		}

		if totalTests == 0 {
			continue
		}

		bar = NewBar(int64(totalTests))

		// 启动进度条更新协程
		startProgressUpdater()

		// 优化并发数量 - 根据测试数量动态调整
		maxWorkers := 30 // 增加默认最大并发数从20到30
		if Conc > 0 {
			maxWorkers = int(Conc)
		} else if totalTests > 50 {
			maxWorkers = 40 // 测试数量多时增加并发数
		} else if totalTests < 20 {
			maxWorkers = 25 // 测试数量少时适当减少并发数
		}

		// 创建工作池和结果通道
		workerPool := make(chan struct{}, maxWorkers)
		resultChan := make(chan *result, totalTests)
		completedTests := 0

		// 使用WaitGroup确保所有goroutine完成
		var wg sync.WaitGroup

		// 启动测试 - 使用工作池控制并发
		for _, test := range validTests {
			// 处理子分组分隔行
			if test.Func == nil {
				ResultLines = append(ResultLines, &result{Name: test.Name, Divider: true})
				continue
			}

			// 获取工作池槽位
			workerPool <- struct{}{}
			wg.Add(1)

			go func(test testItem) {
				defer func() {
					// 释放工作池槽位
					<-workerPool
					wg.Done()
				}()

				// 标记测试开始
				activeTestsMutex.Lock()
				activeTests[test.Name] = true
				activeTestsMutex.Unlock()

				// 创建单个测试的超时context
				testCtx, testCancel := context.WithTimeout(ctx, testTimeout)
				defer testCancel()

				// 检查缓存
				cacheKey := fmt.Sprintf("%s_%d", test.Name, ipType)
				cacheMutex.RLock()
				if cachedResult, exists := resultCache[cacheKey]; exists {
					cacheMutex.RUnlock()
					r := &result{Name: test.Name, Value: cachedResult}
					select {
					case resultChan <- r:
					case <-testCtx.Done():
						// 如果context已取消，丢弃结果
					}

					// 标记测试完成
					activeTestsMutex.Lock()
					delete(activeTests, test.Name)
					activeTestsMutex.Unlock()
					return
				}
				cacheMutex.RUnlock()

				// 执行测试
				done := make(chan m.Result, 1)
				go func() {
					defer func() {
						if r := recover(); r != nil {
							// 防止panic导致程序崩溃
							select {
							case done <- m.Result{
								Status: m.StatusFailed,
								Err:    fmt.Errorf("测试panic: %v", r),
							}:
							default:
								// 如果done channel已满，丢弃结果
							}
						}
					}()

					m.ResetSessionHeaders()
					result := test.Func(client)
					result.CachedResult = false

					// 缓存结果
					cacheMutex.Lock()
					resultCache[cacheKey] = result
					cacheMutex.Unlock()

					select {
					case done <- result:
					default:
						// 如果done channel已满，说明已经超时，丢弃结果
					}
				}()

				// 等待测试完成或超时
				select {
				case testResult := <-done:
					// 测试完成
					r := &result{Name: test.Name, Value: testResult}
					select {
					case resultChan <- r:
					case <-testCtx.Done():
						// 如果context已取消，丢弃结果
					}

				case <-testCtx.Done():
					// 测试超时
					timeoutResult := m.Result{
						Status: m.StatusFailed,
						Err:    fmt.Errorf("测试超时 (%v)", testTimeout),
					}
					r := &result{Name: test.Name, Value: timeoutResult}
					select {
					case resultChan <- r:
					default:
						// 如果resultChan已满，丢弃结果
					}
				}

				// 标记测试完成
				activeTestsMutex.Lock()
				delete(activeTests, test.Name)
				activeTestsMutex.Unlock()
			}(test)
		}

		// 启动一个goroutine来等待所有测试完成，然后关闭结果通道
		go func() {
			wg.Wait()
			close(resultChan)
		}()

		// 收集结果
		for completedTests < totalTests {
			select {
			case r, ok := <-resultChan:
				if !ok {
					// 通道已关闭，所有结果都已收集
					break
				}
				ResultLines = append(ResultLines, r)
				completedTests++

				if bar != nil {
					// 不再在进度条中显示测试结果，只更新计数
					bar.Add(1)
				}

			case <-ctx.Done():
				// 整个地区测试超时
				fmt.Printf("警告：%s (%s) 测试超时，已完成 %d/%d 个测试\n",
					region.Name, ipTypeStr, completedTests, totalTests)
				goto timeoutExit
			}
		}

	timeoutExit:

		regionDuration := time.Since(regionStartTime)
		fmt.Printf("%s (%s) 检测完毕 (%d/%d) - 耗时: %v\n",
			region.Name, ipTypeStr, completedTests, totalTests, regionDuration)
		if bar != nil {
			bar.Finish()
			// 确保进度条完全清理，防止空白行
			fmt.Fprint(os.Stderr, "\r\033[K")
		}

		// 清理活动测试状态
		activeTestsMutex.Lock()
		for testName := range activeTests {
			delete(activeTests, testName)
		}
		activeTestsMutex.Unlock()
	}

	// 显示总体性能统计
	totalDuration := time.Since(startTime)
	totalTests := 0
	for _, region := range regions {
		if region.Enabled {
			for _, test := range region.Tests {
				if test.Func != nil {
					totalTests++
				}
			}
		}
	}

	if totalTests > 0 {
		avgTime := totalDuration / time.Duration(totalTests)
		fmt.Printf("\n性能统计:\n")
		fmt.Printf("总测试数量: %d\n", totalTests)
		fmt.Printf("总耗时: %v\n", totalDuration)
		fmt.Printf("平均每个测试耗时: %v\n", avgTime)
		fmt.Printf("测试速度: %.2f 测试/秒\n", float64(totalTests)/totalDuration.Seconds())
	}
}

var GlobeTests = []testItem{
	{"Amazon Prime Video", m.PrimeVideo, false},
	{"Apple", m.Apple, true},
	{"Bing", m.Bing, true},
	{"Dazn", m.Dazn, false},
	{"Disney+", m.DisneyPlus, true},
	{"Google Play Store", m.GooglePlayStore, true},
	{"iQiYi", m.IQiYi, false},
	{"Netflix", m.NetflixRegion, true},
	{"Netflix CDN", m.NetflixCDN, true},
	{"Reddit", m.Reddit, false},
	{"Spotify", m.Spotify, true},
	{"Steam", m.Steam, false},
	{"TikTok", m.TikTok, false},
	{"TVBAnywhere+", m.TVBAnywhere, false},
	{"Viu.com", m.ViuCom, false},
	{"Wikipedia", m.WikipediaEditable, true},
	{"Youtube CDN", m.YoutubeCDN, true},
	{"Youtube Premium", m.YoutubeRegion, true},
}

var HongKongTests = []testItem{
	{"Bahamut Anime", m.BahamutAnime, false},
	{"Bilibili HongKong/Macau Only", m.BilibiliHKMO, false},
	{"Hoy TV", m.HoyTV, true},
	{"Max", m.Max, true},
	{"MyTVSuper", m.MyTvSuper, false},
	{"NBA TV", m.NBA_TV, true},
	//{"Now TV", m.NowTV, false},
	{"SonyLiv", m.SonyLiv, false},
	{"Viu.TV", m.ViuTV, false},
}

var TaiwanTests = []testItem{
	{"4GTV", m.TW4GTV, false},
	{"Bahamut Anime", m.BahamutAnime, false},
	{"Bilibili Taiwan Only", m.BilibiliTW, false},
	{"CatchPlay+", m.Catchplay, false},
	{"Friday Video", m.FridayVideo, false},
	{"Hami Video", m.HamiVideo, false},
	{"KKTV", m.KKTV, false},
	{"LiTV", m.LiTV, false},
	{"LineTV", m.LineTV, false},
	{"Max", m.Max, true},
	{"MyVideo", m.MyVideo, false},
	{"Ofiii", m.Ofiii, false},
}

var JapanTests = []testItem{
	{"Abema", m.Abema, false},
	{"AnimeFesta", m.AnimeFesta, false},
	{"D Anime Store", m.DAnimeStore, false},
	{"DMM", m.DMM, false},
	{"DMM TV", m.DMMTV, true},
	{"EroGameSpace", m.EroGameSpace, false},
	{"FOD(Fuji TV)", m.FOD, false},
	{"Hulu Japan", m.HuluJP, false},
	{"J:COM On Demand", m.J_COM_ON_DEMAND, false},
	{"Kancolle", m.Kancolle, false},
	{"Karaoke@DAM", m.Karaoke, false},
	{"Lemino", m.Lemino, true},
	{"MGStage", m.MGStage, false},
	{"Mora", m.Mora, false},
	{"Music.jp", m.MusicJP, false},
	{"NHK+", m.NHKPlus, true},
	{"Niconico", m.Niconico, false},
	{"Pretty Derby Japan", m.PrettyDerbyJP, true},
	{"Princess Connect Re:Dive Japan", m.PCRJP, false},
	{"Project Sekai: Colorful Stage", m.PJSK, false},
	{"Radiko", m.Radiko, false},
	{"Rakuten Magazine", m.RakutenMagazine, false},
	{"Rakuten TV JP", m.RakutenTV_JP, false},
	{"Telasa", m.Telasa, true},
	{"TVer", m.TVer, false},
	{"U-NEXT", m.U_NEXT, true},
	{"VideoMarket", m.VideoMarket, false},
	{"Watcha", m.Watcha, false},
	{"Wowow", m.Wowow, false},
}

var KoreaTests = []testItem{
	{"Afreeca", m.Afreeca, false},
	{"Coupang Play", m.CoupangPlay, false},
	{"KBS", m.KBS, false},
	{"Naver TV", m.NaverTV, false},
	{"Panda TV", m.PandaTV, false},
	{"Spotv Now", m.SpotvNow, false},
	{"Tving", m.Tving, false},
	{"Watcha", m.Watcha, false},
	{"Wavve", m.Wavve, false},
}

var NorthAmericaTests = []testItem{
	{"A&E TV", m.AETV, false},
	{"Acorn TV", m.AcornTV, false},
	{"AMC+", m.AMCPlus, true},
	{"BritBox", m.BritBox, true},
	{"CBC Gem", m.CBCGem, false},
	{"Crave", m.Crave, false},
	{"Crunchyroll", m.Crunchyroll, false},
	{"CW TV", m.CW_TV, true},
	{"DirecTV Stream", m.DirectvStream, true},
	{"Discovery+", m.DiscoveryPlus, false},
	{"encoreTVB", m.EncoreTVB, false},
	{"ESPN+", m.ESPNPlus, true},
	{"FOX", m.Fox, true},
	{"Fubo TV", m.FuboTV, false},
	{"FXNOW", m.FXNOW, false},
	{"Hotstar", m.Hotstar, true},
	{"Hulu", m.Hulu, true},
	{"KOCOWA+", m.KOCOWA, false},
	{"MGM+", m.MGMPlus, false},
	{"MathsSpot Roblox", m.MathsSpotRoblox, false},
	{"Max", m.Max, true},
	{"NBC TV", m.NBC_TV, true},
	{"NFL+", m.NFLPlus, false},
	{"NBA TV", m.NBA_TV, true},
	{"Paramount+", m.ParamountPlus, true},
	{"Peacock TV", m.PeacockTV, true},
	{"Philo", m.Philo, false},
	{"Pluto TV", m.PlutoTV, false},
	{"SHOWTIME", m.SHOWTIME, true},
	{"Shudder", m.Shudder, true},
	{"Sling TV", m.SlingTV, true},
	{"SonyLiv", m.SonyLiv, true},
	{"Starz", m.Starz, false},
	{"TLC GO", m.TlcGo, true},
	{"Tubi TV", m.TubiTV, true},
	{"Viaplay", m.Viaplay, false},
}

var SouthAmericaTests = []testItem{
	{"DirecTV GO", m.DirecTVGO, false},
	{"Max", m.Max, true},
}

var EuropeTests = []testItem{
	{"Rakuten TV EU", m.RakutenTV_EU, false},
	{"Sky Show Time", m.SkyShowTime, true},
	{"Viaplay", m.Viaplay, true},
	{"TNTSports", m.TNTSports, false},
	{"Eurosport RO", m.EurosportRO, false},
	{"Setanta Sports", m.SetantaSports, true},
	{"KOCOWA+", m.KOCOWA, false},
	{"MathsSpot Roblox", m.MathsSpotRoblox, false},
	{"Max", m.Max, true},
	{"SonyLiv", m.SonyLiv, true},
	{"GB", nil, true},
	{"BBC iPlayer", m.BBCiPlayer, false},
	{"BritBox", m.BritBox, true},
	{"ITVX", m.ITVX, false},
	{"Channel 4", m.Channel4, false},
	{"Channel 5", m.Channel5, false},
	{"Discovery+ UK", m.DiscoveryPlus_UK, false},
	{"Sky Go", m.SkyGo, false},
	{"FR", nil, true},
	{"Canal+", m.CanalPlus, false},
	{"Molotov", m.Molotov, true},
	{"France TV", m.FranceTV, true},
	{"DE", nil, false},
	{"Joyn", m.Joyn, false},
	{"Sky DE", m.Sky_DE, false},
	{"ZDF", m.ZDF, false},
	{"NL", nil, true},
	{"NLZIET", m.NLZIET, false},
	{"Video Land", m.VideoLand, true},
	{"NPO Start Plus", m.NPOStartPlus, false},
	{"ES", nil, false},
	{"Movistar Plus+", m.MoviStarPlus, false},
	{"IT", nil, false},
	{"Rai Play", m.RaiPlay, false},
	{"CH", nil, false},
	{"Sky CH", m.Sky_CH, false},
	{"RU", nil, false},
	{"Amediateka", m.Amediateka, false},
}

var AfricaTests = []testItem{
	{"DSTV", m.DSTV, false},
	{"Showmax", m.Showmax, true},
}

var SouthEastAsiaTests = []testItem{
	{"Max", m.Max, true},
	{"Hotstar", m.Hotstar, true},
	{"NBA TV", m.NBA_TV, true},
	{"Bilibili SouthEastAsia Only", m.BilibiliSEA, false},
	{"SG", nil, false},
	{"MeWatch", m.MeWatch, false},
	{"CatchPlay+", m.Catchplay, false},
	{"TH", nil, false},
	{"AIS Play", m.AISPlay, false},
	{"TrueID", m.TrueID, false},
	{"Bilibili Thailand Only", m.BilibiliTH, false},
	{"ID", nil, false},
	{"Bilibili Indonesia Only", m.BilibiliID, false},
	{"VN", nil, false},
	{"Clip TV", m.ClipTV, false},
	{"Galaxy Play", m.GalaxyPlay, false},
	{"K+", m.KPlus, false},
	{"Bilibili Vietnam Only", m.BilibiliVN, false},
	{"MY", nil, false},
	{"Sooka", m.Sooka, false},
	{"IN", nil, true},
	{"Tata Play", m.TataPlay, true},
	{"SonyLiv", m.SonyLiv, true},
	{"JioCinema", m.JioCinema, true},
	{"MX Player", m.MXPlayer, false},
	{"Zee5", m.Zee5, true},
}

var OceaniaTests = []testItem{
	{"10 Play", m.Channel10, false},
	{"7Plus", m.SevenPlus, true},
	{"ABC iView", m.ABCiView, false},
	{"Acorn TV", m.AcornTV, false},
	{"AMC+", m.AMCPlus, true},
	{"Binge", m.Binge, true},
	{"BritBox", m.BritBox, true},
	{"Channel 9", m.Channel9, true},
	{"Doc Play", m.DocPlay, false},
	{"Kayo Sports", m.KayoSports, false},
	{"KOCOWA+", m.KOCOWA, false},
	{"Maori TV", m.MaoriTV, false},
	{"NBA TV", m.NBA_TV, true},
	{"Neon TV", m.NeonTV, false},
	{"Optus Sports", m.OptusSports, true},
	{"Paramount+", m.ParamountPlus, true},
	{"SBS on Demand", m.SBSonDemand, false},
	{"Sky Go NZ", m.SkyGo_NZ, false},
	{"SonyLiv", m.SonyLiv, true},
	{"Stan", m.Stan, false},
	{"Three Now", m.ThreeNow, false},
}

var AITests = []testItem{
	{"ChatGPT", m.ChatGPT, true},
	{"Claude", m.Claude, true},
	{"Copilot", m.Copilot, true},
	{"Google Gemini", m.Gemini, true},
	{"Meta AI", m.MetaAI, true},
	{"Sora", m.Sora, true},
}

func GetIPInfo(url string, ipType int, isCloudflare bool) (string, error) {
	timeout := 6
	if ipType == 6 {
		timeout = 3
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var client http.Client
	switch ipType {
	case 6:
		client = m.Ipv6HttpClient
	case 4:
		client = m.Ipv4HttpClient
	case 0:
		client = m.AutoHttpClient
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
	if isCloudflare {
		s := string(b)
		i := strings.Index(s, "ip=")
		s = s[i+3:]
		i = strings.Index(s, "\n")
		return s[:i], nil
	} else {
		return strings.TrimSpace(string(b)), nil
	}
}

func ReadSelect() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signalChan)

	fmt.Println("请选择检测项目：")
	fmt.Println(Green("直接按回车进行全部检测") + "，" + Yellow("按 Ctrl+C 取消检测") + "。")
	fmt.Println("")
	fmt.Println("[0]  :   跨国平台")
	fmt.Println("[1]  :   台湾平台")
	fmt.Println("[2]  :   香港平台")
	fmt.Println("[3]  :   日本平台")
	fmt.Println("[4]  :   韩国平台")
	fmt.Println("[5]  :   北美平台")
	fmt.Println("[6]  :   南美平台")
	fmt.Println("[7]  :   欧洲平台")
	fmt.Println("[8]  :   非洲平台")
	fmt.Println("[9]  : 东南亚平台")
	fmt.Println("[10] : 大洋洲平台")
	fmt.Println("[11] :   ＡＩ平台")
	fmt.Println("")
	fmt.Print("请输入对应数字，空格分隔，回车确认: ")

	inputChan := make(chan string, 1)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-signalChan:
				fmt.Println("")
				fmt.Println(Yellow("输入中止，检测已取消。"))
				os.Exit(0)
			default:
				if input, err := reader.ReadString('\n'); err == nil {
					inputChan <- strings.TrimSpace(input)
					return
				}
			}
		}
	}()

	select {
	case <-signalChan:
		fmt.Println("")
		fmt.Println(Yellow("输入中止，检测已取消。"))
		os.Exit(0)
	case input := <-inputChan:
		for _, c := range strings.Split(input, " ") {
			switch c {
			case "0":
				M = true
			case "1":
				TW = true
			case "2":
				HK = true
			case "3":
				JP = true
			case "4":
				KR = true
			case "5":
				NA = true
			case "6":
				SA = true
			case "7":
				EU = true
			case "8":
				AFR = true
			case "9":
				SEA = true
			case "10":
				OCEA = true
			case "11":
				AI = true
			default:
				M, TW, HK, JP, KR, NA, SA, EU, AFR, SEA, OCEA, AI = true, true, true, true, true, true, true, true, true, true, true, true
			}
		}
	}
}

func (d *Downloader) Read(p []byte) (n int, err error) {
	n, err = d.Reader.Read(p)
	d.Current += uint64(n)
	if d.done {
		return
	}
	d.Pb.Add(n)
	if d.Current == d.Total {
		d.done = true
		d.Pb.Describe("unlock-test下载完成")
		d.Pb.Finish()
	}
	return
}

func checkUpdate() {
	resp, err := http.Get("https://unlock.icmp.ing/test/latest/version")
	if err != nil {
		log.Println("[ERR] 获取版本信息时出错:", err)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERR] 读取版本信息时出错:", err)
		return
	}

	parts := strings.Split(string(b), "-")
	if len(parts) != 2 {
		log.Println("[ERR] 版本号格式错误:", err)
		return
	}
	version := parts[0]

	if version == m.Version {
		fmt.Println("已经是最新版本")
		return
	}

	timestampInt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Println("[ERR] 版本号时间戳错误:", err)
		return
	}
	timestamp := time.Unix(timestampInt, 0)

	fmt.Println("最新版本：", version)
	fmt.Println("发布时间：", timestamp.Format("2006-01-02 15:04:05"))

	OS, ARCH := runtime.GOOS, runtime.GOARCH
	fmt.Println("运行系统：", OS)
	fmt.Println("运行架构：", ARCH)

	if OS == "android" && strings.Contains(os.Getenv("PREFIX"), "com.termux") {
		target_path := os.Getenv("PREFIX") + "/bin"
		out, err := os.Create(target_path + "/unlock-test_new")
		if err != nil {
			log.Fatal("[ERR] 创建文件出错:", err)
			return
		}
		defer out.Close()
		log.Println("下载unlock-test中 ...")
		url := "https://unlock.icmp.ing/test/latest/unlock-test_" + OS + "_" + ARCH
		resp, err = http.Get(url)
		if err != nil {
			log.Fatal("[ERR] 下载unlock-test时出错:", err)
		}
		defer resp.Body.Close()
		downloader := &Downloader{
			Reader: resp.Body,
			Total:  uint64(resp.ContentLength),
			Pb:     pb.DefaultBytes(resp.ContentLength, "下载进度"),
		}
		if _, err := io.Copy(out, downloader); err != nil {
			log.Fatal("[ERR] 下载unlock-test时出错:", err)
		}
		if err := os.Chmod(target_path+"/unlock-test_new", 0777); err != nil {
			log.Fatal("[ERR] 更改unlock-test后端权限出错:", err)
		}
		if _, err := os.Stat(target_path + "/unlock-test"); err == nil {
			if err := os.Remove(target_path + "/unlock-test"); err != nil {
				log.Fatal("[ERR] 删除unlock-test旧版本时出错:", err.Error())
			}
		}
		if err := os.Rename(target_path+"/unlock-test_new", target_path+"/unlock-test"); err != nil {
			log.Fatal("[ERR] 更新unlock-test后端时出错:", err)
		}
	} else {
		url := "https://unlock.icmp.ing/test/latest/unlock-test_" + OS + "_" + ARCH
		if OS == "windows" {
			url += ".exe"
		}

		resp, err = http.Get(url)
		if err != nil {
			log.Fatal("[ERR] 下载unlock-test时出错:", err)
			return
		}
		defer resp.Body.Close()

		bar := pb.DefaultBytes(
			resp.ContentLength,
			"下载进度",
		)

		body := io.TeeReader(resp.Body, bar)

		if resp.StatusCode != http.StatusOK {
			log.Fatal("[ERR] 下载unlock-test时出错: 非预期的状态码", resp.StatusCode)
			return
		}

		err = selfUpdate.Apply(body, selfUpdate.Options{})
		if err != nil {
			log.Fatal("[ERR] 更新unlock-test时出错:", err)
			return
		}
	}

	fmt.Println("[OK] unlock-test后端更新成功")
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
	fmt.Printf("当天运行共%s次, 本月运行共%s次, 共计运行%s次\n", SkyBlue(d), Yellow(m), Green(t))
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

var setSocketOptions = func(network, address string, c syscall.RawConn, interfaceName string) (err error) {
	return
}

func main() {
	var (
		Interface   string
		DNSServers  string
		HTTPProxy   string
		SocksProxy  string
		ShowVersion bool
		CheckUpdate bool
		NF          bool
		TestMode    bool
		IPMode      int
		IP4_1       string
		IP6_1       string
		IP4_2       string
		IP6_2       string
		err         error
		IsProxy     bool
	)
	flag.StringVar(&Interface, "I", "", "Source IP or network interface to use for connections")
	flag.StringVar(&DNSServers, "dns-servers", "", "Custom DNS servers (format: ip:port)")
	flag.StringVar(&HTTPProxy, "http-proxy", "", "HTTP proxy URL (format: http://user:pass@host:port)")
	flag.StringVar(&SocksProxy, "socks-proxy", "", "SOCKS5 proxy URL (format: socks5://user:pass@host:port)")
	flag.BoolVar(&ShowVersion, "v", false, "Show version information and exit")
	flag.BoolVar(&CheckUpdate, "u", false, "Update to latest version")
	flag.BoolVar(&NF, "nf", false, "Only test Netflix availability")
	flag.BoolVar(&TestMode, "test", false, "Run in test mode (only checks Viaplay)")
	flag.BoolVar(&Debug, "debug", false, "Enable debug mode for verbose output")
	flag.IntVar(&IPMode, "m", 0, "Connection mode: 0=auto (default), 4=IPv4 only, 6=IPv6 only")
	flag.Uint64Var(&Conc, "conc", 0, "Max concurrent tests (0=unlimited)")
	flag.BoolVar(&ShowActive, "show-active", true, "Show active tests in progress bar (default: true)")     // 新增参数
	flag.BoolVar(&Cache, "cache", false, "Enable caching and sequential region execution (default: false)") // 新增：缓存参数
	flag.Parse()
	if ShowVersion {
		fmt.Println(m.Version)
		return
	}
	if CheckUpdate {
		checkUpdate()
		return
	}
	if Interface != "" {
		if IP := net.ParseIP(Interface); IP != nil {
			m.Dialer.LocalAddr = &net.TCPAddr{IP: IP}
		} else {
			m.Dialer.Control = func(network, address string, c syscall.RawConn) error {
				return setSocketOptions(network, address, c, Interface)
			}
		}
	}
	if DNSServers != "" {
		m.Dialer.Resolver = &net.Resolver{
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "udp", DNSServers)
			},
		}
		m.Ipv4Transport.Resolver = m.Dialer.Resolver
		m.Ipv6Transport.Resolver = m.Dialer.Resolver
		m.AutoHttpClient.Transport.(*m.CustomTransport).Resolver = m.Dialer.Resolver
	}
	if HTTPProxy != "" {
		if u, err := url.Parse(HTTPProxy); err == nil {
			m.ClientProxy = http.ProxyURL(u)
			m.Ipv4Transport.Proxy = m.ClientProxy
			m.Ipv6Transport.Proxy = m.ClientProxy
			m.AutoHttpClient.Transport.(*m.CustomTransport).Proxy = m.ClientProxy
		}
	}
	if SocksProxy != "" {
		proxyURL, err := url.Parse(SocksProxy)
		if err != nil {
			log.Fatal("SOCKS5 地址不合法：", err)
		}
		var auth *proxy.Auth
		if proxyURL.User != nil {
			username := proxyURL.User.Username()
			password, _ := proxyURL.User.Password()
			auth = &proxy.Auth{
				User:     username,
				Password: password,
			}
		}
		dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, auth, proxy.Direct)
		if err != nil {
			log.Fatal("创建 SOCKS5 连接失败：", err)
		}

		m.Ipv4Transport.SocksDialer = dialer
		m.Ipv6Transport.SocksDialer = dialer
		m.AutoHttpClient.Transport.(*m.CustomTransport).SocksDialer = dialer
	}
	if Conc > 0 {
		sem = make(chan struct{}, Conc)
	}

	if NF {
		fmt.Println("Netflix", ShowSingleResult(m.NetflixRegion(m.AutoHttpClient)))
		return
	}

	if TestMode {
		//GetIpv4Info()
		//GetIpv6Info()

		fmt.Println("itvx", ShowSingleResult(m.ITVX(m.AutoHttpClient)))

		return
	}

	fmt.Println("")
	fmt.Println("[ 项目地址: " + SkyBlue("https://github.com/HsukqiLee/MediaUnlockTest") + " ]")
	fmt.Println("[ 使用方式: " + Yellow("bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)") + " ]")
	fmt.Println()

	if !Debug {
		IP4_2, err := GetIPInfo("https://1.1.1.1/cdn-cgi/trace", 4, true)
		if err != nil {
			fmt.Println(Red("无法获取 IPv4 地址"))
			IPV4 = false
		} else {
			fmt.Println(SkyBlue("IPv4 地址： ") + Green(IP4_2))
			IPV4 = true
		}
		IP6_2, err := GetIPInfo("https://[2606:4700:4700::1111]/cdn-cgi/trace", 6, true)
		if err != nil {
			fmt.Println(Red("无法获取 IPv6 地址"))
			IPV6 = false
		} else {
			fmt.Println(SkyBlue("IPv6 地址： ") + Green(IP6_2))
			IPV6 = true
		}
	} else {
		fmt.Println("[ 正在获取国内分流 IP... ]")
		if IPMode == 0 || IPMode == 4 {
			IP4_1, err = GetIPInfo("https://ipv4.tsinbei.cn", IPMode, false)
			if err != nil {
				if Debug {
					fmt.Println(Red("无法获取国内分流 IPv4 地址 (") + Yellow(err.Error()) + Red(")"))
				} else {
					fmt.Println(Red("无法获取国内分流 IPv4 地址"))
				}
			} else {
				fmt.Println(SkyBlue("IPv4 地址： ") + Green(IP4_1))
			}
		}
		if IPMode == 0 || IPMode == 6 {
			IP6_1, err = GetIPInfo("https://ipv6.tsinbei.cn", IPMode, false)
			if err != nil {
				if Debug {
					fmt.Println(Red("无法获取国内分流 IPv6 地址 (") + Yellow(err.Error()) + Red(")"))
				} else {
					fmt.Println(Red("无法获取国内分流 IPv6 地址"))
				}
			} else {
				fmt.Println(SkyBlue("IPv6 地址： ") + Green(IP6_1))
			}
		}
		fmt.Println("")
		fmt.Println("[ 正在获取国外分流 IP... ]")
		if IPMode == 0 || IPMode == 4 {
			IP4_2, err = GetIPInfo("https://1.1.1.1/cdn-cgi/trace", IPMode, true)
			if err != nil {
				if Debug {
					fmt.Println(Red("无法获取国外 IPv4 地址 (") + Yellow(err.Error()) + Red(")"))
				} else {
					fmt.Println(Red("无法获取国外 IPv4 地址"))
				}
			} else {
				fmt.Println(SkyBlue("IPv4 地址： ") + Green(IP4_2))
			}
		}
		if IPMode == 0 || IPMode == 6 {
			IP6_2, err = GetIPInfo("https://[2606:4700:4700::1111]/cdn-cgi/trace", IPMode, true)
			if err != nil {
				if Debug {
					fmt.Println(Red("无法获取国外 IPv6 地址 (") + Yellow(err.Error()) + Red(")"))
				} else {
					fmt.Println(Red("无法获取国外 IPv6 地址"))
				}
			} else {
				fmt.Println(SkyBlue("IPv6 地址： ") + Green(IP6_2))
			}
		}
		fmt.Println("")
		fmt.Println("[ 正在检测系统代理... ]")

		if IPMode == 0 || IPMode == 4 {
			IP4, err := GetIPInfo("https://www.cloudflare.com/cdn-cgi/trace", 4, true)
			if err != nil {
				if IP4_1 != "" || IP4_2 != "" {
					IsProxy = true
					fmt.Println(Yellow("正在使用系统代理，且无法通过 IPv4 连接代理"))
				} else {
					IPV4 = false
					fmt.Println(Red("未使用 IPv4 代理，无 IPv4 网络"))
				}
			} else {
				IPV4 = true
				if IP4_1 != IP4_2 || IP4_1 != IP4 {
					IsProxy = true
					fmt.Println(Yellow("正在使用监听地址为 IPv4 的代理，出口 IP：") + Red(IP4))
				} else if IP4 == IP4_1 {
					fmt.Println(Green("未使用 IPv4 代理，有 IPv4 网络"))
				} else {
					fmt.Println(Red("无法强制使用 IPv4 网络测试，可能使用 IPv4 代理"))
					IPV4 = false
					if IPMode == 4 {
						IPV6 = false
					}
				}
			}
		}
		if IPMode == 0 || IPMode == 6 {
			IP6, err := GetIPInfo("https://www.cloudflare.com/cdn-cgi/trace", 6, true)
			if err != nil {
				if IP6_1 != "" && IP6_2 != "" {
					IsProxy = true
					fmt.Println(Yellow("正在使用系统代理，且无法通过 IPv6 连接代理"))
				} else {
					IPV6 = false
					fmt.Println(Red("未使用 IPv6 代理，无 IPv6 网络"))
				}
			} else {
				IPV6 = true
				if IP6_1 != IP6_2 && IP6_1 != IP6 {
					IsProxy = true
					fmt.Println(Yellow("正在使用监听地址为 IPv6 的代理，出口 IP：") + Red(IP6))
				} else if IP6 == IP6_1 {
					fmt.Println(Green("未使用 IPv6 代理，有 IPv6 网络"))
				} else {
					fmt.Println(Red("无法强制使用 IPv6 网络测试，可能使用 IPv6 代理"))
					IPV6 = false
					if IPMode == 6 {
						IPV4 = false
					}
				}
			}
		}
	}

	if IsProxy {
		fmt.Println(Yellow("提示：正在使用系统代理，此时连接行为全部受代理控制"))
	}
	if IPMode != 0 {
		switch IPMode {
		case 4:
			IPV6 = false
		case 6:
			IPV4 = false
		}
	}
	fmt.Println()

	if IPV4 || IPV6 {
		ReadSelect()
	}
	regions := []regionItem{
		{Enabled: M, Name: "Globe", Tests: GlobeTests},
		{Enabled: TW, Name: "Taiwan", Tests: TaiwanTests},
		{Enabled: HK, Name: "HongKong", Tests: HongKongTests},
		{Enabled: JP, Name: "Japan", Tests: JapanTests},
		{Enabled: KR, Name: "Korea", Tests: KoreaTests},
		{Enabled: NA, Name: "NorthAmerica", Tests: NorthAmericaTests},
		{Enabled: SA, Name: "SouthAmerica", Tests: SouthAmericaTests},
		{Enabled: EU, Name: "Europe", Tests: EuropeTests},
		{Enabled: AFR, Name: "Africa", Tests: AfricaTests},
		{Enabled: SEA, Name: "SouthEastAsia", Tests: SouthEastAsiaTests},
		{Enabled: OCEA, Name: "Oceania", Tests: OceaniaTests},
		{Enabled: AI, Name: "AI", Tests: AITests},
	}
	if IsProxy {
		if Cache {
			ExecuteTests(regions, m.AutoHttpClient, 0)
		} else {
			ExecuteTestsParallel(regions, m.AutoHttpClient, 0)
		}
	} else {
		if IPV4 {
			if Cache {
				ExecuteTests(regions, m.Ipv4HttpClient, 4)
			} else {
				ExecuteTestsParallel(regions, m.Ipv4HttpClient, 4)
			}
		}
		if IPV6 {
			if Cache {
				ExecuteTests(regions, m.Ipv6HttpClient, 6)
			} else {
				ExecuteTestsParallel(regions, m.Ipv6HttpClient, 6)
			}
		}
	}
	fmt.Println()
	ShowFinalResult()
	fmt.Println()
	fmt.Println("检测完毕，感谢您的使用！")
	ShowCounts()
	fmt.Println()
	ShowAD()
}
