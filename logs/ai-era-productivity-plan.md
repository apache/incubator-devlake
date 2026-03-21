# Plan: AI 时代效能分析增强

## Context

DevLake 已有 Kiro (Q-Dev) 和 GitHub Copilot 两套 AI 工具的数据采集，也有 Copilot+DORA 的关联分析 dashboard。但存在三个关键缺口：
1. Kiro 的新格式数据（credits/messages）没有和 DORA 关联
2. 没有跨 AI 工具的统一对比视图
3. 没有 AI 成本效益指标（credits per PR/deployment）

本计划分三个阶段实现，P0 纯 Grafana dashboard 不改后端。

## Phase 1 (P0): 三个新 Grafana Dashboard — 无后端改动

### Dashboard A: `grafana/dashboards/KiroCreditsDORA.json`
**Kiro Credits + DORA Correlation**

复用 `GithubCopilotDORACorrelation.json` 的模式（周级聚合 JOIN + Pearson's r）。

关键面板：
- Weekly Credits vs PR Cycle Time（双轴时序图）
- Pearson's r 相关系数（stat panel）
- High AI Usage vs Low AI Usage 周的 cycle time 对比
- Credits vs Deployment Frequency / CFR / MTTR

数据源：`_tool_q_dev_user_report`（周聚合 credits）JOIN `project_pr_metrics`（周聚合 cycle time）on `week_start`

模板变量：`project` (from `project_pr_metrics`)

### Dashboard B: `grafana/dashboards/AICostEfficiency.json`
**AI Cost-Efficiency**

关键面板：
- Credits per Merged PR（周趋势）
- Credits per Deployment（周趋势）
- Credits per Issue Resolved（周趋势）
- Summary Stats（总计 credits, credits/PR, credits/deploy）

数据源：`_tool_q_dev_user_report` 周聚合 LEFT JOIN `pull_requests` / `cicd_deployment_commits` / `issues` 周聚合 on `week_start`

### Dashboard C: `grafana/dashboards/MultiAIComparison.json`
**Multi-AI Tool Comparison (Copilot vs Kiro)**

关键面板：
- Active Users 并排对比（周趋势）
- Code Generation Activity 对比
- LOC Accepted 对比
- Acceptance Rate 对比（柱状图/表格）

数据源：`_tool_q_dev_user_report` + `_tool_q_dev_user_data` vs `_tool_copilot_enterprise_daily_metrics`

模板变量：`project`, `connection_id`/`scope_id`（Copilot 侧）

## Phase 2 (P1): 用户身份映射 — 需要后端改动

### 核心问题：Kiro userId 如何绑定 PR author

**Kiro 侧数据：**
- `user_id`: AWS Identity Store UUID（如 `6478a4a8-60a1-70d9-37bc-6aae85f6746a`）
- `display_name`: 通过 AWS Identity Store API 解析的名字（如 `Yingchu Chen`）

**PR 侧数据：**
- `author_name`: git 提交者名字（如 `Yingchu Chen`）
- `author_id`: 平台格式（如 `github:GithubAccount:1:12345`）

**可行方案：**

| 方案 | 做法 | 优点 | 缺点 |
|---|---|---|---|
| A: DisplayName 匹配 | Kiro display_name == PR author_name | 简单 | 名字可能不一致，多人同名 |
| B: 手动配置映射 | 在 connection/scope config 里让用户配置映射 | 精确 | 需要手动维护 |
| C: AWS SSO → 邮箱 | 扩展 Identity Client 获取用户邮箱，匹配 git commit 邮箱 | 自动且精确 | 需要额外 AWS API 权限 |
| D: 聚合级别（不绑定用户） | P0 已用的方式：按周聚合整个团队 | 零改动 | 无法做 per-developer 分析 |

**推荐路线：D（P0 先上）→ B+C（P1 实现）**

### 新模型
- `backend/core/models/domainlayer/crossdomain/ai_tool_user_mapping.go` — 映射 Kiro userId/DisplayName ↔ git AuthorName ↔ Copilot UserLogin

### 新 Migration
- `backend/core/models/migrationscripts/20260321_add_ai_tool_user_mapping.go`

### 新 Task
- `backend/plugins/q_dev/tasks/identity_mapper.go` — 从 `_tool_q_dev_user_report.display_name` 匹配 `pull_requests.author_name`（精确 + 模糊）
- `backend/plugins/gh-copilot/tasks/identity_mapper.go` — 类似逻辑

### 新 Dashboard
- `grafana/dashboards/KiroUserProductivity.json` — 基于映射表的 per-developer AI productivity vs PR metrics

## Phase 3 (P2): AI + Code Quality 关联

### 新 Dashboard
- `grafana/dashboards/AICodeQuality.json`
- AI Usage Intensity vs 新增 Code Issues（周趋势）
- Per-Developer AI Usage vs Code Quality（需 P1 映射表）
- AI Code "Survival Rate"（高 AI 使用周 vs 低使用周的 PR revert 率）

数据源：`cq_issues`（有 `commit_author_email` 可作为 join key）

## 关键参考文件

| 用途 | 文件路径 |
|---|---|
| Copilot+DORA 模板 | `grafana/dashboards/GithubCopilotDORACorrelation.json` |
| QDev+DORA 现有 | `grafana/dashboards/QDevDORA.json` |
| Kiro credits 模型 | `backend/plugins/q_dev/models/user_report.go` |
| PR metrics 模型 | `backend/core/models/domainlayer/crossdomain/project_pr_metric.go` |
| 身份解析参考 | `backend/plugins/q_dev/tasks/identity_client.go` |
| 用户映射表参考 | `backend/core/models/domainlayer/crossdomain/user_account.go` |

## 实施顺序

| Step | Phase | 产出 | 依赖 |
|------|-------|------|------|
| 1 | P0 | `KiroCreditsDORA.json` | 无 |
| 2 | P0 | `AICostEfficiency.json` | 无 |
| 3 | P0 | `MultiAIComparison.json` | Copilot + Kiro 数据 |
| 4 | P1 | `ai_tool_user_mapping.go` 模型 + migration | 无 |
| 5 | P1 | `q_dev/tasks/identity_mapper.go` | Step 4 |
| 6 | P1 | `gh-copilot/tasks/identity_mapper.go` | Step 4 |
| 7 | P1 | `KiroUserProductivity.json` | Steps 5-6 |
| 8 | P2 | `AICodeQuality.json` | Steps 5-6 + SonarQube 数据 |

Steps 1-3 可并行。

## 验证方法

### P0 验证
1. Import dashboard JSON 到 Grafana
2. 确认 `$project` 变量从 `project_pr_metrics` 正确加载
3. 确认 Pearson's r 值在 [-1, 1] 范围且至少有 4 个数据点
4. Playwright E2E 截图每个 dashboard

### P1 验证
1. 单元测试精确匹配和模糊匹配逻辑
2. 集成测试：插入测试数据 → 运行 mapper → 验证映射表
3. 验证 per-developer dashboard 数据正确

### P2 验证
1. 需要 SonarQube 数据在 `cq_issues` 中
2. 验证 `commit_author_email` 与映射表匹配
