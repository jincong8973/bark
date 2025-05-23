# Bark - GitLab MR Review Bot

## 环境变量

以下环境变量可以在运行时配置：

- `BARK_GITLAB_TOKEN`: GitLab 访问令牌
- `BARK_DEEPSEEK_TOKEN`: DeepSeek API 令牌
- `BARK_GITLAB_URL`: GitLab 服务器地址
- `BARK_DEEPSEEK_URL`: DeepSeek API 地址

## 配置文件

复制 `config.yaml.example` 到 `config.yaml` 并根据需要修改配置：

```bash
cp config.yaml.example config.yaml
```

## Docker 运行

```bash
# 构建镜像
docker build -t bark .

# 运行容器
docker run -d \
  -p 18080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  bark
```