# WSL 环境下生成 go.sum 文件

## 问题说明

在 GitHub Actions 中打包 Go 项目时，需要 `go.sum` 文件来验证依赖的完整性。如果缺少这个文件，会报错：

```
Error: missing go.sum entry for module providing package fyne.io/fyne/v2
```

## 解决方案

在 WSL 环境下运行以下命令来生成 `go.sum` 文件：

### 步骤 1: 进入项目目录

```bash
# 在 WSL 中，Windows 的 F 盘通常挂载在 /mnt/f
cd /mnt/f/麒麟打印机安装/printer-installer-go
```

### 步骤 2: 下载依赖并生成 go.sum

```bash
# 下载所有依赖
go mod download

# 整理依赖并生成 go.sum
go mod tidy
```

### 步骤 3: 验证文件生成

```bash
# 检查 go.sum 文件是否生成
ls -lh go.sum

# 查看 go.sum 内容（可选）
head -20 go.sum
```

### 步骤 4: 提交到 Git

```bash
# 添加 go.sum 文件到 Git
git add go.sum

# 提交
git commit -m "Add go.sum file"

# 推送到远程仓库
git push
```

## 完整命令（一键执行）

```bash
cd /mnt/f/麒麟打印机安装/printer-installer-go && \
go mod download && \
go mod tidy && \
ls -lh go.sum && \
echo "✓ go.sum 文件已生成"
```

## 验证编译

生成 `go.sum` 后，可以测试编译是否正常：

```bash
# 测试编译
go build -o printer_installer main.go

# 如果编译成功，清理生成的文件
rm printer_installer

echo "✓ 编译测试通过"
```

## 注意事项

1. **Go 版本**: 确保 WSL 中安装了 Go 1.21 或更高版本
   ```bash
   go version
   ```

2. **如果没有安装 Go**:
   ```bash
   # 下载 Go
   wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
   
   # 解压
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   
   # 添加到 PATH
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   
   # 验证安装
   go version
   ```

3. **代理设置**（如果网络受限）:
   ```bash
   # 设置 Go 代理
   export GOPROXY=https://goproxy.cn,direct
   
   # 或者添加到 ~/.bashrc
   echo 'export GOPROXY=https://goproxy.cn,direct' >> ~/.bashrc
   source ~/.bashrc
   ```

## 故障排除

### 问题 1: 找不到 go 命令

**解决方案**: 安装 Go（见上面的安装步骤）

### 问题 2: 网络超时

**解决方案**: 设置 Go 代理
```bash
export GOPROXY=https://goproxy.cn,direct
go mod download
```

### 问题 3: 权限错误

**解决方案**: 检查文件权限
```bash
# 修改文件所有者（如果需要）
sudo chown -R $USER:$USER /mnt/f/麒麟打印机安装/printer-installer-go

# 或者使用 sudo 运行
sudo go mod download
sudo go mod tidy
```

### 问题 4: go.sum 内容不正确

**解决方案**: 清理并重新生成
```bash
# 删除旧的 go.sum
rm go.sum

# 清理缓存
go clean -modcache

# 重新生成
go mod download
go mod tidy
```

## GitHub Actions 配置

确保你的 GitHub Actions 工作流包含以下步骤：

```yaml
- name: Setup Go
  uses: actions/setup-go@v4
  with:
    go-version: '1.21'

- name: Download dependencies
  run: |
    cd printer-installer-go
    go mod download
    go mod tidy

- name: Build
  run: |
    cd printer-installer-go
    go build -o printer_installer main.go
```

## 成功标志

当你看到以下输出时，说明成功了：

```
✓ go.sum 文件已生成
-rw-r--r-- 1 user user 12345 Nov 30 14:00 go.sum
```

## 后续步骤

1. ✅ 生成 `go.sum` 文件
2. ✅ 提交到 Git 仓库
3. ✅ 推送到 GitHub
4. ✅ GitHub Actions 会自动使用这个文件进行构建

---

**最后更新**: 2025-11-30
**适用版本**: Go 1.21+
