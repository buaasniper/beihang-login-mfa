# 智御万象：融合多元威胁情报的Agent驱动自适应MFA认证系统

2026 年 4 月

## 摘要

面向统一身份认证场景，传统基于静态规则和人工维护的登录防护体系正面临由大模型、自动化框架与智能代理共同驱动的攻击升级压力。攻击流量不仅具备更强的伪装能力，而且能够跨周期复用环境、调整访问节奏并绕过固定校验逻辑，导致常规验证码、黑名单和单点设备识别方案逐步失效。针对上述问题，本文提出一套融合多元威胁情报的 Agent 驱动自适应 MFA 认证系统，围绕浏览器指纹、设备锚定、风险驱动认证和离线智能反馈构建统一防御闭环。

该系统以浏览器端高维环境信号为基础，将 Canvas、WebGL、Fonts 等核心特征统一标准化为事件事实层，在在线链路中首先执行 Anti-Bot 门禁，再对通过门禁的请求实施 Device-ID 与风险驱动认证，从而形成“先过滤自动化流量、再判断真实身份风险”的两阶段决策机制。在离线链路中，系统围绕多时间窗口特征、设备关系分析、研究验证数据集和规则挖掘平台开展 intelligence 生产，并通过 Agent 控制层将高价值规则、阈值和画像动态回灌至在线执行体系，实现小时级策略迭代。

在工程实现上，本文构建了 Shared Foundation、Online System、Offline Intelligence 和 Agent Control Layer 四层架构，并提出 device-id v1 到 device-id v2 的可演进设备识别路径。前者基于确定性指纹链接实现稳定落地，后者面向模糊匹配与图结构推理提供升级接口。系统已在北京航空航天大学统一身份认证平台相关场景中完成部署与验证，累计处理数据超过 200 万条，识别并锁定 Bot 流量 6000 余条，验证了方案在真实业务中的有效性与可落地性。

本文工作的意义不仅在于提出一套可解释、可扩展的高阶 MFA 认证方案，更在于构建了一个可持续演进的安全智能底座，为高校、政企和金融等高安全需求场景提供国产化、工程化和研究化兼顾的身份风控能力。

**关键词：** 浏览器指纹；Anti-Bot；设备锚定；Agent；MFA；风险驱动认证

## Abstract

This paper presents an Agent-driven adaptive multi-factor authentication system for real-world login protection under the evolving threat landscape created by large language models, browser automation frameworks, and intelligent agents. In contemporary identity systems, attackers can imitate legitimate interaction patterns, rotate execution environments, and continuously adjust their strategies according to observed defenses. Consequently, conventional protection mechanisms based on static rules, manually maintained blacklists, and isolated device checks are becoming increasingly inadequate for high-value scenarios such as university identity gateways and enterprise login services.

To address this problem, the proposed system builds a unified architecture around browser fingerprinting, anti-bot detection, device anchoring, and risk-based authentication. The central idea is to transform heterogeneous browser-side signals into a canonical event layer and then organize online decisions into two stages. The first stage performs low-latency anti-bot filtering through basic attribute validation, device-level risk inspection, and environment-consistency analysis. The second stage is activated only for traffic that passes the anti-bot gate, where the system resolves device identity and applies risk-based authentication. This design reduces online computation cost and prevents obvious bot traffic from contaminating identity reasoning.

At the system level, the architecture is organized into four coordinated layers: Shared Foundation, Online System, Offline Intelligence, and Agent Control Layer. The Shared Foundation standardizes incoming signals and maintains canonical event facts; the Online System performs real-time anti-bot judgment, device identity resolution, and final decision output; the Offline Intelligence layer aggregates multi-window features, generates anti-bot and RBA intelligence, and supports research validation; and the Agent Control Layer publishes validated rules, thresholds, and risk profiles back to online execution while keeping offline evaluation and research artifacts aligned with the newest strategies.

To support long-term evolution, the paper further proposes a versioned device identity path. The deployable baseline, device-id v1, relies on deterministic browser-fingerprint linking, while the reserved upgrade path, device-id v2, is designed for fuzzy matching and graph-based identity reasoning. This interface-based design allows the production system to remain stable while more advanced methods are validated and gradually promoted. The system has been deployed in a real university identity-authentication scenario related to Beihang University, where it has processed more than two million records, identified over six thousand bot flows, and achieved hour-level strategy iteration. These results demonstrate the practical value of integrating online execution, offline intelligence, and research-driven evolution into a unified adaptive MFA platform.

**Keywords:** Browser Fingerprinting; Anti-Bot; Device Anchoring; Agent; MFA; Risk-Based Authentication

## 第一章 绪论

### 研究背景与意义

随着人工智能技术、浏览器自动化框架以及智能代理系统的快速发展，登录认证入口正在成为高价值攻击最集中的防线之一。高校统一身份认证、政务门户和企业账号体系普遍承载着个人信息、教学资源、办公应用和业务权限，一旦认证侧被突破，后续权限扩散和数据泄露风险将被显著放大。与传统脚本攻击不同，当前自动化攻击已具备更强的页面理解、行为模仿和环境伪装能力，攻击者能够借助大模型生成自适应访问策略，基于自动化框架批量发起低成本、高频次、长周期的试探与渗透行为。

在这一背景下，传统依赖验证码、IP 黑名单、固定阈值和人工维护策略的防护体系逐渐暴露出明显缺陷。一方面，静态规则容易被学习和绕过；另一方面，单点设备识别对浏览器环境漂移十分敏感，难以在长时间跨度内稳定追踪同一攻击主体。此外，从离线分析到线上策略发布往往依赖人工流转，导致策略迭代速度无法匹配攻击演化速度。由此可见，构建兼顾实时拦截、持续设备锚定和动态 intelligence 反馈的一体化认证系统，已成为智能时代身份安全建设的重要需求。

浏览器指纹技术为解决该问题提供了新的切入点。与行为特征相比，浏览器指纹更强调设备底层执行环境的稳定差异，如图形渲染链路、字体集合、系统接口响应和浏览器能力组合等；与单纯网络特征相比，其在跨网络、跨账号、跨会话关联方面具有更强的可持续性。若进一步将浏览器指纹与风险驱动认证、研究验证平台和 Agent 控制机制结合，就有可能构建一个既能在线快速判定、又能离线持续进化的统一认证风险平台。

### 国内外研究现状

围绕浏览器指纹、反自动化检测和设备识别，国内外学术界与工业界已经积累了较多研究与实践成果。国外研究较早关注浏览器指纹的唯一性、稳定性与隐私影响，并围绕 Canvas、WebGL、字体、插件、音频上下文等高熵特征开展了测量分析和模型研究。工业界方面，Cloudflare、FingerprintJS、DataDome 等厂商已将浏览器指纹纳入反 Bot 和风险控制体系，通过与网络侧、行为侧和威胁情报侧特征联动，提高复杂流量检测能力。

然而，现有方案仍存在若干共性问题。首先，部分系统偏向黑盒检测，缺乏对判定依据、规则来源和决策逻辑的可解释表达，这会降低认证场景中的可信度和可审计性。其次，设备识别往往过度依赖静态特征，一旦浏览器升级、插件变化或系统环境轻微调整，便容易出现识别漂移。再次，许多系统只强调“检测”而忽视“更新”，即便离线侧发现了新的攻击模式，也难以快速、安全地将其转化为在线可执行规则。

国内相关研究近年来亦在持续推进，尤其在浏览器指纹稳定性建模、环境伪装识别、多源风险融合和账号安全防控方面取得了明显进展。但从完整平台角度看，国内仍较少出现同时覆盖高阶 Anti-Bot、设备锚定、动态 intelligence 更新、研究验证闭环与可解释展示能力的统一系统。特别是在高校统一身份认证等真实业务场景中，如何将研究能力、工程能力与国产化落地能力融合到同一平台，仍然具有较高探索价值。

### 研究内容与贡献

针对上述问题，本文围绕“融合多元威胁情报的 Agent 驱动自适应 MFA 认证系统”展开研究，主要工作与贡献包括以下几个方面。

第一，提出“Anti-Bot 优先、Device-ID 与 RBA 递进”的两阶段在线认证框架。系统将明显自动化流量在第一阶段快速筛除，仅对通过门禁的流量开展设备身份与风险驱动认证，有效降低后续高成本分析的负担，并提升身份推理质量。第二，设计多层递进的 Anti-Bot 引擎，围绕基础属性校验、设备级风险识别和环境一致性交叉验证构建在线快速防御机制。

第三，构建基于浏览器指纹链接的设备锚定路径，提出可工程落地的 device-id v1，并在架构层预留面向模糊匹配和图结构推理的 device-id v2 升级接口。第四，提出 Agent 驱动的离线到在线 intelligence 闭环，使规则、阈值、画像和模型结果能够通过 research platform 与控制层统一验证、统一发布。第五，在真实的北航统一身份认证相关场景中完成部署与验证，证明了该体系在工程可落地性、在线实时性、离线可研究性和系统可解释性上的综合价值。

除上述贡献外，本文还强调研究与工程的统一：一方面以可部署、可运行、可追踪为基本目标完成系统落地；另一方面通过研究验证数据集、影子评估和版本化升级路径，为后续规则、模型和图算法持续迭代提供标准化实验底座。

进一步而言，本文的创新并非单点算法改进，而是围绕三个可检验痛点展开：其一，针对“前端特征扩展慢、后端接入成本高”的问题，提出以统一事实层和 feature registry 为中心的可演进特征体系；其二，针对“离线 intelligence 难以快速上线”的问题，提出 Agent 控制下的 research-to-production 发布闭环；其三，针对“设备识别难以在工程系统中平滑升级”的问题，提出 device-id v1 到 v2 的接口化演进路径。上述创新均直接对应系统中的具体模块与数据资产，因此具备明确的工程可验证性。

## 第二章 系统总体设计

### 2.1 浏览器指纹概述

浏览器指纹是指在不依赖显式身份标识的前提下，通过采集浏览器执行环境中的多维属性来描述设备差异的一类技术。其核心思想不是寻找单一的永久标识，而是通过一组具备区分度和相对稳定性的环境特征，构建对设备执行上下文的联合表达。典型信号包括 User-Agent、Canvas 渲染输出、WebGL 参数、字体集合、屏幕信息、语言配置、时区表现以及部分接口能力响应等。

在认证安全场景中，浏览器指纹的价值主要体现在三个方面。其一，它能够在无感知条件下补充设备真实性信息，减少对 Cookie、IP 等易变信号的过度依赖。其二，它可以为设备级风控和跨账号关联提供底座，使风控对象从“单次请求”扩展到“持续出现的设备实体”。其三，在自动化攻击广泛采用行为伪装的前提下，底层执行环境往往仍会暴露难以完全伪造的一致性缺口，这为高对抗检测提供了重要突破口。

需要指出的是，浏览器指纹并非静态不变。浏览器升级、插件变化、字体增减和系统配置调整都会引起一定程度的指纹漂移。因此，指纹技术若要真正服务于认证系统，就必须与在线决策、离线建模和动态更新机制结合，形成面向演进场景的统一设计，而非停留在一次性设备标识层面。

### 2.2 系统整体架构设计

本文提出的系统由 Shared Foundation、Online System、Offline Intelligence 与 Agent Control Layer 四个部分构成，并围绕统一事件层形成贯通在线执行与离线 intelligence 的整体框架。系统首先通过前端 SDK 采集浏览器环境特征，再由后端标准化引擎完成特征提取、哈希生成与事件落库，形成统一事实表 `bfp_event`。在此基础上，在线系统执行低延迟风险判定，离线系统执行多窗口分析、设备关系构建和规则挖掘，控制层则负责 intelligence 的双向反馈与动态编排。

从决策顺序看，系统遵循“先 Anti-Bot，后 Device-ID / RBA”的原则。Online Anti-Bot Engine 负责识别明显的自动化流量并快速终止其后续流程；只有在 Anti-Bot 阶段通过的请求，才进入设备身份解析与风险驱动认证。这一设计既减少了无效计算，又避免将明显 bot 流量引入设备关系分析链路，从源头降低后续 identity graph 被污染的风险。

离线侧的核心目标不是替代在线决策，而是为在线决策持续提供更强 intelligence。通过 5 分钟和 1 小时等多时间窗口特征、设备关系分析、研究验证数据集和影子评估机制，离线系统可以不断发现更高价值的规则、阈值和画像，再经由 Agent 控制层发布到在线配置与画像资产中。由此，系统不再是单次检测器，而成为可持续进化的认证风险平台。

需要强调的是，上述分层并非任意拆分，而是对若干备选方案进行权衡后的结果。若将 Anti-Bot、Device-ID 与 RBA 合并为单阶段统一决策，则在线链路需要同时承担自动化流量过滤、设备解析与身份推理，既增加了延迟，也会让明显 bot 流量进入高成本分析流程；若完全依赖离线模型再统一下发，则在线系统将失去对新攻击的即时门禁能力；若仅保留在线规则而缺少离线 intelligence 和研究层，则系统又会迅速退化为静态防御工具。因此，本文采用“共享底座 + 在线分阶段执行 + 离线 intelligence 沉淀 + 控制层发布”的结构，以平衡实时性、可解释性、可演进性与研究可复现性。

#### 2.2.1 数据采集与接入机制

系统入口由前端 SDK 与后端 Go Ingest API 共同构成。SDK 不仅负责采集指纹与上下文特征，还承担统一 payload schema、SDK 版本控制和 feature flag 管理等职责，从而将前端特征采集从零散脚本提升为可演进、可治理的产品化组件。相比于简单的页面级 JS 上报方式，SDK 可以更稳定地支持字段扩展、浏览器兼容处理和后端口径统一。

Go Ingest API 则作为在线接入的统一网关，负责请求校验、身份鉴别、流量限速、日志记录与链路编排。所有来自浏览器端的 JSON payload 都通过该入口进入后端体系，避免 SDK 直接耦合底层存储或在线判定逻辑。该设计为高并发接入、统一监控和后续微服务拆分提供了基础。

#### 2.2.2 特征标准化与统一数据基座

在接收到原始 payload 后，系统通过 Normalize + Hash Engine 提取在线核心特征，并重点生成 `canvas_hash`、`webgl_hash` 与 `fonts_hash` 三类标准化标识。为了避免大对象直接进入高频判定链路，系统将原始特征压缩为结构化、可索引、可关联的 hash 资产，并沉淀到对应 hash library 中。这一层既降低了在线计算成本，也为离线统计、画像构建和设备锚定提供了统一表示。

标准化完成后，系统将事件写入 `bfp_event`。该表是全系统的统一事实表，也是 online 与 offline 的共同入口。在线侧所有实时判定都围绕 `bfp_event` 展开，离线侧的特征聚合、规则挖掘和研究验证同样从该事实层派生。通过这一设计，系统避免了多模块重复解析大 JSON、字段口径不一致以及特征逻辑碎片化的问题，确保平台具备清晰的数据中轴。

#### 2.2.3 在线分阶段决策机制

在线决策链路分为两个连续阶段。第一阶段由 Online Anti-Bot Engine 实现，其目标是在极低延迟下完成 Bot 门禁。该阶段综合读取 `bfp_event`、`online_rule_config` 以及来自离线系统的风险画像，在请求抵达时快速进行规则匹配与环境校验，并输出 `online_bot_result`。一旦命中 bot / challenge 条件，请求即可在此阶段结束。

第二阶段为 Device Identity 与 Online RBA。系统通过 Device Identity Interface 将设备识别过程抽象为统一接口：P1 中以 `device-id v1` 为主，基于确定性指纹链接快速生成设备身份；P2 中引入 `device-id v2`，在保持接口稳定的前提下实现模糊匹配、图谱推理和更高阶的实体解析。随后，Online RBA Engine 结合设备身份、账号画像、设备画像和请求上下文，输出 `online_rba_result`，并与 Anti-Bot 结果共同构成最终的 MFA 风险决策。

这一分阶段机制的优势在于：既保留了在线链路的实时性，又使身份推理资源只服务于更高价值的流量；既满足当前规则驱动的工程落地需求，又为后续更复杂的设备识别与 RBA 演进保留了清晰接口。

#### 2.2.4 离线智能建模与能力沉淀

离线智能系统承担多时间尺度建模、深层模式识别和高阶关系分析任务。首先，系统通过 `job_offline_prepare_5min` 与 `job_offline_prepare_1h` 生成短窗口与长窗口聚合特征，用于刻画账号、UA、设备和指纹实体在不同时间尺度下的访问强度与行为模式。随后，`job_offline_detect` 基于这些特征执行离线 Anti-Bot 判断，产出 `offline_bot_result`，并同步生成 online 可消费的规则候选与画像资产。

在设备与风险关系分析方面，离线侧进一步通过 linking / graph / RBA 任务构建 `offline_rba_result`、`risk_link` 与 `risk_entity_score`。这些结果不只是辅助报表，而是构成整个系统更高阶 intelligence 的沉淀层：一方面，它们为 device-id v2 的训练与验证提供基础；另一方面，它们也是 online 风险画像、规则刷新和实验比较的关键输入。

为了支持后续研究与算法验证，系统还从离线决策层中派生出 `research_validation_dataset`，将 `offline_bot_result` 与 `offline_rba_result` 等结果快照化、版本化，使新规则、新模型和新链接算法能够在统一验证集上进行公平比较。这样，离线系统就从单纯分析模块提升为兼具 intelligence 生产、研究验证和策略供给能力的统一平台。

其中，`research_validation_dataset` 不被设计为新的复杂多源事实层，而是定位为离线决策结果的研究副本。其核心样本由 `offline_bot_result`、`offline_rba_result` 及必要的事件上下文字段派生而来，并通过版本号、时间窗和实验切分标记进行管理。这样的设计有两个直接好处：一是避免研究平台再次与生产明细表深度耦合，二是使规则实验、模型实验和 linking 实验能够在同一批稳定快照上反复复现。为了降低数据泄漏风险，研究数据集在构造时应按时间段和实体维度进行切分，使同一设备链或同一攻击批次不会同时进入训练与验证集合。

#### 2.2.5 智能控制与动态反馈机制

Agent Control Layer 是系统中的双更新控制中枢。一方面，它负责接收 `job_offline_detect`、离线 RBA 分析以及 research platform 产生的候选 intelligence，将其转化为可执行的 `online_rule_config` 与 `risk_profiles`，并通过在线 refresh 机制更新 Anti-Bot 和 RBA 的执行资产。另一方面，它也会将最新的在线配置版本和发布结果反向同步到离线检测与评估路径中，确保研究验证、shadow evaluation 和离线 intelligence 与线上状态保持一致。

换言之，该层并非简单的“配置分发器”，而是连接在线执行、离线 intelligence 与研究平台的智能编排器。它既管理规则、阈值和画像的上线节奏，也管理哪些 intelligence 需要进入影子验证、哪些模型结论只能在 research platform 内部保留、哪些结果可以回流更新离线检测逻辑。正是这一双更新机制，使系统能够在保证在线稳定性的同时持续吸收离线发现，形成真正可演进的认证安全闭环。

从控制逻辑上看，Agent Control Layer 同时承担“向前发布”和“向后约束”两类职责。所谓向前发布，是指将通过验证的离线 intelligence 动态注入在线规则与画像系统；所谓向后约束，是指要求任何候选规则、候选模型和候选链接算法在进入线上前必须先经过研究数据集验证与影子评估，从而避免未验证 intelligence 直接影响生产判定。由此，系统中的 Agent 不是泛化意义上的自动化脚本，而是连接 research、offline decision 与 online execution 的策略治理核心。

#### 2.2.6 模块协同与闭环机制

整体而言，系统各模块围绕统一事件事实层和统一 intelligence 资产层形成稳定协同。Shared Foundation 负责统一采集、标准化和状态管理；Online System 负责实时检测、设备识别和认证决策；Offline Intelligence 负责多窗口特征分析、关系建模和研究验证；Agent Control Layer 负责 intelligence 的筛选、发布与回灌。各层既分工明确，又通过共享表结构和版本化资产紧密耦合。

这种协同关系使系统具备三个显著特点：其一，在线和离线不再是割裂的两套体系，而是通过统一事实层和画像层构成闭环；其二，产品执行面和研究验证面可以在同一底座上并行演进，既保障生产可用，也支撑算法创新；其三，系统天然适合从 device-id v1 逐步升级到 v2、从规则驱动逐步扩展到模型驱动，形成面向长期发展的技术路线。

### 2.3 指纹链接技术原理

在认证安全场景中，设备识别的难点不在于一次性区分设备，而在于在环境轻微变化和对抗伪装同时存在的条件下保持跨周期稳定关联。为此，系统将浏览器指纹从静态标识视角转化为动态链接视角，通过特征分层、规则约束与版本化解析接口实现设备锚定能力的工程落地。

#### 2.3.1 指纹属性分类与稳定性建模

系统首先将指纹特征划分为稳定底层特征和可演化环境特征两类。稳定底层特征主要反映设备硬件、驱动或渲染能力层面的差异，例如部分 WebGL 渲染表现、Canvas 绘制差异和与底层执行路径强相关的信号；可演化环境特征则包括浏览器版本、字体集合、语言配置、时区与部分功能支持状态等，它们更容易随正常使用过程发生变化。

通过这一分类，系统能够在设备识别中显式区分“可以容忍的自然漂移”和“值得警惕的风险变化”。也就是说，设备链接不再要求所有字段严格一致，而是要求稳定特征维持足够一致性，同时允许部分上层特征在可解释范围内发生演化。该思想为后续确定性链接和模糊图谱推理提供了统一建模基础。

#### 2.3.2 规则决策逻辑与链接算法

在 P1 阶段，系统采用规则型 `device-id v1` 作为主方案。该方案基于核心 hash 值组合和稳定环境约束生成确定性设备标识，强调快速、可解释和可工程化部署。在决策逻辑上，系统对关键底层特征设置较强一致性要求，对部分上层特征设置有限容忍区间，从而在保留稳定性的同时降低自然漂移带来的误拆分。

进入 P2 之后，系统将在统一 Device Identity Interface 下引入 `device-id v2`。该版本不再局限于确定性匹配，而是通过模糊相似度、关系边权重和图结构传播来完成更高阶的设备链接。由于接口保持不变，在线系统可以逐步从 v1 迁移到 v2，而无需重构上层 RBA 和决策逻辑。这种“先落地、后升级”的分层设计兼顾了生产稳定性与研究前沿性。

#### 2.3.3 工程实现的优势

当前以规则型 `device-id v1` 为基础的设备链接方案在工程上具有明显优势。首先，它适合在高并发在线环境中以较低延迟生成设备身份；其次，其判定逻辑可解释、可回溯，便于与在线风控和可视化展示结合；再次，它可以作为 device-id v2 的稳定基线，为后续模型化、图谱化算法提供对照与回退机制。

更重要的是，通过统一 Device Identity Interface 和持久化 identity tables，系统从架构上避免了“每次升级算法都要重写上层流程”的问题。也就是说，设备识别能力的增强被封装在身份解析层内部，而 Online Anti-Bot、Online RBA、Research Platform 等上层模块只依赖统一输出结果，从而获得更好的长期可维护性。

### 2.4 Anti-Bot 引擎与防御算法设计

Anti-Bot 引擎承担系统第一阶段门禁职责，其设计目标是在保证低延迟的前提下，对不同复杂度的自动化流量实施层级化识别。系统采用“三层递进”思路，将基础属性异常、设备级风险异常与环境一致性异常纳入同一防御体系，兼顾粗粒度快速过滤和高强度对抗识别。

#### 2.4.1 基础属性校验

第一层为基础属性校验，主要针对明显的自动化环境和低阶伪装。典型信号包括 `webdriver` 状态、无头浏览器特征、异常窗口参数、明显不合理的浏览器能力组合等。该层优先识别最容易暴露的自动化请求，在较小成本下完成第一轮剔除，为后续更精细的识别保留系统资源。

#### 2.4.2 基于风险的授权机制

第二层以设备和实体为中心构建风险判断逻辑。系统综合 `canvas_hash`、`webgl_hash`、`fonts_hash` 以及时间窗口聚合特征，识别短时异常高频访问、跨账号复用、典型攻击访问节奏和异常行为爆发等模式。当某一设备、账号或 UA 在窗口内表现出明显偏离正常分布的特征时，系统可触发更严格的 MFA 动作，包括 challenge、额外验证或直接拒绝。

#### 2.4.3 不一致性检测

第三层用于识别高伪装、高对抗流量。高级 Bot 可能在单一字段上表现“正常”，但在多维环境上出现互相矛盾，例如宣称为移动端环境却暴露桌面自动化行为特征，或在 UA、渲染能力、字体集合之间出现组合不一致。系统通过多维特征交叉校验捕获此类“伪一致”现象，从而提升对复杂自动化攻击的识别能力。与前两层相比，该层更强调精细化判定与高置信告警，是高阶在线 Anti-Bot 的关键部分。

### 2.5 基于 Agent 的全链路动态更新机制

本文系统并不将规则与模型视为一次性配置，而是将其纳入 Agent 驱动的持续更新框架。系统更新能力覆盖前端采集、后端数据、离线 intelligence、研究验证和在线发布五个方面，目标是在不打断主业务链路的情况下实现小时级策略迭代。

#### 2.5.1 前端特征采集的动态注入

当前端面对新的伪装方法或新的攻击载体时，控制层可以推动 SDK 扩展新的采集项，例如新的环境一致性特征、新的底层渲染信号或新的交互上下文字段。借助 feature registry 和版本管理机制，新增特征不必立刻进入在线主表，而可以先进入研究和离线验证流程，在确认价值后再晋升为在线核心特征。

#### 2.5.2 后端数据库的实时对齐

系统通过统一事实层、画像层和规则层实现数据快速对齐。离线系统一旦发现新的高风险 hash、账号画像、设备画像或关系边，便可以通过标准化 publisher 写入在线可消费资产。同时，pipeline state 负责保障离线任务的增量执行与窗口重算，使 mock 触发、定时运行和异常恢复都能在统一机制下完成。这种数据库层对齐保证了 offline intelligence 能够稳定沉淀为 online execution assets。

#### 2.5.3 分析算法的闭环更新

在更高层面，research platform 与 Agent Control Layer 共同构成“发现 - 验证 - 发布”的闭环。新规则、新阈值和新模型先在 `research_validation_dataset` 上完成离线实验，再通过 shadow evaluation 评估其与当前线上版本的差异和潜在风险，最后由 control layer 按策略选择是否发布至在线。由此，系统不仅具备更新能力，更具备可验证、可回退、可治理的更新能力，从而使动态演进成为受控过程而非无序试验。

## 第三章 实验验证与分析

### 3.1 实验平台搭建

为验证本文系统的工程可行性与实际效果，研究团队将其部署于北航统一身份认证相关场景中进行观察与验证。实验平台由前端 SDK、Go 在线服务、Python 离线智能任务、共享数据库和 dashboard 展示层组成，能够覆盖真实登录请求、异常访问流量和多类型浏览器环境。整个实验体系既支持在线实时判定，也支持离线 intelligence 生成、影子评估与研究复现。

在实验组织方式上，平台同时保留生产执行链路与研究验证链路。前者用于承载在线检测、设备识别和 MFA 决策，后者则围绕 `research_validation_dataset`、离线结果快照和候选规则集合开展对比分析。通过这一设计，系统既能够观察真实场景中的运行效果，又能够在统一验证集上复现实验过程，从而保证评估结论具备可解释性与可重复性。

在评估协议上，本文将系统性能划分为四类指标：第一类为 Anti-Bot 检测指标，包括高风险流量命中率、误报率、漏报率以及多时间窗口下的告警稳定性；第二类为在线执行指标，包括判定延迟、配置刷新时延和高并发下的吞吐稳定性；第三类为设备识别与 RBA 指标，包括同设备跨周期关联能力、账号-设备关系一致性以及风险分层表现；第四类为演进能力指标，包括新规则从发现到发布的时间、影子评估周期和动态更新前后的效果变化。若条件允许，后续实验还将进一步围绕 `device-id v1` 与 `device-id v2` 的对比、关键模块消融以及动态更新前后效果差异开展系统性验证。

### 3.2 系统效能分析

从实际运行结果看，系统已累计处理数据 200 余万条，识别并锁定 Bot 流量 6000 余条，且能够将离线 intelligence 在小时级别内反馈到在线执行体系。这表明系统已经初步形成“在线检测、离线建模、动态发布”的完整闭环。

从论文分析视角看，上述结果至少说明三个问题。第一，系统具备现实业务中的持续运行能力，而非仅停留在实验室原型阶段；第二，离线 intelligence 已能够转化为在线可执行资产，说明 Agent 控制层具备实际价值；第三，浏览器指纹与设备关系分析在真实认证入口中具有足够的工程承载能力，为后续 device-id v2 和更高阶 RBA 提供了可扩展基础。尽管当前版本尚未给出完整的 precision / recall 或消融对比结果，但从业务指标与链路闭环状态来看，系统已经达到进入更精细研究评估阶段的前提条件。

#### 3.2.1 流量纯净度监控

系统能够持续跟踪真实业务流量中的正常访问与异常访问比例变化，辅助运维人员观察认证入口的整体纯净度趋势，从而评估攻击活动是否出现阶段性上升或规则更新是否取得预期效果。

#### 3.2.2 拦截战报与实时日志

在线系统能够记录每次 Anti-Bot 判定、设备身份解析和最终认证决策，并将其沉淀为结构化日志与战报信息，为后续追踪、回溯和运营分析提供基础。

#### 3.2.3 指纹特征分布与攻击者画像

通过对核心 hash 资产和设备关系的持续聚合，系统可以逐步勾勒攻击者设备画像，识别重复出现的风险设备、异常环境组合和跨账号活动模式，使防御视角从“单次请求”扩展为“持续实体”。

#### 3.2.4 风险评分与离线关联

离线侧通过多窗口特征和关系分析实现从单点事件判定到实体级风险表达的拓展。风险评分、设备画像和关系边的联合使用，使 online rule 与 RBA 决策能够获得更稳定的 intelligence 支撑。

进一步地，离线 risk score 并不直接替代在线判定，而是通过控制层转化为画像资产、阈值候选和规则候选进入在线系统。这种“离线建模、在线消费”的协同方式，既避免了将复杂模型直接压入低延迟链路，也确保了线上决策可以持续吸收离线高价值发现。

#### 3.2.5 交互性与可解释性验证

系统能够将检测依据、设备身份信息、环境一致性特征和风险画像以可视化方式展示给工程人员与业务管理者，从而增强判定结果的可理解性、可审计性和可运营性，这也是该方案区别于传统黑盒检测系统的重要特征。

此外，可解释性并不仅体现在 dashboard 展示层，还体现在系统的版本化与结构化设计中。例如，online rule、risk profile、device-id 版本和研究验证数据集版本均可被独立记录与回溯，这使得具体判定结果能够被映射到明确的规则版本、画像版本和身份解析版本之上。对于高校认证平台这类需要兼顾安全、管理和沟通的场景，这种解释能力具有直接实际价值。

## 结论

本文围绕 AI 时代登录认证场景中的自动化攻击、设备识别与动态风险控制问题，提出并实现了一套融合多元威胁情报的 Agent 驱动自适应 MFA 认证系统。该系统以浏览器指纹为核心基础，以 Anti-Bot 作为第一阶段门禁，以 Device-ID 与 RBA 作为第二阶段身份推理能力，并通过离线智能、研究验证和控制发布机制构建持续演进的 intelligence 闭环。

从方法与架构层面看，本文的主要价值在于将在线执行、离线 intelligence、研究平台和动态发布统一到同一工程体系中，使系统兼具实时性、可解释性、可演进性和研究可复现性。从实践层面看，系统已在北航统一身份认证相关场景中完成部署验证，展示了良好的工程可落地性与场景适配能力。

未来，随着 device-id v2、graph linking、shadow evaluation 和更丰富的 Agent 控制策略进一步成熟，该平台有望从高校认证场景扩展到更广泛的政企与金融安全场景，形成兼具研究创新和产品价值的新一代认证安全基础设施。进一步而言，本文所提出的体系并不限于某一具体业务入口，而是为“在线执行系统 + 离线 intelligence 系统 + 研究验证平台 + 动态发布控制层”的统一建模提供了可复制范式，具有进一步推广到更多身份安全与风控场景的潜力。

当然，当前工作仍然以规则驱动和工程可部署为主，围绕更高阶设备关系推理、反馈标签沉淀和自动化策略治理仍有进一步深化空间。后续研究可继续围绕更强的图结构建模、更稳定的 ground-truth 构造机制以及更精细的在线-离线协同发布流程展开，从而推动该体系从“可落地系统”进一步演进为“可持续学习平台”。

## 参考文献

1. Park, J. S., O'Brien, J., Cai, C. J., Morris, M. R., Liang, P., & Bernstein, M. S. Generative agents: Interactive simulacra of human behavior.
2. Eckersley, P. How unique is your web browser?
3. Venugopalan, H., Munir, S., Ahmed, S., Wang, T., King, S. T., & Shafiq, Z. Fp-inconsistent: Detecting evasive bots using browser fingerprint inconsistencies.
4. Hölbl, M., Zadorozhny, V., Welzer Družovec, T., Komapara, M., & Nemec Zlatolas, L. Browser Fingerprinting: Overview and Open Challenges.
5. Laperdrix, P., Rudametkin, W., & Baudry, B. Beauty and the beast: Diverting modern web browsers to build unique browser fingerprints.