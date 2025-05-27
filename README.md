# Bark - GitLab MR Review Bot


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

## Pre-commit Hook 设置

1. 安装 pre-commit：
```bash
pip install pre-commit
```

2. 在项目根目录创建 `.pre-commit-config.yaml`：
```yaml
repos:
  - repo: https://github.com/jincong8973/bark
    rev: v0.0.1
    hooks:
      - id: bark-code-review
        name: Bark Code Review
        description: 使用 Bark 进行代码审查
        args: ["--url", "https://your-bark-service.com", "--token", "your-bark-token"]
```

3. 安装 pre-commit hooks：
```bash
pre-commit install
```

现在，每次提交代码时，pre-commit hook 都会自动检查修改的文件，并调用 bark 服务进行代码审查。如果发现问题，提交将被阻止。