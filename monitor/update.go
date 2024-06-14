package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
	"strings"
	"strconv"

	"github.com/kardianos/service"
	selfUpdate "github.com/inconshreveable/go-update"
	pb "github.com/schollz/progressbar/v3"
)

func checkUpdate() {
	resp, err := http.Get("https://unlock.icmp.ing/monitor/latest/version")
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
	
	if version == Version {
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
	    out, err := os.Create(target_path + "/unlock-monitor_new")
	    if err != nil {
	    	log.Fatal("[ERR] 创建文件出错:", err)
	    	return
	    }
	    defer out.Close()
	    log.Println("下载unlock-monitor中 ...")
	    url := "https://unlock.icmp.ing/monitor/latest/unlock-monitor_" + runtime.GOOS + "_" + runtime.GOARCH
    	resp, err = http.Get(url)
    	if err != nil {
	    	log.Fatal("[ERR] 下载unlock-monitor时出错:", err)
	    }
	    defer resp.Body.Close()
	    downloader := &Downloader{
	    	Reader: resp.Body,
	    	Total:  uint64(resp.ContentLength),
	    	Pb:     pb.DefaultBytes(resp.ContentLength, "下载进度"),
	    }
	    if _, err := io.Copy(out, downloader); err != nil {
	    	log.Fatal("[ERR] 下载unlock-monitor时出错:", err)
	    }
	    if os.Chmod(target_path + "/unlock-monitor_new", 0777) != nil {
	    	log.Fatal("[ERR] 更改unlock-monitor后端权限出错:", err)
	    }
	    if _, err := os.Stat(target_path + "/unlock-monitor"); err == nil {
    		if os.Remove(target_path + "/unlock-monitor") != nil {
	    		log.Fatal("[ERR] 删除unlock-monitor旧版本时出错:", err.Error())
	    	}
	    }
	    if os.Rename(target_path + "/unlock-monitor_new", target_path + "/unlock-monitor") != nil {
	    	log.Fatal("[ERR] 更新unlock-monitor后端时出错:", err)
	    }
	} else {
	    url := "https://unlock.icmp.ing/monitor/latest/unlock-monitor_" + OS + "_" + ARCH
	    if OS == "windows" {
	        url += ".exe"
	    }

	    resp, err = http.Get(url)
	    if err != nil {
	    	log.Fatal("[ERR] 下载unlock-monitor时出错:", err)
	    	return
	    }
	    defer resp.Body.Close()
	    
	    bar := pb.DefaultBytes(
	    	resp.ContentLength,
	    	"下载进度",
	    )

	    body := io.TeeReader(resp.Body, bar)

	    if resp.StatusCode != http.StatusOK {
	    	log.Fatal("[ERR] 下载unlock-monitor时出错: 非预期的状态码", resp.StatusCode)
	    	return
	    }

	    err = selfUpdate.Apply(body, selfUpdate.Options{})
	    if err != nil {
	    	log.Fatal("[ERR] 更新unlock-monitor时出错:", err)
	    	return
	    }
	}

	log.Println("[OK] unlock-monitor后端更新成功")

	serviceConfig := &service.Config{
		Name:        "unlock-monitor",
		DisplayName: "unlock-monitor",
		Description: "Service to monitor media unlock status.",
	}

	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		log.Fatal(err)
	}

	if serviceExists(s) && serviceIsActive(s) {
		restartService(s)
	} else {
		log.Println("[INFO] unlock-monitor 服务未运行")
	}
}

type Downloader struct {
	io.Reader
	Total   uint64
	Current uint64
	Pb      *pb.ProgressBar
	done    bool
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
		d.Pb.Describe("unlock-monitor下载完成")
		d.Pb.Finish()
	}
	return
}