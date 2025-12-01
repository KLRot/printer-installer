#!/bin/bash

# DEB 打包脚本
# 用于将编译好的可执行文件打包成 .deb 安装包

set -e

VERSION="1.0.0"
PACKAGE_NAME="printer-installer"
MAINTAINER="Your Name <your.email@example.com>"
DESCRIPTION="麒麟系统打印机自动安装程序"

echo "========================================="
echo "开始打包 DEB 安装包"
echo "========================================="

# 检查可执行文件是否存在
if [ ! -f "dist/printer-installer-amd64" ]; then
    echo "错误: 找不到 dist/printer-installer-amd64"
    echo "请先运行编译命令生成可执行文件"
    exit 1
fi

if [ ! -f "dist/printer-installer-arm64" ]; then
    echo "错误: 找不到 dist/printer-installer-arm64"
    echo "请先运行编译命令生成可执行文件"
    exit 1
fi

# 函数：创建 DEB 包
create_deb() {
    ARCH=$1
    BINARY_NAME=$2
    
    echo ""
    echo "正在打包 ${ARCH} 版本..."
    
    # 创建临时目录
    BUILD_DIR="build/deb-${ARCH}"
    rm -rf "${BUILD_DIR}"
    mkdir -p "${BUILD_DIR}"
    
    # 创建目录结构
    mkdir -p "${BUILD_DIR}/DEBIAN"
    mkdir -p "${BUILD_DIR}/usr/bin"
    mkdir -p "${BUILD_DIR}/usr/share/applications"
    mkdir -p "${BUILD_DIR}/usr/share/pixmaps"
    
    # 复制可执行文件
    cp "dist/${BINARY_NAME}" "${BUILD_DIR}/usr/bin/printer-installer"
    chmod 755 "${BUILD_DIR}/usr/bin/printer-installer"
    
    # 复制图标（如果存在）
    if [ -f "printer_icon.png" ]; then
        cp printer_icon.png "${BUILD_DIR}/usr/share/pixmaps/printer-installer.png"
    fi
    
    # 创建桌面快捷方式
    cat > "${BUILD_DIR}/usr/share/applications/printer-installer.desktop" << EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=打印机安装程序
Name[zh_CN]=打印机安装程序
Comment=麒麟系统打印机自动安装工具
Comment[zh_CN]=麒麟系统打印机自动安装工具
Exec=/usr/bin/printer-installer
Icon=printer-installer
Terminal=false
Categories=System;Settings;
Keywords=printer;install;打印机;安装;
EOF
    
    # 创建 control 文件
    cat > "${BUILD_DIR}/DEBIAN/control" << EOF
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: ${MAINTAINER}
Description: ${DESCRIPTION}
 麒麟系统打印机自动安装程序，支持批量安装和管理打印机。
 .
 主要功能：
  - 从配置服务器自动获取打印机列表
  - 支持批量选择和安装打印机
  - 自动下载和配置 PPD 文件
  - 友好的图形界面
Depends: cups
EOF
    
    # 创建 postinst 脚本（安装后执行）
    cat > "${BUILD_DIR}/DEBIAN/postinst" << 'EOF'
#!/bin/bash
set -e

# 更新桌面数据库
if [ -x /usr/bin/update-desktop-database ]; then
    update-desktop-database -q
fi

# 更新图标缓存
if [ -x /usr/bin/gtk-update-icon-cache ]; then
    gtk-update-icon-cache -q -t -f /usr/share/pixmaps 2>/dev/null || true
fi

echo "打印机安装程序已成功安装"
echo "您可以从应用菜单启动，或在终端运行: printer-installer"

exit 0
EOF
    
    chmod 755 "${BUILD_DIR}/DEBIAN/postinst"
    
    # 创建 prerm 脚本（卸载前执行）
    cat > "${BUILD_DIR}/DEBIAN/prerm" << 'EOF'
#!/bin/bash
set -e
exit 0
EOF
    
    chmod 755 "${BUILD_DIR}/DEBIAN/prerm"
    
    # 构建 DEB 包
    DEB_FILE="${PACKAGE_NAME}_${VERSION}_${ARCH}.deb"
    dpkg-deb --build "${BUILD_DIR}" "dist/${DEB_FILE}"
    
    echo "✓ 成功创建: dist/${DEB_FILE}"
    
    # 显示包信息
    echo ""
    echo "包信息:"
    dpkg-deb --info "dist/${DEB_FILE}"
    
    # 清理临时目录
    rm -rf "${BUILD_DIR}"
}

# 创建 build 目录
mkdir -p build
mkdir -p dist

# 打包 amd64 版本
create_deb "amd64" "printer-installer-amd64"

# 打包 arm64 版本
create_deb "arm64" "printer-installer-arm64"

echo ""
echo "========================================="
echo "打包完成！"
echo "========================================="
echo ""
echo "生成的 DEB 包:"
ls -lh dist/*.deb
echo ""
echo "安装方法:"
echo "  sudo dpkg -i dist/${PACKAGE_NAME}_${VERSION}_amd64.deb"
echo ""
echo "卸载方法:"
echo "  sudo dpkg -r ${PACKAGE_NAME}"
echo ""
