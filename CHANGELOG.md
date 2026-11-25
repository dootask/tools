# Changelog

DooTask Tools 仓库的版本变更记录。包含前端工具库 `@dootask/tools`、命令行 `@dootask/cli`（含 doo 二进制）、Python SDK `dootask-tools`、Node SDK 与 Go SDK，统一版本号发布。

## [2.3.0]

前端工具库 `onNativeBack` 支持控制侧滑返回指示器，返回的默认语义统一为「离开」。

### Features

- `onNativeBack(cb, options?)` 新增 `options.indicator`（默认 `true`）：侧滑返回时是否跟手画返回指示器（iOS App、手机浏览器；Android 由系统绘制），多次登记以最后登记的为准。返回的默认行为随之明确为「离开应用（转到后台，状态保留）」，并补充电脑窗口里按 Esc 也算返回

### Documentation

- `onBeforeClose` 说明澄清：离开（胶囊 ×、返回）不会询问；后台应用过多需要挤掉最早的时也会询问，但被拦住时不弹确认框，只保留原应用

## [2.2.0]

前端工具库新增 `setCapsulePosition`，沉浸式页面可以临时移动胶囊。

### Features

- 新增 `setCapsulePosition({top, right} | null)`：进到子页面、打开抽屉时把胶囊移到别处（有过渡动画），缺的一边用默认位置，`null` 恢复默认，返回实际位置；超出页面时贴边，换打开方式、转屏、重新加载后自动恢复默认。需要 DooTask 主程序支持，可以先用 `isSupported("setCapsulePosition")` 判断

## [2.1.0]

前端工具库新增 `selectImage`，可直接从 DooTask 的图片空间选图。

### Features

- 新增 `selectImage({width, height, crop})`：打开 DooTask 的图片空间选一张图片（点选上传过的，或上传新图片，上传时按给的尺寸缩放、`crop` 为 `true` 时裁满），返回 `{url, path}`，用户关掉时抛 `cancelled`。需要 DooTask 主程序支持，可以先用 `isSupported("selectImage")` 判断

## [2.0.0]

前端工具库按新的应用通信协议（第 2 版）重写，不兼容 1.x；服务端 SDK 支持应用令牌。需要 DooTask 主程序支持协议第 2 版。

### Breaking Changes

- 不再依赖 micro-app：应用的窗口页面都在 iframe 里打开，通过 `postMessage` 和 DooTask 握手（`hello` / `welcome`），1.x 的 `window.microApp`、`methodTryParent`、`getAppData` 等全部移除
- 接口重新设计，见 README「API 文档」：调用类接口返回 Promise，登记类接口（`on*`）返回取消函数；`popoutWindow` → `popout`，`closeApp` → `close`，`interceptBack` → `onBeforeClose`，`requestAPI` → `request`（返回接口的 `data`，只能调用应用 `permissions` 里声明过的接口），`selectUsers` 参数改为 `selected` 并返回 `{users, groups}`；移除 `isMicroApp`、`isElectron`、`isEEUIApp` 等环境探测函数，改用 `ready()` / `getEnv()` / `isSupported()`；窗口里不再需要 `nextZIndex`（只在扩展里提供）
- 不再提供用户登录令牌：应用调用自己的后端时用 `getAppToken()` 取应用令牌，放在 `Authorization: Bearer` 里
- 调用失败统一抛 `DooTaskError`（带 `code`）

### Features

- 新增窗口相关接口：`getPresentation` / `onPresentationChange`、`getCapsuleRect` / `onCapsuleChange`、`setCapsuleMenus` / `onCapsuleMenu`、`openWindow`、`openApp`
- 外部打开的页面也能用工具库：桌面客户端、App 里由客户端注入通道（`electron` / `app`），浏览器里取走地址上的启动码（`dootask_launch`），经隐藏的 DooTask 桥接页换取应用令牌（`bridge`）
- 主题、语言、安全距离握手时就带过来，`getTheme` 等直接返回，变化时通过 `onThemeChange` 等通知
- Go / Node / Python 三端 SDK：`token` 传应用令牌（`dta_` 开头）时用 `Authorization: Bearer` 发送，应用后端可以直接拿页面请求里的应用令牌调用主程序接口；个人访问令牌（`dtp_` 开头）同样用 `Authorization: Bearer` 发送

## [1.3.6]

`doo` 命令行文件操作补齐增删改与历史版本能力。

### Features

- `doo file` 补齐文件 CRUD：新增 `create`（新建文件/文件夹，`--type` 指定类型）、`save`（保存文本内容，支持 `--content` 直传或 `--file` 从本地文件读取）、`rename`（重命名）、`delete`（可批量删除）、`history`（查看文件内容的修改历史）、`restore`（恢复到指定历史版本）。此前 `file` 仅支持 list/search/view/fetch，现覆盖完整文件生命周期；`delete`/`restore` 等写操作带 `--yes` 确认护栏，非交互环境不加则拒绝执行

### Documentation

- `@dootask/cli` 新增独立 README 的「For AI Agents」章节：提供可直接粘贴进 system prompt 的接入片段（环境变量认证、`--json` 机器可读输出、退出码约定、危险操作 `--yes`），方便把 `doo` 作为工具交给 LLM Agent 调用；同步更新命令一览中 `file` 的子命令列表

## [1.3.5]

Go / Node / Python 三端 SDK 新增设置应用菜单角标能力。

### Features

- 三端 SDK 新增 `setAppBadge`：通过 `POST /api/apps/badge/set` 设置应用在 DooTask 菜单上的角标。以 `appid + secret`（默认读取环境变量 `APP_ID` / `APP_SECRET`，也可在构造时覆盖）鉴权；`userid` 支持单个或数组、`count`（数字角标）/ `dot`（小红点）/ `menu_key`（指定菜单项）均为可选。每次调用为绝对覆盖、幂等，方便从插件后端主动提醒用户。三端 README 同步补充用法

## [1.3.4]

`doo` 命令行新增「调用应用自定义指令」能力。

### Features

- 新增 `doo app call <应用ID> [指令]`：基于应用暴露的 OpenAPI/Swagger 规范，直接调用插件自己的后端接口。无指令时列出该应用全部可用指令；`--describe` 查看某指令的参数；`K=V`（或 `--arg K=V`）传参、`--data '<json>'` 传完整请求体、`--raw <METHOD> <path>` 调用未在规范中声明的端点、`--refresh` 跳过规范缓存。调用经 AppStore 内网代理到应用容器，自动透传当前用户身份、由应用后端自鉴权，无需关心容器地址、端口与鉴权细节

### Miscellaneous

- 构建产物语法目标从 ES5 升级到 ES2017：`@dootask/tools` 不再向下转译为 ES5，贴合 Electron 与现代浏览器运行环境，同时规避 TypeScript 7.0 对 ES5 target 的弃用警告（如需兼容老旧引擎请自行转译）

## [1.3.3]

`doo` 命令行支持上传本地应用压缩包。

### Features

- 新增 `doo app upload <压缩包路径>`：把本地应用压缩包（`.zip` / `.tar.gz`）上传导入应用市场，等同网页的「上传本地应用」。仅做导入与合规校验、不自动部署，导入成功后用 `doo app install <ID>` 安装；`--appid` 留空时按文件名自动识别应用 ID（去掉结尾版本段与扩展名，如 `myapp-1.0.0.tar.gz` → `myapp`），文件名与目标 ID 不一致时可显式指定

## [1.3.2]

`doo` 命令行补齐「应用升级」查询能力。

### Features

- 新增 `doo app updates`：一条命令列出**已安装且有新版本可升级**的应用（含 community 应用），输出应用 ID、名称、当前已装版本与可升级到的最新版本。与网页应用市场的「可升级」标记同源——此前只有 `doo app list` 查看已装版本，无法判断哪些应用有新版可升，现在「有哪些应用可以更新」可直接得到答案

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
