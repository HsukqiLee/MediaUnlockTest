package main

import (
	core "MediaUnlockTest/pkg/core"
	m "MediaUnlockTest/pkg/providers"
	"context"
	"fmt"
	"sync"
	"time"
)

// ExecuteTestsParallel 并行执行所有测试（不按地区分组）
func ExecuteTestsParallel(regions []regionItem, client core.HttpClient, ipType int) {
	startTime := time.Now()

	// 收集所有启用的测试，并记录地区信息
	type testWithRegion struct {
		test   m.TestItem
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

				regionItems = append(regionItems, test)
			} else if !(ipType == 6 && !test.SupportsV6) {

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

	bar = newProgressBar(int64(len(allTests)), "正在测试...")
	startProgressUpdater()

	maxWorkers := 30
	if Conc > 0 {
		maxWorkers = int(Conc)
	} else if len(allTests) > 50 {
		maxWorkers = 40
	} else if len(allTests) < 20 {
		maxWorkers = 25
	}

	workerPool := make(chan struct{}, maxWorkers)
	resultChan := make(chan *result, len(allTests))
	completedTests := 0

	var wg sync.WaitGroup

	for _, testWithRegion := range allTests {
		workerPool <- struct{}{}
		wg.Add(1)

		test := testWithRegion.test

		go func() {
			defer func() {
				<-workerPool
				wg.Done()
			}()

			activeTestsMutex.Lock()
			activeTests[test.Name] = true
			activeTestsMutex.Unlock()

			done := make(chan core.Result, 1)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						select {
						case done <- core.Result{
							Status: core.StatusFailed,
							Err:    fmt.Errorf("测试panic: %v", r),
						}:
						default:
						}
					}
				}()

				core.ResetSessionHeaders()
				client = core.NewHttpClient(ipType)
				result := test.Func(client)
				result.CachedResult = false

				select {
				case done <- result:
				default:
				}
			}()

			select {
			case testResult := <-done:
				r := &result{Name: test.Name, Value: testResult}
				select {
				case resultChan <- r:
				case <-time.After(1 * time.Second):
				}
			case <-time.After(testTimeout):
				timeoutResult := core.Result{
					Status: core.StatusFailed,
					Err:    fmt.Errorf("测试超时 (%v)", testTimeout),
				}
				r := &result{Name: test.Name, Value: timeoutResult}
				select {
				case resultChan <- r:
				default:
				}
			}

			activeTestsMutex.Lock()
			delete(activeTests, test.Name)
			activeTestsMutex.Unlock()
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	resultsByRegion := make(map[string][]*result)
	for completedTests < len(allTests) {
		r, ok := <-resultChan
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

	ipTypeStr := fmt.Sprintf("IPv%d", ipType)
	if ipType == 0 {
		ipTypeStr = "Auto"
	}

	for _, regionName := range regionOrder {
		if regionItems, exists := regionSubGroups[regionName]; exists && len(regionItems) > 0 {

			ResultLines = append(ResultLines, &result{Name: fmt.Sprintf("%s (%s)", regionName, ipTypeStr), Divider: true})

			for _, item := range regionItems {
				if test, ok := item.(m.TestItem); ok {
					if test.Func == nil {

						ResultLines = append(ResultLines, &result{Name: test.Name, Divider: true})
					} else {

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

	stopProgressUpdater()
	if bar != nil {
		bar.Finish()
	}

	activeTestsMutex.Lock()
	for testName := range activeTests {
		delete(activeTests, testName)
	}
	activeTestsMutex.Unlock()

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

func ExecuteTests(regions []regionItem, client core.HttpClient, ipType int) {
	startTime := time.Now()

	for _, region := range regions {
		if !region.Enabled {
			continue
		}
		regionStartTime := time.Now()

		ipTypeStr := fmt.Sprintf("IPv%d", ipType)
		if ipType == 0 {
			ipTypeStr = "Auto"
		}
		fmt.Printf("\n正在检测 %s (%s) ...\n", region.Name, ipTypeStr)

		ResultLines = append(ResultLines, &result{Name: fmt.Sprintf("%s (%s)", region.Name, ipTypeStr), Divider: true})
		if region.Name == "Korea" && ipType == 6 {
			ResultLines = append(ResultLines, &result{Name: "No Korean platform supports IPv6", Divider: false})
		}

		ctx, cancel := context.WithTimeout(context.Background(), regionTimeout)
		defer cancel()

		totalTests := 0
		var validTests []m.TestItem
		for _, test := range region.Tests {
			if test.Func != nil && !(ipType == 6 && !test.SupportsV6) {
				totalTests++
				validTests = append(validTests, test)
			}
		}

		if totalTests == 0 {
			continue
		}

		bar = newProgressBar(int64(totalTests), "正在测试...")

		startProgressUpdater()

		maxWorkers := 30
		if Conc > 0 {
			maxWorkers = int(Conc)
		} else if totalTests > 50 {
			maxWorkers = 40
		} else if totalTests < 20 {
			maxWorkers = 25
		}

		workerPool := make(chan struct{}, maxWorkers)
		resultChan := make(chan *result, totalTests)
		completedTests := 0

		// 使用WaitGroup确保所有goroutine完成
		var wg sync.WaitGroup

		for _, test := range validTests {

			if test.Func == nil {
				ResultLines = append(ResultLines, &result{Name: test.Name, Divider: true})
				continue
			}

			workerPool <- struct{}{}
			wg.Add(1)

			go func(test m.TestItem) {
				defer func() {

					<-workerPool
					wg.Done()
				}()

				activeTestsMutex.Lock()
				activeTests[test.Name] = true
				activeTestsMutex.Unlock()

				testCtx, testCancel := context.WithTimeout(ctx, testTimeout)
				defer testCancel()

				cacheKey := fmt.Sprintf("%s_%d", test.Name, ipType)
				cacheMutex.RLock()
				if cachedResult, exists := resultCache[cacheKey]; exists {
					cacheMutex.RUnlock()
					r := &result{Name: test.Name, Value: cachedResult}
					select {
					case resultChan <- r:
					case <-testCtx.Done():

					}

					activeTestsMutex.Lock()
					delete(activeTests, test.Name)
					activeTestsMutex.Unlock()
					return
				}
				cacheMutex.RUnlock()

				done := make(chan core.Result, 1)
				go func() {
					defer func() {
						if r := recover(); r != nil {

							select {
							case done <- core.Result{
								Status: core.StatusFailed,
								Err:    fmt.Errorf("测试panic: %v", r),
							}:
							default:

							}
						}
					}()

					core.ResetSessionHeaders()
					client = core.NewHttpClient(ipType)
					result := test.Func(client)
					result.CachedResult = false

					cacheMutex.Lock()
					resultCache[cacheKey] = result
					cacheMutex.Unlock()

					select {
					case done <- result:
					default:

					}
				}()

				select {
				case testResult := <-done:

					r := &result{Name: test.Name, Value: testResult}
					select {
					case resultChan <- r:
					case <-testCtx.Done():

					}

				case <-testCtx.Done():

					timeoutResult := core.Result{
						Status: core.StatusFailed,
						Err:    fmt.Errorf("测试超时 (%v)", testTimeout),
					}
					r := &result{Name: test.Name, Value: timeoutResult}
					select {
					case resultChan <- r:
					default:

					}
				}

				activeTestsMutex.Lock()
				delete(activeTests, test.Name)
				activeTestsMutex.Unlock()
			}(test)
		}

		go func() {
			wg.Wait()
			close(resultChan)
		}()

		for completedTests < totalTests {
			select {
			case r, ok := <-resultChan:
				if !ok {

					break
				}
				ResultLines = append(ResultLines, r)
				completedTests++

				if bar != nil {

					bar.Add(1)
				}

			case <-ctx.Done():

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
			stopProgressUpdater()
			bar.Finish()
		}

		activeTestsMutex.Lock()
		for testName := range activeTests {
			delete(activeTests, testName)
		}
		activeTestsMutex.Unlock()
	}

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
