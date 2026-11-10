# Agent Delivery Kit

从需求到代码交付的 AI Agent 协作工具包，支持跨仓库执行、人工评审与持久交接。

当一个需求涉及多个仓库和独立 Agent 会话时，上下文传递、依赖确认和验证汇总往往需要人工反复协调。本项目通过 Markdown Skills、模板和本地 CLI `coordctl`，保存任务边界、依赖、设计决策与交接证据，让新的 Agent 会话能够继续已有工作。

## 核心能力

- **完整交付流程**：覆盖需求澄清、方案确认、仓库准备、任务开发、逐项评审和集成验收。
- **任务认领与依赖检查**：用文件锁保护认领，拒绝启动前置任务未完成的任务。
- **跨仓库交接**：保存固定修订、验证结果和下游说明，供后续任务读取。
- **跨会话恢复**：任务状态与设计偏离记录保存在本地文件中，可由独立 CLI 进程读取。

## 适用场景

适合需要组织 AI Agent 在多个独立仓库和 IDE / CLI 会话中完成同一需求的技术负责人或工程师，尤其是需要交接接口版本、保留人工评审和恢复执行上下文的场景。

项目与模型、IDE 和企业需求系统解耦。CLI 管理持久状态，Skills 提供交付方法；Agent 由使用者所在的环境运行。

## 运行要求

- **CLI**：Go 1.23 或兼容工具链，支持 Unix `flock` 的系统，以及可读写的共享本地文件系统。Windows 尚未验证。
- **交付流程**：能阅读 Markdown、读写授权文件和执行命令的编码 Agent，以及参与评审和确认的用户。
- **业务仓库**：Git、对应语言工具链、可获取的依赖和任务所需的隔离验证环境。

无需特定模型或 IDE，也无需安装外部 Skills 或自定义 Agent。支持通用子 Agent 和独立会话；只有一个执行者时可按依赖顺序串行处理。

## 快速开始

### 1. 构建并试运行 CLI

克隆仓库后，在项目根目录执行：

```bash
CLI_DIR="$(mktemp -d "${TMPDIR:-/tmp}/agent-delivery-kit-cli.XXXXXX")" &&
GOCACHE=/private/tmp/agent-coordinator-go-cache go build -o "$CLI_DIR/coordctl" ./cmd/coordctl &&
"$CLI_DIR/coordctl" summary --project templates/project.template.json
```

空任务模板的输出为：

```json
{
  "counts": {},
  "ready": []
}
```

这一步确认 CLI 可运行。试运行二进制位于 `$CLI_DIR`，用完可删除该临时目录；正式需求使用流程建立的独立控制目录和二进制。全部命令及字段说明见 [CLI 文档](docs/cli.md)。

### 2. 将需求交给协调 Agent

在能够访问本仓库及业务仓库的 Agent 会话中，提供本项目的绝对路径、需求资料和仓库路径。例如：

> 用 Agent Delivery Kit 处理这个需求。工具包路径：`/absolute/path/to/agent_delivery_kit`。原始需求：……。涉及仓库：……。请先读取工具包的 `AGENTS.md` 和 `skills/delivery-flow/SKILL.md`，整理需求供我确认；确认需求后编写技术方案，方案确认后再准备仓库并开展开发。

Skills 是 Agent 阅读的指令文件，`coordctl` 不会自动执行它们。协调 Agent 按 [delivery-flow](skills/delivery-flow/SKILL.md) 组织以下流程：

```text
需求澄清与确认 → 技术方案与确认 → 仓库准备
             → 任务执行与逐项评审 → 集成验收与交付
```

材料默认保存在被 Git 忽略的 `workspaces/<需求名>/`：`requirement.md` 和 `design.md` 用于人工评审，`delivery.json` 保存整体进度，`control/` 保存任务图与执行上下文。完成后从 `delivery-report.md` 查看各仓库修订、验证结果及交接说明。

## 任务与交接示例

以下用虚构需求说明任务图与交接格式。路径需替换为实际授权 worktree；示例不代表需求、方案或提交已经获批。

```json
{
  "version": 1,
  "name": "item-contract",
  "tasks": [
    {
      "id": "T1",
      "repo": "/absolute/path/to/requirement/repos/api",
      "title": "发布 Item v1 契约",
      "objective": "定义稳定的 id 和 name JSON 字段",
      "acceptance": ["契约测试通过，JSON 字段名保持稳定"],
      "handoff_requirements": ["固定提交、契约版本及获取方式"],
      "depends_on": [],
      "status": "pending"
    },
    {
      "id": "T2",
      "repo": "/absolute/path/to/requirement/repos/sdk",
      "title": "适配 SDK 的 Item 类型",
      "objective": "按 T1 的固定契约版本实现公开类型",
      "acceptance": ["SDK 序列化与上游契约一致"],
      "handoff_requirements": ["固定提交、验证结果及下游使用说明"],
      "depends_on": ["T1"],
      "status": "pending"
    }
  ]
}
```

初始 `summary` 仅返回 `T1` 为 ready。T1 认领、实现、评审并获提交授权后，使用实际修订和验证结果运行 `complete`。CLI 保存的任务字段片段如下，尖括号内容需替换为真实值：

```json
{
  "id": "T1",
  "status": "completed",
  "owner": "api-agent",
  "result": {
    "commit": "<已验收的完整 Git SHA>",
    "validation": ["<实际运行的契约测试及结果>"],
    "handoff": "<Item v1 的固定依赖版本、获取方式及验证证据>"
  }
}
```

随后重新运行 `summary`，T2 成为 ready；`context --task T2` 返回其任务边界及 T1 的完成证据。状态由 CLI 维护，不能直接编辑活跃文件套用这段示例。存在 Git commit 本身不解锁依赖，必须成功完成 `complete`。

## 当前状态与限制

已通过一个真实跨仓库需求验证，覆盖需求与方案确认、任务执行、逐项评审、交接和本地集成验证。现有自动化回归覆盖认领、依赖、完成证据、上下文、决策、跳过和状态汇总；独立 CLI 进程的状态恢复与并发认领尚未纳入自动化回归。

- 只支持共享本地文件系统上的协调，跨机器协作尚未验证。
- 没有任务租约、心跳、超时回收和自动重试；中断后按 [恢复约定](skills/delivery-flow/references/recovery-and-corrections.md) 继续。
- CLI 不会启动或托管 Agent，也不自动生成文档或落实人工评审、提交及发布授权。
- `needs_approval` 只记录审批要求，尚未形成强制门禁；修订和证据字段的真实性仍需实际核验。
- 任务认领不解决同一仓库内的所有修改冲突，仓库写入仍需协调。
- 完全不配置自定义 Agent、也不安装外部 Skills 的全新环境完整交付，以及新产物管理约定的落实效果，尚待验证。

后续重点是不同使用者和全新环境中的复用、第二个真实用例，以及效率收益的量化评估。当前范围保持为本地交付工具包。

## 文档导航

| 文档 | 用途 |
|---|---|
| [CLI 文档](docs/cli.md) | 命令、控制文件字段、输出和退出码 |
| [delivery-flow](skills/delivery-flow/SKILL.md) | 完整交付、调度与跨会话恢复入口 |
| [requirement-analysis](skills/requirement-analysis/SKILL.md) / [需求模板](skills/requirement-analysis/assets/requirement_template.md) | 需求整理与确认 |
| [technical-design](skills/technical-design/SKILL.md) / [方案模板](skills/technical-design/assets/design_template.md) | 技术设计与方案评审 |
| [repository-preparation](skills/repository-preparation/SKILL.md) | 基线核查与 worktree 准备 |
| [delivery-review](skills/delivery-review/SKILL.md) / [integration-validation](skills/integration-validation/SKILL.md) | 单任务评审与跨仓库验收 |
| [Agent 委派与可选角色配置](skills/delivery-flow/references/agent-delegation.md) | 执行者职责、调用输入与专业角色配置 |
| [验证与证据](skills/delivery-flow/references/validation-and-evidence.md) / [产物管理](skills/delivery-flow/references/artifact-management.md) | 验证标准、替代路径与资源管理 |
| [控制文件模板](templates/project.template.json) | version 1 的空任务图 |
| [AGENTS.md](AGENTS.md) | 仓库协作与开发约束 |

## 反馈与贡献

欢迎通过仓库 Issues 提交使用问题、缺陷和真实复用反馈，通过 Pull Request 提交改进。问题报告请包含工具链、复现步骤、预期与实际结果，并使用可公开的虚构任务数据。

贡献前阅读 [AGENTS.md](AGENTS.md)。当前优先处理正确性、兼容性、已有行为的测试和使用文档改进；功能扩展需有实际使用依据。

## 许可证

采用 [MIT License](LICENSE)。
