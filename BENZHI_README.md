# BENZHI_README

这是一个面向多租户具身机器人舰队运营的 Go 后端服务，负责机器人注册、任务审批与调度、遥测采集、远程交接、告警审计及可靠事件投递。

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
./build_benzhi_docker.sh benzhi-task-233-amd64 linux/amd64
./build_benzhi_docker.sh benzhi-task-233-arm64 linux/arm64
docker run -it benzhi-task-233-amd64:latest
docker run -it --platform linux/arm64 benzhi-task-233-arm64:latest
```
