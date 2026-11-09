# 单仓编码任务执行

供通用编码 subAgent、独立仓库会话及串行执行者使用。协调者按 `agent-delegation.md` 提供调用输入；本文件负责编码执行步骤，不要求自定义编码 Agent，也不授予整体协调权。

没有子 Agent 能力时，mainAgent 可按本文件兼任编码执行者，仍以真实 owner 逐任务认领/恢复和交接；整体协调继续按 delivery-flow 执行，评审遵循委派约定的自检与独立审查边界。不伪造子会话，也不因串行执行绕过任务状态或授权。

## 1. 恢复任务与固定输入

读取本项目及目标仓库适用指令、委派范围、确认需求/方案和 `control/execution.md`。旧控制目录缺少执行说明时沿用已有基线和记录，不为恢复重建控制面。

核对唯一授权 worktree、基线、实际分支/HEAD、未提交修改及 writer。运行：

```bash
"$COORD" context --project "$PROJECT" --task "$TASK_ID"
```

随后读取 `delivery.json` 当前交付修订、corrections 和有效批准决策；历史 handoff 不自动包含后续修复。按固定版本核对上游依赖及其获取/解析证据，不从浮动 HEAD 推断交付。

以上变量均来自真实调用输入，路径使用绝对路径并正确引用。检查每条命令的退出码；不发明 CLI 命令或绕过依赖。

## 2. 确定新任务或续接

- `pending` 新任务先执行下述 claim，只有成功后才能修改业务代码。
- `claimed` 恢复任务先核对持久 owner 和旧 writer 已停止，续接同一 owner；不再次 claim，也不另起 owner 强抢。
- `completed` 任务的修复按 `recovery-and-corrections.md` 处理，默认复用原需求 worktree，保留原 `result`，不重新 claim/complete。
- 非本次范围、依赖未满足、来源或 writer 无法确认时返回具体阻塞，不自行换仓库、分支或任务。

```bash
"$COORD" claim --project "$PROJECT" --task "$TASK_ID" --owner "$OWNER"
```

CLI 没有租约、超时回收或重分配命令。同 owner 重复 claim 也会失败；失败或工具中断先核对实际持久状态，再只补未完成步骤。

## 3. 实现与验证

在唯一授权仓库按 acceptance、确认方案和有效决策实施最小改动。行为变化先加一个失败行为测试并确认失败，再最小实现并运行必要完整检查；涉及所有权、状态、决策或持久化时补并发验证。低影响文档修改不机械新增测试。

按 `validation-and-evidence.md` 核查初始化、TestMain、脚本、工具链、依赖、外部目标与隔离条件后运行验证。保留真实失败和未运行范围，冻结最终源码、差异、依赖、配置与结果；不能为错误的测试夹具改动正确业务逻辑。

读取 [产物管理](artifact-management.md)，使用分配的当前证据、attempt 和 runtime。必要运行脚本、测试副本和缓存放 runtime；日志与关键失败/最终结果按 attempt 组织。已有仓库工具优先复用，不每轮另造 freeze 或记录脚本，不把完整运行环境默认归入长期证据。

按需读取本次范围内的其他 Skills。Skill 不扩大修改或发布授权；业务选择、跨仓库合同变化和验收缺口返回 mainAgent。需记录偏离时按委派安排执行 `decide` 或交由 mainAgent 记录，需审批的变化在确认前暂停受影响部分。

## 4. 交给协调者评审

实现和必要验证就绪后，按产物管理检查材料完整性并返回实际差异、冻结摘要、验证与下游交接材料，以及新增材料和遗留资源。mainAgent 按 `delivery-review` 核验或委派评审，再单独向用户展示本任务；编码者只修复返回的问题，不自行派生 Reviewer 或把自检称为独立审查。

缺少覆盖本次差异的提交授权时保持 `claimed` 和原 owner。等待授权不是任务完成，也不需要重新认领；沿用已经明确覆盖的授权，不由方案确认或另一任务批准推导提交权限。

准备接口/数据合同、固定依赖版本或产物及下游实际获取/解析证据。提交尚未形成时先交付冻结内容和可行获取方式；真实提交形成后补齐最终固定输入，不用浮动源码路径作为最终交接。

## 5. 授权提交与完成交接

mainAgent 确认本次验收和有效授权满足后，当前 writer 按仓库与宿主 Git 规范创建新增 commit，核对完整 SHA、父提交、源 blob 与被批准的冻结差异一致。变化影响证据时重验相关范围；仅提交或路径变化可在证明等价后复用证据。

`--validation` 包含当前 Review/验证摘要与入口，`--handoff` 包含下游合同、确切固定版本/产物、获取验证、兼容与集成方法。新任务在仍为 `claimed` 且 owner 正确时执行：

```bash
"$COORD" complete --project "$PROJECT" --task "$TASK_ID" --owner "$OWNER" \
  --commit "$COMMIT" --validation "$VALIDATION" --handoff "$HANDOFF"
```

检查实际结果并立即通知 mainAgent。由 mainAgent 核验、同步当前交付指针并读取 summary 调度下游；Git commit 本身不解锁依赖。修复不执行上面的 complete，返回真实新提交及修复证据，由 mainAgent 更新 correction 和当前交付指针。

工具部分失败先核对 Git、CLI、进程与证据；成功提交或 complete 不重放，不改写活跃控制文件。按产物管理保存必要证据后，仅处理本次创建且归属明确、不再使用的临时资源，分别核验停止、目录删除和必要保留，报告真实遗留与条件；不清理正式 worktree 或必要固定交接产物。
