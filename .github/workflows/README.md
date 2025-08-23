# GitHub Actions Workflows

本目录包含项目的GitHub Actions自动化工作流配置。

## Release.yml

### 概述
`Release.yml` 是一个用于自动化构建和发布的工作流，当创建预发布或正式发布时自动触发。

### 主要特性

#### 🚀 并行构建
- 使用矩阵构建策略，同时编译多个平台和架构
- 支持 `cli` 和 `monitor` 两个项目
- 大幅提升构建效率，减少总体构建时间

#### 🌍 多平台支持
**主要平台：**
- **Linux**: amd64, arm64, 386, arm7/6/5, loong64, mips系列, ppc64系列, riscv64, s390x
- **macOS**: amd64, arm64
- **Windows**: amd64, arm64, 386, arm7/6/5
- **FreeBSD**: amd64, arm64, 386, arm7/6/5, riscv64
- **NetBSD**: amd64, arm64, 386, arm7/6/5
- **OpenBSD**: amd64, arm64, 386, arm7/6/5, ppc64

**Android平台（需要CGO）：**
- arm64, amd64, arm7/6/5, 386

#### 🔧 智能构建
- 自动检测平台类型，设置相应的Go环境变量
- Android平台自动下载和配置NDK工具链
- 根据项目类型设置不同的构建标志
- 支持ARM架构的GOARM版本设置

#### 📦 自动化发布
- 自动更新版本号到源代码
- 构建完成后自动上传到GitHub Release
- 使用GitHub Artifacts进行文件管理
- 提供详细的构建和上传日志

### 触发条件
```yaml
on:
  release:
    types: [prereleased, released]
```

### 工作流程

#### 1. 构建阶段 (build)
- **并行执行**: 每个平台和项目的组合在独立的runner上并行构建
- **环境配置**: 自动设置Go环境变量和Android NDK（如需要）
- **交叉编译**: 支持跨平台编译，生成目标平台的二进制文件
- **文件命名**: 自动生成符合规范的二进制文件名

#### 2. 发布阶段 (release)
- **依赖关系**: 等待所有构建任务完成后开始
- **版本更新**: 自动更新源代码中的版本号
- **文件收集**: 下载所有构建产物
- **批量上传**: 自动上传所有二进制文件到GitHub Release

### 构建产物命名规则

#### CLI项目
- `unlock-test_<平台>_<架构>[.exe]`
- 例如: `unlock-test_linux_amd64`, `unlock-test_windows_amd64.exe`

#### Monitor项目
- `unlock-monitor_<平台>_<架构>[.exe]`
- 例如: `unlock-monitor_linux_amd64`, `unlock-monitor_windows_amd64.exe`

### 环境要求
- **Go版本**: 1.24+
- **运行环境**: Ubuntu Latest
- **权限**: 需要 `contents: write` 和 `actions: read` 权限

### 性能优化
- **矩阵构建**: 最多可同时运行 80+ 个构建任务
- **缓存策略**: 使用GitHub Artifacts进行文件缓存
- **并行执行**: 充分利用GitHub Actions的并行能力
- **资源优化**: 只在需要时下载Android NDK

### 使用建议
1. **发布前检查**: 确保所有代码已合并到主分支
2. **版本管理**: 使用语义化版本号（如 v1.2.3）
3. **测试验证**: 在预发布阶段验证构建产物
4. **资源监控**: 注意GitHub Actions的使用配额

### 故障排除
- **构建失败**: 检查Go版本兼容性和依赖项
- **上传失败**: 验证GitHub Token权限
- **版本更新失败**: 检查分支保护和权限设置
- **NDK下载失败**: 检查网络连接和NDK版本可用性


