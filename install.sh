#!/bin/bash

# 定义文件下载地址和目标文件名
DOWNLOAD_URL="https://github.com/xifo-wu/post-bee/releases/download/V0.0.1.BETA/postbee.zip"
TARGET_FILE="postbee.zip"

# 定义解压后的目标文件夹路径
TARGET_FOLDER="~/.postbee"

# 下载文件
echo "开始下载postbee.zip文件..."
curl -L -o $TARGET_FILE $DOWNLOAD_URL

# 检查文件是否下载成功
if [ $? -ne 0 ]; then
  echo "下载文件失败！"
  exit 1
fi

echo "下载完成！"

# 解压文件
echo "开始解压postbee.zip文件..."
unzip $TARGET_FILE -d $TARGET_FOLDER

# 检查解压是否成功
if [ $? -ne 0 ]; then
  echo "解压文件失败！"
  exit 1
fi

echo "解压完成！"

# 删除下载的zip文件
echo "删除下载的zip文件..."
rm $TARGET_FILE

echo "安装完成！"