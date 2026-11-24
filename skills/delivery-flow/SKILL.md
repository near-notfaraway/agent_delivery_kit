---
name: delivery-flow
description: Coordinate analysis, design review, worktrees, and multi-agent delivery with coordctl. Use when users ask to develop a requirement or resume delivery. Not for isolated edits.
---

# 完整交付流程

由通用 mainAgent 承担协调职责，维护整体交付记录，按依赖组织执行者并核验交付证据。CLI 管理任务持久状态，Agent 的启动与运行由宿主提供。

## 1. 确认本次请求与执行边界

- 遵守宿主及实际读写位置的适用指令，核对目标仓库的 `AGENTS.md`（如有）。本文件是需求交付入口。
- 定位本 Skill 所属项目的绝对路径 `KIT_DIR`，确认 `cmd/coordctl`、`templates/` 和五个阶段 Skill 存在。Skill 是要实际读取的指令文件；缺失时先定向查找，仍找不到则记录受影响阻塞。
- 按第 2 节核验本次交付入口，再恢复当前阶段与有效确认。参考输入的收录和读取按 [来源材料管理](references/source-materials.md)；来源存档不推进阶段。
- 按阶段读取 [requirement-analysis](../requirement-analysis/SKILL.md)、[technical-design](../technical-design/SKILL.md)、[repository-preparation](../repository-preparation/SKILL.md)、[delivery-review](../delivery-review/SKILL.md) 和 [integration-validation](../integration-validation/SKILL.md)。
- 按 [Agent 委派约定](references/agent-delegation.md) 选择执行者并提供完整输入。mainAgent 直接准备仓库，编码默认委派通用 subAgent 执行 [单仓任务执行](references/task-execution.md)；无子 Agent 能力时按同一入口串行执行。
- 按 [产物管理](references/artifact-management.md) 分配当前入口、可写范围和必要运行位置；材料按阶段生成。中断与部分失败按 [恢复与返工](references/recovery-and-corrections.md) 核对真实状态，只补未完成步骤。
- 按 [README 的运行要求](../../README.md#运行要求) 和 [运行能力与替代路径](references/validation-and-evidence.md#运行能力与替代路径) 核查当前阶段能力。外部 Skills 和专业角色可选，验收范围与证据标准不随工具替换缩减。

用户要求完整交付时，方案确认后继续已授权的准备和开发，不因阶段 Skill 的职责终点重复请求相同授权。仅分析或设计不自动扩展为开发；Git commit、push、PR、合并、部署及对外发送分别核对实际授权。

推送是整体验收后的收尾动作，执行条件是对应仓库与目标已有明确授权。用户明确“验收通过后 push 所有 repo”即为条件式授权，记录目标与依据，条件满足后按第 9 节直接执行。未获推送授权时保留待确认动作及当前交付证据；明确仅本地交付时按该范围收尾，不由使用本 Skill 或模板更新补造授权。

## 2. 定位交付入口并判定新建或恢复

从用户输入和现有材料提取需求名、所有开发仓库共用的分支名，以及用户指定的保存位置。只询问影响当前阶段的必要信息，未明确的分支和执行配置可以留到执行准备阶段；用户未指定位置时直接使用默认目录，不额外请求位置确认。

定义：

- `WORKSPACES_ROOT`：默认 `$KIT_DIR/workspaces`；用户指定其他根目录时沿用。
- `DELIVERY_DIR`：默认 `$WORKSPACES_ROOT/<需求名>`；用户直接指定需求目录时使用该目录，不重复追加需求名。
- `CONTROL_DIR`：`$DELIVERY_DIR/control`，只在方案确认后创建。

真实需求默认保存在本项目的 `workspaces/`，确认根 `.gitignore` 包含 `/workspaces/`。不要强制暂存该目录或把其中的私有材料复制到公开跟踪文件。Git 可用时检查忽略规则及是否已有被跟踪内容；忽略规则不会移除已有跟踪。恢复已有需求时沿用记录的位置，不因默认值变化自动迁移。

需求名必须是单一目录名，拒绝空值、`.`、`..`、路径分隔符及目录逃逸；需要改名时先对齐实际名称，不静默替换。使用 `git check-ref-format --branch` 校验分支名。所有路径、变量都要正确引用。

进入阶段前，先检查指定或默认 `DELIVERY_DIR` 是否存在，再读取该入口的 `delivery.json`、人工评审文档及真实确认依据。记忆、外部同名目录、业务仓库分支及资料里的任务图仅用于定向定位，不能替代当前交付核验：

| 实际情况 | 处理与阶段依据 |
|---|---|
| 当前需求目录不存在，也没有已核实的同一交付记录；用户只提供需求描述与参考文档 | 新建需求，从 `init` → `analyzing_requirement` 开始；旧需求或预调研可提取信息，不能直接进入设计、仓库准备或编码 |
| 目录存在，但只有 sources、参考文件或空目录 | 不视为需求已经整理或确认；保留已有文件，在当前入口建立需求稿及初始化记录 |
| 当前入口有匹配的 delivery.json、文档和确认依据 | 恢复同一交付；逐项核对需求身份、实际版本、确认范围、控制材料和授权，再继续首个未完成阶段 |
| 当前入口有本次需求/方案正文，但缺少 delivery.json 或确认依据 | 核实确为当前交付材料后按第 3 节兼容接入；没有真实确认时恢复到首个未确认阶段，不能因文件存在而跳过评审 |
| 只在其他目录发现同名历史需求、已评审方案或控制文件，用户未指定续接 | 仅作参考，不将其目录、版本、批准、owner、result 或任务图导入本次交付 |
| 用户明确续接指定旧交付，或已有本次交付记录明确指向旧位置 | 沿原入口恢复，不因默认目录不存在而重建；找不到或无法核验关键材料时保留对应缺口，不新建目录伪装恢复成功 |
| 用户明确将某份旧文档作为本次基线 | 核对确切原文及版本，并记录本次采用的范围；只有明确覆盖该整稿的真实确认才可沿用批准，不自动采用旧任务状态或发布授权 |

“开发这个需求，之前有需求信息和预调研方案”表示开发意图及参考输入，不等于确认了当前需求和方案。既有批准仅在同一交付、同一版本和实际范围内恢复；用户明确采用并确认旧整稿时，不再重复询问同一批准。整理、重写或补全后发生内容变化，按所属阶段 Skill 记录受影响的确认。

写入前检查父目录和目标位置，冲突目录不得覆盖。初始化或恢复整体记录前，读取 [整体交付记录](references/delivery-records.md)。确认为新建后，从 [delivery.template.json](assets/delivery.template.json) 初始化 `delivery.json`，填写真实名称、绝对路径、已明确的分支、时间和授权依据，初始 `status` 为 `init`。新需求的文档版本按各模板起步，不沿用参考资料的版本号；尚未明确但不影响需求分析的执行配置留到相应阶段核查。

材料布局见整体交付记录，按阶段和实际需要创建；`requirement.md` / `design.md` 保持人工评审入口，执行快照和任务图在仓库全部就绪后生成。恢复沿用原位置，来源、证据与临时资源分别按第 1 节引用的共用约定处理。

## 3. 管理状态与确认依据

状态转移、批准登记及记录恢复前，按 [整体交付记录](references/delivery-records.md) 核对进入条件、文档状态映射、当前与已确认版本、字段和写入方式。整体阶段与真实确认由 mainAgent 唯一维护，`project.json` 保持 version 1，任务状态只通过 CLI 修改。

`delivery.json` 写入仍须使用文件系统锁、同目录临时文件刷新后原子替换并保留未知字段；交接协调权先确认旧 writer 停止。只记录真实确认和阶段转移，局部答复、Agent 返回或文件存在不等于整稿获批，参考资料不自动成为当前批准。

需求和方案评审等待使用对应 awaiting 状态；任务 Review 等待保持 `claimed` 和原 owner。局部阻塞不停止独立工作，当前阶段所有工作受阻时才进入 `blocked` 并保存续接条件。恢复核验当前文档、Git、CLI 与证据，保留成功副作用和有效授权；基线变化显式同步受影响确认和历史，不静默回退或重建旧交付。

## 4. 需求分析与技术评审

1. 新需求先进入 `analyzing_requirement`，按来源材料管理登记相关输入并收录必要原文，再按委派约定提供来源身份、实际读取范围、名称、当前文档位置及本次已有结论，由 mainAgent 或需求设计 subAgent 执行 `requirement-analysis`。材料无法读取时登记缺口，继续可独立整理的部分。
2. 核验返回的 `requirement.md` 及真实状态。mainAgent 组织需要的澄清和整稿 Review，反馈由文档 writer 按需求 Skill 写回；稿件就绪进入 `awaiting_requirement_review`，实际整稿确认后同步 `delivery.json` 并进入 `requirement_confirmed`。
3. 进入 `designing`，提供确认需求及仓库只读入口，由 mainAgent 或需求设计 subAgent 执行 `technical-design`。该 Skill 负责模板、仓库核查、模块契约与验证计划，不提前创建 worktree 或控制目录。
4. 核验返回的 `design.md` 和未决事项，并按 [技术方案写作与交付检查](../technical-design/references/document-writing.md) 检查评审入口、必要核查和真实状态；“整理中”或“修订中”的稿件继续在 `designing` 完成，交审或交确认后才进入 `awaiting_design_review`。mainAgent 展示具体版本并组织反馈，按所属 Skill 同步受影响文档。整稿确认后同步 `delivery.json` 并进入 `design_confirmed`，继续授权范围内的准备。

需求和方案服务人工评审，Agent 命令和执行任务放入后续控制目录。会话已有确认时保留授权直接推进，不能要求用户重复确认。

## 5. 调用仓库准备

先核对需求和方案都存在对应版本的真实确认，且 `delivery.json` 已同步到 `design_confirmed`。进入 `preparing_repos`，从确认方案提取仓库清单，区分需要修改、只读参考和外部团队负责的仓库。只为授权修改的仓库创建 worktree。

mainAgent 直接执行 [repository-preparation](../repository-preparation/SKILL.md)，使用确认版本、需求目录/分支、逐仓授权和已有准备记录，形成准备证据入口，再以协调者身份同步整体记录。

逐仓来源、基线和准备事实按该 Skill 保存。未解决仓库保留局部结果及成功 worktree，整体状态按第 3 节处理；所有授权修改仓库均创建并核验成功后进入 `repos_ready`，才生成控制材料。执行环境缺口仍按影响范围限制后续任务，不将 worktree 就绪写成测试或验收已通过。

后续开发和需求内修复默认使用原需求 worktree，同一 worktree 最多一个写入者，不覆盖活跃源码或二进制。具体复用检查和基线纠正按仓库准备 Skill，不由协调者另设一套准备流程。

## 6. 构造 coordctl 控制文件

进入 `preparing_control`，创建 `CONTROL_DIR`：

1. 将已确认需求和方案复制为完整快照，并在 `baselines/` 保存不可覆盖基线，记录版本、确认依据及内容摘要。复制去重按产物管理，当前快照与评审入口保持一致；快照内相对引用的原解析位置按来源材料管理记录，不能因复制而假设链接仍可导航。
2. 根据确认方案生成 `execution.md`，写清绝对路径、确认版本、需求验收/模块/验证用例与任务的对应关系、依赖原因、仓库边界、验证命令和交接要求。按 [单仓任务执行](references/task-execution.md) 为各任务划分初始里程碑，确定目标、验收边界及检查时点。使用实际名称或已有编号，不要求人工评审稿先提供编号映射表。
3. 从 `templates/project.template.json` 生成 `project.json`，保留 `version: 1`。根据实际方案生成任务，不复制 README 的示例任务图。
4. 每项任务填写唯一 `id`、唯一授权 worktree 的绝对路径 `repo`、`title`、`objective`、`acceptance`、`handoff_requirements`、`depends_on` 和初始 `status: pending`。一个仓库只有一个任务时，`id` 直接使用仓库别名，`title` 以仓库名开头；同一仓库确有多个任务时，使用 `<仓库别名>-<范围>`，并在标题中保留仓库名。不要只用 `contract`、`core`、`forwarding` 等抽象能力名隐藏任务对应的仓库。不添加整体阶段、审批或 Agent 调度字段。
5. 检查任务 ID 唯一、引用存在、无自依赖或环、仓库授权正确，且方案范围被完整覆盖。同一 worktree 的多个写任务用依赖或协调调度串行化。
6. 验收须包括符合设计、必要代码 Review、仓库要求的检查与相关测试；交接须包括接口/数据契约、不可变修订、可消费的固定依赖/产物及解析验证、兼容及下游使用方式，按 [验证与证据约定](references/validation-and-evidence.md)。
7. 在 `KIT_DIR` 为本需求构建独立 CLI，不覆盖活跃二进制。使用 Go 默认缓存或本会话明确配置的 `GOCACHE`，沿用已有资源记录，并将实际缓存位置和配置提供给后续执行者：

```bash
go build -o "$CONTROL_DIR/bin/coordctl" ./cmd/coordctl
```

设置 `COORD="$CONTROL_DIR/bin/coordctl"`、`PROJECT="$CONTROL_DIR/project.json"`。运行 `summary`，逐任务运行 `context`，核对输出和预期任务图一致。CLI 不负责校验所有图结构、文档确认、Git 修订真实性或证据质量，以上检查由协调者完成。

命令、字段和退出码查询 [CLI 文档](../../docs/cli.md)，阶段执行与门禁以本 Skill 及相关执行约定为准。

控制文件初始化后，任务状态通过 CLI 修改。不要将真实活跃 `project.json` 或 `decisions.jsonl` 当成开发/测试夹具。检查通过后进入 `control_ready`。

准备脚本只实现本次初始化时放本次 runtime；结果写入既有准备/控制验证入口，不要求每个需求另建初始化报告、准备 JSON 或记录脚本。状态写入辅助必须符合恢复约定，不把一次性脚本作为可盲目重跑的标准工具。

## 7. 调度多个仓库 Agent

使用当前环境实际提供的子 Agent 能力，为可并行的独立仓库分派通用编码 subAgent；执行依赖来自 `summary.ready` 与真实方案。CLI 只管理持久状态，不会启动或托管 Agent。

进入 `developing`。按 [委派约定](references/agent-delegation.md) 提供完整输入和唯一 owner，指定读取 [单仓任务执行](references/task-execution.md)。在 `delivery.json.agents` 记录实际会话与 owner 的对应关系，不伪造宿主会话 ID。

以单任务为单位推进交付，任务内按 [单仓任务执行](references/task-execution.md) 逐里程碑开发、Review 和续接。委派前提供初始划分及本轮里程碑，执行者只读核查后提出细化建议，由 mainAgent 核定并写回执行入口；不增加独立的里程碑方案审批。独立任务可以并行启动；任一里程碑达到可审阅条件后立即处理，不等待同批其他任务全部结束，也不将多个仓库的提交确认捆绑为默认门槛。用户对本批任务另有明确评审或授权安排时沿用其安排。

执行者不可用时按委派约定选择已有通用执行者；没有子 Agent 能力时，mainAgent 按依赖串行处理。用户已安排独立 IDE / CLI 会话时可提供相同任务输入，不将配置 Agent 或启动额外会话作为串行交付的前置。缺少角色不构成调度阻塞；实际缺失的执行能力或验收条件仍按影响范围记录。

编码者按单仓任务执行入口读取 context、认领/恢复、实现、验证和返回材料。mainAgent 核验其真实执行状态，不因 Agent 已返回就当作任务完成，也不在同一 worktree 并行派发 writer。

每个里程碑就绪后按下列第 1–3 步完成审查、用户 Review 和真实确认；通过后由 mainAgent 放行下一里程碑，尚未通过时只处理本轮修复和补验。已有明确覆盖本轮的评审与授权安排直接沿用。里程碑通过不等于提交授权，也不执行 `complete` 或解锁下游；全部里程碑通过并满足整项任务最终门禁后，才执行第 4–5 步。

1. mainAgent 执行或委派 [delivery-review](../delivery-review/SKILL.md)，提供任务及本轮里程碑范围、审阅基准、实际增量与累计 diff、确认版本、有效决策和验证材料，注明是中间里程碑还是整项任务最终评审。编码者按反馈修复，Reviewer 复核并返回该 Skill 规定的证据；核查自检与独立审查的真实边界，并按产物管理核验当前入口、历史、可读取证据和资源处理。
2. 材料就绪就按 [delivery-review](../delivery-review/SKILL.md) 第 4 节的用户阅读顺序交审：先展示本轮改动、解决的问题和主要代码，再简述验证及影响评审决定的限制，附本次范围的 Review diff 与评审入口；不要只发送路径、哈希或通过请求。缺少本轮真实确认时请求 Review，需要提交且缺少覆盖范围的授权时再请求提交确认，不等待独立任务。沿用已明确覆盖的评审与授权安排，流程安排、方案确认或另一任务批准不能替代本轮所需依据。
3. 将真实批准按稳定 ID、实际范围、冻结摘要及依据记录到 Review 和 `delivery.json`。检查失败默认阻塞该项验收，仅有真实批准的具体验收例外时按范围交接，原失败及补验条件保留。
4. mainAgent 核对全部里程碑的有效通过依据、整项任务累计 diff、最终检查与可消费交接，确认验收和有效授权满足后，当前 writer 按单仓任务执行约定创建新增提交、核验冻结内容并执行 `complete`，返回真实修订和 CLI 结果；已有获授权的里程碑提交时保留并核验，不重复提交。
5. mainAgent 核对实际 Git/CLI 与最终 Review、验证及可消费 handoff 一致，核对新增材料及遗留资源，再更新当前 `delivery_commit` 和证据入口。中途失败按恢复约定补齐缺失步骤，不重放已成功动作；实际 complete 已成功时不因补证据问题重复 complete。

每次收到里程碑、单任务完成报告或用户 Review 反馈，先处理该任务的可执行后续步骤。向用户请求 Review 时暂停该任务的下一里程碑及未获授权的提交，其依赖仍不解锁；不停止其他独立 Agent，也不等待它们结束，继续推进其他已授权的 ready 任务。

有任务或里程碑可 Review 且需要用户反馈或提交授权时，mainAgent 核验并展示该范围的具体材料，保存当前交付状态、证据入口、持久 owner、待确认范围和下一里程碑的续接信息，然后结束本轮回复以释放交互。不得为了等待其他任务而等待整批完成，也不持续阻塞在等待调用或轮询中。任务保持 `claimed`，仍有独立工作可推进时整体保持 `developing`；释放交互不代表停止其他执行者或结束交付。收到授权或实际完成通知后，先核对持久状态、最终差异和原 writer，再推进对应任务；已有覆盖范围的授权则直接继续，不为释放交互重复请求批准。其他任务能否在本轮回复结束后独立运行，以宿主实际能力为准，不假定支持。

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

流程中仅对仍为 `pending`、未执行且判断依据充分的条件任务使用 `skip`，不能借此结束失败、遗漏或仍有人执行的任务。CLI 将 skipped 依赖视为满足，必须确认下游不依赖被跳过任务的实际产物；本流程约束不代表 CLI 已强制全部条件。

已完成任务的修复按 [恢复与返工约定](references/recovery-and-corrections.md) 执行，保留原任务 `result`，更新当前交付修订并核查下游影响。任务边界/依赖未变时不启动计划维护；问题或必需复验未解决前不得标记整体完成。

## 9. 整体验收与交付

所有任务已成为 `completed` 或有充分依据的 `skipped` 后进入 `verifying`，由 mainAgent 执行或委派 [integration-validation](../integration-validation/SKILL.md)。已具备前置的局部联调可以提前执行，最终验收仍核对完整确认范围。

按委派约定提供确认版本、对应需求验收/模块/验证用例（名称或已有编号）、有效决策、当前交付修订与固定依赖、Review/验证入口和必需环境/外部条件，并分配历史 attempt 和 runtime。mainAgent 核验该 Skill 返回的验收结果、材料完整性及实际资源处理；发现缺陷交回第 8 节修复。

协调者核对 summary 无 pending/claimed、非空任务图覆盖范围、原完成与当前修订的证据均真实、交付代码未遗漏、必需合同及环境验收成立、需审批决策有实际确认。已批准的检查例外保留其阶段范围和待补事项，不自动替代必需整体验收。证据不足不能由测试名称或工具成功返回补成通过。

各业务仓库执行其适用检查和已确认验收；不将活跃控制文件作为测试夹具。

生成 `delivery-report.md`，记录解决的问题、需求与方案版本、各仓库路径/分支/提交、外部可见变化、Review、测试及整体验收结果、兼容影响、必要迁移/回退、已确认偏离、明确延期范围和后续发布事项。

证据不足或必需联调不可运行时，报告“代码阶段已交付，整体验收待完成”等准确状态，整体不进入 `completed`。延期范围必须有真实确认，不能为结束任务自行缩减验收。按产物管理分别核验自有服务停止、临时目录处理与必要产物保留，在现有交付入口说明遗留及处理条件，不额外生成目录清洁报告。

### 整体验收通过后的推送

整体验收通过且有对应仓库、目标和范围的明确推送授权时，mainAgent 在执行或恢复推送前读取 [逐仓推送步骤](references/push-delivery.md)，普通推送固定 `delivery_commit` 并核验远端 SHA。沿用已有条件式授权，保留逐仓成功与失败，不 force push 或推导 PR、合并、部署授权。

全部必需验收和本次要求的推送核验完成后，更新报告和整体状态为 `completed`。部分推送失败按真实可推进范围保留 `verifying` 或 `blocked`；明确仅本地交付时按其确认范围收尾。
