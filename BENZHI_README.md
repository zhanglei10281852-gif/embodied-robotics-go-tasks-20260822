# BENZHI_README

这是一个基于 Go 实现的后端服务，用于承载 embodied-robotics-go-tasks-20260822 的业务处理、数据管理与稳定运行。

## 项目说明

- 项目：zhanglei10281852-gif/embodied-robotics-go-tasks-20260822
- 项目用途：This repository is a self-built production-style Go backend for multi-tenant embodied robot fleet operations. It coordinates robot registration, capability versions, mission approval, scheduling, telemetry ingestion, policy review, remote handoff, alerts, audit trails and durable outbox delivery.
- Go 工具链：`golang:1.26`
- 前端工具链：无

## 标准构建、运行和测试命令

进入容器后执行：

```bash
# 编译
cd '/app' && GOTOOLCHAIN=local go build ./...

# 启动
cd '/app' && GOTOOLCHAIN=local go run ./cmd/server

# 测试
cd '/app' && GOTOOLCHAIN=local go test ./...
```

## Docker 构建和进入容器

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh benzhi-task-231-amd64 linux/amd64
./build_benzhi_docker.sh benzhi-task-231-arm64 linux/arm64
docker run -it benzhi-task-231-amd64:latest
docker run -it --platform linux/arm64 benzhi-task-231-arm64:latest
```

## 题目验证命令

1. 预期退出码 0：`go test ./internal/httpapi -run '^TestMissionHandlerStopsOnCancelledRequest$' -count=1`
2. 预期退出码 0：`go test ./...`
3. 预期退出码 0：`GOTOOLCHAIN=local go build -buildvcs=false ./... && GOTOOLCHAIN=local go vet ./...`
