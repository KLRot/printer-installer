#!/bin/bash

echo "=== 中文字体诊断工具 ==="
echo ""

echo "1. 检查系统中所有中文字体："
echo "-----------------------------------"
fc-list :lang=zh family | sort -u

echo ""
echo "2. 检查程序搜索的字体路径："
echo "-----------------------------------"

# 程序中定义的字体路径
FONT_PATHS=(
    "/usr/share/fonts/truetype/ukui/ukui-default.ttf"
    "/usr/share/fonts/ukui/ukui-default.ttf"
    "/usr/share/fonts/truetype/kylin-font/kylin-font.ttf"
    "/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc"
    "/usr/share/fonts/noto-cjk/NotoSansCJK-Regular.ttc"
    "/usr/share/fonts/truetype/wqy/wqy-microhei.ttc"
    "/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc"
    "/usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf"
    "/usr/share/fonts/truetype/arphic/uming.ttc"
    "/usr/share/fonts/truetype/arphic/ukai.ttc"
)

FOUND=0
for path in "${FONT_PATHS[@]}"; do
    if [ -f "$path" ]; then
        echo "✓ 找到: $path"
        FOUND=1
    else
        echo "✗ 未找到: $path"
    fi
done

echo ""
if [ $FOUND -eq 0 ]; then
    echo "❌ 警告：没有找到任何程序支持的中文字体！"
    echo ""
    echo "请安装中文字体："
    echo "  sudo apt-get install fonts-wqy-microhei"
    echo "  或"
    echo "  sudo apt-get install fonts-noto-cjk"
else
    echo "✓ 至少找到一个可用的中文字体"
fi

echo ""
echo "3. 查找系统中第一个可用的中文字体文件："
echo "-----------------------------------"
FIRST_FONT=$(fc-list :lang=zh file | head -n 1 | cut -d: -f1 | xargs)
if [ -n "$FIRST_FONT" ]; then
    echo "找到: $FIRST_FONT"
    echo ""
    echo "你可以使用环境变量强制使用这个字体："
    echo "  FYNE_FONT=\"$FIRST_FONT\" ./printer-installer-amd64"
else
    echo "未找到任何中文字体文件"
fi

echo ""
echo "=== 诊断完成 ==="
