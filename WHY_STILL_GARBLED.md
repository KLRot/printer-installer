# 为什么系统有字体还是乱码？

## 问题原因

**关键点**：不是所有字体都支持中文！

### 常见误区

很多 Linux 系统默认安装的字体（如 DejaVu、Liberation、Noto Sans）**只支持拉丁字符**，不包含中文字符。如果程序加载了这些字体，中文会显示为：
- 方框 □□□
- 问号 ???
- 或者其他乱码

### 支持中文的字体 vs 不支持中文的字体

| 字体名称 | 支持中文 | 说明 |
|---------|---------|------|
| DejaVu Sans | ❌ | 只有拉丁字符 |
| Liberation Sans | ❌ | 只有拉丁字符 |
| Noto Sans | ❌ | 只有拉丁字符 |
| **Noto Sans CJK** | ✅ | 支持中日韩文字 |
| **WQY Micro Hei** | ✅ | 文泉驿微米黑 |
| **WQY Zen Hei** | ✅ | 文泉驿正黑 |
| **AR PL UMing** | ✅ | 文鼎 PL 明体 |

## 诊断步骤

### 1. 运行诊断脚本

```bash
chmod +x diagnose_fonts.sh
./diagnose_fonts.sh
```

这个脚本会：
- 列出系统中所有中文字体
- 检查程序搜索的字体路径是否存在
- 找到第一个可用的中文字体文件

### 2. 查看程序输出

运行程序时，注意控制台输出：

```bash
./printer-installer-amd64
```

**正确的输出**：
```
✓ 成功加载系统字体: /usr/share/fonts/truetype/wqy/wqy-microhei.ttc
```

**错误的输出**：
```
! 警告: 未找到中文字体
```

### 3. 手动检查字体

```bash
# 查找所有中文字体
fc-list :lang=zh

# 查找特定字体
ls -l /usr/share/fonts/truetype/wqy/
```

## 解决方案

### 方案 1：安装推荐的中文字体

```bash
# 文泉驿微米黑（推荐，体积小）
sudo apt-get install fonts-wqy-microhei

# 或者 Noto CJK（Google 字体）
sudo apt-get install fonts-noto-cjk

# 或者文鼎字体
sudo apt-get install fonts-arphic-uming
```

### 方案 2：使用环境变量指定字体

如果你已经有中文字体但程序找不到，可以手动指定：

```bash
# 1. 找到字体文件
fc-list :lang=zh file | head -n 1

# 2. 使用环境变量运行
FYNE_FONT="/path/to/your/chinese/font.ttc" ./printer-installer-amd64
```

### 方案 3：检查文件编码

如果安装了字体还是乱码，可能是源代码文件编码问题：

1. 在 VS Code 中打开 `main.go`
2. 查看右下角编码（应该是 UTF-8）
3. 如果不是，点击编码 → "Save with Encoding" → 选择 UTF-8
4. 重新编译

## 常见问题

### Q: 我安装了字体，但程序还是提示找不到？

A: 可能是字体安装路径不在程序搜索列表中。运行诊断脚本找到实际路径，然后：
- 使用 `FYNE_FONT` 环境变量
- 或者修改 `main.go` 中的 `fontPaths` 数组

### Q: 程序加载了字体，但中文还是方框？

A: 检查加载的是哪个字体。如果是 DejaVu、Liberation 等，说明程序加载了不支持中文的字体。确保安装了 `fonts-wqy-microhei` 或 `fonts-noto-cjk`。

### Q: 在 Windows 上测试正常，Linux 上乱码？

A: Windows 默认有中文字体（如微软雅黑），但 Linux 需要手动安装。

## 验证字体是否支持中文

```bash
# 查看字体包含的字符范围
fc-query /usr/share/fonts/truetype/wqy/wqy-microhei.ttc | grep lang

# 应该包含：
# lang: zh-cn|zh-sg|zh-tw|zh-hk|...
```

## 总结

1. ✅ 确保安装了**支持中文**的字体
2. ✅ 运行诊断脚本检查字体路径
3. ✅ 查看程序启动时的字体加载信息
4. ✅ 如果需要，使用 `FYNE_FONT` 环境变量

**记住**：不是所有字体都支持中文！必须安装 CJK（中日韩）字体。
