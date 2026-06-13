# Changelog

DooTask Tools 仓库的版本变更记录。包含前端工具库 `@dootask/tools`、命令行 `@dootask/cli`（含 doo 二进制）、Python SDK `dootask-tools`、Node SDK 与 Go SDK，统一版本号发布。

## [1.3.1]

`doo` 命令行围绕「消息处理、任务看板流转、工作报告 AI」补齐三条完整操作闭环，并修复一批接口参数错配问题；Go / Node / Python 三端 SDK 同步更新。

### Features

- 消息处理闭环：新增 `doo dialog inbox`（发现列表外仍有未读/待办的对话）、`doo dialog unread`（查看对话未读与提及统计）、`doo dialog read`（将对话标记为已读，可只标记某消息及之后）、`doo dialog mytodo`（列出我未完成的待办）、`doo message todolist`（列出消息的待办记录）与 `doo message todo-remind`（设置/取消待办提醒时间），让「找到待办 → 处理 → 标记完成」可在命令行一站完成
- 任务看板流转闭环：新增 `doo flow list`（查看项目工作流及各状态）、`doo task move`（跨项目/看板列移动任务或推进工作流状态）、`doo tag` 任务标签增删改查；`doo task update` 扩展 `--tag` / `--flow` / `--visibility`，支持直接为任务设置标签、工作流状态与可见性
- 工作报告 AI 闭环：新增 `doo report unread`（我收到的未读报告总数）、`doo report analyze`（将 AI 分析结果写回报告）、`doo report share`（把报告以分享链接发送到对话或成员）
- 任务对象新增 `task_tag` 标签字段，Go / Node / Python 三端 SDK 的 `ProjectTask` 同步暴露，便于读取任务已绑定的标签

### Bug Fixes

- 修复用户搜索按关键词（如「老李」）搜不到结果的问题：接口参数由 `keyword` 修正为后端实际期望的 `keys[key]`；工作报告列表的类型/状态/关键词筛选同样修正
- 修复消息搜索与「标记待办完成」：消息搜索改用正确的 `/api/search/message` 接口，完成待办按待办记录 ID（来自 `message todolist`）提交
- 修复文件搜索与查看：搜索关键词参数、文件详情的文本内容参数修正，`doo file view --content` 可正确返回文本内容
- 三端 SDK 请求统一补发 Version 头，修复后端接口版本校验失败的问题

### Documentation

- 同步 Go / Node / Python 三端 README 与示例：补充消息搜索、待办列表与完成待办等用法

## [1.3.0]

首次三包（`@dootask/tools` / `@dootask/cli` / `dootask-tools`）统一版本号发布，并接入 GitHub Actions 自动发布流程。

### Features

- 新增 `doo` 命令行工具：DooTask 的 gh 式 CLI，覆盖任务/项目/对话/消息/页面操作等常用操作，作为独立 npm 包 `@dootask/cli` 发布，支持 Linux/macOS/Windows × x64/arm64 五平台原生二进制
- `doo app` 命令：管理应用插件（AppStore），支持参数定义、必填校验、资源限额、sticky 等完整能力；`doo app catalog` 支持 `--search` 关键词模糊匹配 id / name / description / tags
- 新增 Node 版 `DooTaskClient` 与类型定义，与 Go、Python 三端 SDK 对齐方法命名与返回结构
- 新增 `callExtraEmitter` 方法，用于跨 micro-app 调用 `emitter.emit`
- 新增 `focusParentWindowIfIframe` 函数，便于在 iframe 场景中将焦点切回父窗口；`selectUsers` 等弹层流程已使用
- 新增微应用胶囊（capsule）配置能力
- 接入自动发布流程：GitHub Actions 一键完成 npm（`@dootask/tools` + `@dootask/cli` 主包与 5 个平台子包）、PyPI（`dootask-tools`）、GitHub Release（doo 五平台二进制）的统一发布

### Bug Fixes

- `doo app` 请求自动带上 Version 头，修复后端 `require_version` 校验失败的问题
- 调整 `MicroAppProps` 接口字段 `urlType` 为更通用的 `type`，并同步更新 `isIframe` 的引用

### Documentation

- 同步 `doo app` 命令文档：fields / --search / --param / sticky / 资源限额等参数完整说明
- 补充 `DooTaskClient` 客户端默认服务地址为 `http://nginx` 的说明
