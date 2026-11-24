---
name: agent-delivery-kit
description: Set up Agent Delivery Kit and open its delivery workflow. Use when the user asks to initialize the kit, check coordctl, configure its optional agents, or start or resume a requirement with Agent Delivery Kit.
---

# Agent Delivery Kit

本入口负责首次配置，以及将新建或续接请求转交给 `delivery-flow`。交付阶段、确认和写入规则以实际读取的 `delivery-flow` 为准。

## 定位工具包

公开仓库为 `https://github.com/near-notfaraway/agent_delivery_kit.git`，默认本地位置为 `~/.local/share/agent-delivery-kit`。依次使用用户本次指定的工具包目录、安装副本中已有的“工具包目录”配置（若有）、当前已知的仓库源位置，最后检查默认位置。安装副本的父目录和当前业务仓库不一定是工具包目录，不从两者猜测。

使用前确认目录包含 `go.mod`（模块为 `example.com/agent_delivery_kit`）、`README.md`、`cmd/coordctl`、`templates/project.template.json` 和 `skills/delivery-flow/SKILL.md`。用户明确指定的路径或已有配置失效时，说明如何修正，不静默切换到其他检出。正常交付不自动拉取或更新工具包；缺少有效检出时，提示先初始化或提供现有检出位置。

## 初始化

仅在用户明确要求初始化、检查或配置时执行对应步骤。没有明确动作时，提示用户可发“初始化”、提供新需求，或提供原需求目录续接。

1. 显式初始化且没有指定路径、已配置位置或可复用检出时，确认默认目标尚不存在，再通过 `git clone` 从上述公开仓库检出到默认位置，按宿主管理文件写入权限。默认位置若被其他文件占用、是符号链接或并非有效检出，不覆盖它，请用户指定已有工具包目录。检出后按上节核验，后续会话直接复用。不要对已有检出执行 `pull`、`reset` 或 `checkout`。
2. 初始化或检查 CLI 时，读取工具包 `README.md` 的“运行要求”和“构建并试运行 CLI”。在工具包根目录按 README 的命令编译到新的临时目录，再用仓库的空任务模板运行 `summary`；只有退出成功且返回空 `counts` 和 `ready` 才算通过。按 README 处理本次试运行目录，不安装全局二进制，也不覆盖已有需求里的二进制。
3. 默认使用宿主现有的通用 Agent。用户要配置专业 Agent 时，读取 `skills/delivery-flow/references/agent-delegation.md` 的“可选专业角色的配置”，使用 `agents/` 中的源指令，按宿主支持的方式配置和核验。宿主需要用户点击时，只指出具体操作；安装文件已写入不等于宿主已发现或角色已实际运行。专业角色缺失不阻塞通用流程。
4. 报告实际工具包位置、`coordctl` 检查结果和可选角色状态。Git、Go、权限、下载、构建或试运行未通过时，说明实际失败与恢复方式，不声称初始化完成；有有效检出时，可按正常授权继续文档整理。初始化本身不创建需求目录、任务图或业务 worktree。

## 新建或续接需求

实际读取已核验工具包目录中的 `skills/delivery-flow/SKILL.md`。以工具包根目录的绝对路径作为 `KIT_DIR`，传递用户的需求、涉及仓库或原需求目录。由 `delivery-flow` 辨别新建和恢复、核对已有效的确认及授权，并按当前阶段读取仓库内的其他 Skills。

正常启动或续接不重复安装、配置角色或试编译。正式交付使用需求独立的 `coordctl`：新需求在执行准备阶段构建，已有需求按恢复约定检查并沿用其二进制。仅选择本 Skill 不构成需求、方案、提交或发布授权。
