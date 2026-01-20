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

- `src/` - 前端工具库源码，入口 `index.ts`，核心逻辑在 `utils.ts`，类型定义在 `types.ts`
- `server/go/` - Go SDK
- `server/python/dootask/` - Python SDK，核心类 `DooTaskClient`
- `server/node/src/` - Node.js SDK
- `example/` - Vue 3 + Vite 示例应用，用于本地验证

## 架构要点

### 微前端通信模式

前端工具通过 `window.microApp.getData()` 获取主应用注入的数据。所有对主应用的方法调用必须经由：
- `methodTryParent` - 底层通信封装
- 已有高层封装（`requestAPI`、`openWindow`、`closeApp`、`getUserInfo` 等）

**禁止在新代码中直接使用 `postMessage` 或自定义全局字段。**

### 运行环境判断

统一使用已有函数：`isMicroApp`、`isElectron`、`isEEUIApp`、`isIframe`。避免重复实现环境探测逻辑。

### 错误处理模式

- 环境检测方法（如 `isMicroApp`）返回 `boolean`，不抛异常
- API 方法在不支持的环境下抛出 `UnsupportedError`
- 网络请求失败抛出 `ApiError`

## 扩展指引

### 扩展前端 API

1. 在 `src/types.ts` 补充类型定义
2. 在 `src/utils.ts` 复用 `getAppData` / `methodTryParent` 实现新函数
3. 在 `README.md` 补充使用示例

### 扩展后端 SDK

在 `server/go`、`server/python`、`server/node` 中为同一后端接口建模，保持方法命名、请求参数与返回结构一致。

## 代码风格

- 使用 Prettier 格式化（无分号、双引号、2 空格缩进）
- 遵循现有 Promise / `error` 返回模式和命名风格
- 避免引入与现有栈不兼容的新依赖

## 语言偏好

- 总体说明与重要总结优先使用简体中文
- 如用户明确要求其他语言则遵从
- 完成改动后附带推荐的 Git 提交 message
