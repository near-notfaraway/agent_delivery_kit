# Agent Delivery Kit 开发指南

这是一个 Go CLI 和一组交付 Skills。初次接手先看 [README](README.md)；使用工具包交付业务需求时，按 [delivery-flow](skills/delivery-flow/SKILL.md) 和目标仓库的指令执行。

## 代码结构

```text
cmd/coordctl/      命令入口和 CLI 测试
internal/control/  任务状态、依赖判断、文件锁和持久化
skills/            安装入口、交付流程、阶段指令和模板
agents/            可选专业角色的提示词
templates/         CLI 任务图模板
docs/              CLI 参考和 README 配图
```

**主要入口**

- 改命令、参数、输出或退出码：从 [`cmd/coordctl/main.go`](cmd/coordctl/main.go) 的 `run` 入手，同时检查 [`main_test.go`](cmd/coordctl/main_test.go) 和 [CLI 参考](docs/cli.md)。
- 改任务认领、依赖或状态存储：看 [`internal/control/store.go`](internal/control/store.go) 和 [`store_test.go`](internal/control/store_test.go)。涉及 `project.json` 格式时，核对现有版本 1 文件和[任务图模板](templates/project.template.json)。
- 改安装或初始化：看 [入口 Skill](skills/agent-delivery-kit/SKILL.md) 和 README 的「快速开始」。改交付行为时，从 [delivery-flow](skills/delivery-flow/SKILL.md) 找到对应的阶段指令、参考文件或模板。

`project.json` 是 CLI 管理的任务状态，`decisions.jsonl` 由 CLI 追加；`delivery.json` 中的整体进度和确认记录由协调 Agent 维护。

## 本地开发

使用 Go 1.23 或兼容工具链，在当前工作树开发即可。构建和试运行按 README 的「构建并试运行 CLI」执行。`workspaces/` 存放本地交付材料，已被 Git 忽略；测试使用临时目录和虚构数据。

## 测试

迭代时先跑改动所属的包；交接 Go 改动前运行 `go test ./... -count=1` 和 `go vet ./...`。如果 `workspaces/` 下的 Go 源码会被扫到，用包含当前改动、不含该目录的临时副本验证。

只改文档或 Skills 时，检查格式、链接和相关说明是否一致。可运行示例变更时，验证受影响的示例。
