# 整体交付记录

初始化、恢复或更新 `delivery.json` 前读取本文件。它约定整体状态、确认字段和材料布局，不改变 CLI 的 version 1 任务状态；阶段推进按 delivery-flow，当前交付修订及事件去重按恢复与返工约定。

## 状态与确认依据

从初始化开始，整体阶段写入 `delivery.json`。文档状态及确认记录遵循需求分析/技术设计 Skill 和各自模板，协调者同步真实确认版本、范围与依据。恢复旧文档时沿用已有真实记录，不因模板更新重建状态。

整体状态只由协调者维护，`requirement-analysis` 和 `technical-design` 不创建或修改 `delivery.json`。不能将整体状态混入 `project.json`；后者保持 version 1，任务状态仍为 `pending / claimed / completed / skipped`。

| 整体状态 | 进入条件及下一步 |
|---|---|
| `init` | 名称、目录及已明确的执行配置已记录，初始化状态已保存 |
| `analyzing_requirement` | 开始执行需求分析 Skill |
| `awaiting_requirement_review` | 可审阅需求稿已完成，等待用户澄清或确认具体版本 |
| `requirement_confirmed` | 需求关键问题已解决，存在对应版本的真实确认 |
| `designing` | 基于确认需求执行技术设计 Skill |
| `awaiting_design_review` | 技术评审稿已完成，等待用户 Review |
| `design_confirmed` | 技术评审问题已解决，方案版本得到真实确认 |
| `preparing_repos` | 开始核对开发基线并准备各仓库 worktree |
| `repos_ready` | 所有改动仓库的来源、路径、分支、基线及写入占用核验通过；执行环境缺口已单独记录 |
| `preparing_control` | 开始生成快照、执行说明及 version 1 任务图 |
| `control_ready` | 文档和任务图一致、CLI 检查通过、可分发任务 |
| `developing` | 仓库 Agent 按依赖认领、开发、评审和验证 |
| `verifying` | 所有任务已结束，进行跨仓库整体验收与证据复核 |
| `completed` | 全部完成条件满足，最终报告已落盘 |
| `blocked` | 环境、权限、提交授权或执行故障阻止继续推进 |

正常按表顺序推进，需求/方案评审等待使用对应 awaiting 状态，不冒充已确认或执行失败。开发阶段的单任务 Review 或提交等待保留 `claimed`，在该任务的评审入口和 `delivery.json` 中记录，不新增 CLI 任务状态或退回文档评审阶段。

局部阻塞写入 `blockers`，其他工作继续：准备阶段有仓库可推进时保持 `preparing_repos`，开发阶段有执行者或 ready 任务可推进时保持 `developing`，验收/推送阶段仍有可执行项时保持 `verifying`。所有当前阶段可推进工作均受阻时才进入 `blocked`，保存 `resume_status`、影响范围和解除条件；需求/方案的正常整稿评审等待仍使用 awaiting 状态。恢复只记录真实接入和后续转移，不补造早期历史。

每次转移同时保存 `updated_at` 和 `history` 条目，包含前后状态、时间、原因与证据路径。确认记录必须包含文档版本、Reviewer、时间和实际确认依据；草稿生成、局部问题回答、任务完成均不能替代整份文档确认。

## 文档状态与版本

`documents.requirement` 和 `documents.design` 分别对应当前人工评审入口。文档 writer 按所属 Skill 维护正文状态、版本和确认依据；协调者核验后同步以下字段，不由子 Agent 修改整体记录：

- `path`：当前文档入口；模板中的相对路径从需求目录解析。
- `version`：当前文档的实际版本，尚未形成文档时为 `null`；不从参考资料或整体阶段推断。
- `review_status`：当前整稿状态，使用下表的机器值；不是某个评审事项、任务或验证结果的状态。
- `confirmed_version`：最近一次有真实整稿确认依据的文档版本；没有已核实的整稿确认时为 `null`。修订时保留该历史版本，当前是否获批须结合 `version`、`review_status`、实际确认范围和依据判断。

新记录中的 `review_status` 按下表与正文状态对应。状态进入条件及反馈处理以所属 Skill 为准：

| 适用文档 | 正文状态 | `review_status` |
|---|---|---|
| 需求、方案 | 尚未开始，无当前文档 | `not_started` |
| 需求 | 待澄清 | `clarifying` |
| 需求、方案 | 待确认 | `awaiting_confirmation` |
| 需求、方案 | 已确认 | `confirmed` |
| 方案 | 整理中 | `drafting` |
| 方案 | 待评审 | `awaiting_review` |
| 方案 | 修订中 | `revising` |

`not_started` 只用于尚未开始的文档；当前稿已存在时不能继续以模板初始值代表真实状态，各文档只使用上表适用的状态。

当前整稿确认成立时，`review_status` 为 `confirmed`，`version` 与 `confirmed_version` 相同，且有覆盖当前版本及范围的真实确认依据；字段相等本身不能证明批准。局部确认按实际范围保存在正文和 `approvals`，不将整稿标记为 `confirmed`。需求或方案发生实质修订时，更新当前版本及所属文档状态，保留历史批准和 `confirmed_version`，明确受影响范围；新版本获得覆盖整稿的真实确认后才更新 `confirmed_version`。等价的表达或布局调整按所属 Skill 沿用有效确认，不因编辑动作自动失效，也不把旧批准静默改成新版本批准。

整体 `status` 按当前实际工作维护，不由 `review_status` 一对一推导：

- 正在整理或澄清需求使用 `analyzing_requirement`；可审阅稿已交给用户澄清或确认、正在等待反馈时使用 `awaiting_requirement_review`，正文仍可能是“待澄清”或“待确认”。
- 正在形成方案或处理实质修订使用 `designing`；交审或交确认后等待用户反馈时使用 `awaiting_design_review`，方案状态为“待评审”或“待确认”。影响交审的必要事实未核查完成时不能进入该等待阶段。
- 当前需求整稿确认成立后才能进入 `requirement_confirmed` 并开始设计；当前需求及方案的有效确认都成立后才能进入 `design_confirmed` 并准备仓库。后续准备、开发、验收或局部阻塞不改变仍有效的文档确认状态。

需求修订同时核查方案所依据的需求版本和受影响契约，不能以方案仍标记“已确认”为由跨过新的需求确认；必要时按恢复约定暂停受影响工作并显式退回所属阶段，独立已授权工作继续。`delivery.json` 中的文档字段、整体阶段及批准更新由协调者在同次同步中核对一致，不因单个字段缺失就补造批准或阶段历史。

旧记录缺少字段或采用其他状态值时，保留未知字段和已有历史，不要求迁移或批量替换。恢复时读取正文及真实证据核验实际状态，正常维护时只补已核实的当前信息；不能仅凭缺失、未知或相互矛盾的值推定已确认。模板更新不重写活跃状态，也不改变 CLI 的 version 1 协议。

## 字段与写入

模板中的 `null` 表示尚未确认或尚未生成，不能当成已授权。时间采用带时区的 ISO 8601 格式，记录结构保持一致：

- `kit_dir`：工具包根目录的绝对路径，对应 `KIT_DIR`。新记录只写 `kit_dir`。恢复旧记录时，若 `kit_dir` 缺失或为空，则读取旧字段 `poc_dir`，并按入口规则核验目录。两个字段均有值时以 `kit_dir` 为准；核验失败按现有恢复规则处理。保留旧字段，无需仅因改名重写已有记录。
- `history[]`：`from`、`to`、`at`、`reason`、`evidence`；初始记录的 `from` 为 null。
- `approvals[]`：稳定 `id`、`stage`、`version`、`reviewer`、`at`、`evidence`，决策确认另含 `decision_id`。任务 Review/提交确认另注明 `task_ids`，修复确认含 `correction_id`，在证据中记录审阅的 commit 或工作树摘要；批准一项任务不扩展到其他任务。旧记录无 ID 时兼容恢复，不重复追加同一批准。
- `authorization` 的非空项：`scope`、`at`、`evidence`；只记录明确授权的动作和范围。
- `sources[]`：协调者登记本次有效输入。新记录用来源 `id` 和来源索引入口关联，正文、版本、读取范围及副本信息集中在 `sources/index.md`；已有 url/path 等记录兼容保留。具体收录规则见来源材料管理，不把来源字段当作批准或阶段状态。
- `blockers[]`：`id`、`reason`、`task_ids`、`unblock_condition`、`resolved_at`。
- `agents[]`：`session_id`、`owner`、`task_ids`、`repo`；没有真实会话时不填造假的 ID。
- `repositories[]`：保存逐仓来源、准备结果和当前交付修订；`delivery_commit`、证据指针及 `corrections[]` 的字段与恢复规则统一见 [恢复约定](recovery-and-corrections.md)。
- `pushes[]`：逐仓普通推送及远端核验记录，包含稳定 `id`、`repo`、`remote`、`ref`、固定 `commit`、实际 `remote_commit`、`result`、`verified_at`、`evidence`；失败或未核验字段如实保留，远端地址与证据不包含凭证。缺少该可选字段的旧实例仍可恢复。

协调者独占写入 `delivery.json`，使用文件系统锁、同目录临时文件完整写入和刷新后原子替换，保留未知字段；这些是协调者的写入要求，CLI 不负责该文件。交接协调权前确认旧协调者停止写入，不能只依赖记录中的名字推断互斥。批准、修复和重试去重按恢复约定，history 只记录真实阶段转移。

恢复时先通过 delivery-flow 的交付入口判定核验身份，再检查记录、文件、worktree、CLI 状态及证据是否一致，修复中断处再继续。只有已核实属于本次交付的需求或方案而没有 `delivery.json` 时，才沿用原入口兼容接入，由协调者初始化状态并据真实记录恢复到首个未获确认的阶段。参考文档不适用这条恢复路径。基线变更可以显式退回需求或方案阶段，保留历史及受影响确认记录，不静默回退或跳过门禁。

## 材料布局

`delivery.json` 是 delivery-flow 管理的整体状态文件，从初始化开始维护。`requirement-analysis` 的人工评审产物为 `requirement.md`，来源收录按分工执行，版本、澄清和确认都在文档内；由协调者同步整体状态。技术设计阶段再生成 `design.md`，方案确认后才生成执行控制材料。下图是完整交付阶段的可能目录，除初始化状态外的材料按阶段生成：

```text
workspaces/<需求名>/
├── delivery.json          # 初始化时创建，保存整体阶段与恢复信息
├── requirement.md         # 人工评审入口
├── design.md              # 人工评审入口
├── sources/               # 来源索引与确有必要的原文/导出，按来源材料管理创建
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

相关文档、日志和报告集中在需求目录。来源登记、必要正文/附件保留及版本引用按 [来源材料管理](source-materials.md)；稳定可取的外部原文可保留原位置并引用。不复制凭证、无关数据或整个私人会话。

图示为新执行的默认布局，目录按需创建。已有需求沿用其入口和证据位置；一次性脚本、中间载荷和下载缓存不堆放根目录，非标准材料按产物管理约定说明用途。
