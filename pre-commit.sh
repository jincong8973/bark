#!/bin/bash

# 默认配置
BARK_URL="http://localhost:8080"
BARK_TOKEN=""

# 显示进度条动画
show_progress() {
    local pid=$1
    local delay=0.1
    local spinstr='|/-\'
    local i=0
    while ps -p $pid > /dev/null; do
        i=$(( (i+1) %4 ))
        printf "\r正在检查代码... [%c] " "${spinstr:$i:1}"
        sleep $delay
    done
    printf "\r"
}

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

# 获取暂存区的差异
DIFF=$(git diff --cached)
if [ -z "$DIFF" ]; then
    echo "No changes to check"
    exit 0
fi

# 构建 JSON 数据
JSON_DATA=$(echo "$DIFF" | jq -R -s '{"diff": .}')

# 验证 JSON 格式
if ! echo "$JSON_DATA" | jq . > /dev/null 2>&1; then
    echo "Error: Invalid JSON format"
    echo "JSON data:"
    echo "$JSON_DATA"
    exit 1
fi

# 将 JSON 数据写入临时文件
TEMP_JSON_DATA=$(mktemp)
echo "$JSON_DATA" > "$TEMP_JSON_DATA"

# 构建 curl 命令
CURL_CMD="curl -s -X POST -H \"Content-Type: application/json\""

# 如果提供了 token，添加到请求头
if [ -n "$BARK_TOKEN" ]; then
    CURL_CMD="$CURL_CMD -H \"Authorization: Bearer $BARK_TOKEN\""
fi

# 发送请求到 bark 服务
echo "Sending request to: $BARK_URL/precommit"
RESPONSE=$(eval "$CURL_CMD -d @\"$TEMP_JSON_DATA\" \"$BARK_URL/precommit\"") &
CURL_PID=$!

# 显示进度条
show_progress $CURL_PID

# 等待请求完成
wait $CURL_PID
CURL_EXIT_CODE=$?

# 清理临时文件
rm -f "$TEMP_JSON_DATA"

# 检查 curl 命令是否成功
if [ $CURL_EXIT_CODE -ne 0 ]; then
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