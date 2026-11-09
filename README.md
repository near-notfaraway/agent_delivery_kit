# Agent Delivery Kit

从需求到代码交付的 AI Agent 协作工具包，支持跨仓库执行、人工评审与持久交接。

通过仓库自带的 Skills、模板和 `coordctl`，组织需求分析、技术设计、任务执行、逐项评审与集成验收，适用于多个独立仓库和 IDE 会话。仓库名称采用 `agent_delivery_kit`，CLI 名称保持 `coordctl`。

面向需要使用 AI Agent 在多个独立仓库和 IDE 会话中完成同一需求的技术负责人或高级工程师。跨仓库协作中的上下文传递、依赖确认、方案同步和验证汇总仍需要人工协调，本项目用持久记录支撑这些工作。

新的 Agent 无需历史聊天，也能从任务上下文和交付记录中识别当前需求、责任仓库与任务、前置依赖、上游接口和固定修订、验收与交接要求、有效方案及偏离决策，以及当前阻塞。

## 核心能力

通过持久化的本地控制面，组织独立 IDE / 仓库中的 Agent 协作，不依赖同一聊天会话或特定模型的内建 Agent Team。

工具包支持：

- 显式绑定需求、仓库、任务和 Agent。
- 原子任务认领，避免重复执行。
- 依赖门禁。
- 带验证证据的完成状态。
- 标准化仓库交接。
- 设计决策与代码偏离记录。
- 跨会话恢复状态。

## 差异定位

本项目面向独立仓库和 IDE 会话之间的交付协作，重点是：

1. 多个相互独立的 IDE 工作区，而非同一仓库内的并行 worktree。
2. 跨仓库接口、版本和验证证据的持久交接。
3. 人工评审参与下的方案变更与代码漂移治理。
4. 与模型、IDE和企业需求系统解耦的本地协议。

## 当前状态

项目已通过真实跨仓库需求验证，覆盖需求与方案确认、任务执行、逐项评审、交接和本地集成验证。现面向开源分享与复用，用持久化控制面减少跨仓库协作中的人工上下文传递。

当前包含以下能力：

- 文件锁保护的原子任务认领。
- 依赖门禁和可开始任务汇总。
- 必填 Git 修订、验证证据和交接说明。
- 面向全新 Agent 会话的任务与上游上下文读取。
- 设计偏离决策日志和有依据的任务跳过。

现有自动化回归覆盖任务认领、依赖门禁、完成证据、上下文、决策、跳过和状态汇总。独立 CLI 进程之间的状态恢复与并发认领尚未纳入自动化回归。

后续验证重点是不同使用者与全新环境中的复用，以及效率收益的量化评估。

## 运行要求与可选能力

本项目提供可下载的本地交付流程与控制工具。核心流程由仓库自带的六个 Skills、模板和 coordctl 组成，不要求安装仓库之外的 Skill，也不要求配置自定义 Agent。自定义 Agent 与外部 Skills 都是可选增强项。

| 使用范围 | 最低要求 |
|---|---|
| 构建和运行 coordctl | Go 1.23 或兼容工具链；支持 Unix `flock` 的系统和共享本地文件系统；控制目录可读写。当前未验证 Windows |
| 按流程交付需求 | 可阅读本仓库 Markdown 指令、读写授权文件并执行命令的编码 Agent，以及参与评审和确认的用户；不要求特定模型或 IDE |
| 准备和开发业务仓库 | Git、授权仓库及相应语言工具链、可获取的固定依赖和任务所需的隔离验证环境 |
| 多会话执行 | 原生子 Agent 或用户启动的独立 IDE / CLI 会话。没有并行能力时可由同一执行者按依赖顺序处理，仍遵守单仓唯一写入者和任务认领 |

使用本项目自带 Skills 时，各 Agent 环境按阶段直接读取 `skills/<名称>/SKILL.md` 及其引用材料；Skill 是指令文件，coordctl 不会自动执行它。

按实际需求选择场景能力，不要求每个使用者安装同一套工具：

| 必需能力及触发场景 | 推荐实现 | 替代路径及必须保留的结果 |
|---|---|---|
| 读取需求和参考资料 | 本地 Markdown；在线来源可使用对应平台的文档 Skill / API 工具 | 用户提供导出文件或完整正文，注明来源、版本及未读取的评论/附件；搜索摘要不能代替正文 |
| 行为变化的 TDD | 执行本仓库规定的先失败测试、再最小实现与验证；可选 `test-driven-development` | 编码 Agent 直接按相同步骤执行并保存失败/通过证据，无需外部 TDD Skill |
| 页面实现 | 按确认方案和仓库既有组件、样式开发；可选 `frontend-design` | 使用既有设计系统，仍需满足功能、交互和验收要求 |
| 实际代码评审 | 协调者或只读 Reviewer 按 `delivery-review` 审查；可选代码审查 Skill | 人工审查同一固定差异并记录问题、复核和审阅范围；注明真实 Reviewer，不能将自检称为独立审查 |
| 需求要求的浏览器验收 | 浏览器自动化工具，如 Playwright、agent-browser；可选配套浏览器 Skill | 人工在实际浏览器按相同 Case 执行，记录修订、环境、步骤、实际结果及必要截图/请求响应与数据证据 |
| 数据库、并发或目标环境验收 | 授权的隔离真实资源及其客户端/测试工具；可选对应平台 Skill | 使用另一种可证明相同语义的隔离真实环境；Mock 仅覆盖其模拟范围，不替代真实存储或目标环境验收 |

`lark-doc` 等平台 Skill 只适用于相应在线来源，`skill-creator` 仅辅助维护 Skill，均不属于普通交付的必装项。用户可按所用 Agent 环境的安装方式选择外部 Skills；本项目不自动安装或托管它们。

开始相应阶段前核查能力和运行条件，将选择的实现、替代路径及证据要求写入既有执行/验证入口。替换工具且验收范围与证据标准不变时无需重新批准流程；能力缺失或证据不完整则保留对应待验/阻塞，继续独立可执行工作。不能以缺少外部 Skill 为理由跳过必需验收，缩减范围仍须真实确认。替代证据遵循 [验证与证据约定](skills/delivery-flow/references/validation-and-evidence.md)。

以上定义了可移植执行约定；尚未在不配置自定义 Agent、完全不安装外部 Skills 的新环境中验证六个 Skills 的完整交付，CLI 的本地测试不能代替该验证。

## 下次需求如何使用

当一个新需求涉及多个仓库和多个独立 IDE 时，先在协调会话中说明：

> 用 Agent Delivery Kit 处理这个需求。原始需求文档：……。涉及仓库：……。先整理标准化需求，不确定的问题跟我确认；需求确认后编写技术评审方案，方案确认后再在我指定的目录建立 Agent 控制面，不要立即开发。

通用 mainAgent 按 `delivery-flow` 建立任务图、维护整体状态和组织确认；通用编码 subAgent 或独立仓库会话只负责分配的任务和仓库。

```text
通用 mainAgent → 读取 AGENTS.md，按 delivery-flow 协调
├── 需求/设计 → 直接执行阶段 Skill，或委派需求设计 subAgent
├── 仓库准备 → 直接执行 repository-preparation
├── 控制材料与调度 → mainAgent 维护
├── 单仓编码 → 通用 subAgent；没有子 Agent 时 mainAgent 串行执行
├── 单任务评审 → mainAgent 或评审 subAgent 执行 delivery-review
└── 集成验收 → mainAgent 或验收 subAgent 执行 integration-validation
```

整个准备过程分为三个阶段。用户未指定保存位置时，需求材料和状态默认放在本项目的 `workspaces/<需求名>/`，根 `.gitignore` 忽略整个 `/workspaces/`；已存在的需求沿用原位置。前两个阶段的文档优先让人容易理解和评审；第三个阶段在需求目录的 `control/` 中组织 Agent 执行材料。用户指定其他位置时沿用指定位置。原始需求不能直接当作执行计划。

本地流程分工如下，均为 Agent 阅读的指令文件，CLI 不会自动调用：

| Skill | 职责 |
|---|---|
| [delivery-flow](skills/delivery-flow/SKILL.md) | 恢复阶段、维护整体记录、调度、授权和完成交接 |
| [requirement-analysis](skills/requirement-analysis/SKILL.md) | 标准化需求、澄清和整稿确认 |
| [technical-design](skills/technical-design/SKILL.md) | 仓库事实核查、技术设计及方案评审 |
| [repository-preparation](skills/repository-preparation/SKILL.md) | 核实基线、创建或复用 worktree、运行准备 |
| [delivery-review](skills/delivery-review/SKILL.md) | 单任务实际差异审查、验证冻结、例外和依赖交接 |
| [integration-validation](skills/integration-validation/SKILL.md) | 当前固定修订的真实跨仓库合同与整体验收 |

`./skills/` 是本项目自带的六个 Skill 及后续新增项目 Skill 的唯一源目录和读取入口，均以真实目录和文件维护。外部及用户级（全局）Skills 按所用 Agent 环境的安装与发现方式管理。共享规则见 [验证与证据](skills/delivery-flow/references/validation-and-evidence.md) 和 [当前修订、恢复与返工](skills/delivery-flow/references/recovery-and-corrections.md)。

`delivery-flow` 在初始化时建立 `delivery.json`，持续记录整体阶段、确认和恢复信息。`requirement-analysis` 的单文件输出要求只约束需求分析职责，其输出为 `requirement.md`；整体状态文件由协调者维护。

### 交付材料与临时资源

材料按实际阶段生成。需求根目录维护需求、方案、整体状态和最终交付入口；控制文件在 control/，正式 worktree 在 repos/，评审与验证保持每个范围一个当前入口。历史证据默认放 `validation/attempts/<scope>/<attempt>/`，可重建脚本、源码副本、数据库、二进制和缓存放 `.runtime/<scope>/<attempt>/`；固定下游产物可集中保存在 validation/artifacts/。

创建前优先复用已有入口和工具，不为每轮生成另一套 final、before、freeze 或中间状态清单。mainAgent 分配位置与 writer，交接核查链接、固定修订、历史覆盖和资源实际处理；进程停止、目录删除和必要保留分别说明。需要保留的失败证据、固定产物及 worktree 支撑镜像不因整理目录被丢弃。

完整规则见 [产物管理](skills/delivery-flow/references/artifact-management.md)。这是后续执行的默认约定，不要求迁移或清理已有需求，也不新增 CLI 命令或 version 1 状态字段。规则由 Agent 执行，尚未通过新的完整交付验证其落实效果。

### Agent 角色与安装

角色与执行方法分别维护：`AGENTS.md` 保存稳定约束和关键门禁，Skills 保存具体流程，`agents/` 保存可选专业角色的职责、输入输出及权限范围。角色 prompt 引用 Skill，不复述阶段步骤；本文的使用示例不构成另一份流程规范。CLI 执行其已有的依赖、认领、证据必填和持久化检查，不自动落实文档确认或宿主权限。

协调使用通用 mainAgent，编码使用通用 subAgent，仓库准备由 mainAgent 直接执行对应 Skill，不要求安装这三类自定义 Agent。可按需安装的专业角色如下：

| 可选角色 | Prompt 源文件 | 使用的 Skill |
|---|---|---|
| 需求与技术方案 Agent | [requirement-design-agent.md](agents/requirement-design-agent.md) | `requirement-analysis` / `technical-design` |
| 交付评审 Agent | [delivery-review-agent.md](agents/delivery-review-agent.md) | `delivery-review` |
| 集成验收 Agent | [integration-validation-agent.md](agents/integration-validation-agent.md) | `integration-validation` |

用户未配置自定义 Agent、平台不支持配置或角色不可用时，流程直接使用通用 subAgent；没有子 Agent 能力时，mainAgent 按依赖串行执行同一套 Skills。无需先设置角色或更换平台。专业角色调用部分成功时先恢复实际状态、确认旧 writer 停止，再继续执行，避免重复认领或写入。

执行者选择、输入输出、writer、反馈及默认路径按 [Agent 委派约定](skills/delivery-flow/references/agent-delegation.md)，编码方法按 [单仓任务执行](skills/delivery-flow/references/task-execution.md)。执行方式变化不减免文档确认、任务认领、实际评审、验证或交接；mainAgent 编码后的自检不能标为独立审查。

本仓库使用 `agents/` 保存 prompt 源文件，`skills/` 保存项目 Skill 源文件，不维护 `.agents/` 或其他宿主发现目录。不同宿主的自定义 Agent 配置位置和格式按其当前文档确认；这些 Markdown 文件没有统一的自动导入协议，coordctl 不负责安装或启动 Agent。

仅在用户希望启用增强项、且宿主支持时安装或手动配置专业角色；这不是运行交付流程的准备步骤：

1. 在所用宿主支持的 Agent 配置中映射源文件的名称、调用描述和“Agent 指令”；核实工具及写入权限，模型沿用宿主默认。平台不支持或配置方式尚未确认时沿用通用执行路径，不等待配置；不把文件复制成功当作已注册。
2. 每次委派提供本项目绝对 `POC_DIR` 和适用指令入口，保证执行者可以读取项目 Skills 与材料。安装后的 prompt 不依赖安装目录的相对路径，也不假设跨仓库会话自动继承主会话历史或 AGENTS。
3. 更新源 prompt 后同步对应已安装配置；更新 Skills 后核对安装角色实际读取的项目来源仍正确。保留源版本或内容摘要，避免旧安装副本与当前流程矛盾。
4. 分别确认配置已保存、宿主已发现角色和代表任务的行为。只有实际启动并观察过才能报告行为已验证；未验证阶段如实保留。

subAgent 可在任务范围内使用其他 Skills，但进一步派生 Agent、修改整体状态或扩大权限须交回 mainAgent。编码者的自检不自动构成独立评审；角色名称和模型选择不代替实际证据。

### 1. 整理并确认标准化需求

采用 [requirement-analysis 自带模板](skills/requirement-analysis/assets/requirement_template.md) 作为标准模板，需求分析只输出一份 `requirement.md`，版本记录在文件内部。使用模板的元信息、适用章节和当前确认记录，删除填写提示及不适用内容，读取原始需求后整理：

- 背景、目标、用户场景和现有问题。
- 本次范围、排除项和约束。
- 带稳定编号（如 R1）的需求行为与验收条件，覆盖异常及边界场景。
- 实际提出的非功能要求。
- 待确认问题及其对范围、行为或验收的影响。

存在影响范围、业务行为或验收的不确定事项时，直接在对话中询问用户，取得结论后写回正文并保留确认依据；没有此类事项就交付这一份文档，不要求需求分析增加澄清问题、评审副本、接入报告或状态文件。正文表达当前有效规则，附录保留历史澄清及替代关系。技术核查与执行配置按所属阶段记录，只有影响业务选择或验收时才升级为需求确认。

需求澄清完成不自动代表整稿批准；记录实际确认版本、范围及依据后再进入技术方案阶段。会话中已有明确确认时沿用该结论。需求和方案确认保存在各自文档中，由协调者同步到初始化时建立的 `delivery.json`；恢复已有记录时不重建或重置。

### 2. 编写并评审技术方案

采用 [technical-design 自带标准模板](skills/technical-design/assets/design_template.md)，以已确认的需求版本为基线，核查涉及仓库的现有实现，输出 `design.md`。按模板组织内容：

- 背景与目标：本次问题、目标、范围与既有约束。
- 当前现状：当前链路、可复用能力、缺口及核查依据。
- 整体方案：推荐方案、核心变化、模块职责和关键流程。
- 模块改动：按模块或仓库说明数据与接口变化、核心处理和关键边界。
- 关键取舍与影响：选择依据、实质性备选、兼容影响及成立条件。
- 实施与验证：改造范围与真实依赖、具体验证重点、必要发布与回退。
- 评审事项：待决定或待验证事项、建议与影响、实际结论及整稿确认记录。

保留模板元信息中的需求链接及确认版本、作者、更新日期与方案版本，状态如实使用 `待评审 / 修订中 / 待确认 / 已确认`。先讲清方案和选择依据，正文保留影响评审结论的设计、接口边界和风险；详细接口、测试用例、SOP 与评审历史按需链接到附录。提交评审前删除填写提示、占位内容和不适用内容，历史方案可用于参考表达方式和图示，用户明确调整时按用户要求组织。

使用需求编号 R、组件改动编号 C 和测试用例编号 V 建立对应关系，映射较多时放入附录。Agent 的任务拆分、运行命令与状态放到下一阶段。第 7 节保留真实整稿确认版本、确认人、日期及依据；未确认时填写“待确认”，解决局部问题不等于整稿批准。取得方案确认后再生成控制目录。已有方案及执行记录沿用原位置、版本和真实确认，不因模板更新迁移。

### 3. 根据确认方案建立独立 Agent 控制目录

准备好已确认的需求、技术方案后，再生成执行材料。真实需求也默认放在被 Git 忽略的 `workspaces/<需求名>/`，用户指定位置时覆盖默认值；不覆盖其他需求的状态，不将私有材料提交到公开仓库。

开发所需的各仓库 worktree 放在需求目录的 `repos/` 下，默认从对应仓库的 `master` 创建需求分支。master 缺失或不适合作为基线时，暂停该仓库并询问用户应使用哪个分支；用户指定的分支也不存在时继续询问正确分支或远端，不自行替代。选定来源后先记录确切 ref、不可变 commit 及确认依据，再创建 worktree。保留已成功准备的仓库和已有工作，直到所有改动仓库的 worktree 创建并验证成功后才生成控制材料、启动开发。

按 repository-preparation 完成所有仓库后，在新的、未被其他需求使用的目录建立控制面，并构建本需求独立 CLI。以下路径替换为实际授权值，构建在项目根目录执行；先核查工具链和命令副作用。恢复既有控制目录时沿用其材料和二进制，不重放初始化或覆盖活跃实例。

```bash
POC_DIR=/absolute/path/to/agent_delivery_kit
DELIVERY_DIR="$POC_DIR/workspaces/<需求名>"
CONTROL_DIR="$DELIVERY_DIR/control"
COORD="$CONTROL_DIR/bin/coordctl"
REQUIREMENT_DOC="$DELIVERY_DIR/requirement.md"
DESIGN_DOC="$DELIVERY_DIR/design.md"
PROJECT="$CONTROL_DIR/project.json"
ls -ld "$DELIVERY_DIR" &&
test ! -e "$CONTROL_DIR" &&
mkdir -p "$CONTROL_DIR/bin" &&
cp "$REQUIREMENT_DOC" "$CONTROL_DIR/requirement.md" &&
cp "$DESIGN_DOC" "$CONTROL_DIR/design.md" &&
cp "$POC_DIR/templates/project.template.json" "$PROJECT" &&
GOCACHE=/private/tmp/agent-coordinator-go-cache go build -o "$COORD" ./cmd/coordctl
```

`POC_DIR` 沿用既有变量名，值为本项目实际绝对路径；已有目录仍可使用，无需因项目改名迁移需求材料或控制状态。Go 模块名更新为 `example.com/agent_delivery_kit`，CLI 命令和 version 1 控制文件格式保持不变。

文档在外部系统维护时，导出已确认版本作为快照。不要把空白需求、方案模板复制进去后才开始澄清，也不要把本文的示例任务图当作真实需求的初始计划。

每个需求的控制目录包含：

| 文件 | 用途 |
|---|---|
| `requirement.md` | 已确认的标准化需求完整快照，保留版本及 Review 记录 |
| `design.md` | 已确认的技术评审方案完整快照，保留版本及 Review 记录 |
| `execution.md` | 原文档位置、确认版本、任务映射、仓库边界、依赖原因及 Agent 操作指引 |
| `project.json` | 可执行任务图、任务状态、负责人、完成证据和交接说明 |
| `decisions.jsonl` | 方案偏离的追加日志，首次运行 `decide` 时创建 |

协调 Agent 根据方案和下述字段说明生成 `execution.md` 和 `project.json`；上述命令只准备快照、JSON 模板和 CLI。通用 JSON 模板保留 version 1，初始任务列表为空；将需求名称与任务填完整后才能交付仓库 Agent。为每个任务填写：

- `id`：稳定的任务编号。
- `repo`：被授权修改的唯一仓库。
- `title` 和 `objective`：任务名称与边界。
- `acceptance`：Agent 必须完成的验证。
- `handoff_requirements`：下游需要的交接证据。
- `depends_on`：前置任务编号。
- `status`：新任务从 `pending` 开始。

在 `execution.md` 中把任务 T 映射到方案中的 R / C / V 编号，解释每条依赖及下游交接要求。它的任务、仓库、依赖和验收方式必须与 `project.json` 一致。仓库使用授权的绝对路径，避免控制目录位置改变后指向错误的仓库。任务状态及首次完成 `result` 由 CLI 维护；整体阶段、真实批准、后续修复及当前交付修订由协调者维护于 `delivery.json`，执行说明引用其入口。

交付前核对任务图是否忠实覆盖确认方案，不能在生成执行计划时引入新的需求或架构选择。发现新问题时回到相应阶段确认。需求与方案原文档是维护入口，控制目录保存确认版本快照；基线变化时同步原文档、快照及受影响的计划，并记录确认依据。

当前文档生成、版本同步与 Review 检查由协调 Agent 按流程显式执行，CLI 不会自动生成、读取或校验 Markdown，也不强制 Review 门禁。现有控制目录没有 `execution.md` 时仍可按原有需求与方案记录使用，无需迁移。

### 4. Agent 读取任务上下文

每个仓库 Agent 开始前先阅读 `execution.md` 及同目录的需求、方案快照，确认基线和相关 Review 结论，再运行：

```bash
"$COORD" context \
  --project "$PROJECT" \
  --task T1
```

输出包含本任务目标、仓库边界、验收条件，以及上游任务的修订、验证结果和交接说明。新 Agent 会话应以这里的持久状态为准，不依赖历史聊天。

### 5. 认领任务

```bash
"$COORD" claim \
  --project "$PROJECT" \
  --task T1 \
  --owner api-agent
```

如果前置依赖尚未完成，认领会失败。多个 Agent 同时认领同一任务时，文件锁保证只有一个能够成功。

只有认领成功后，Agent 才开始修改自己的仓库。

### 6. 开发、验证并提交 Git 修订

Agent 按 `acceptance` 完成代码与测试，由协调者审查实际差异后，单独交给用户 Review 并确认是否允许提交该任务。材料就绪就进入该任务的评审，不等待其他独立任务全部结束；已明确覆盖该任务的授权继续有效。满足验收且获得提交授权后，由当前仓库 Agent 创建真实 Git commit。协调器 CLI 不自动修改代码或启动、托管 Agent。

等待一项任务 Review 时，其他独立 Agent 继续执行。单项检查失败或缺少授权只阻塞该任务和依赖它的任务，按实际范围记录；不能据此假设其他任务也需要暂停。

### 7. 提交完成证据和交接

```bash
"$COORD" complete \
  --project "$PROJECT" \
  --task T1 \
  --owner api-agent \
  --commit "<git-revision>" \
  --validation "go test ./... passed" \
  --handoff "Item v1 使用 id 和 name 字段，下游按该结构消费"
```

以下三项缺少任何一个都不能完成：

- 不可变 Git 修订。
- 最小验证结果。
- 下游可理解的交接说明。

每项任务 `complete` 成功后，协调者立即重新运行 `summary`，启动所有依赖均已满足的 ready 任务；下游 Agent 再运行 `context`，并读取 `delivery.json` 当前交付指针与有效批准决策。CLI 的历史任务交接不自动包含后续修复。commit 本身不解锁依赖。

已完成任务需要修复时，默认复用原需求分支/worktree，先核对唯一写入者及活跃验证。原 `result.commit` 保留，修复单独评审、验证并按实际授权新增提交，由协调者维护 `repositories[].delivery_commit` 和 `corrections[]`。旧记录缺当前指针时兼容读取原任务提交；无任务图变化的小修复不需要改图或重新认领。

例如 T1 与 T4 无依赖并行开发，T2 仅依赖 T1：T1 先完成并获用户 Review/提交确认后，提交 T1 并执行 `complete`，T2 即可与仍在执行的 T4 并行。具备前置的局部联调可先行，最终整体验收在全部任务结束后执行。

整体验收通过后，默认由 mainAgent 将所有交付仓库的最新已验收提交普通 push 到对应远端需求分支，并逐仓核验远端 SHA。用户明确“验收通过后 push 所有 repo”可提前授权这个步骤，记录仓库/远端/分支后条件满足直接执行，无需再次确认。只要求分析、设计或明确本地交付时不推导推送授权。

推送前，最新分支 HEAD 必须与当前 delivery_commit 和验收证据一致；新改动先处理受影响评审/验证。远端指向本地镜像时核查源仓库实际远端，目标不明或非快进等失败不强推。结果记录于 delivery.json.pushes 和现有交付入口，部分成功保留，恢复只补未完成项；全部本次要求的推送核验通过后整体 completed。推送不代表已创建 PR、合并或部署，具体方法见 [delivery-flow](skills/delivery-flow/SKILL.md#整体验收通过后的推送)。

### 8. 记录设计变化

开发中发现原方案不合理时，不要只改代码。先保存决策记录：

```bash
"$COORD" decide \
  --project "$PROJECT" \
  --id D1 \
  --task T2 \
  --original "原计划直接引用上游仓库" \
  --problem "会导致两个仓库的发布周期耦合" \
  --choice "改用本地 DTO，并通过契约测试保护" \
  --needs-approval
```

决策会追加到控制目录中的 `decisions.jsonl`。当前 `--needs-approval` 只记录审批要求，尚未形成强制门禁；高风险、跨仓库或影响上线方式的变化仍需人工确认。

协调者根据实际确认更新原技术方案及变更记录，同步 `design.md` 当前及版本快照、`execution.md`；涉及需求范围或验收变化时，同时更新原需求及快照。保留决策编号与真实批准。确需修改活跃任务边界或依赖时，暂停受影响调度，另行形成安全维护方案；当前 CLI 无计划编辑命令，不直接手改活跃 `project.json` 或让用户手工修改。

### 9. 明确结束不需要的条件任务

```bash
"$COORD" skip \
  --project "$PROJECT" \
  --task T4 \
  --reason "现有公开类型已兼容新契约，无需修改"
```

条件任务不能长期停在 `pending`。不需要修改时必须给出可审计的判断依据。

### 10. 协调会话查看整体状态

```bash
"$COORD" summary --project "$PROJECT"
```

协调者重点检查：

- `ready` 中有哪些任务可以开始。
- 是否存在长期处于 `claimed` 的任务。
- 每个完成任务是否具有修订、验证和交接证据。
- `decisions.jsonl` 中是否有待人工确认的变化。
- 最终是否只剩 `completed` 或有依据的 `skipped`。

### 仓库 Agent 启动话术

可以把以下入口指令交给通用 subAgent 或不同 IDE 中的 Agent，并按 [委派约定](skills/delivery-flow/references/agent-delegation.md) 提供真实输入。编码步骤以 [单仓任务执行](skills/delivery-flow/references/task-execution.md) 为准：

```text
你负责单仓编码任务 T2，本次调用附有真实 POC_DIR、owner、唯一 worktree、
确认文档、控制/CLI 路径、授权范围、验证条件和证据入口。
先读取适用 AGENTS.md，
再读取 POC_DIR/skills/delivery-flow/references/agent-delegation.md
和 POC_DIR/skills/delivery-flow/references/task-execution.md，按其执行。
向 mainAgent 返回实际差异、验证、阻塞与交接结果；不承担整体调度。
```

全部任务结束后通过 integration-validation 使用当前交付修订验收真实链路。只编译、Mock、本地通过和目标环境验收分别记录；检查例外必须有实际范围批准。证据按一个当前入口链接历史尝试，保留关键失败和最终结果，清理仅限本次自有临时资源。

## 当前边界

- 只支持共享文件系统上的本地协调，尚未验证跨机器协作。
- 使用 Unix `flock`，尚未验证 Windows。
- 没有任务租约、心跳、超时回收和自动重试。
- `needs_approval` 尚未强制阻止任务继续执行。
- 不自动调用模型、不托管 Agent，也不替代 Git、CI 和代码评审。
- 不解决同一仓库内并行修改的所有文件冲突；执行仍须遵守唯一写入者约束。
- 第二个真实案例验证通过前，不扩展 Web UI、数据库或云服务。

## 模板与示例

- [标准需求模板](skills/requirement-analysis/assets/requirement_template.md)：生成 `requirement.md`。
- [标准技术方案模板](skills/technical-design/assets/design_template.md)：生成 `design.md`。
- [version 1 控制文件模板](templates/project.template.json)：初始任务列表为空，按确认方案填写任务。

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
