#!/bin/bash

# 默认配置
BARK_URL="http://localhost:8080"
BARK_TOKEN=""

# 显示帮助信息
show_help() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -u, --url URL       Bark 服务地址 (默认: http://localhost:8080)"
    echo "  -t, --token TOKEN   Bark 服务访问令牌"
    echo "  -h, --help         显示帮助信息"
    exit 0
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -u|--url)
            BARK_URL="$2"
            shift 2
            ;;
        -t|--token)
            BARK_TOKEN="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            ;;
        *)
            echo "未知参数: $1"
            show_help
            ;;
    esac
done

# 获取暂存的文件列表
FILES=$(git diff --cached --name-only --diff-filter=ACMR | tr '\n' ' ')

if [ -z "$FILES" ]; then
    echo "No files to check"
    exit 0
fi

# 构建请求体
JSON_DATA=$(echo "$FILES" | jq -R -s -c 'split(" ") | map(select(length > 0)) | {files: .}')

# 构建 curl 命令
CURL_CMD="curl -s -X POST -H \"Content-Type: application/json\""

# 如果提供了 token，添加到请求头
if [ -n "$BARK_TOKEN" ]; then
    CURL_CMD="$CURL_CMD -H \"Authorization: Bearer $BARK_TOKEN\""
fi

# 发送请求到 bark 服务
RESPONSE=$(eval "$CURL_CMD -d \"$JSON_DATA\" \"$BARK_URL/precommit\"")

# 检查 curl 命令是否成功
if [ $? -ne 0 ]; then
    echo "Error: Failed to connect to Bark service at $BARK_URL"
    exit 1
fi

# 解析响应
SUCCESS=$(echo "$RESPONSE" | jq -r '.success')
MESSAGE=$(echo "$RESPONSE" | jq -r '.message')

if [ "$SUCCESS" = "true" ]; then
    echo "Pre-commit check passed"
    exit 0
else
    echo "Pre-commit check failed:"
    echo "$MESSAGE"
    exit 1
fi 