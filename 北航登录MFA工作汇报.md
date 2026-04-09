# 北航登录 MFA 工作汇报

> 使用说明
>
> - 这个文件长期持续使用，不每周新建
> - 每周按日期追加一块内容
> - 每个模块负责人按实际情况补充
> - 重点记录可讨论、可决策、可推进的内容

========================================================================
## 2026-04-09
========================================================================

### 日期

- 周次：第 1 周
- 时间范围：2026-04-03 至 2026-04-09

### 本周总览

- 本周重点：完成北航登录 MFA 智能风险平台的统一系统架构设计，明确 Online、Offline、Research 三层关系，以及 Anti-Bot 先于 Device-ID / RBA 的顺序逻辑。
- 本周整体结果：已经形成完整的工程系统图、详细系统说明、分阶段路线和汇报模板，产品主线已经从“单点 anti-bot”提升为“统一登录风险平台”。
- 下周整体目标：继续把 P1 的核心表设计、Go 在线工程结构、Python offline 工程结构和数据流约束收敛为可直接开工的方案。

### 成员 1：系统架构与产品定义

- 本周做了什么：重新定义了产品边界，明确 shared foundation、online execution、offline intelligence 三层结构。
- 当前结果：系统主干稳定为 `Frontend SDK -> bfp_event -> Online Anti-Bot -> Device-ID / RBA -> Final Decision`。
- 下周计划：继续压缩 P1 MVP 范围，明确最小上线模块。
- 需要讨论的问题：对外发布时采用“登录 MFA 风控平台”还是“BFP + MFA 智能认证平台”的命名方式。

### 成员 2：Online / Offline 闭环设计

- 本周做了什么：明确了 Offline intelligence 通过 `online_rule_config`、`risk_profiles`、`online rule refresh` 动态更新 Online 的闭环。
- 当前结果：已经确认 `device-id v1` 与 `device-id v2` 都放在同一 `Device Identity Interface` 下，v2 在 P2 上线，v1 继续保留。
- 下周计划：进一步细化 online anti-bot、online RBA、offline detect 和 publish 层的输入输出。
- 需要讨论的问题：P2 上线时是否需要给 v2 一段 shadow 观察期。

### 成员 3：Research / 评估治理

- 本周做了什么：增加了 `research_validation_dataset`、`research platform`、`agent / control layer`、`shadow evaluation`、`rule lifecycle` 等模块设计。
- 当前结果：research 层已明确主要从 offline decision 派生，不直接依赖生产查询逻辑。
- 下周计划：继续把 research dataset、evaluation tables 和 experiment control 再规范化。
- 需要讨论的问题：groundtruth 是否先以高置信 offline decision 代理，还是提前预留人工确认接口。

### 成员 4：dashboard 设计

- 本周做了什么：完善了dashboard的展示内容，加入了饼状图，趋势图等可视化形式，同时对检测到的用户的id进行了hash处理，保护用户隐私。
- 当前结果：dashboard已经能够展式线上检测遇到的风险的趋势和分布情况并作出相应的分析。
- 下周计划：继续完善dashboard的功能，删除不必要的展示内容，引入`agent`和用户交互，根据用户的需求做个性化展示。
- 需要讨论的问题：用户的个性化需求具体有多宽泛，是选项类的还是对话类;

### 本周共性问题

- 当前文档层已经清楚，但 P1 表结构和服务接口还需要进一步收敛到“直接可编码”的程度。
- 对外发布仓库需要控制内容范围，以文档和系统蓝图为主，不包含实验性实现代码。

### 本周新增资料 / 文件

- `工程化系统蓝图.md`
- `研究内容与创新蓝图.md`
- `北航登录MFA工作汇报.md`
