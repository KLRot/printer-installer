# 中文字体解决方案

## 最终方案：使用系统字体

经过测试，我们采用了最简单可靠的方案：**程序运行时自动加载系统字体**。

## 为什么不嵌入字体？

虽然嵌入字体理论上可以解决所有问题，但在 CI 环境中 `fyne bundle` 工具存在兼容性问题，会导致编译失败。因此我们选择了更稳定的系统字体方案。

## 使用说明

### 1. 安装中文字体

在麒麟系统上运行程序前，请先安装中文字体：

```bash
sudo apt-get install fonts-wqy-microhei
```

### 2. 运行程序

```bash
./printer-installer-amd64
```

程序启动时会自动搜索并加载系统字体，你会看到：

```
✓ 成功加载系统字体: /usr/share/fonts/truetype/wqy/wqy-microhei.ttc
```

### 3. 如果仍然乱码

如果安装字体后仍然乱码，可以通过环境变量手动指定字体：

```bash
FYNE_FONT=/usr/share/fonts/truetype/wqy/wqy-microhei.ttc ./printer-installer-amd64
```

## 支持的字体路径

程序会按顺序搜索以下路径：

1. 环境变量 `FYNE_FONT` 指定的路径
2. 麒麟/UKUI 系统字体
   - `/usr/share/fonts/truetype/ukui/ukui-default.ttf`
   - `/usr/share/fonts/ukui/ukui-default.ttf`
   - `/usr/share/fonts/truetype/kylin-font/kylin-font.ttf`
3. 通用 Linux 字体
   - `/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc`
   - `/usr/share/fonts/truetype/wqy/wqy-microhei.ttc`
   - `/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc`

## 优势

- ✅ 编译简单，无需额外步骤
- ✅ 可执行文件体积小
- ✅ 用户可以选择自己喜欢的字体
- ✅ 稳定可靠，不依赖 CI 工具

## 注意事项

- 程序**不会**自动安装字体，需要用户手动安装
- 如果系统没有中文字体，程序会显示警告信息
- 图标已嵌入到可执行文件中，无需额外文件
