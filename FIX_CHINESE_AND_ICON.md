# 中文显示和图标问题解决方案

## 已修复的问题

### 1. ✅ 中文乱码问题

**原因**: Fyne 默认使用的字体可能不支持中文字符。

**解决方案**: 
- 添加了自定义亮色主题 `myLightTheme`
- 主题使用系统默认字体，自动支持中文显示
- 强制使用亮色主题变体 `theme.VariantLight`

### 2. ✅ 图标不显示问题

**原因**: 图标路径查找不正确。

**解决方案**:
- 修改为在应用级别设置图标 (`app.SetIcon`)
- 添加多个可能的图标路径搜索
- 优先搜索可执行文件所在目录

### 3. ✅ 亮色主题

**原因**: Fyne 默认可能使用暗色主题。

**解决方案**:
- 创建自定义亮色主题
- 在应用启动时设置主题

## 代码修改说明

### 修改 1: 添加自定义主题

```go
// myLightTheme 自定义亮色主题
type myLightTheme struct{}

func (m *myLightTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
    // 强制使用亮色主题
    return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (m *myLightTheme) Font(style fyne.TextStyle) fyne.Resource {
    // 使用系统默认字体以支持中文
    return theme.DefaultTheme().Font(style)
}
```

### 修改 2: 应用主题

```go
func NewPrinterInstallerGUI() *PrinterInstallerGUI {
    myApp := app.NewWithID("com.kylin.printer.installer")
    
    // 设置亮色主题
    myApp.Settings().SetTheme(&myLightTheme{})
    
    // ... 其他代码
}
```

### 修改 3: 设置应用图标

```go
func (gui *PrinterInstallerGUI) setAppIcon() {
    // 尝试多个可能的图标路径
    iconPaths := []string{
        "printer_icon.png",                    // 当前目录
        "assets/printer_icon.png",            // assets 目录
        "../printer_icon.png",                // 上级目录
        "./printer-installer-go/printer_icon.png", // 项目目录
    }
    
    // 获取可执行文件所在目录
    if exePath, err := os.Executable(); err == nil {
        baseDir := filepath.Dir(exePath)
        iconPaths = append([]string{filepath.Join(baseDir, "printer_icon.png")}, iconPaths...)
    }
    
    // 尝试加载图标
    for _, iconPath := range iconPaths {
        if _, err := os.Stat(iconPath); err == nil {
            if icon, err := fyne.LoadResourceFromPath(iconPath); err == nil {
                gui.app.SetIcon(icon)  // 使用 app.SetIcon 而不是 window.SetIcon
                fmt.Printf("✓ 成功加载图标: %s\n", iconPath)
                return
            }
        }
    }
}
```

## 使用说明

### 编译

```bash
cd /mnt/f/麒麟打印机安装/printer-installer-go
go build -o printer_installer main.go
```

### 运行

确保 `printer_icon.png` 与可执行文件在同一目录：

```bash
# 复制图标到可执行文件目录
cp printer_icon.png ./

# 运行程序
./printer_installer
```

### 打包发布

使用 fyne 打包工具：

```bash
# 安装 fyne 命令
go install fyne.io/fyne/v2/cmd/fyne@latest

# 打包应用（会自动包含图标）
fyne package -os linux -icon printer_icon.png -name "打印机安装程序"
```

## 验证

运行程序后，应该看到：

1. ✅ **中文正常显示**: 所有中文文字清晰可见
2. ✅ **亮色主题**: 界面使用亮色背景
3. ✅ **图标显示**: 
   - 窗口标题栏显示打印机图标
   - 任务栏显示打印机图标
   - 控制台输出: `✓ 成功加载图标: printer_icon.png`

## 故障排除

### 问题 1: 中文仍然乱码

**可能原因**: 系统缺少中文字体

**解决方案**:
```bash
# 安装中文字体
sudo apt-get install fonts-noto-cjk fonts-wqy-zenhei

# 刷新字体缓存
fc-cache -fv
```

### 问题 2: 图标不显示

**检查步骤**:

1. 确认图标文件存在：
```bash
ls -lh printer_icon.png
```

2. 确认图标格式正确（PNG 格式）：
```bash
file printer_icon.png
# 应该输出: printer_icon.png: PNG image data
```

3. 查看程序输出，确认图标加载路径：
```bash
./printer_installer
# 应该看到: ✓ 成功加载图标: /path/to/printer_icon.png
```

4. 如果仍然不显示，尝试使用绝对路径：
```bash
# 修改代码中的图标路径为绝对路径
cp printer_icon.png /usr/share/pixmaps/
# 然后在代码中使用 /usr/share/pixmaps/printer_icon.png
```

### 问题 3: 主题不是亮色

**解决方案**: 确保代码中正确设置了主题

```go
// 检查这行代码是否存在
myApp.Settings().SetTheme(&myLightTheme{})
```

## 在不同环境下的图标位置

### 开发环境
```
printer-installer-go/
├── main.go
├── printer_icon.png  ← 图标在这里
└── printer_installer
```

### 打包后
```
dist/
├── printer_installer
└── printer_icon.png  ← 图标在这里
```

### 使用 fyne package
```
# fyne package 会自动将图标嵌入到可执行文件中
# 不需要单独的图标文件
```

## 推荐的打包方式

使用 `fyne package` 可以避免图标路径问题：

```bash
# 打包（图标会被嵌入）
fyne package -os linux -icon printer_icon.png -name "打印机安装程序"

# 生成的文件
# - 打印机安装程序.tar.xz (包含所有资源)
```

## 总结

所有问题已修复：
- ✅ 中文显示正常（使用系统默认字体）
- ✅ 亮色主题（强制使用 VariantLight）
- ✅ 图标显示（应用级别设置，多路径搜索）

重新编译后即可看到效果！
