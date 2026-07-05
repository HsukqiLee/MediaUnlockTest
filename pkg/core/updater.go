package core

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/charmbracelet/glamour"
	selfUpdate "github.com/inconshreveable/go-update"
	"github.com/schollz/progressbar/v3"
)

func init() {
	// Clean up any .old files left behind by go-update after a successful self-update.
	// Windows locks the currently running executable, so go-update renames it to .old.
	// We delete it upon the next startup.
	if exe, err := os.Executable(); err == nil {
		oldPath := filepath.Join(filepath.Dir(exe), "."+filepath.Base(exe)+".old")
		_ = os.Remove(oldPath)
	}
}

type UpdateConfig struct {
	AppName         string
	VersionURL      string
	BinaryURLPrefix string
	Silent          bool
	ForceUpdate     bool
	JustCheck       bool
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
	// Detect `go run` by checking if the executable resides in the temp directory
	exe, err := os.Executable()
	if err == nil {
		tmpDir := os.TempDir()
		if strings.HasPrefix(filepath.ToSlash(exe), filepath.ToSlash(tmpDir)) {
			if !cfg.JustCheck {
				log.Println("[ERR] 检测到 go run 环境，不支持自动更新，请编译后再使用 -u")
			}
			return false
		}
	}

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

		if cfg.JustCheck {
			notesURL := "https://unlock.icmp.ing/api/release-notes"
			if resp, err := http.Get(notesURL); err == nil {
				if b, err := io.ReadAll(resp.Body); err == nil && len(b) > 0 {
					printReleaseNotes(string(b))
				}
				resp.Body.Close()
			}
			fmt.Printf("\n提示: 发现新版本，请运行 %s -u 进行更新\n", cfg.AppName)
			return false
		}
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
			if resp.StatusCode == http.StatusNotFound {
				log.Println("[ERR] 下载失败: GitHub Actions 构建可能仍在运行中，请稍候再试")
			} else {
				log.Println("[ERR] 下载时出错: 非预期的状态码", resp.StatusCode)
			}
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
		notesURL := "https://unlock.icmp.ing/api/release-notes"
		if resp, err := http.Get(notesURL); err == nil {
			if b, err := io.ReadAll(resp.Body); err == nil && len(b) > 0 {
				printReleaseNotes(string(b))
			}
			resp.Body.Close()
		}
	} else {
		log.Println("[OK]", cfg.AppName, "后台更新成功")
	}
	return true
}

func printReleaseNotes(md string) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width < 40 {
		width = 100 // Fallback to a reasonable width if we can't get terminal size or it's unreasonably small
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)

	var rendered string
	if err == nil {
		rendered, err = r.Render(md)
	}

	if err != nil {
		// Fallback: strip common markdown syntax manually
		md = strings.ReplaceAll(md, "**", "")
		md = strings.ReplaceAll(md, "`", "")
		fmt.Println("\n更新日志：\n" + strings.TrimSpace(md))
	} else {
		// Collapse 3+ consecutive newlines → 2 (avoids double blank lines between heading and list)
		rendered = regexp.MustCompile(`\n{3,}`).ReplaceAllString(rendered, "\n\n")
		fmt.Println("\n更新日志：")
		fmt.Print(rendered)
	}
	fmt.Println("--------------------------------")
}
