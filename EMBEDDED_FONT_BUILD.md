# 嵌入字体构建指南

## 方案说明

使用 Go 的 `embed` 指令将中文字体直接编译进可执行文件，彻底解决跨平台字体问题。

## 优势

- ✅ **零依赖**：无需用户安装任何字体
- ✅ **跨平台一致**：所有系统显示效果完全相同
- ✅ **简化部署**：单个文件即可运行
- ✅ **体积可控**：使用精简版字体，增加约 2-3MB

## 本地构建步骤

### 1. 下载字体

```bash
cd /mnt/f/麒麟打印机安装/printer-installer-go

# 运行下载脚本
chmod +x download_font.sh
./download_font.sh
```

这会下载 `NotoSansSC-Regular.otf` 文件（约 2.5MB）。

### 2. 编译程序

```bash
go build -o printer_installer main.go
```

Go 编译器会自动读取 `//go:embed` 指令，将字体文件编译进可执行文件。

### 3. 验证

```bash
./printer_installer
```

启动时应该看到：
```
✓ 使用嵌入的中文字体 (Noto Sans SC)
```

## GitHub Actions 自动构建

GitHub Actions 会自动：
1. 下载 Noto Sans SC 字体
2. 使用 `go:embed` 编译进程序
3. 生成包含字体的可执行文件

## 文件大小对比

| 版本 | 可执行文件大小 | 说明 |
|------|---------------|------|
| 无嵌入字体 | ~20 MB | 依赖系统字体 |
| 嵌入字体 | ~22-23 MB | 完全独立 |

增加约 2-3MB，换来完全的独立性和一致性。

## 字体选择

使用 **Noto Sans SC**（简体中文）的原因：
- Google 开源字体，授权友好（OFL 1.1）
- 精简版体积小（2.5MB）
- 覆盖常用汉字（约 6000+ 字）
- 显示效果优秀

## 故障排除

### 问题：编译时提示找不到 NotoSansSC-Regular.otf

**解决**：
```bash
# 确保字体文件在项目根目录
ls -l NotoSansSC-Regular.otf

# 如果不存在，运行下载脚本
./download_font.sh
```

### 问题：编译后文件很大

**说明**：这是正常的。嵌入字体后文件会增加 2-3MB。如果文件超过 30MB，可能是 debug 信息未去除。

**解决**：
```bash
# 使用 -ldflags 去除调试信息
go build -ldflags="-s -w" -o printer_installer main.go
```

### 问题：运行时仍然提示找不到字体

**检查**：
```bash
# 查看程序输出
./printer_installer

# 应该看到：
# ✓ 使用嵌入的中文字体 (Noto Sans SC)
```

如果看到其他信息，说明 embed 没有生效。确保：
1. `NotoSansSC-Regular.otf` 在项目根目录
2. `main.go` 中有 `//go:embed NotoSansSC-Regular.otf`
3. 使用 Go 1.16+ 版本

## 许可证

Noto Sans SC 字体使用 SIL Open Font License 1.1，允许商业使用和嵌入。

## 总结

使用 Go embed 嵌入字体是最可靠的解决方案：
- 无需用户安装字体
- 无需担心系统差异
- 一次编译，到处运行
- 体积增加可接受

这是生产环境的最佳选择！🎉
