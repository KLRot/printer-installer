# 字体、窗口和 Core Dump 问题解决方案

## 1. ✅ 字体乱码问题

**原因**: Fyne 默认字体不支持中文。

**解决方案**:
- 代码已修改为**自动搜索系统字体**。
- 程序启动时会按顺序检查以下路径：
  1. `/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc`
  2. `/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc`
  3. `/usr/share/fonts/truetype/wqy/wqy-microhei.ttc`
  4. `/usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf`
  5. `/usr/share/fonts/noto-cjk/NotoSansCJK-Regular.ttc`
  6. Windows 字体路径 (兼容性)

**验证**:
启动程序时，控制台会输出：
```
✓ 已加载系统字体: /usr/share/fonts/...
```

## 2. ✅ 窗口缩放问题

**原因**: 窗口初始化顺序不当，导致布局计算错误。

**解决方案**:
- 优化了初始化顺序：
  1. `SetMaster()` - 设为主窗口
  2. `initUI()` - 先设置内容
  3. `Resize()` - 再调整大小
  4. `CenterOnScreen()` - 最后居中
  5. `ShowAndRun()` - 显示

## 3. ✅ .core 文件问题

**原因**: 程序退出时底层图形驱动可能发生异常，导致生成 core dump 文件。

**解决方案 A: 代码层面 (已实施)**
- 添加了 `recover()` 捕获 panic
- 添加了生命周期钩子 `SetOnStopped` 进行优雅退出

**解决方案 B: 系统层面 (推荐)**
如果仍然生成 `.core` 文件，建议在启动脚本中禁止生成。

修改启动脚本 (例如 `run.sh`)：

```bash
#!/bin/bash

# 禁止生成 core dump 文件
ulimit -c 0

# 运行程序
./printer-installer-amd64
```

或者在终端中直接运行：
```bash
ulimit -c 0
./printer-installer-amd64
```

## 重新编译

请在 WSL 中重新编译程序以应用这些更改：

```bash
cd /mnt/f/麒麟打印机安装/printer-installer-go
go build -o printer_installer main.go
```

## 故障排除

如果中文仍然乱码，请检查系统是否安装了中文字体：

```bash
# 查找系统中文字体
fc-list :lang=zh

# 如果没有，安装字体
sudo apt-get install fonts-noto-cjk fonts-wqy-zenhei
```

如果找到了字体但路径不在我们的列表中，请告诉我字体的实际路径，我会添加到代码中。
