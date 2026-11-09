---
name: delivery-flow
description: Coordinate analysis, design review, worktrees, and multi-agent delivery with coordctl. Use when users ask to develop a requirement or resume delivery. Not for isolated edits.
---

# 完整交付流程

由通用 mainAgent 承担协调职责，维护需求从初始化到开发完成的持久记录，按依赖组织各仓库执行者，核验交付证据。复用现有 CLI 和宿主实际提供的 Agent 能力，不要求自定义协调/编码 Agent，不增加模型托管、数据库或通用工作流引擎。

用户只要求创建或修改本 Skill 时，仅处理 Skill 文件，不为这次编辑启动需求交付流程。

## 1. 恢复上下文与执行边界

- 读取当前工作区与涉及仓库的 `AGENTS.md`。在 Agent Delivery Kit 项目中修改内容前，读取 `README.md`，并按改动范围读取相关材料。
- 定位本 Skill 所属项目的绝对路径 `POC_DIR`，确认 `cmd/coordctl`、`templates/` 和下述依赖 Skill 存在。通过链接读取时先解析真实项目路径，不依赖终端当前目录。
- 按阶段读取并执行 [requirement-analysis](../requirement-analysis/SKILL.md)、[technical-design](../technical-design/SKILL.md)、[repository-preparation](../repository-preparation/SKILL.md)、[delivery-review](../delivery-review/SKILL.md) 和 [integration-validation](../integration-validation/SKILL.md)。Skill 是指令文件，不是假想的可调用 API。若缺失，先定向查找；无法找到时记录阻塞，不冒充已调用。
- 执行者选择、完整调用输入、权限与 writer、反馈及返回结果按 [Agent 委派约定](references/agent-delegation.md)，包括自定义 Agent 未配置、平台不支持及角色调用失败时的默认路径。自定义角色均为增强项，缺失不阻塞流程；使用通用 subAgent，或由 mainAgent 直接执行同一 Skill。仓库准备由 mainAgent 直接执行；编码默认交给通用 subAgent，读取 [单仓任务执行](references/task-execution.md)，没有子 Agent 能力时按该入口串行执行。
- 产物清单、尝试编号、临时目录与阶段交接检查遵循 [产物管理](references/artifact-management.md)。mainAgent 分配当前入口和可写范围，材料按实际阶段生成，不提前初始化目录或额外状态清单。
- 恢复已有文档、状态、真实确认和用户授权。继续同一需求时沿用记录；不重建目录、不重置任务、不重复请求已获得的确认。
- 恢复及部分失败按 [当前修订、恢复与返工约定](references/recovery-and-corrections.md) 比对 Git、CLI、整体记录和证据；已成功的副作用只补齐缺失步骤，不盲目重放。
- 用户要求完整交付时，方案确认后继续执行已授权的 worktree、控制面和开发步骤；不因设计 Skill 的职责终点再索取同一份开发授权。仅分析或设计的请求不自动扩展为开发。
- Git commit、push、PR、合并、部署及对外发送分别遵循用户实际授权。使用本 Skill 不等于用户已明确要求提交或发布。
- 整体验收通过后，推送本需求所有交付仓库的最新已验收提交到对应远端需求分支，作为默认收尾步骤。用户明确“验收通过后 push 所有 repo”即为该范围的条件式推送授权，记录目标与依据，条件满足后直接执行，不再次索取相同授权。仅分析、设计或明确限定本地交付的请求不扩展为推送；旧记录不能由模板更新补造授权。
- 按 [README 的运行要求与可选能力](../../README.md#运行要求与可选能力) 核查当前阶段所需能力。核心依赖是本仓库提供的五个阶段 Skill；外部 Skill 可选，不要求特定文档平台、浏览器或原生 Agent Team。按相同验收与证据标准选择工具或人工路径，能力缺失只阻塞相关工作，不自动 skip 或缩减验收；仍遵守适用仓库与宿主环境指令。

## 2. 对齐名称、分支与保存位置

从用户输入和现有材料提取需求名、所有开发仓库共用的分支名，以及用户指定的保存位置。只询问影响当前阶段的必要信息，未明确的分支和执行配置可以留到执行准备阶段；用户未指定位置时直接使用默认目录，不额外请求位置确认。

定义：

- `WORKSPACES_ROOT`：默认 `$POC_DIR/workspaces`；用户指定其他根目录时沿用。
- `DELIVERY_DIR`：默认 `$WORKSPACES_ROOT/<需求名>`；用户直接指定需求目录时使用该目录，不重复追加需求名。
- `CONTROL_DIR`：`$DELIVERY_DIR/control`，只在方案确认后创建。

真实需求默认保存在本项目的 `workspaces/`，确认根 `.gitignore` 包含 `/workspaces/`。不要强制暂存该目录或把其中的私有材料复制到公开跟踪文件。Git 可用时检查忽略规则及是否已有被跟踪内容；忽略规则不会移除已有跟踪。恢复已有需求时沿用记录的位置，不因默认值变化自动迁移。

需求名必须是单一目录名，拒绝空值、`.`、`..`、路径分隔符及目录逃逸；需要改名时先对齐实际名称，不静默替换。使用 `git check-ref-format --branch` 校验分支名。所有路径、变量都要正确引用。

写入前检查父目录和目标是否存在。若同名目录存在，核对需求身份及已有记录后恢复；冲突目录不得覆盖。名称与位置确定后，创建需求目录，从 [delivery.template.json](assets/delivery.template.json) 生成 `delivery.json`，填写真实名称、绝对路径、已明确的分支、初始化时间和授权依据，初始 `status` 为 `init`。尚未明确但不影响需求分析的执行配置保留未确定值，在相应阶段核查。

`delivery.json` 是本 Skill 管理的整体状态文件，从初始化开始维护。调用 `requirement-analysis` 时，该 Skill 只产出 `requirement.md`，版本、澄清和确认都在文档内；由协调者同步整体状态。技术设计阶段再生成 `design.md`，方案确认后才生成执行控制材料。下图是完整交付阶段的可能目录，除初始化状态外的材料按阶段生成：

```text
workspaces/<需求名>/
├── delivery.json          # 初始化时创建，保存整体阶段与恢复信息
├── requirement.md         # 人工评审入口
├── design.md              # 人工评审入口
├── sources/               # 协调者按需保存必要原始材料
├── repos/<仓库别名>/       # 方案确认后建立的 Git worktree
├── control/               # worktree 准备完成后生成
│   ├── requirement.md     # 已确认完整快照，包含版本与确认依据
│   ├── design.md
│   ├── execution.md
│   ├── project.json       # version 1 任务状态与证据
│   ├── decisions.jsonl    # 首次 decide 时由 CLI 创建
│   ├── baselines/         # 历次确认版本的不可覆盖快照
│   └── bin/coordctl       # 本需求独立二进制
├── reviews/<scope>.md      # 当前评审入口
├── validation/            # 当前验证入口
│   ├── attempts/<scope>/<attempt>/ # 必要历史证据
│   └── artifacts/         # 按需固定可消费产物
├── .runtime/<scope>/<attempt>/    # 可重建运行资源，使用后核验处理
└── delivery-report.md     # 最终交付说明
```

相关文档、日志和报告集中在需求目录。仓库本身、外部来源文档和工具可以保留原位置，在记录中引用；不复制凭证、无关数据或整个私人会话。

图示为新执行的默认布局，目录按需创建。已有需求沿用其入口和证据位置；一次性脚本、中间载荷和下载缓存不堆放根目录，非标准材料按产物管理约定说明用途。

## 3. 管理状态与确认依据

从初始化开始，整体阶段写入 `delivery.json`。文档状态及确认记录遵循需求分析/技术设计 Skill 和各自模板，协调者同步真实确认版本、范围与依据。恢复旧文档时沿用已有真实记录，不因模板更新重建状态。

整体状态只由协调者维护，`requirement-analysis` 和 `technical-design` 不创建或修改 `delivery.json`。不能将整体状态混入 `project.json`；后者保持 version 1，任务状态仍为 `pending / claimed / completed / skipped`。

| 整体状态 | 进入条件及下一步 |
|---|---|
| `init` | 名称、目录及已明确的执行配置已记录，初始化状态已保存 |
| `analyzing_requirement` | 开始执行需求分析 Skill |
| `awaiting_requirement_review` | 可审阅需求稿已完成，等待用户确认具体版本 |
| `requirement_confirmed` | 需求关键问题已解决，存在对应版本的真实确认 |
| `designing` | 基于确认需求执行技术设计 Skill |
| `awaiting_design_review` | 技术评审稿已完成，等待用户 Review |
| `design_confirmed` | 技术评审问题已解决，方案版本得到真实确认 |
| `preparing_repos` | 开始核对开发基线并准备各仓库 worktree |
| `repos_ready` | 所有改动仓库的路径、分支与基线验证通过 |
| `preparing_control` | 开始生成快照、执行说明及 version 1 任务图 |
| `control_ready` | 文档和任务图一致、CLI 检查通过、可分发任务 |
| `developing` | 仓库 Agent 按依赖认领、开发、评审和验证 |
| `verifying` | 所有任务已结束，进行跨仓库整体验收与证据复核 |
| `completed` | 全部完成条件满足，最终报告已落盘 |
| `blocked` | 环境、权限、提交授权或执行故障阻止继续推进 |

正常按表顺序推进，需求/方案评审等待使用对应的 awaiting 状态，不冒充已确认或执行失败。开发阶段的单任务人工 Review 或提交等待保留 `claimed`，在该任务的 `reviews/` 材料和 `delivery.json` 中记录；不增加 CLI 任务状态，也不退回需求/方案评审阶段。进入 `blocked` 时保存 `resume_status`、阻塞原因、影响任务及解除条件。只要有独立 Agent 正在执行或还有可推进任务，保留 `developing`，将局部等待记录在 `blockers` 中；所有可推进工作都已受阻时才将整体置为 `blocked`。恢复已有材料时只记录真实发生的接入和后续转移，不补造早期历史。

每次转移同时保存 `updated_at` 和 `history` 条目，包含前后状态、时间、原因与证据路径。确认记录必须包含文档版本、Reviewer、时间和实际确认依据；草稿生成、局部问题回答、任务完成均不能替代整份文档确认。

模板中的 `null` 表示尚未确认或尚未生成，不能当成已授权。时间采用带时区的 ISO 8601 格式，记录结构保持一致：

- `history[]`：`from`、`to`、`at`、`reason`、`evidence`；初始记录的 `from` 为 null。
- `approvals[]`：稳定 `id`、`stage`、`version`、`reviewer`、`at`、`evidence`，决策确认另含 `decision_id`。任务 Review/提交确认另注明 `task_ids`，修复确认含 `correction_id`，在证据中记录审阅的 commit 或工作树摘要；批准一项任务不扩展到其他任务。旧记录无 ID 时兼容恢复，不重复追加同一批准。
- `authorization` 的非空项：`scope`、`at`、`evidence`；只记录明确授权的动作和范围。
- `blockers[]`：`id`、`reason`、`task_ids`、`unblock_condition`、`resolved_at`。
- `agents[]`：`session_id`、`owner`、`task_ids`、`repo`；没有真实会话时不填造假的 ID。
- `repositories[].delivery_commit`：当前已评审、获对应授权并提交的仓库修订；`delivery_review` / `delivery_validation` 链接对应证据。缺失时读取相关任务 `result.commit`，旧字段含义不变。
- `corrections[]`：按稳定 ID 保存修复及其原任务、父提交、当前提交、证据、批准和下游影响，字段见 [恢复约定](references/recovery-and-corrections.md)。未形成的事实保留 null，不伪造批准或提交。
- `pushes[]`：逐仓普通推送及远端核验记录，包含稳定 `id`、`repo`、`remote`、`ref`、固定 `commit`、实际 `remote_commit`、`result`、`verified_at`、`evidence`；失败或未核验字段如实保留，远端地址与证据不包含凭证。缺少该可选字段的旧实例仍可恢复。

协调者独占写入 `delivery.json`，遵循现有锁、同目录临时文件完整写入、刷新后原子替换，保留未知字段；子 Agent 不改整体状态。交接协调权前确认旧协调者停止写入，不能只依赖记录中的名字推断互斥。批准和修复按稳定身份去重，history 只记录真实阶段转移。

恢复时检查已有记录、文件、worktree、CLI 状态及证据是否一致，修复中断处再继续。只有需求或方案而没有 `delivery.json` 时，沿用原文档入口，由协调者初始化状态并据真实记录恢复到首个未获确认的阶段，不为接入新建报告或迁移材料。文件存在不是阶段完成依据。基线变更可以显式退回需求或方案阶段，保留历史及受影响确认记录，不静默回退或跳过门禁。

## 4. 需求分析与技术评审

1. 进入 `analyzing_requirement`，按委派约定提供原始来源、名称、文档位置及已有结论，由 mainAgent 或需求设计 subAgent 执行 `requirement-analysis`。
2. 核验返回的 `requirement.md` 及真实状态。mainAgent 组织需要的澄清和整稿 Review，反馈由文档 writer 按需求 Skill 写回；稿件就绪进入 `awaiting_requirement_review`，实际整稿确认后同步 `delivery.json` 并进入 `requirement_confirmed`。
3. 进入 `designing`，提供确认需求及仓库只读入口，由 mainAgent 或需求设计 subAgent 执行 `technical-design`。该 Skill 负责模板、仓库核查、模块契约与验证计划，不提前创建 worktree 或控制目录。
4. 核验返回的 `design.md` 和未决事项，进入 `awaiting_design_review`，mainAgent 展示具体版本并组织反馈；按所属 Skill 同步受影响文档。整稿确认后同步 `delivery.json` 并进入 `design_confirmed`，继续授权范围内的准备。

需求和方案服务人工评审，Agent 命令和执行任务放入后续控制目录。会话已有确认时保留授权直接推进，不能要求用户重复确认。

## 5. 调用仓库准备

先核对需求和方案都存在对应版本的真实确认，且 `delivery.json` 已同步到 `design_confirmed`。进入 `preparing_repos`，从确认方案提取仓库清单，区分需要修改、只读参考和外部团队负责的仓库。只为授权修改的仓库创建 worktree。

mainAgent 直接执行 [repository-preparation](../repository-preparation/SKILL.md)，使用确认版本、需求目录/分支、逐仓授权和已有准备记录，形成该 Skill 规定的准备证据入口。

协调者逐仓保存源路径、职责、ref/commit、分析修订、基线授权、需求分支/worktree、指令/命令及实际结果。来源未解决时保留成功 worktree，整体留在 `preparing_repos`；所有授权修改仓库均创建并核验成功后进入 `repos_ready`，才生成控制材料。

后续开发和需求内修复默认使用原需求 worktree，同一 worktree 最多一个写入者，不覆盖活跃源码或二进制。具体复用检查和基线纠正按仓库准备 Skill，不由协调者另设一套准备流程。

## 6. 构造 coordctl 控制文件

进入 `preparing_control`，创建 `CONTROL_DIR`：

1. 将已确认需求和方案复制为完整快照，并在 `baselines/` 保存不可覆盖基线，记录版本、确认依据及内容摘要。相同内容引用已有基线并在索引记录实际确认，不重复创建 before/copy 副本；新内容才保存新基线。当前快照与评审入口保持一致。
2. 根据确认方案生成 `execution.md`，写清绝对路径、确认版本、R / C / V 到任务映射、依赖原因、仓库边界、验证命令和交接要求。
3. 从 `templates/project.template.json` 生成 `project.json`，保留 `version: 1`。根据实际方案生成任务，不复制 README 的示例任务图。
4. 每项任务填写唯一 `id`、唯一授权 worktree 的绝对路径 `repo`、`title`、`objective`、`acceptance`、`handoff_requirements`、`depends_on` 和初始 `status: pending`。不添加整体阶段、审批或 Agent 调度字段。
5. 检查任务 ID 唯一、引用存在、无自依赖或环、仓库授权正确，且方案范围被完整覆盖。同一 worktree 的多个写任务用依赖或协调调度串行化。
6. 验收须包括符合设计、必要代码 Review、仓库要求的检查与相关测试；交接须包括接口/数据契约、不可变修订、可消费的固定依赖/产物及解析验证、兼容及下游使用方式，按 [验证与证据约定](references/validation-and-evidence.md)。
7. 在 `POC_DIR` 为本需求构建独立 CLI，不覆盖活跃二进制：

```bash
GOCACHE=/private/tmp/agent-coordinator-go-cache go build -o "$CONTROL_DIR/bin/coordctl" ./cmd/coordctl
```

设置 `COORD="$CONTROL_DIR/bin/coordctl"`、`PROJECT="$CONTROL_DIR/project.json"`。运行 `summary`，逐任务运行 `context`，核对输出和预期任务图一致。CLI 不负责校验所有图结构、文档确认、Git 修订真实性或证据质量，以上检查由协调者完成。

控制文件初始化后，任务状态通过 CLI 修改。不要将真实活跃 `project.json` 或 `decisions.jsonl` 当成开发/测试夹具。检查通过后进入 `control_ready`。

准备脚本只实现本次初始化时放本次 runtime；结果写入既有准备/控制验证入口，不要求每个需求另建初始化报告、准备 JSON 或记录脚本。状态写入辅助必须符合恢复约定，不把一次性脚本作为可盲目重跑的标准工具。

## 7. 调度多个仓库 Agent

使用当前环境实际提供的子 Agent 能力，为可并行的独立仓库分派通用编码 subAgent；执行依赖来自 `summary.ready` 与真实方案。CLI 只管理持久状态，不会启动或托管 Agent。

进入 `developing`。按 [委派约定](references/agent-delegation.md) 提供完整输入和唯一 owner，指定读取 [单仓任务执行](references/task-execution.md)。在 `delivery.json.agents` 记录实际会话与 owner 的对应关系，不伪造宿主会话 ID。

以单任务为单位推进开发、人工 Review、提交与依赖解锁。独立任务可以并行启动；任一任务达到可审阅条件后立即处理该任务，不等待同批其他任务全部结束，也不将多个仓库的提交确认捆绑为默认门槛。用户对本批任务另有明确评审或授权安排时沿用其安排。

执行者不可用时按委派约定选择已有通用执行者；没有子 Agent 能力时，mainAgent 按依赖串行处理。用户已安排独立 IDE / CLI 会话时可提供相同任务输入，不将配置 Agent 或启动额外会话作为串行交付的前置。缺少角色不构成调度阻塞；实际缺失的执行能力或验收条件仍按影响范围记录。

编码者按单仓任务执行入口读取 context、认领/恢复、实现、验证和返回材料。mainAgent 核验其真实执行状态，不因 Agent 已返回就当作任务完成，也不在同一 worktree 并行派发 writer。

每个任务在 `complete` 前按以下步骤推进：

1. mainAgent 执行或委派 [delivery-review](../delivery-review/SKILL.md)，提供任务、实际 diff、确认版本、有效决策和验证材料。编码者按反馈修复，Reviewer 复核并返回该 Skill 规定的单任务证据；核查自检与独立审查的真实边界，并按产物管理核验当前入口、历史、可读取证据和资源处理。
2. 材料就绪就向用户交付本任务的具体差异、Review、验证及限制，缺少覆盖范围的授权时请求 Review/提交确认，不等待独立任务。沿用已明确覆盖的授权，流程安排、方案确认或另一任务批准不能替代本任务授权。
3. 将真实批准按稳定 ID、实际范围、冻结摘要及依据记录到 Review 和 `delivery.json`。检查失败默认阻塞该项验收，仅有真实批准的具体验收例外时按范围交接，原失败及补验条件保留。
4. mainAgent 确认验收和有效授权满足后，当前 writer 按单仓任务执行约定创建新增提交、核验冻结内容并执行 `complete`，返回真实修订和 CLI 结果。
5. mainAgent 核对实际 Git/CLI 与最终 Review、验证及可消费 handoff 一致，核对新增材料及遗留资源，再更新当前 `delivery_commit` 和证据入口。中途失败按恢复约定补齐缺失步骤，不重放已成功动作；实际 complete 已成功时不因补证据问题重复 complete。

每次收到单任务完成报告或用户 Review 反馈，先处理该任务的可执行后续步骤。向用户请求 Review 只暂停该任务的提交及依赖，不停止其他独立 Agent，也不等待它们结束；继续推进其他已授权的 ready 任务。

有任务可 Review 且需要用户反馈或提交授权时，mainAgent 核验并展示该任务的具体材料，保存当前交付状态、证据入口、持久 owner、待确认范围和续接信息，然后结束本轮回复以释放交互。不得为了等待其他任务而等待整批完成，也不持续阻塞在等待调用或轮询中。任务保持 `claimed`，仍有独立工作可推进时整体保持 `developing`；释放交互不代表停止其他执行者或结束交付。收到授权或实际完成通知后，先核对持久状态、最终差异和原 writer，再推进对应任务；已有覆盖范围的授权则直接继续，不为释放交互重复请求批准。其他任务能否在本轮回复结束后独立运行，以宿主实际能力为准，不假定支持。

单任务提交且 `complete` 成功后，立即重新读取 `summary`，为所有依赖均已满足的 ready 任务分派 Agent。下游重新读取 `context`，并核对当前交付指针、corrections 与有效批准决策，按确切修订集成。CLI context 的历史 handoff 不自动包含后续修复。Git commit 本身不解锁依赖，必须提交真实验证和 handoff 并成功 `complete`。不在共享 worktree 并行派发写入者。

例如 T1、T4 无依赖并行开发，T2 仅依赖 T1：T1 先完成实现、协调者审查和验证，就单独请用户 Review T1 并确认提交；确认后提交 T1、执行 `complete`、启动 T2，T4 继续执行。T2 若还有其他未完成依赖，则仍等待对应依赖。具备前置的局部联调可以先行，最终整体验收在全部任务结束后执行。

## 8. 处理设计偏离、阻塞与返工

执行中发现方案偏离，先通过 CLI 记录：

```bash
"$COORD" decide --project "$PROJECT" --id "$DECISION_ID" --task "$TASK_ID" \
  --original "$ORIGINAL" --problem "$PROBLEM" --choice "$CHOICE" --needs-approval
```

确实无需审批的局部实现选择可省略 `--needs-approval`，仍记录依据。影响需求、跨仓库契约、兼容、权限或发布方式的变化必须暂停受影响任务并取得确认；独立任务可继续。当前 `needs_approval` 只是日志标记，协调者负责执行门禁。

实际审批写入文档 Review 记录和 `delivery.json.approvals`，引用决策 ID；不改写历史 `decisions.jsonl` 来伪装已审批。确认后的变更同步原文档、当前快照、版本快照和执行说明。

需要修改活跃任务图时先停止受影响调度，保留控制文件副本，核对无并发写入。当前 CLI 无安全计划编辑/迁移命令：记录所需变更并请求独立的维护方案，不直接手改活跃状态，也不让用户手工编辑控制文件。不要为继续交付临时扩展 v1 协议。

条件任务确实不需要实施时，协调者核对范围与下游所需证据，再运行：

```bash
"$COORD" skip --project "$PROJECT" --task "$TASK_ID" --reason "$REASON"
```

仅对未执行且判断依据充分的条件任务使用 `skip`，不能借此结束失败、遗漏或仍有人执行的任务。CLI 将 skipped 依赖视为满足，必须确认下游不依赖被跳过任务的实际产物。

已完成任务发现问题时执行 [恢复与返工约定](references/recovery-and-corrections.md)：默认复用原需求 worktree，先核对唯一 writer、基线及活跃验证；保存稳定 correction、原因、原任务 commit 与下游影响，按本次实际范围评审、验证和取得提交授权，新增 commit 后更新当前交付指针。原 `result` 保留，不重新 claim/complete。不改变任务图的小修复无需升级为计划维护；只有真实任务边界/依赖变化才走上述维护流程。

修复提交后立即检查 summary 和下游影响，重验相关合同，继续可推进的独立工作。工具部分成功按实际 Git/CLI 状态续接，不重复记录批准或重放副作用；问题或必需复验未解决前不得标记整体完成。

## 9. 整体验收与交付

所有任务已成为 `completed` 或有充分依据的 `skipped` 后进入 `verifying`，由 mainAgent 执行或委派 [integration-validation](../integration-validation/SKILL.md)。已具备前置的局部联调可以提前执行，最终验收仍核对完整确认范围。

按委派约定提供确认版本、R / C / V、有效决策、当前交付修订与固定依赖、Review/验证入口和必需环境/外部条件，并分配历史 attempt 和 runtime。mainAgent 核验该 Skill 返回的验收结果、材料完整性及实际资源处理；发现缺陷交回第 8 节修复。

协调者核对 summary 无 pending/claimed、非空任务图覆盖范围、原完成与当前修订的证据均真实、交付代码未遗漏、必需合同及环境验收成立、需审批决策有实际确认。已批准的检查例外保留其阶段范围和待补事项，不自动替代必需整体验收。证据不足不能由测试名称或工具成功返回补成通过。

对于协调器仓库的修改，至少执行：

```bash
GOCACHE=/private/tmp/agent-coordinator-go-cache go test ./... -count=1
GOCACHE=/private/tmp/agent-coordinator-go-cache go vet ./...
```

涉及协调器持久化行为时，使用虚构的 v1 控制文件验证新旧工作流兼容，并通过独立 CLI 进程验证状态恢复与并发认领；不修改活跃实例。

生成 `delivery-report.md`，记录解决的问题、需求与方案版本、各仓库路径/分支/提交、外部可见变化、Review、测试及整体验收结果、兼容影响、必要迁移/回退、已确认偏离、明确延期范围和后续发布事项。

证据不足或必需联调不可运行时，报告“代码阶段已交付，整体验收待完成”等准确状态，整体不进入 `completed`。延期范围必须有真实确认，不能为结束任务自行缩减验收。按产物管理分别核验自有服务停止、临时目录处理与必要产物保留，在现有交付入口说明遗留及处理条件，不额外生成目录清洁报告。

### 整体验收通过后的推送

整体验收通过并有对应范围的明确推送授权后，由 mainAgent 推送全部交付仓库；只读参考仓库或无需产物的 skipped 任务不因此产生无关分支。已有条件式授权直接沿用，用户明确暂停或限定本地范围时按其最新要求执行。

1. 从当前交付记录恢复逐仓 worktree、最新 `delivery_commit`、远端和目标分支，核对真实 Git 状态、完整 SHA、Review/授权与验收证据一致。需求分支 HEAD 必须等于当前已验收修订，工作树不能遗漏本次未提交改动；出现新提交或变化时先处理受影响评审/验证，不把浮动 HEAD 直接推送，也不自行补提交。
2. 核实各仓库对应远端与完整 `refs/heads/<需求分支>`。优先使用用户已指定或已有记录的目标；未指定时核查源仓库 origin 是否为实际交付远端。worktree origin 指向本地 clone/mirror 时，追溯已确认源仓库的实际远端，不将推到本地路径当作远端交付。来源或目标有歧义时只暂停受影响仓库，不猜测、初始化新远端或修改 Git 配置。
3. 检查实际远端分支 SHA；已等于目标修订时记录“已同步”并保留原成功结果。否则从该 worktree 普通推送固定提交，变量替换为已核实值：

```bash
git -C "$WORKTREE_PATH" push --porcelain "$PUSH_REMOTE" \
  "$DELIVERY_COMMIT:refs/heads/$BRANCH"
git -C "$WORKTREE_PATH" ls-remote --heads "$PUSH_REMOTE" "refs/heads/$BRANCH"
```

逐条检查退出码；远端查询返回的完整 SHA 必须等于被验收的 `DELIVERY_COMMIT`。非快进、权限、保护分支或网络失败保留实际结果，不 force push、跳过 hooks、改写历史或自行合并冲突。

4. 将逐仓目标、固定修订、推送与核验结果保存到 `delivery.json.pushes` 和现有交付报告/验证入口；必要输出按产物管理归档，不默认另建独立当前推送清单。全部仓库结果逐项维护，允许独立推送继续；部分成功不是回滚理由。
5. 中断或部分失败后先比对实际远端、Git 和已有记录，成功项不重复推送，只补未完成步骤。实际已经推送但记录失败时先核验远端，再补同一稳定事件。尚有可推进项保留 `verifying`，所有剩余项受阻时记录 `blocked`、`resume_status: verifying` 及解除条件，明确“验收已通过，推送未全部完成”。

全部必需验收及本次要求的逐仓推送核验完成后，更新报告和 `delivery.json.status` 为 `completed`，附各仓库远端分支及固定修订。推送不等于 PR、合并或部署，这些动作及正式 worktree 清理继续遵循各自实际授权；明确仅本地交付的范围按其确认完成条件收尾。
