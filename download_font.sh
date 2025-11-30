#!/bin/bash

# 下载精简版中文字体用于嵌入
# 使用 Noto Sans SC 的子集版本（约 2-3MB）

echo "正在下载精简版中文字体..."

# 方案1：从 Google Fonts 下载 Noto Sans SC（推荐）
wget -O NotoSansSC-Regular.otf "https://github.com/googlefonts/noto-cjk/raw/main/Sans/SubsetOTF/SC/NotoSansSC-Regular.otf"

if [ $? -eq 0 ]; then
    echo "✓ 字体下载成功"
    ls -lh NotoSansSC-Regular.otf
else
    echo "✗ 下载失败，尝试备用源..."
    
    # 备用方案：使用文泉驿微米黑（约 4MB）
    wget -O wqy-microhei.ttc "https://github.com/anthonyfok/fonts-wqy-microhei/raw/master/wqy-microhei.ttc"
    
    if [ $? -eq 0 ]; then
        echo "✓ 备用字体下载成功"
        ls -lh wqy-microhei.ttc
    else
        echo "✗ 所有下载源都失败"
        exit 1
    fi
fi

echo ""
echo "下一步：运行 go generate 生成嵌入代码"
