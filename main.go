package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "sync"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"
)

type Printer struct {
    Name  string `json:"name"`
    Model string `json:"model"`
    IP    string `json:"ip"`
    PPD   string `json:"ppd"`
    URI   string `json:"uri"`
}

type PrinterModelCfg struct {
    PPDUrl string `json:"ppd_url"`
}

type Config struct {
    Locations     map[string][]Printer       `json:"locations"`
    PrinterModels map[string]PrinterModelCfg `json:"printer_models"`
}

// TappableBox 是一个自定义组件，允许容器响应点击事件
type TappableBox struct {
    widget.BaseWidget
    content  *fyne.Container
    OnTapped func()
}

func newTappableBox(content *fyne.Container, onTapped func()) *TappableBox {
    b := &TappableBox{content: content, OnTapped: onTapped}
    b.ExtendBaseWidget(b)
    return b
}

func (b *TappableBox) CreateRenderer() fyne.WidgetRenderer {
    return widget.NewSimpleRenderer(b.content)
}

func (b *TappableBox) Tapped(_ *fyne.PointEvent) {
    if b.OnTapped != nil {
        b.OnTapped()
    }
}

func main() {
    configURL := "http://10.245.93.86/printer/printer_config.json"

    a := app.New()
    w := a.NewWindow("打印机安装工具 (Go)")
    w.Resize(fyne.NewSize(1000, 720))

    locationSelect := widget.NewSelect([]string{}, func(s string) {})
    refreshBtn := widget.NewButton("刷新配置", nil)
    installBtn := widget.NewButton("安装选中的打印机", nil)
    installBtn.Disable()

    status := widget.NewLabel("就绪")

    printerListContainer := container.NewVBox()
    scroll := container.NewVScroll(printerListContainer)
    scroll.SetMinSize(fyne.NewSize(950, 520))

    var cfg Config
    var printers []Printer
    var checks []*widget.Check

    loadConfig := func() {
        status.SetText("正在下载配置...")
        resp, err := http.Get(configURL)
        if err != nil {
            dialog.ShowError(fmt.Errorf("无法下载配置: %v", err), w)
            status.SetText("配置加载失败")
            return
        }
        defer resp.Body.Close()
        body, _ := io.ReadAll(resp.Body)

        if err := json.Unmarshal(body, &cfg); err != nil {
            dialog.ShowError(fmt.Errorf("解析配置失败: %v", err), w)
            status.SetText("配置解析失败")
            return
        }

        locs := make([]string, 0, len(cfg.Locations))
        for k := range cfg.Locations {
            locs = append(locs, k)
        }
        locationSelect.Options = locs
        if len(locs) > 0 {
            locationSelect.SetSelected(locs[0])
        }
        locationSelect.Refresh()
        status.SetText(fmt.Sprintf("配置加载成功 - %d 个地点", len(locs)))
    }

    refreshBtn.OnTapped = func() { go loadConfig() }

    updatePrinterList := func(loc string) {
        printerListContainer.Objects = nil
        printers = cfg.Locations[loc]
        checks = make([]*widget.Check, len(printers))

        for i, p := range printers {
            check := widget.NewCheck("", func(bool) {})
            checks[i] = check

            nameLbl := widget.NewLabel(p.Name)
            modelLbl := widget.NewLabel(p.Model)
            ipLbl := widget.NewLabel(p.IP)

            rowContainer := container.NewHBox(
                check,
                container.NewVBox(nameLbl, modelLbl),
                layout.NewSpacer(),
                ipLbl,
            )
            
            // 使用自定义的可点击组件包裹行内容
            row := newTappableBox(rowContainer, func() {
                check.SetChecked(!check.Checked)
            })
            
            printerListContainer.Add(row)

        }
        printerListContainer.Refresh()
        installBtn.Enable()
    }

    locationSelect.OnChanged = func(s string) { updatePrinterList(s) }

    installBtn.OnTapped = func() {
        selected := []Printer{}
        for i, c := range checks {
            if c != nil && c.Checked {
                selected = append(selected, printers[i])
            }
        }
        if len(selected) == 0 {
            dialog.ShowInformation("提示", "请先选择至少一台打印机", w)
            return
        }
        confirm := dialog.NewConfirm("确认安装",
            fmt.Sprintf("确定安装 %d 台打印机？", len(selected)),
            func(ok bool) {
                if !ok { return }
                go runInstallFlow(selected, cfg, w, status)
            }, w)
        confirm.Show()
    }

    topBar := container.NewHBox(widget.NewLabel("地点："), locationSelect, refreshBtn, layout.NewSpacer(), status)
    bottomBar := container.NewHBox(layout.NewSpacer(), installBtn)
    content := container.NewBorder(topBar, bottomBar, nil, nil, scroll)
    w.SetContent(content)

    go loadConfig()
    w.ShowAndRun()
}

func runInstallFlow(selected []Printer, cfg Config, w fyne.Window, status *widget.Label) {
    var wg sync.WaitGroup
    mu := sync.Mutex{}
    successCnt := 0
    failList := []string{}

    status.SetText("开始安装...")
    for _, p := range selected {
        wg.Add(1)
        go func(pr Printer) {
            defer wg.Done()
            ok, err := installSingle(pr, cfg)
            mu.Lock()
            defer mu.Unlock()
            if ok {
                successCnt++
            } else {
                failList = append(failList, fmt.Sprintf("%s: %v", pr.Name, err))
            }
            status.SetText(fmt.Sprintf("安装中: %d 成功, %d 失败", successCnt, len(failList)))
        }(p)
    }
    wg.Wait()

    if len(failList) == 0 {
        dialog.ShowInformation("安装完成", fmt.Sprintf("全部安装成功，共 %d 台", successCnt), w)
        status.SetText("安装完成")
    } else {
        preview := strings.Join(failList, "\n")
        dialog.ShowError(fmt.Errorf("成功: %d，失败: %d\n%s", successCnt, len(failList), preview), w)
        status.SetText("安装完成（有失败）")
    }
}

func installSingle(p Printer, cfg Config) (bool, error) {
    modelCfg, ok := cfg.PrinterModels[p.Model]
    if !ok || modelCfg.PPDUrl == "" {
        return false, fmt.Errorf("未配置型号 %s 的 ppd_url", p.Model)
    }
    ppdURL := modelCfg.PPDUrl

    if strings.Contains(ppdURL, " ") {
        ppdURL = strings.ReplaceAll(ppdURL, " ", "%20")
    }

    tmpfile := filepath.Join(os.TempDir(), filepath.Base(ppdURL))
    out, err := os.Create(tmpfile)
    if err != nil { return false, err }
    defer func() { out.Close(); os.Remove(tmpfile) }()

    resp, err := http.Get(ppdURL)
    if err != nil { return false, err }
    defer resp.Body.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil { return false, err }
    out.Sync()

    uri := p.URI
    if uri == "" {
        uri = "ipp://" + p.IP + "/ipp/print"
    }

    checkCmd := exec.Command("lpstat", "-p", p.Name)
    if err := checkCmd.Run(); err == nil {
        _ = exec.Command("lpadmin", "-x", p.Name).Run()
    }

    cmd := exec.Command("lpadmin", "-p", p.Name, "-v", uri, "-P", tmpfile, "-E",
        "-D", fmt.Sprintf("%s (%s)", p.Name, p.Model))
    outb, err := cmd.CombinedOutput()
    if err != nil {
        return false, fmt.Errorf("%v: %s", err, string(outb))
    }
    return true, nil
}
