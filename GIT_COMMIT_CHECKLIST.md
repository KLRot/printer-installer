# Git 提交清单

## 需要提交的文件

### 1. 核心代码
- ✅ `main.go` - 添加了字体嵌入逻辑

### 2. 字体文件（重要！）
- ✅ `NotoSansSC-Regular.otf` - 中文字体文件（2.5MB）

### 3. 配置文件
- ✅ `.gitattributes` - 确保字体文件以二进制方式处理
- ✅ `.github/workflows/build.yml` - 移除了字体下载步骤

### 4. 文档
- ✅ `EMBEDDED_FONT_BUILD.md` - 嵌入字体构建指南
- ✅ `FONT_FILE.md` - 字体文件说明
- ✅ `download_font.sh` - 字体下载脚本（可选）

## 提交命令

```bash
cd /mnt/f/麒麟打印机安装/printer-installer-go

# 查看状态
git status

# 添加所有修改的文件
git add main.go
git add NotoSansSC-Regular.otf
git add .gitattributes
git add .github/workflows/build.yml
git add EMBEDDED_FONT_BUILD.md
git add FONT_FILE.md

# 提交
git commit -m "feat: 嵌入中文字体，彻底解决跨平台字体问题

- 使用 Go embed 嵌入 Noto Sans SC 字体
- 移除对系统字体的依赖
- 更新 GitHub Actions 配置
- 添加字体文件说明文档"

# 推送
git push
```

## 验证

提交后，GitHub Actions 会自动构建。检查：

1. ✅ 构建成功
2. ✅ 可执行文件大小约 22-23MB（比之前增加 2-3MB）
3. ✅ 下载并运行，中文显示正常

## 注意事项

### 字体文件大小

`NotoSansSC-Regular.otf` 约 2.5MB，这会增加：
- Git 仓库大小：+2.5MB
- 可执行文件大小：+2.5MB

这是可接受的代价，换来的是：
- ✅ 零依赖
- ✅ 跨平台一致
- ✅ 用户体验完美

### Git LFS（可选）

如果担心仓库体积，可以使用 Git LFS：

```bash
# 安装 Git LFS
git lfs install

# 跟踪字体文件
git lfs track "*.otf"

# 提交 .gitattributes
git add .gitattributes
git commit -m "chore: 使用 Git LFS 管理字体文件"
```

但对于 2.5MB 的文件，通常不需要 LFS。

## 下一步

提交后：
1. 等待 GitHub Actions 构建完成
2. 下载构建产物
3. 在麒麟系统上测试
4. 验证中文显示正常

完成！🎉
