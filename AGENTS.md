# Repository Guidelines

该仓库为 DooTask 工具库与 SDK，包含前端微前端辅助工具以及 Go/Python/Node 后端客户端。

## 项目结构与核心入口

- 前端工具源码：`src/`，编译产物：`dist/`；入口为 `src/index.ts`，主要逻辑集中在 `src/utils.ts`，类型在 `src/types.ts`。
- 示例前端应用：`example/`，入口为 `example/src/main.ts` + `App.vue`，用于本地验证工具库（根目录可运行 `npm run example`）。
- Go SDK：`server/go/`，主要文件为 `types.go`、`utils.go`。
- Python SDK：`server/python/dootask/`，以 `DooTaskClient` 为核心类。
- Node SDK：源码位于 `server/node/src/`，通过 npm 包 `@dootask/tools` 中的 `DooTaskClient` 使用。

## 微前端交互与 API 约定

- 前端工具通过 `window.microApp.getData()` 获取主应用注入的数据，相关封装集中在 `src/utils.ts` 的 `appReady`、`getAppData`、`methodTryParent` 等方法。
- 对主应用的方法调用应统一经由 `methodTryParent` 或已存在的高层封装（如 `requestAPI`、`openWindow`、`closeApp`、`getUserInfo` 等），不要在新代码中直接使用 `postMessage` 或自定义全局字段。
- 运行环境判断统一使用已有函数：`isMicroApp`、`isElectron`、`isEEUIApp`、`isIframe` 等；避免重复实现环境探测逻辑。

## 扩展与修改指引

- 扩展前端 API 时：先在 `src/types.ts` 补充类型，然后在 `src/utils.ts` 复用 `getAppData` / `methodTryParent` 实现新函数，最后在根 `README.md` 补充使用示例。
- 扩展后端能力时：在 `server/go`、`server/python` 与 `server/node` 中分别为同一后端接口建模，保持方法命名、请求参数与返回结构尽量一致。
- 修改应尽量局部化，遵循现有 Promise / `error` 返回模式和命名风格，避免引入与现有栈不兼容的新依赖或框架。

## AI 回复风格与语言偏好

- 总体说明与重要总结（尤其是最终回答的 recap 部分），在不影响技术表达准确性的前提下，应优先使用简体中文进行回复。
- 如用户在对话中明确要求使用其他语言（例如英文），则以用户的显式指令为最高优先级。
- 当本次协作的改动已经较为完整且自然形成一个提交单元时，应在最终回答中附带一条或数条推荐的 Git 提交 message，方便用户直接复制使用。
