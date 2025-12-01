# 麒麟系统打印机自动安装程序

## 文件说明

- `printer-installer-amd64` - x86_64 架构可执行文件（图标+字体已嵌入）
- `printer-installer-arm64` - ARM64 架构可执行文件（图标+字体已嵌入）

## 使用方法

### 1. 选择对应架构的文件

```bash
# x86_64 系统
chmod +x printer-installer-amd64
./printer-installer-amd64

# ARM64 系统
chmod +x printer-installer-arm64
./printer-installer-arm64
```

**无需安装任何字体！** 程序已内置 Noto Sans SC 中文字体。

### 2. 图标和字体说明

- 图标已嵌入到可执行文件中
- 中文字体（Noto Sans SC）已嵌入到可执行文件中
- **完全独立运行，无需任何外部依赖**

### 3. 安装到系统（可选）

```bash
# 复制到系统目录
sudo cp printer-installer-amd64 /usr/local/bin/printer-installer

# 运行
printer-installer
```

## 系统要求

- 麒麟 Linux 系统
- CUPS 打印服务
- 网络连接（用于下载配置和 PPD 文件）

## 功能特性

- ✅ 从服务器自动加载打印机配置
- ✅ 按地点分类显示打印机
- ✅ 批量选择和安装打印机
- ✅ 自动下载 PPD 驱动文件
- ✅ 实时显示安装进度
- ✅ 亮色主题，完整中文支持
- ✅ 图标已嵌入，无需额外文件

## 注意事项

1. 需要 lpadmin 权限来安装打印机
2. 确保 CUPS 服务正在运行
3. 图标已嵌入到可执行文件中，无需单独的图标文件

---
构建时间: $(date)
