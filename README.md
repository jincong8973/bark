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

1. 安装依赖：
```bash
# macOS
brew install jq

# Ubuntu/Debian
apt-get install jq

# CentOS/RHEL
yum install jq
```

2. 安装 pre-commit：
```bash
pip install pre-commit
```

3. 在项目根目录创建 `.pre-commit-config.yaml`：
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

4. 安装 pre-commit hooks：
```bash
pre-commit install
```

5. 单独运行 pre-commit 检查：
```bash
# 检查所有文件
pre-commit run --all-files

# 检查特定文件
pre-commit run --files file1.go file2.go

# 检查特定 hook
pre-commit run bark-code-review

# 检查特定 hook 的特定文件
pre-commit run bark-code-review --files file1.go file2.go
```

现在，每次提交代码时，pre-commit hook 都会自动检查修改的文件，并调用 bark 服务进行代码审查。如果发现问题，提交将被阻止。