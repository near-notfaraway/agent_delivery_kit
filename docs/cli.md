# coordctl CLI

`coordctl` 读取 version 1 的 `project.json`，管理任务认领、依赖、完成证据、上下文和设计偏离记录。构建与试运行见 [README](../README.md#快速开始)，实际交付准备见 [delivery-flow](../skills/delivery-flow/SKILL.md)。

## 使用已有控制目录

以下命令在需求的控制目录已准备完成后使用。将路径替换为实际值：

```bash
CONTROL_DIR=/absolute/path/to/requirement/control
COORD="$CONTROL_DIR/bin/coordctl"
PROJECT="$CONTROL_DIR/project.json"
```

`POC_DIR` 在流程中表示工具包的绝对路径。需求与方案、整体记录和控制材料的位置由交付流程定义；CLI 只读取命令指定的控制文件，不自动读取 Markdown 或 `delivery.json`。恢复已有需求时沿用其入口和二进制。

## 命令

### summary：查看状态

```bash
"$COORD" summary --project "$PROJECT"
```

返回 JSON：`counts` 按实际任务状态计数，`ready` 列出所有依赖均已 `completed` 或 `skipped` 的 `pending` 任务。没有可开始任务时，`ready` 为 `[]`。输出不包含文档评审状态或待审批决策汇总。

### context：读取任务与上游交接

```bash
"$COORD" context --project "$PROJECT" --task T1
```

返回 JSON：`task` 包含任务目标、仓库、验收与交接要求，`dependencies` 包含直接上游任务及其已保存的完成结果。它不自动纳入后续修复；当前修订和有效批准按 [恢复约定](../skills/delivery-flow/references/recovery-and-corrections.md) 读取。

### claim：认领任务

```bash
"$COORD" claim --project "$PROJECT" --task T1 --owner api-agent
```

任务必须为 `pending`，且全部依赖已 `completed` 或 `skipped`。成功后状态为 `claimed`，`owner` 保存认领者标识。同一任务的并发认领只有一个成功；重复认领已 `claimed` 的任务会失败，不能用于恢复。

### complete：保存完成证据

```bash
"$COORD" complete \
  --project "$PROJECT" --task T1 --owner api-agent \
  --commit "<已验收的完整 Git SHA>" \
  --validation "<实际运行的验证及结果>" \
  --handoff "<固定依赖版本、获取方式及下游说明>"
```

`owner` 必须与记录中的认领者一致。`commit`、`validation` 和 `handoff` 均为必填；成功后状态为 `completed`，证据保存到 `result`。`--validation` 接收一个字符串，保存为含该字符串的数组。

CLI 检查证据字段非空，不检查 Git 修订是否存在、测试是否真实运行或提交是否获批。实际差异评审、验证与授权由 [delivery-review](../skills/delivery-review/SKILL.md) 和交付流程落实。仅有 Git commit 不解锁依赖，成功执行 `complete` 后下游才可认领。

### decide：记录设计偏离

```bash
"$COORD" decide \
  --project "$PROJECT" --id D1 --task T2 \
  --original "原计划直接引用上游仓库" \
  --problem "会导致两个仓库的发布周期耦合" \
  --choice "改用本地 DTO，并通过契约测试保护" \
  --needs-approval
```

除 `--needs-approval` 外，上述参数均为必填。决策追加到 `project.json` 同目录的 `decisions.jsonl`，每行一个 JSON 对象，包含 `id`、`task_id`、`original`、`problem`、`choice` 和 `needs_approval`。

`--needs-approval` 默认为 false，只记录人工审批要求，不会阻止后续命令执行。真实确认及方案同步按 [delivery-flow](../skills/delivery-flow/SKILL.md#8-处理设计偏离阻塞与返工) 执行。

### skip：跳过无需实施的条件任务

仅对实际任务图中无需实施的条件任务使用，将任务 ID 替换为真实值：

```bash
"$COORD" skip \
  --project "$PROJECT" --task "<条件任务 ID>" \
  --reason "现有公开类型已兼容新契约，无需修改"
```

`reason` 必填。成功后状态为 `skipped`，理由保存到 `resolution`；下游依赖检查将该任务视为已满足。CLI 不判断跳过理由和下游产物是否充分，这项核查由协调者负责。

## 控制文件字段

[空任务模板](../templates/project.template.json) 和 [两任务示例](../README.md#任务与交接示例) 展示了 version 1 格式。初始化任务图按确认方案生成；运行后通过 CLI 修改任务状态。

| 项目字段 | 含义 |
|---|---|
| `version` | 当前格式为 `1` |
| `name` | 需求名称 |
| `tasks` | 任务数组 |

| 任务字段 | 含义 |
|---|---|
| `id` | 稳定、唯一的任务编号 |
| `repo` | 本任务唯一授权 worktree 的绝对路径 |
| `title` / `objective` | 任务名称与目标边界 |
| `acceptance` | 验收要求的字符串数组 |
| `handoff_requirements` | 下游交接要求的字符串数组 |
| `depends_on` | 直接前置任务 ID 数组 |
| `status` | `pending`、`claimed`、`completed` 或 `skipped` |
| `owner` | 认领者标识，由 `claim` 保存 |
| `result` | 完成证据，由 `complete` 保存 |
| `resolution` | 跳过理由，由 `skip` 保存 |

`result` 包含 `commit` 字符串、`validation` 字符串数组和 `handoff` 字符串。CLI 不负责校验全部任务图结构、授权路径或格式兼容性，初始化核查由交付流程完成。

整体阶段、文档确认和当前交付修订保存在协调者维护的 `delivery.json` 中，不属于 CLI 任务格式。现有需求的材料位置及版本 1 用法保持有效，无需因文档整理迁移状态。

## 输出与退出码

`summary` 和 `context` 成功时向标准输出写入 JSON；其余命令成功时输出操作结果文本。错误写入标准错误。

| 退出码 | 含义 |
|---|---|
| `0` | 命令成功 |
| `1` | 读写、任务、依赖、所有者或 JSON 输出等操作失败 |
| `2` | 命令不存在、参数解析失败或必填参数缺失 |

CLI 不自动重试。中断、部分成功、已认领任务和已完成任务的修复按 [恢复与返工](../skills/delivery-flow/references/recovery-and-corrections.md) 核对实际状态后续接。
