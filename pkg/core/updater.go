package core

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	selfUpdate "github.com/inconshreveable/go-update"
	"github.com/schollz/progressbar/v3"
)

type UpdateConfig struct {
	AppName         string
	VersionURL      string
	BinaryURLPrefix string
	Silent          bool
	ForceUpdate     bool
}

type Downloader struct {
	io.Reader
	Total   uint64
	Current uint64
	Pb      *progressbar.ProgressBar
	Done    bool
	Silent  bool
}

func (d *Downloader) Read(p []byte) (n int, err error) {
	n, err = d.Reader.Read(p)
	d.Current += uint64(n)
	if d.Done {
		return
	}
	if !d.Silent && d.Pb != nil {
		d.Pb.Add(n)
	}
	if d.Current == d.Total {
		d.Done = true
		if !d.Silent && d.Pb != nil {
			d.Pb.Describe("下载完成")
			d.Pb.Finish()
		}
	}
	return
}

type BarWriter struct {
	bar *progressbar.ProgressBar
}

func (bw *BarWriter) Write(p []byte) (n int, err error) {
	return bw.bar.Write(p)
}

// CheckUpdate checks and applies an update based on the provided configuration.
// Returns true if an update was successfully applied.
func CheckUpdate(cfg UpdateConfig) bool {
	resp, err := http.Get(cfg.VersionURL)
	if err != nil {
		log.Println("[ERR] 获取版本信息时出错:", err)
		return false
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERR] 读取版本信息时出错:", err)
		return false
	}

	parts := strings.Split(string(b), "-")
	if len(parts) != 2 {
		log.Println("[ERR] 版本号格式错误")
		return false
	}
	version := parts[0]

	if !cfg.ForceUpdate && strings.TrimPrefix(version, "v") == strings.TrimPrefix(Version, "v") {
		if !cfg.Silent {
			fmt.Println("已经是最新版本")
		}
		return false
	}

	timestampInt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Println("[ERR] 版本号时间戳错误:", err)
		return false
	}
	timestamp := time.Unix(timestampInt, 0)

	if !cfg.Silent {
		fmt.Println("最新版本：", version)
		fmt.Println("发布时间：", timestamp.Format("2006-01-02 15:04:05"))
		fmt.Println("运行系统：", runtime.GOOS)
		fmt.Println("运行架构：", runtime.GOARCH)
	}

	OS, ARCH := runtime.GOOS, runtime.GOARCH
	if OS == "android" && strings.Contains(os.Getenv("PREFIX"), "com.termux") {
		target_path := os.Getenv("PREFIX") + "/bin"
		out, err := os.Create(target_path + "/" + cfg.AppName + "_new")
		if err != nil {
			log.Println("[ERR] 创建文件出错:", err)
			return false
		}
		defer out.Close()
		if !cfg.Silent {
			log.Println("下载", cfg.AppName, "中 ...")
		}
		url := cfg.BinaryURLPrefix + "_" + OS + "_" + ARCH
		resp, err = http.Get(url)
		if err != nil {
			log.Println("[ERR] 下载时出错:", err)
			return false
		}
		defer resp.Body.Close()
		
		var pbBar *progressbar.ProgressBar
		if !cfg.Silent {
			pbBar = progressbar.DefaultBytes(resp.ContentLength, "下载进度")
		}
		downloader := &Downloader{
			Reader: resp.Body,
			Total:  uint64(resp.ContentLength),
			Pb:     pbBar,
			Silent: cfg.Silent,
		}
		if _, err := io.Copy(out, downloader); err != nil {
			log.Println("[ERR] 下载时出错:", err)
			return false
		}
		if err := os.Chmod(target_path+"/"+cfg.AppName+"_new", 0777); err != nil {
			log.Println("[ERR] 更改后端权限出错:", err)
			return false
		}
		if _, err := os.Stat(target_path + "/" + cfg.AppName); err == nil {
			if err := os.Remove(target_path + "/" + cfg.AppName); err != nil {
				log.Println("[ERR] 删除旧版本时出错:", err.Error())
				return false
			}
		}
		if err := os.Rename(target_path+"/"+cfg.AppName+"_new", target_path+"/"+cfg.AppName); err != nil {
			log.Println("[ERR] 更新后端时出错:", err)
			return false
		}
	} else {
		url := cfg.BinaryURLPrefix + "_" + OS + "_" + ARCH
		if OS == "windows" {
			url += ".exe"
		}

		resp, err = http.Get(url)
		if err != nil {
			log.Println("[ERR] 下载时出错:", err)
			return false
		}
		defer resp.Body.Close()

		var body io.Reader = resp.Body
		if !cfg.Silent {
			bar := progressbar.DefaultBytes(
				resp.ContentLength,
				"下载进度",
			)
			barWrapper := &BarWriter{bar: bar}
			body = io.TeeReader(resp.Body, barWrapper)
		}

		if resp.StatusCode != http.StatusOK {
			log.Println("[ERR] 下载时出错: 非预期的状态码", resp.StatusCode)
			return false
		}

		err = selfUpdate.Apply(body, selfUpdate.Options{})
		if err != nil {
			log.Println("[ERR] 更新时出错:", err)
			return false
		}
	}
	if !cfg.Silent {
		fmt.Println("[OK]", cfg.AppName, "更新成功")
	} else {
		log.Println("[OK]", cfg.AppName, "后台更新成功")
	}
	return true
}
