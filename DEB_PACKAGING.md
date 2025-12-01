# 打印机安装程序 - DEB 打包说明

## 快速开始

### 1. 编译程序

```bash
# 使用 fyne-cross 编译（推荐）
fyne-cross linux -arch=amd64,arm64 -name printer-installer -icon printer_icon.png

# 或者使用 GitHub Actions 自动编译
git push
```

### 2. 打包 DEB

```bash
chmod +x build-deb.sh
./build-deb.sh
```

这会在 `dist/` 目录生成：
- `printer-installer_1.0.0_amd64.deb`
- `printer-installer_1.0.0_arm64.deb`

### 3. 安装

```bash
# AMD64 架构
sudo dpkg -i dist/printer-installer_1.0.0_amd64.deb

# ARM64 架构
sudo dpkg -i dist/printer-installer_1.0.0_arm64.deb
```

### 4. 运行

安装后可以：
- 从应用菜单启动（系统设置 → 打印机安装程序）
- 或在终端运行：`printer-installer`

### 5. 卸载

```bash
sudo dpkg -r printer-installer
```

## DEB 包内容

```
/usr/bin/printer-installer              # 主程序
/usr/share/applications/printer-installer.desktop  # 桌面快捷方式
/usr/share/pixmaps/printer-installer.png          # 图标
```

## 修改版本号

编辑 `build-deb.sh` 文件，修改：

```bash
VERSION="1.0.0"
MAINTAINER="Your Name <your.email@example.com>"
```

## 依赖

- **运行时依赖**: `cups` (CUPS 打印系统)
- **打包工具**: `dpkg-deb`

## 故障排除

### 问题：dpkg-deb: command not found

```bash
sudo apt-get install dpkg-dev
```

### 问题：中文显示乱码

确保系统安装了中文字体：

```bash
sudo apt-get install fonts-noto-cjk
```

程序会自动查找并使用系统中的楷体字体。
