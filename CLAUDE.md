# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

DooTask Tools 是一个微前端工具库和多语言 SDK 集合，用于 DooTask 应用开发。包含前端工具库（TypeScript）以及 Go、Python、Node.js 后端客户端。

## 常用命令

```bash
# 构建项目（编译前端工具库和 Node SDK）
npm run build

# 运行示例应用（Vue 3 + Vite）
npm run example

# 版本更新
npm version patch|minor|major

# 发布到 npm
npm publish
```

## 项目结构

- `src/` - 前端工具库（2.x），入口 `index.ts`
  - `channel.ts` - 通道：只负责收发协议信封（目前有 frame，electron / app / bridge 待补）
  - `client.ts` - 连接：握手、调用、事件、回调、应用令牌续期
  - `api.ts` - 对外接口，和 README「API 文档」一一对应
  - `errors.ts` - `DooTaskError` 和错误码
  - `types.ts` - 对外类型
- `server/go/` - Go SDK
- `server/python/dootask/` - Python SDK，核心类 `DooTaskClient`
- `server/node/src/` - Node.js SDK
- `example/` - Vue 3 + Vite 示例应用，配置成应用的窗口入口后在 DooTask 里逐个试接口

## 架构要点

### 通信协议

前端工具库和 DooTask 之间只走协议第 2 版的信封 `{dootask: 2, kind, ...}`，协议定义在主程序仓库 `docs/microapp-window/contract.md`：

- 加载后自动选通道并发 `hello`，收到 `welcome` 后就绪；`welcome` 带环境、状态（主题、语言、安全距离……）、可用接口列表和应用令牌
- 调用走 `call` / `result`；参数里的函数换成 `{__fn: id}`，DooTask 用 `callback` 调回
- DooTask 主动推的走 `event`，带 `id` 的事件要用 `reply` 回答（如 `before_close`）
- frame 通道只认父窗口发来的消息，`welcome` 之后只发给、只认 DooTask 的 origin

**新接口只在 `api.ts` 里用 `call` / `listen` 实现，不要直接 `postMessage`，也不要读主程序注入的全局变量。** `sdk.` 开头的方法是工具库内部用的，不对应用公开。

### 错误处理

- `ready()` 不在 DooTask 里抛 `unsupported`，协议版本不一致抛 `version_mismatch`；`isHosted()` / `isSupported()` 不抛错
- 调用失败统一抛 `DooTaskError`，`code` 是 `errors.ts` 里的错误码
- 应用令牌只在内存里，离过期不到 10 分钟时 `getAppToken()` 自动换新

## 扩展指引

### 扩展前端 API

1. 主程序宿主端先实现并登记到 capabilities，协议文档同步
2. 在 `src/types.ts` 补充类型，在 `src/api.ts` 用 `call` / `listen` 实现
3. 在 `README.md` 补充说明和示例，示例页加一个按钮

### 扩展后端 SDK

在 `server/go`、`server/python`、`server/node` 中为同一后端接口建模，保持方法命名、请求参数与返回结构一致。

## 代码风格

- 使用 Prettier 格式化（无分号、双引号、2 空格缩进）
- 遵循现有 Promise / `error` 返回模式和命名风格
- 避免引入与现有栈不兼容的新依赖

## 语言偏好

- 回复一律使用简体中文，除非用户明确要求其他语言
- 完成改动后附带推荐的 Git 提交 message
