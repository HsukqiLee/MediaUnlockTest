package main

import (
	"fmt"
	"io"
	"os"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	
	selfUpdate "github.com/inconshreveable/go-update"
	pb "github.com/schollz/progressbar/v3"
)

var (
	Version   = "1.4.5"
	buildTime string
)

func serviceExists(serviceName string) bool {
	cmd := exec.Command("systemctl", "is-enabled", serviceName)
	err := cmd.Run()
	return err == nil
}

func serviceIsActive(serviceName string) bool {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	err := cmd.Run()
	return err == nil
}

func restartService(serviceName string) error {
	cmd := exec.Command("systemctl", "restart", serviceName)
	return cmd.Run()
}

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
	version := string(b)
	if version == Version {
		fmt.Println("已经是最新版本")
		return
	}
	fmt.Println("检测到新版本", version)

	OS, ARCH := runtime.GOOS, runtime.GOARCH
	fmt.Println("OS:", OS)
	fmt.Println("ARCH:", ARCH)
	
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

	fmt.Println("[OK] unlock-monitor后端更新成功")
	
	if serviceExists("unlock-monitor.service") && serviceIsActive("unlock-monitor.service") {
		err = restartService("unlock-monitor.service")
		if err != nil {
			log.Fatal("[ERR] 重启unlock-monitor服务时出错:", err)
			return
		}
		fmt.Println("[OK] unlock-monitor服务已重启")
	} else {
		fmt.Println("[INFO] unlock-monitor服务不存在或未运行，无需重启")
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
	// fmt.Printf("\r正在下载，进度：%.2f%% [%s/%s]", float64(d.Current*10000/d.Total)/100, humanize.Bytes(d.Current), humanize.Bytes(d.Total))
	if d.Current == d.Total {
		d.done = true
		d.Pb.Describe("unlock-monitor下载完成")
		d.Pb.Finish()
	}
	return
}