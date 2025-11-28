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
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/theme"
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

func main() {

    configURL := "http://10.245.93.86/printer/printer_config.json"

    a := app.New()
    w := a.NewWindow("打印机安装工具 (Go)")
    w.Resize(fyne.NewSize(1024, 720))
    w.CenterOnScreen()

    title := canvas.NewText("打印机安装工具", theme.ForegroundColor())
    title.TextStyle = fyne.TextStyle{Bold: true}
    title.Alignment = fyne.TextAlignLeading
    title.TextSize = 20

    // 控件
    status := widget.NewLabel("就绪")
    status.TextStyle = fyne.TextStyle{Bold: true}

    locationSelect := widget.NewSelect([]string{}, func(s string) {})
    locationSelect.PlaceHolder = "选择地点"

    refreshBtn := widget.NewButton("刷新配置", nil)
    refreshBtn.Importance = widget.HighImportance

    installBtn := widget.NewButton("安装选中的打印机", nil)
    installBtn.Importance = widget.HighImportance
    installBtn.Disable()

    printerListContainer := container.NewVBox()
    printerListContainer.Add(widget.NewLabel("请选择地点..."))

    scroll := container.NewVScroll(printerListContainer)
    scroll.SetMinSize(fyne.NewSize(980, 560))

    var cfg Config
    var printers []Printer
    var checks []*widget.Check

    loadConfig := func() {
        status.SetText("正在加载配置...")
        resp, err := http.Get(configURL)
        if err != nil {
            dialog.ShowError(fmt.Errorf("下载配置失败: %v", err), w)
            status.SetText("加载失败")
            return
        }
        defer resp.Body.Close()

        body, _ := io.ReadAll(resp.Body)

        if err := json.Unmarshal(body, &cfg); err != nil {
            dialog.ShowError(fmt.Errorf("配置文件解析失败: %v", err), w)
            status.SetText("解析失败")
            return
        }

        // 更新地点列表
        locs := make([]string, 0, len(cfg.Locations))
        for k := range cfg.Locations {
            locs = append(locs, k)
        }
        locationSelect.Options = locs
        if len(locs) > 0 {
            locationSelect.SetSelected(locs[0])
        }
        locationSelect.Refresh()

        status.SetText("配置加载成功")
    }

    refreshBtn.OnTapped = func() {
        go loadConfig()
    }

    updatePrinterList := func(loc string) {
        printerListContainer.Objects = nil
        printers = cfg.Locations[loc]
        checks = make([]*widget.Check, len(printers))

        for i, p := range printers {
            idx := i

            check := widget.NewCheck("", func(bool) {})
            checks[i] = check

            name := canvas.NewText(p.Name, theme.ForegroundColor())
            name.TextStyle = fyne.TextStyle{Bold: true}
            name.TextSize = 16

            model := widget.NewLabel("型号：" + p.Model)
            ip := widget.NewLabel("IP：" + p.IP)

            leftBox := container.NewVBox(name, model)
            rightBox := container.NewVBox(ip)

            card := container.NewBorder(nil, nil, check, nil,
                container.NewHBox(
                    leftBox,
                    layout.NewSpacer(),
                    rightBox,
                ),
            )

            cardBox := container.NewVBox(card)
            cardBox.Add(canvas.NewRectangle(theme.ShadowColor()))

            printerListContainer.Add(card)

            card.OnTapped = func() {
                check.SetChecked(!check.Checked)
            }
        }

        printerListContainer.Refresh()
        installBtn.Enable()
    }

    locationSelect.OnChanged = func(s string) {
        updatePrinterList(s)
    }

    installBtn.OnTapped = func() {
        selected := []Printer{}
        for i, c := range checks {
            if c != nil && c.Checked {
                selected = append(selected, printers[i])
            }
        }

        if len(selected) == 0 {
            dialog.ShowInformation("提示", "请至少选择一台打印机", w)
            return
        }

        dialog.ShowConfirm("确认安装",
            fmt.Sprintf("确定要安装 %d 台打印机？", len(selected)),
            func(ok bool) {
                if ok {
                    go runInstallFlow(selected, cfg, w, status)
                }
            }, w)
    }

    // 顶栏
    topBar := container.NewVBox(
        title,
        canvas.NewRectangle(theme.ShadowColor()),
        container.NewHBox(
            widget.NewLabel("地点："),
            locationSelect,
            refreshBtn,
            layout.NewSpacer(),
            status,
        ),
        canvas.NewRectangle(theme.ShadowColor()),
    )

    bottomBar := container.NewHBox(
        layout.NewSpacer(),
        installBtn,
    )

    mainLayout := container.NewBorder(topBar, bottomBar, nil, nil, scroll)

    w.SetContent(mainLayout)

    go loadConfig()
    w.ShowAndRun()
}

func runInstallFlow(selected []Printer, cfg Config, w fyne.Window, status *widget.Label) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    success := 0
    failed := []string{}

    status.SetText("正在安装...")

    for _, p := range selected {
        wg.Add(1)
        go func(pr Printer) {
            defer wg.Done()

            ok, err := installSingle(pr, cfg)
            mu.Lock()
            defer mu.Unlock()
            if ok {
                success++
            } else {
                failed = append(failed, fmt.Sprintf("%s: %v", pr.Name, err))
            }
            status.SetText(fmt.Sprintf("安装进度: %d 成功, %d 失败", success, len(failed)))
        }(p)
    }

    wg.Wait()

    if len(failed) == 0 {
        dialog.ShowInformation("完成", fmt.Sprintf("全部成功安装 (%d 台)", success), w)
        status.SetText("安装完成")
    } else {
        dialog.ShowError(
            fmt.Errorf("成功 %d, 失败 %d\n%s", success, len(failed), strings.Join(failed, "\n")),
            w,
        )
        status.SetText("安装完成（部分失败）")
    }
}

func installSingle(p Printer, cfg Config) (bool, error) {
    modelCfg, ok := cfg.PrinterModels[p.Model]
    if !ok {
        return false, fmt.Errorf("未配置型号: %s", p.Model)
    }

    ppdURL := modelCfg.PPDUrl
    if strings.Contains(ppdURL, " ") {
        ppdURL = strings.ReplaceAll(ppdURL, " ", "%20")
    }

    tmpfile := filepath.Join(os.TempDir(), filepath.Base(ppdURL))
    out, err := os.Create(tmpfile)
    if err != nil {
        return false, err
    }
    defer func() { out.Close(); os.Remove(tmpfile) }()

    resp, err := http.Get(ppdURL)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return false, err
    }

    uri := p.URI
    if uri == "" {
        uri = "ipp://" + p.IP + "/ipp/print"
    }

    exec.Command("lpadmin", "-x", p.Name).Run()

    cmd := exec.Command("lpadmin",
        "-p", p.Name,
        "-v", uri,
        "-P", tmpfile,
        "-E",
        "-D", fmt.Sprintf("%s (%s)", p.Name, p.Model),
    )

    b, err := cmd.CombinedOutput()
    if err != nil {
        return false, fmt.Errorf("%v: %s", err, string(b))
    }

    return true, nil
}
