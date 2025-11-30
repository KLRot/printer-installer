package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed NotoSansSC-Regular.otf
var embeddedFontData []byte

// myLightTheme 自定义亮色主题
type myLightTheme struct {
	regular fyne.Resource
	bold    fyne.Resource
}

var (
	// 定义常见的中文字体路径
	// 注意：只包含支持中文的字体！
	fontPaths = []string{
		// 麒麟/UKUI 系统字体 (优先级最高)
		"/usr/share/fonts/truetype/ukui/ukui-default.ttf",
		"/usr/share/fonts/ukui/ukui-default.ttf",
		"/usr/share/fonts/truetype/kylin-font/kylin-font.ttf",
		
		// 通用 Linux 中文字体
		"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",
		"/usr/share/fonts/noto-cjk/NotoSansCJK-Regular.ttc",
		"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
		"/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc",
		"/usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf",
		"/usr/share/fonts/truetype/arphic/uming.ttc",
		"/usr/share/fonts/truetype/arphic/ukai.ttc",
		
		// Windows 兼容
		"C:\\Windows\\Fonts\\msyh.ttc",
		"C:\\Windows\\Fonts\\simhei.ttf",
	}
)

func newLightTheme() *myLightTheme {
	theme := &myLightTheme{}
	theme.loadFonts()
	return theme
}

func (m *myLightTheme) loadFonts() {
	// 1. 优先使用嵌入的字体（最可靠）
	if len(embeddedFontData) > 0 {
		m.regular = fyne.NewStaticResource("NotoSansSC-Regular.otf", embeddedFontData)
		m.bold = m.regular // 使用同一字体
		fmt.Println("✓ 使用嵌入的中文字体 (Noto Sans SC)")
		return
	}

	// 2. 检查环境变量 FYNE_FONT
	if envFont := os.Getenv("FYNE_FONT"); envFont != "" {
		if _, err := os.Stat(envFont); err == nil {
			if fontData, err := os.ReadFile(envFont); err == nil {
				m.regular = fyne.NewStaticResource("regular", fontData)
				m.bold = fyne.NewStaticResource("bold", fontData)
				fmt.Printf("✓ 使用环境变量指定的字体: %s\n", envFont)
				return
			}
		}
	}

	// 3. 尝试加载系统字体（后备方案）
	for _, path := range fontPaths {
		if _, err := os.Stat(path); err == nil {
			if fontData, err := os.ReadFile(path); err == nil {
				m.regular = fyne.NewStaticResource("regular", fontData)
				m.bold = fyne.NewStaticResource("bold", fontData)
				fmt.Printf("✓ 成功加载系统字体: %s\n", path)
				return
			}
		}
	}
	
	// 4. 如果都没找到，使用默认字体（可能不支持中文）
	fmt.Println("! 警告: 未找到中文字体")
	fmt.Println("! 请安装字体: sudo apt-get install fonts-wqy-microhei")
}

// 自定义颜色
var (
	kylinBlue   = color.RGBA{R: 40, G: 102, B: 255, A: 255} // 麒麟蓝
	lightBg     = color.RGBA{R: 248, G: 250, B: 252, A: 255} // 浅灰背景
	headerColor = color.RGBA{R: 30, G: 41, B: 59, A: 255}    // 深色标题
)

func (m *myLightTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return kylinBlue
	case theme.ColorNameBackground:
		return lightBg
	case theme.ColorNameInputBackground:
		return color.White
	}
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (m *myLightTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *myLightTheme) Font(style fyne.TextStyle) fyne.Resource {
	if m.regular != nil {
		if style.Bold {
			return m.bold
		}
		return m.regular
	}
	return theme.DefaultTheme().Font(style)
}

func (m *myLightTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 14 // 稍微增大默认字体
	}
	return theme.DefaultTheme().Size(name)
}

// ... (PrinterConfig 等结构体定义保持不变) ...

// initUI 初始化用户界面
func (gui *PrinterInstallerGUI) initUI() {
	// 1. 标题区域 (使用 canvas.Text 实现大字体)
	titleText := canvas.NewText("麒麟系统打印机自动安装工具", kylinBlue)
	titleText.TextSize = 28 // 大字体
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.Alignment = fyne.TextAlignCenter
	
	subTitle := widget.NewLabel("快速 • 智能 • 自动")
	subTitle.Alignment = fyne.TextAlignCenter
	
	headerBox := container.NewVBox(
		container.NewCenter(titleText),
		subTitle,
		widget.NewSeparator(),
	)
	
	// 2. 地点选择部分
	locationLabel := widget.NewLabel("📍 选择安装地点:")
	locationLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	gui.locationSelect = widget.NewSelect([]string{}, gui.onLocationChanged)
	gui.locationSelect.PlaceHolder = "请选择您的办公区域..."
	
	gui.refreshBtn = widget.NewButtonWithIcon("刷新配置", theme.ViewRefreshIcon(), func() {
		go gui.loadConfig()
	})
	
	locationBox := container.NewBorder(
		nil, nil,
		locationLabel,
		gui.refreshBtn,
		gui.locationSelect,
	)
	
	// 给地点选择加一个带边框的卡片效果
	locationCard := widget.NewCard("", "", container.NewPadded(locationBox))
	
	// 3. 打印机列表（使用 List + 复选框）
	gui.printerTable = widget.NewList(
		func() int {
			return len(gui.printerData)
		},
		func() fyne.CanvasObject {
			// CreateItem: 创建列表项模板
			check := widget.NewCheck("", nil)
			check.Resize(fyne.NewSize(30, 20))
			
			// 使用 canvas.Text 可以设置颜色
			nameText := canvas.NewText("打印机名称", headerColor)
			nameText.TextSize = 16
			nameText.TextStyle = fyne.TextStyle{Bold: true}
			
			modelLabel := widget.NewLabel("型号")
			ipLabel := widget.NewLabel("IP")
			
			// 布局: [Check] [Name]
			//               [Model] - [IP]
			infoBox := container.NewVBox(
				nameText,
				container.NewHBox(modelLabel, widget.NewLabel("-"), ipLabel),
			)
			
			return container.NewHBox(check, infoBox)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// UpdateItem: 更新数据
			if id >= len(gui.printerData) {
				return
			}
			
			printer := gui.printerData[id]
			
			// item 是 HBox
			box := item.(*fyne.Container)
			
			// 1. 复选框
			if len(box.Objects) > 0 {
				if check, ok := box.Objects[0].(*widget.Check); ok {
					check.Checked = gui.checkedItems[id]
					check.OnChanged = func(checked bool) {
						gui.mutex.Lock()
						gui.checkedItems[id] = checked
						gui.mutex.Unlock()
						gui.updateInstallBtnState()
					}
					check.Refresh()
				}
			}
			
			// 2. 信息区域
			if len(box.Objects) > 1 {
				if infoBox, ok := box.Objects[1].(*fyne.Container); ok {
					// infoBox [0]NameText, [1]DetailBox
					if len(infoBox.Objects) > 0 {
						if nameText, ok := infoBox.Objects[0].(*canvas.Text); ok {
							nameText.Text = printer.Name
							nameText.Refresh()
						}
					}
					
					if len(infoBox.Objects) > 1 {
						if detailBox, ok := infoBox.Objects[1].(*fyne.Container); ok {
							if len(detailBox.Objects) > 0 {
								detailBox.Objects[0].(*widget.Label).SetText(printer.Model)
							}
							if len(detailBox.Objects) > 2 {
								detailBox.Objects[2].(*widget.Label).SetText(printer.IP)
							}
						}
					}
				}
			}
		},
	)
	
	printerCard := widget.NewCard("可用打印机", "", gui.printerTable)

	
	// 4. 全选/全不选按钮
	gui.selectAllBtn = widget.NewButton("全选", gui.selectAll)
	gui.deselectAllBtn = widget.NewButton("全不选", gui.deselectAll)
	
	selectBtnBox := container.NewHBox(
		gui.selectAllBtn,
		gui.deselectAllBtn,
	)
	
	// 5. 进度条（默认隐藏）
	gui.progressBar = widget.NewProgressBar()
	gui.progressBar.Hide()
	
	// 6. 底部操作按钮
	gui.statusLabel = widget.NewLabel("")
	gui.statusLabel.Bind(gui.statusText)
	
	gui.installBtn = widget.NewButtonWithIcon("安装选中的打印机", theme.ConfirmIcon(), gui.installPrinters)
	gui.installBtn.Importance = widget.HighImportance
	gui.installBtn.Disable()
	
	exitBtn := widget.NewButton("退出", func() {
		gui.app.Quit()
	})
	
	actionBox := container.NewBorder(
		nil, nil,
		gui.statusLabel,
		container.NewHBox(gui.installBtn, exitBtn),
	)
	
	// 组合所有组件
	content := container.NewBorder(
		container.NewVBox(
			headerBox,
			locationCard,
		),
		container.NewVBox(
			selectBtnBox,
			gui.progressBar,
			actionBox,
		),
		nil, nil,
		printerCard,
	)
	
	gui.window.SetContent(content)
}

// loadConfig 从服务器加载配置文件
func (gui *PrinterInstallerGUI) loadConfig() {
	gui.statusText.Set("正在加载配置文件...")
	gui.refreshBtn.Disable()
	
	resp, err := http.Get(gui.configURL)
	if err != nil {
		gui.refreshBtn.Enable()
		gui.statusText.Set("配置加载失败")
		dialog.ShowError(fmt.Errorf("无法加载配置: %v", err), gui.window)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		gui.refreshBtn.Enable()
		gui.statusText.Set("配置加载失败")
		dialog.ShowError(fmt.Errorf("读取配置失败: %v", err), gui.window)
		return
	}
	
	var config PrinterConfig
	if err := json.Unmarshal(body, &config); err != nil {
		gui.refreshBtn.Enable()
		gui.statusText.Set("配置加载失败")
		dialog.ShowError(fmt.Errorf("解析配置失败: %v", err), gui.window)
		return
	}
	
	gui.config = &config
	gui.updateLocations()
	gui.refreshBtn.Enable()
}

// updateLocations 更新地点列表
func (gui *PrinterInstallerGUI) updateLocations() {
	if gui.config == nil {
		return
	}
	
	locations := make([]string, 0, len(gui.config.Locations))
	for location := range gui.config.Locations {
		locations = append(locations, location)
	}
	sort.Strings(locations)
	
	gui.locationSelect.Options = locations
	
	if len(locations) > 0 {
		gui.locationSelect.SetSelected(locations[0])
		gui.statusText.Set(fmt.Sprintf("配置加载成功 - 共 %d 个地点", len(locations)))
		
		// 移除成功弹窗，避免打扰用户
		// dialog.ShowInformation("成功", "配置文件加载成功", gui.window)
	} else {
		gui.statusText.Set("配置加载成功 - 但没有地点数据")
		dialog.ShowInformation("警告", "配置文件中没有找到任何地点信息", gui.window)
	}
}

// onLocationChanged 地点选择变化时更新打印机列表
func (gui *PrinterInstallerGUI) onLocationChanged(location string) {
	if location == "" || gui.config == nil {
		return
	}
	
	gui.mutex.Lock()
	gui.printerData = gui.config.Locations[location]
	gui.checkedItems = make(map[int]bool)
	gui.mutex.Unlock()
	
	gui.printerTable.Refresh()
	gui.updateInstallBtnState()
	gui.statusText.Set(fmt.Sprintf("已加载 %d 台打印机", len(gui.printerData)))
}



// selectAll 全选
func (gui *PrinterInstallerGUI) selectAll() {
	gui.mutex.Lock()
	for i := range gui.printerData {
		gui.checkedItems[i] = true
	}
	gui.mutex.Unlock()
	
	gui.printerTable.Refresh()
	gui.updateInstallBtnState()
}

// deselectAll 全不选
func (gui *PrinterInstallerGUI) deselectAll() {
	gui.mutex.Lock()
	gui.checkedItems = make(map[int]bool)
	gui.mutex.Unlock()
	
	gui.printerTable.Refresh()
	gui.updateInstallBtnState()
}

// updateInstallBtnState 更新安装按钮状态
func (gui *PrinterInstallerGUI) updateInstallBtnState() {
	count := 0
	for _, checked := range gui.checkedItems {
		if checked {
			count++
		}
	}
	
	if count > 0 {
		gui.installBtn.Enable()
		gui.installBtn.SetText(fmt.Sprintf("安装选中的打印机 (%d)", count))
	} else {
		gui.installBtn.Disable()
		gui.installBtn.SetText("安装选中的打印机")
	}
}

// installPrinters 安装选中的打印机
func (gui *PrinterInstallerGUI) installPrinters() {
	selectedPrinters := make([]Printer, 0)
	
	gui.mutex.Lock()
	for i, checked := range gui.checkedItems {
		if checked && i < len(gui.printerData) {
			selectedPrinters = append(selectedPrinters, gui.printerData[i])
		}
	}
	gui.mutex.Unlock()
	
	if len(selectedPrinters) == 0 {
		return
	}
	
	// 确认对话框
	confirmMsg := fmt.Sprintf("确定要安装 %d 台打印机吗?", len(selectedPrinters))
	dialog.ShowConfirm("确认安装", confirmMsg, func(confirmed bool) {
		if confirmed {
			go gui.installProcess(selectedPrinters)
		}
	}, gui.window)
}

// installProcess 安装过程
func (gui *PrinterInstallerGUI) installProcess(printers []Printer) {
	// 显示进度条
	gui.progressBar.Show()
	gui.progressBar.Max = float64(len(printers))
	gui.progressBar.SetValue(0)
	gui.installBtn.Disable()
	
	successCount := 0
	failedPrinters := make([]string, 0)
	
	for i, printer := range printers {
		// 更新进度
		gui.statusText.Set(fmt.Sprintf("正在安装: %s...", printer.Name))
		gui.progressBar.SetValue(float64(i))
		
		success, errMsg := gui.installSinglePrinter(printer)
		if success {
			successCount++
		} else {
			failedPrinters = append(failedPrinters, fmt.Sprintf("%s: %s", printer.Name, errMsg))
		}
	}
	
	// 完成
	gui.progressBar.Hide()
	gui.updateInstallBtnState()
	gui.statusText.Set(fmt.Sprintf("安装完成 - 成功: %d, 失败: %d", successCount, len(failedPrinters)))
	
	// 显示结果
	resultMsg := fmt.Sprintf("安装完成!\n\n成功: %d 台\n失败: %d 台", successCount, len(failedPrinters))
	if len(failedPrinters) > 0 {
		resultMsg += "\n\n失败详情:\n"
		displayCount := len(failedPrinters)
		if displayCount > 5 {
			displayCount = 5
		}
		resultMsg += strings.Join(failedPrinters[:displayCount], "\n")
		if len(failedPrinters) > 5 {
			resultMsg += fmt.Sprintf("\n... 还有 %d 台", len(failedPrinters)-5)
		}
		dialog.ShowInformation("安装结果", resultMsg, gui.window)
	} else {
		dialog.ShowInformation("安装结果", resultMsg, gui.window)
	}
}

// installSinglePrinter 安装单台打印机
func (gui *PrinterInstallerGUI) installSinglePrinter(printer Printer) (bool, string) {
	// 获取 PPD URL
	ppdURL := ""
	if gui.config != nil {
		if modelInfo, ok := gui.config.PrinterModels[printer.Model]; ok {
			ppdURL = modelInfo.PPDURL
		}
	}
	
	if ppdURL == "" {
		return false, fmt.Sprintf("配置文件中未找到型号 '%s' 的ppd_url，请在服务器的printer_config.json中配置", printer.Model)
	}
	
	// 对URL中的非ASCII字符进行编码
	if idx := strings.LastIndex(ppdURL, "/"); idx != -1 {
		baseURL := ppdURL[:idx]
		filename := ppdURL[idx+1:]
		encodedFilename := url.PathEscape(filename)
		ppdURL = baseURL + "/" + encodedFilename
	}
	
	// 下载 PPD 文件
	tempFile, err := os.CreateTemp("", "printer-*.ppd")
	if err != nil {
		return false, fmt.Sprintf("创建临时文件失败: %v", err)
	}
	tempPPDPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPPDPath)
	
	resp, err := http.Get(ppdURL)
	if err != nil {
		return false, fmt.Sprintf("下载PPD文件失败 (%s): %v", ppdURL, err)
	}
	defer resp.Body.Close()
	
	outFile, err := os.Create(tempPPDPath)
	if err != nil {
		return false, fmt.Sprintf("创建PPD文件失败: %v", err)
	}
	
	_, err = io.Copy(outFile, resp.Body)
	outFile.Close()
	if err != nil {
		return false, fmt.Sprintf("保存PPD文件失败: %v", err)
	}
	
	// 检查打印机是否已存在
	checkCmd := exec.Command("lpstat", "-p", printer.Name)
	if err := checkCmd.Run(); err == nil {
		// 打印机已存在，先删除
		deleteCmd := exec.Command("lpadmin", "-x", printer.Name)
		deleteCmd.Run()
	}
	
	// 设置打印机 URI
	printerURI := printer.URI
	if printerURI == "" {
		printerURI = fmt.Sprintf("ipp://%s/ipp/print", printer.IP)
	}
	
	// 安装打印机
	installCmd := exec.Command(
		"lpadmin",
		"-p", printer.Name,
		"-v", printerURI,
		"-P", tempPPDPath,
		"-E",
		"-D", fmt.Sprintf("%s (%s)", printer.Name, printer.Model),
	)
	
	output, err := installCmd.CombinedOutput()
	if err != nil {
		errMsg := string(output)
		if errMsg == "" {
			errMsg = "未知错误"
		}
		return false, errMsg
	}
	
	return true, ""
}

func main() {
	// 捕获 Panic 并写入日志文件
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("程序发生严重错误: %v\n堆栈信息:\n%s", r, string(debug.Stack()))
			fmt.Println(err)
			
			// 写入 crash.log
			f, _ := os.OpenFile("crash.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if f != nil {
				f.WriteString(fmt.Sprintf("\n[%s] %v\n", time.Now().Format(time.RFC3339), err))
				f.Close()
			}
			
			// 尝试显示错误对话框（如果 UI 还没死）
			// 注意：如果 Fyne 驱动已经崩溃，这可能不起作用
			os.Exit(1)
		}
	}()

	gui := NewPrinterInstallerGUI()
	
	// 设置退出时的清理工作
	gui.app.Lifecycle().SetOnStopped(func() {
		// 这里可以添加清理代码
		fmt.Println("程序正在退出...")
	})
	
	gui.Run()
}
