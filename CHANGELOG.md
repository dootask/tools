# Changelog

DooTask Tools 仓库的版本变更记录。包含前端工具库 `@dootask/tools`、命令行 `@dootask/cli`（含 doo 二进制）、Python SDK `dootask-tools`、Node SDK 与 Go SDK，统一版本号发布。

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
