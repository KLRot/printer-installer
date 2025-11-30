# 解决中文乱码的终极指南

## 1. 确认文件编码（非常重要！）

如果源代码文件本身不是 UTF-8 编码，那么无论怎么设置字体，中文都会乱码。

**检查步骤：**
1. 在 VS Code 中打开 `main.go`
2. 查看窗口右下角的状态栏
3. 如果显示 **UTF-8**，则正常
4. 如果显示 **GBK** 或 **GB2312**：
   - 点击该编码
   - 选择 **"Save with Encoding" (通过编码保存)**
   - 选择 **UTF-8**

## 2. 检查运行时输出

重新编译并运行程序，观察控制台输出：

```bash
./printer-installer-amd64
```

**情况 A: 找到字体**
```
✓ 成功加载系统字体: /usr/share/fonts/...
```
如果显示这个但仍然乱码，说明是**文件编码问题**（见第1点）。

**情况 B: 未找到字体**
```
! 警告: 未找到任何中文字体！
! 已搜索的路径:
  - /usr/share/fonts/...
```
这说明程序没找到任何可用的字体。

## 3. 解决方案

### 方法一：设置环境变量（推荐）

找到系统中任意一个支持中文的字体文件（例如 `simsun.ttc` 或 `DroidSansFallback.ttf`），然后运行：

```bash
# 临时运行
FYNE_FONT=/usr/share/fonts/truetype/arphic/uming.ttc ./printer-installer-amd64

# 或者永久设置
export FYNE_FONT=/usr/share/fonts/truetype/arphic/uming.ttc
./printer-installer-amd64
```

### 方法二：安装常用字体

麒麟系统通常可以通过以下命令安装常用中文字体：

```bash
sudo apt-get update
sudo apt-get install fonts-noto-cjk fonts-wqy-microhei fonts-wqy-zenhei
```

安装后重新运行程序即可。

### 方法三：查找本机字体

如果不知道字体在哪里，可以使用命令查找：

```bash
# 查找所有中文字体
fc-list :lang=zh

# 输出示例：
# /usr/share/fonts/truetype/ukui/ukui-default.ttf: UKUI-Default:style=Regular
```

将找到的路径告诉程序（通过 `FYNE_FONT` 环境变量）。

## 4. 重新编译

```bash
cd /mnt/f/麒麟打印机安装/printer-installer-go
go build -o printer_installer main.go
```
