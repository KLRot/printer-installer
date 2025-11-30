# 字体文件说明

## 文件信息

- **文件名**: `NotoSansSC-Regular.otf`
- **字体**: Noto Sans SC (简体中文)
- **大小**: 约 2.5 MB
- **许可**: SIL Open Font License 1.1

## 用途

该字体文件通过 Go 的 `embed` 指令嵌入到可执行文件中，确保程序在任何系统上都能正确显示中文。

## 嵌入方式

在 `main.go` 中：

```go
//go:embed NotoSansSC-Regular.otf
var embeddedFontData []byte
```

编译时，Go 编译器会自动将字体文件内容读取到 `embeddedFontData` 变量中。

## Git 管理

- ✅ 字体文件已提交到仓库
- ✅ 使用 `.gitattributes` 标记为二进制文件
- ✅ GitHub Actions 无需下载，直接使用

## 更新字体

如果需要更新字体文件：

1. 下载新的字体文件
2. 替换 `NotoSansSC-Regular.otf`
3. 提交到 Git
4. 重新编译

## 许可证

Noto Sans SC 使用 SIL Open Font License 1.1，允许：
- ✅ 商业使用
- ✅ 修改
- ✅ 分发
- ✅ 嵌入到软件中

详见：https://scripts.sil.org/OFL
