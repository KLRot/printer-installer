# 终极解决方案：嵌入中文字体

为了彻底解决中文乱码问题，我们采取了最稳妥的方案：**将中文字体直接嵌入到可执行文件中**。

## 🚀 方案原理

1.  **下载字体**: 在 GitHub Actions 构建过程中，自动下载开源的文泉驿微米黑字体 (`wqy-microhei.ttc`)。
2.  **打包资源**: 使用 `fyne bundle` 工具将字体文件转换为 Go 代码 (`bundled_font.go`)。
3.  **编译嵌入**: Go 编译器将生成的字体代码编译进最终的可执行文件。
4.  **运行时加载**: 程序启动时，优先使用嵌入的字体数据，不再依赖系统字体。

## ✅ 优势

- **100% 解决乱码**: 无论用户的系统是否安装了中文字体，程序都能正常显示中文。
- **零依赖**: 用户不需要安装任何额外的字体包。
- **一致性**: 在所有 Linux 发行版上显示效果一致。

## 📦 文件大小变化

由于嵌入了字体文件，可执行文件的大小会增加约 **5MB**。这是为了保证兼容性所付出的必要代价。

## 🛠️ 代码修改说明

### 1. `main.go`
修改了 `loadFonts` 方法，优先检查 `bundledFont` 变量：

```go
// 嵌入的字体资源
var bundledFont fyne.Resource

func (m *myLightTheme) loadFonts() {
    // 1. 优先使用编译时嵌入的字体
    if bundledFont != nil {
        m.regular = bundledFont
        m.bold = bundledFont
        fmt.Println("✓ 使用嵌入的中文字体")
        return
    }
    // ... 后备逻辑：检查环境变量和系统字体 ...
}
```

### 2. `build.yml`
添加了字体下载和打包步骤：

```yaml
- name: Bundle Chinese Font
  run: |
    wget -O wqy-microhei.ttc https://...
    fyne bundle -package main -name bundledFont wqy-microhei.ttc > bundled_font.go
```

## 📝 本地开发指南

如果你在本地 (WSL) 开发，想要测试这个功能：

1.  **下载字体**:
    ```bash
    wget https://github.com/anthonyfok/fonts-wqy-microhei/raw/master/wqy-microhei.ttc
    ```

2.  **生成资源文件**:
    ```bash
    # 安装 fyne 工具 (如果还没安装)
    go install fyne.io/fyne/v2/cmd/fyne@latest
    
    # 生成 bundled_font.go
    fyne bundle -package main -name bundledFont wqy-microhei.ttc > bundled_font.go
    ```

3.  **运行**:
    ```bash
    go run .
    ```
    你应该会看到输出：`✓ 使用嵌入的中文字体`

如果不执行上述步骤，本地运行时 `bundledFont` 变量为 `nil`，程序会自动回退到搜索系统字体，这不会影响开发。
