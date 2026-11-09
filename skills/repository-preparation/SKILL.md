---
name: repository-preparation
description: Verify repository baselines and prepare or reuse requirement worktrees before coding. Use for delivery preparation and recovery. Not for design, implementation, or publishing.
---

# 仓库准备

将已确认方案中的授权仓库准备为可开发、可恢复的工作区。输出真实基线和准备证据，不生成执行任务图、不启动编码 Agent。

在完整交付流程中由通用 mainAgent 直接执行，无需自定义仓库准备 Agent。输入与 writer 边界遵循 [委派约定](../delivery-flow/references/agent-delegation.md)，具体准备方法以本 Skill 为准。

## 1. 恢复输入与阶段

读取适用 `AGENTS.md`、已确认 `requirement.md` / `design.md`、已有仓库准备记录和用户授权。与 [delivery-flow](../delivery-flow/SKILL.md) 协作时，由协调者维护 `delivery.json`；本 Skill 返回准备结果，不与协调者并发写状态。

必要输入为需求目录、需求分支、各仓库职责与源路径、确认版本及授权范围。技术设计未确认时仅做只读核查，不创建 worktree。恢复已有需求沿用其位置、分支与确认，不能因缺少新版模板而重建。

按仓库区分授权修改、只读参考、外部团队负责。只准备授权修改仓库。需求名与分支按协调流程校验，路径使用绝对路径，创建目录前核对父目录和已有目标身份。

读取 [产物管理](../delivery-flow/references/artifact-management.md)，准备结果集中到既有准备入口和 delivery.json 的权威仓库记录，不要求额外准备 JSON 或逐仓报告。必要一次性脚本与可重建验证资源使用分配的 runtime。

## 2. 核实每个仓库的来源

读取源仓库的指令、`git status --short`、`git worktree list --porcelain`、配置远端及相关引用。未提交代码不属于基线；只读分析修订不自动成为开发基线。

默认使用该仓库 `refs/heads/master`。master 缺失、异常陈旧或缺少方案依赖的既有接口时，先核实并保存证据，暂停该仓库，直接询问用户应使用哪个分支。不能自行替换为 `HEAD`、当前功能分支、`main` 或发布分支。

用户指定替代分支时，在该仓库本地及配置远端核实，必要时 fetch。分支仍不存在则再次暂停该仓库，询问正确分支或远端，直到来源可解析。该授权只适用于这个仓库，重试保留已有有效授权。

创建前记录确切 ref、完整不可变 commit、来源授权及与分析修订的差异；差异影响确认设计时交回协调者处理。不能把浮动分支名当作固定基线。

## 3. 创建或复用需求 worktree

默认位置为 `<需求目录>/repos/<仓库别名>`。来源已核实、分支未存在且目标未占用时，在该源仓库执行：

```bash
BASE_COMMIT="$(git -C "$SOURCE_REPO" rev-parse --verify "$BASE_REF^{commit}")" &&
git -C "$SOURCE_REPO" worktree add -b "$BRANCH" "$WORKTREE_PATH" "$BASE_COMMIT"
```

变量替换为已核实值。需求分支不能与源分支同名。命令失败先检查实际 worktree 和分支是否已创建，不盲目重放。

分支已存在且未占用时，验证需求身份、历史与记录基线后执行 `git worktree add "$WORKTREE_PATH" "$BRANCH"`。目标 worktree 已存在时，核对真实源仓库、分支、基线、HEAD、未提交代码和占用者后复用，不再次创建。

需求内修复默认复用原需求分支/worktree，包括任务已完成后的修复。先确认旧写入者已停止、没有活跃验证占用和未处理的他人修改，按 [恢复与返工约定](../delivery-flow/references/recovery-and-corrections.md) 续接。只在实际并发、基线或隔离需要成立时另建目录，说明理由和修订关系。

非 Git 源目录、身份不符的非空目标、分支被其他 worktree 占用或基线不能证明时保存冲突，不使用 `--force`、删除目录、重置分支或初始化 Git 来绕过。

已核实的环境限制确需本地镜像等支撑资源时，记录真实源、基线、Git 元数据引用和保留条件。承载正式 worktree 的支撑镜像属于仓库资源，不因名字含 mirror 或位于临时位置而可直接清理；新增或变更隔离方式仍需满足现有授权和基线规则。

用户明确纠正干净且未开发 worktree 的基线时，核对其没有额外提交，保留旧分支作为备份，再按实际授权处理；已有实现或额外提交必须保留。删除、重置等动作仍需明确授权，不能由基线纠正推导任意丢弃工作权限。

## 4. 核查可执行环境

在启动测试或构建前读取 [验证与证据约定](../delivery-flow/references/validation-and-evidence.md) 的运行准备部分：

- 读取实际 `go.mod`、工具链、包初始化、`TestMain` 和将执行测试的副作用；前端核对 lockfile、脚本及环境变量用途。
- 识别网络、数据库和外部资源目标，确认隔离配置与资源归属；不能因命令名称带 test 就假设安全。
- 记录工具链、依赖获取条件、必需命令和已知基线失败。只编译、兼容构建与运行测试分别标注。
- 构建输出使用需求独立目录，不覆盖活跃源码、服务或协调器二进制。

环境未就绪不需要销毁成功 worktree；记录其准备事实与验证阻塞，交给协调者判断当前阶段是否可推进。业务验收和检查例外不能由本 Skill 自动批准。

## 5. 交付准备结果

使用 [preparation_template.md](assets/preparation_template.md) 形成一个准备入口，默认 `validation/preparation.md`；已有入口继续使用。记录各仓库的源路径、ref/commit、需求分支、worktree、实际 HEAD、确认依据、工具链与阻塞，不机械生成逐仓库副本。

部分成功逐仓保存，重试只处理未完成部分。所有授权修改仓库均创建并核验成功后，返回“仓库就绪”；仍有来源不明或 worktree 未完成时返回局部结果，整体留在仓库准备阶段。

协调者据结果更新 `delivery.json.repositories` 及阶段，然后生成控制材料。目录存在、命令已发出或部分仓库成功都不能作为全部准备完成依据。

返回前按产物管理核对当前入口、新增材料、支撑资源和临时资源的实际处理，写入同一准备记录。保留有用途的资源并说明条件，不能将部分准备或仅停止进程描述为资源已清理。

> 为这个已确认方案准备仓库。逐仓核对 master 或用户授权的替代基线，复用符合记录的需求 worktree，保存固定修订和阻塞证据，不启动编码。
