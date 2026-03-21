# AI 时代效能分析 — 现状评估与优化方向

日期: 2026-03-20

## 一、现有数据资产全景

### Kiro (Q-Dev) 插件
| 表名 | 数据内容 | 粒度 |
|---|---|---|
| `_tool_q_dev_user_report` | credits, messages, conversations, subscription tier, overage | 每天/每用户/每 client_type |
| `_tool_q_dev_user_data` | 旧格式 43 个 feature-level 指标 (inline, chat, code fix, code review, dev, doc gen, test gen, transformation) | 每天/每用户 |
| `_tool_q_dev_chat_log` | 逐条 chat 交互 (prompt/response 长度, model, conversation, steering, spec mode, code refs, web links) | 秒级/每请求 |
| `_tool_q_dev_completion_log` | 逐条 inline completion (file, context length, completions count) | 秒级/每请求 |

### GitHub Copilot 插件
| 表名 | 数据内容 | 粒度 |
|---|---|---|
| `_tool_copilot_enterprise_daily_metrics` | DAU/WAU/MAU, code gen/acceptance activity, LOC suggested/added, PR created/reviewed by Copilot | 每天/企业级 |
| `_tool_copilot_user_daily_metrics` | user_login, code gen/acceptance activity | 每天/每用户 |
| `_tool_copilot_metrics_by_ide` | 按 IDE 分的使用量 | 每天 |
| `_tool_copilot_metrics_by_feature` | 按功能 (chat_panel, inline_chat 等) 分 | 每天 |
| `_tool_copilot_metrics_by_language_feature` | 按编程语言+功能分 | 每天 |
| `_tool_copilot_seats` | license 分配和最后活跃时间 | 快照 |

### 传统 DevOps 数据 (Domain Layer)
| 数据域 | 表数 | 关键表 |
|---|---|---|
| Code | 11 | repos, commits, commit_files, pull_requests, refs |
| DevOps/CICD | 8 | cicd_pipelines, cicd_deployments, cicd_deployment_commits |
| Tickets | 18 | issues, boards, sprints |
| Code Quality | 5 | cq_projects, cq_issues, cq_file_metrics |
| QA | 4 | qa_test_cases, qa_test_case_executions |
| Cross-domain | 18 | users, accounts, project_pr_metrics |

### 已有 Dashboard (50+)
- DORA 系列 (5 个): 四指标 + 明细
- Copilot 系列 (2 个): Adoption + DORA Correlation
- Kiro 系列 (4 个): Usage, Legacy Feature, Logging, Executive
- 传统 DevOps: GitHub, GitLab, Jira, Jenkins, SonarQube 等

---

## 二、现有能力 vs AI 时代效能分析的差距

### 已经有的
- Kiro 使用数据 (credits, messages, feature-level metrics, prompt logs)
- GitHub Copilot 采纳数据 (seats, acceptance rates, LOC)
- Copilot + DORA 关联分析 (已有 dashboard, 含 Pearson's r 相关系数)
- 传统 DevOps 指标 (commits, PRs, DORA 四指标, SonarQube)

### 关键缺失：AI 对结果的影响闭环

目前所有 AI 指标都是"输入侧"的 — 用了多少、接受了多少。但缺乏"输出侧"的闭环：

| 维度 | 现状 | 缺失 |
|---|---|---|
| AI 代码命运追踪 | 知道 acceptance rate | 不知道接受后是否被 revert、是否引入 bug |
| AI 对代码质量的影响 | 有 SonarQube 数据 | 没有 AI 使用量 vs code smell 的关联 |
| AI 对效率的因果关系 | Copilot+DORA 有相关性(r值) | Kiro 没有同等的 DORA 关联 |
| 多 AI 工具统一视图 | Copilot 和 Kiro 各自独立 | 无法对比/汇总跨工具的 AI 总效能 |
| AI ROI (成本效益) | 有 credits_used | 没有 credits per PR merged / per deployment |
| 用户身份统一 | 各插件独立 user_id | Kiro userId -> git author -> Copilot userLogin 无法打通 |

---

## 三、优化方向 (按优先级)

### P0 — 短期可做 (利用现有数据，无需新数据源)

#### 1. Kiro + DORA 关联 Dashboard
- 类似已有的 Copilot+DORA correlation dashboard
- 将 Kiro credits/messages 与 PR cycle time 关联
- 按周聚合，计算 Pearson's r 相关系数
- 分桶对比：AI 重度使用周 vs 轻度使用周的 DORA 表现
- 数据来源: `_tool_q_dev_user_report` + `project_pr_metrics`

#### 2. AI 成本效益 Dashboard
- Credits per PR merged
- Credits per deployment
- Credits per accepted line of code (已有)
- Credits per issue resolved
- 趋势：成本效益是否随时间改善
- 数据来源: `_tool_q_dev_user_report` + `pull_requests` + `cicd_deployments` + `issues`

#### 3. 多 AI 工具对比 Dashboard
- Copilot vs Kiro 并排对比
- 统一指标: adoption rate, acceptance rate, LOC generated, active users
- 趋势对比: 两个工具的采纳曲线
- 数据来源: `_tool_copilot_enterprise_daily_metrics` + `_tool_q_dev_user_report`

### P1 — 中期 (需要新的数据打通)

#### 4. AI Code Quality 闭环
- SonarQube findings 按 "AI 重度用户 vs 轻度用户" 分组
- AI 辅助 code review vs 人工 review 的 finding 数量对比
- 需要打通: AI 使用量 -> commit author -> SonarQube findings
- 数据来源: `_tool_q_dev_user_data` + `commits` + `cq_issues`

#### 5. 用户身份统一层
- 建立 Kiro userId -> domain users -> Copilot userLogin 的映射
- 可以通过 email 或 display name 做 fuzzy matching
- 一旦打通，所有 AI 指标都可以跟 git/PR/issue 指标关联
- 实现方式: 新的 cross-domain 映射表或扩展现有 user_accounts 表

### P2 — 长期 (需要新数据源或复杂分析)

#### 6. AI 代码存活率
- 追踪 AI 生成代码从接受到被修改/删除的时间
- 需要: commit diff 分析 + AI acceptance 时间戳关联
- 指标: AI 代码平均存活天数, AI 代码 revert 率

#### 7. AI 辅助 Code Review 效能
- 衡量 AI 对 review 速度和质量的影响
- 有 AI review (Kiro code review) vs 纯人工 review 的对比
- 指标: review time, comments count, approval rate, rework rate

#### 8. Developer Flow 与 AI 的关系
- 分析开发者何时使用 AI (时段、任务类型)
- AI 使用的上下文切换频率
- AI prompt 复杂度随时间的演变 (学习曲线)
- 数据来源: `_tool_q_dev_chat_log` (已有 timestamp, prompt_length, model_id)

---

## 四、可连接的数据节点

```
Kiro user_id ──────┐
                    ├──> domain users.id ──> commits.author_id ──> PR, Issues, DORA
Copilot user_login ─┘                   ──> cq_issues (SonarQube)
                                         ──> cicd_deployments
                                         ──> qa_test_case_executions

_tool_q_dev_user_report.date ──> project_pr_metrics (weekly join for DORA correlation)
_tool_q_dev_chat_log.timestamp ──> commits.created_date (intra-day correlation)
_tool_copilot_enterprise_daily_metrics.day ──> project_pr_metrics (already done in copilot_impact)
```

---

## 五、技术实现建议

1. P0 的三个 dashboard 可以纯 Grafana SQL 实现，无需后端改动
2. P1 的用户身份统一需要一个新的 convertor subtask 或者 domain layer 扩展
3. P2 的代码存活率需要 git diff 分析能力 (git extractor 已有 commit_files 数据)
4. 建议先做 P0 中的 "Kiro + DORA 关联"，因为已有 Copilot+DORA 的模板可以参考
