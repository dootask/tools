# DooTask 应用工具库

[![npm version](https://img.shields.io/npm/v/@dootask/tools.svg?style=flat-square)](https://www.npmjs.com/package/@dootask/tools)
[![npm downloads](https://img.shields.io/npm/dm/@dootask/tools.svg?style=flat-square)](https://www.npmjs.com/package/@dootask/tools)
[![license](https://img.shields.io/npm/l/@dootask/tools.svg?style=flat-square)](https://github.com/dootask/tools/blob/main/LICENSE)

DooTask 应用的前端工具库。应用的窗口页面、扩展、外部打开的页面都用它和 DooTask 交互：拿用户信息和应用令牌、调用主程序接口、选人、弹提示、打开任务和会话等。

应用的配置写法（`config.yml` 的 `menu_items`、`extensions`、`permissions` 等）见 [DooTask 应用开发指南](https://appstore.dootask.com/development/manual)。

> 2.x 是重新设计的版本，接口和通信协议都和 1.x 不兼容，需要 DooTask 1.11.0 以上。从 1.x 升级见文末「从 1.x 升级」。

## 特点

- **一套接口，各端通用**：网页、桌面客户端、App 里的窗口页面，外部打开的页面，扩展，都用同一套接口。工具库自动选择通信方式。
- **不支持就报错**：当前环境不支持的接口一律抛出 `unsupported` 错误，不会静默失败。
- **类型完整**：完整的 TypeScript 类型定义。

## 安装

```bash
npm install @dootask/tools@^2
```

## 快速开始

### 窗口页面

```typescript
import { ready, getUserInfo, getTheme, onThemeChange, request, toast } from "@dootask/tools"

const env = await ready()                 // 等待和 DooTask 握手完成
const user = await getUserInfo()

document.documentElement.dataset.theme = await getTheme()
onThemeChange(theme => {
  document.documentElement.dataset.theme = theme
})

const projects = await request("project/lists", { data: { page: 1 } })
toast.success(`你好，${user.nickname}`)
```

### 扩展

```typescript
import { defineExtension, registerHandler, defineSlot } from "@dootask/tools"

export default defineExtension({
  activate() {
    registerHandler("convertToApprove", async ctx => {
      // ctx 里是挂载点的上下文，比如 {msg_id, dialog_id, msg_type}
    })
    defineSlot("task.detail.block", "okr-link", {
      mount(el, ctx) { /* 在 el 里渲染，ctx.task_id */ },
      update(ctx) {},
      unmount() {},
    })
  },
  deactivate() {
    // 清掉定时器、监听和其他副作用
  },
})
```

## 运行环境

工具库加载后自动判断所在的环境，选择通信方式（通道）：

| 通道 | 环境 |
| --- | --- |
| `frame` | DooTask 里的窗口页面（iframe），包括新标签打开、桌面客户端的应用窗口 |
| `inline` | 扩展，和 DooTask 在同一个页面里运行 |
| `electron` | 桌面客户端里外部打开的页面，由客户端注入通道 |
| `app` | App 里外部打开的页面和独立 WebView 页面，由 App 注入通道 |
| `bridge` | 浏览器里外部打开的页面，通过隐藏的 DooTask 桥接页通信 |

各组接口在不同通道里的支持情况：

| 接口组 | frame | inline | electron / app | bridge |
| --- | --- | --- | --- | --- |
| 环境 | ✓ | ✓ | ✓ | ✓ |
| 身份和请求 | ✓ | ✓ | ✓ | ✓ |
| 界面 | ✓ | ✓ | ✓ | ✓（桥接页临时铺满页面显示） |
| 操作主程序 | ✓ | ✓ | ✓ | ✓（转给用户已打开的 DooTask 桌面客户端） |
| 窗口 | ✓ | ✗ | ✓ | ✗ |
| 扩展 | ✗ | ✓ | ✗ | ✗ |

不确定某个接口能不能用时，先用 `isSupported(name)` 判断。

## 接口约定

- **调用类接口都返回 Promise**。
- **登记类接口是同步的**：`on*`、`register*`、`defineSlot` 立即返回一个取消登记的函数。
- **出错一律抛 `DooTaskError`**，用 `code` 区分（见「错误」）。
- **文字由应用自己翻译**：传给 `toast`、`dialog` 等接口的文字原样显示，请按 `getLanguage()` 翻译好再传。

## API 文档

### 环境

| 接口 | 返回 | 说明 |
| --- | --- | --- |
| `ready()` | `Promise<Env>` | 等待握手完成。不在 DooTask 里时抛 `unsupported`；协议版本不一致时抛 `version_mismatch` |
| `isHosted()` | `Promise<boolean>` | 是否运行在 DooTask 里，不抛错 |
| `getEnv()` | `Promise<Env>` | 当前环境，见下方类型 |
| `capabilities()` | `Promise<string[]>` | 当前环境支持的接口名列表 |
| `isSupported(name)` | `Promise<boolean>` | 某个接口能不能用，`name` 就是接口名，比如 `"popout"`、`"toast.success"` |
| `getSystemInfo()` | `Promise<SystemInfo>` | DooTask 的版本、标题、访问地址 |
| `getTheme()` | `Promise<"light" \| "dark">` | 当前主题 |
| `onThemeChange(cb)` | `() => void` | 主题变化 |
| `getLanguage()` | `Promise<Language>` | 当前语言 |
| `getLanguageList()` | `Promise<Record<Language, string>>` | 支持的语言和显示名 |
| `onLanguageChange(cb)` | `() => void` | 语言变化 |
| `getSafeArea()` | `Promise<SafeArea>` | 移动端顶部和底部的安全距离 |
| `onSafeAreaChange(cb)` | `() => void` | 安全距离变化 |

```typescript
interface Env {
  channel: "frame" | "inline" | "electron" | "app" | "bridge"
  platform: "web" | "electron" | "app"
  presentation: "page" | "window" | null   // 只有窗口有，其他为 null；手机上总是 "page"
  mobile: boolean
  protocol: number                                  // 通信协议版本，2.x 为 2
}

interface SystemInfo {
  version: string      // DooTask 版本
  title: string        // 系统名称
  baseUrl: string      // DooTask 的访问地址
}

type Language = "zh" | "zh-CHT" | "en" | "ko" | "ja" | "de" | "fr" | "id" | "ru"

interface SafeArea {
  top: number
  bottom: number
}
```

### 身份和请求

| 接口 | 返回 | 说明 |
| --- | --- | --- |
| `getAppToken()` | `Promise<string>` | 应用令牌，调用应用自己的后端时放在 `Authorization: Bearer` 里。令牌会自动续期，每次请求时取，不要长期保存 |
| `getUserInfo()` | `Promise<UserInfo>` | 当前用户 |
| `fetchUsers(ids)` | `Promise<UserBasic[]>` | 一批用户的基本信息 |
| `request(path, options?)` | `Promise<any>` | 调用主程序接口，返回接口的 `data`。只能调用 `permissions` 里声明过的接口 |

```typescript
interface UserInfo {
  userid: number
  nickname: string
  email: string
  userimg: string         // 头像地址
  identity: string[]      // 身份，比如 ["admin"]
  department: number[]    // 部门 ID
  profession: string
  tel: string
  bot: boolean
  created_at: string
}

interface UserBasic {
  userid: number
  nickname: string
  email: string
  userimg: string
  department: number[]
  profession: string
  bot: boolean
  online: boolean
  line_at: string         // 最近在线时间
}

interface RequestOptions {
  method?: "GET" | "POST" // 默认 GET
  data?: Record<string, any>
  headers?: Record<string, string>
}
```

- `path` 是主程序接口路径，不带 `api/` 前缀，比如 `"project/lists"`。
- 接口返回失败时抛 `api_error`，`error.data` 是 `{ret, msg, data}`；没有权限时抛 `permission_denied`。

### 界面

| 接口 | 返回 | 说明 |
| --- | --- | --- |
| `toast.success(text)` / `toast.error(text)` / `toast.warning(text)` / `toast.info(text)` | `Promise<void>` | 轻提示 |
| `dialog.alert(options)` / `dialog.info(options)` / `dialog.success(options)` / `dialog.warning(options)` / `dialog.error(options)` | `Promise<void>` | 提示框，关闭后 resolve |
| `dialog.confirm(options)` | `Promise<boolean>` | 确认框，确定返回 `true`，取消返回 `false` |
| `selectUsers(options?)` | `Promise<{users: number[], groups: number[]}>` | 选人，用户关掉选择器时抛 `cancelled` |
| `selectImage(options?)` | `Promise<{url: string, path: string}>` | 选一张图片：打开 DooTask 的图片空间，可以点选上传过的图片，也可以上传新图片；用户关掉时抛 `cancelled` |
| `previewImage(urls, index?)` | `Promise<void>` | 预览图片 |
| `download(url)` | `Promise<void>` | 下载文件 |

```typescript
interface DialogOptions {
  title: string
  content?: string        // 纯文本
  okText?: string
  cancelText?: string
  width?: number
}

interface SelectUsersOptions {
  selected?: number[]                 // 已选的用户
  locked?: number[]                   // 已选、不能取消的用户
  disabled?: number[]                 // 不能选的用户
  max?: number                        // 最多选几个
  projectId?: number                  // 只在这个项目的成员里选
  dialogId?: number                   // 只在这个会话的成员里选
  includeBots?: boolean               // 显示机器人，默认 false
  includeDisabled?: boolean           // 显示已停用的用户，默认 false
  includeGroups?: boolean             // 同时可以选群组，默认 false
  onlyGroups?: boolean                // 只选群组，默认 false
  title?: string
  placeholder?: string
  beforeSubmit?: (ids: {users: number[], groups: number[]}) => boolean | Promise<boolean>  // 返回 false 阻止提交
}

interface SelectImageOptions {
  width?: number                      // 上传新图片时缩放到的宽度，不传按 DooTask 默认
  height?: number                     // 上传新图片时缩放到的高度
  crop?: boolean                      // true 时上传的图片按宽高裁满（如头像），默认等比缩小
}
```

```typescript
// 换头像：选好后拿到图片地址，用户关掉时什么也不做
try {
  const { url } = await selectImage({ width: 512, height: 512, crop: true })
  setAvatar(url)
} catch (error) {
  if (!(error instanceof DooTaskError && error.code === "cancelled")) throw error
}
```

### 操作主程序

| 接口 | 返回 | 说明 |
| --- | --- | --- |
| `openTask(id)` | `Promise<void>` | 打开任务 |
| `addTask(params?)` | `Promise<void>` | 打开新建任务表单并预填。`params`：`{projectId?: number, columnId?: number, name?: string, content?: string, times?: [string, string], owner?: number[]}`，`content` 是 HTML，`times` 是开始、结束时间（`YYYY-MM-DD HH:mm`） |
| `openDialog(id)` | `Promise<void>` | 打开会话 |
| `openDialogWithUser(userid)` | `Promise<void>` | 打开和某个用户的单聊 |
| `createGroup(userids)` | `Promise<void>` | 打开创建群组，预选这些成员 |
| `startMeeting(params?)` | `Promise<void>` | 打开新会议。`params`：`{userids?: number[], dialogId?: number}`，预选成员、会议消息发到哪个会话 |
| `joinMeeting()` | `Promise<void>` | 加入会议 |
| `openSearch(keyword?)` | `Promise<void>` | 打开全局搜索 |
| `openAI(params?)` | `Promise<void>` | 打开 AI 助手。`params`：`{prompt?: string, send?: boolean}`，`prompt` 填进输入框，`send: true` 直接发送 |
| `openApp(appid, key, query?)` | `Promise<void>` | 打开某个应用的窗口入口，`key` 是 `menu_items` 的 `key`，`query` 拼到地址上 |
| `getHostPreference(key)` / `setHostPreference(key, value)` | `Promise<any>` / `Promise<void>` | 读写 DooTask 的偏好设置，只开放登记过的键，目前只有 `ai.float_button` |
| `getTaskFields(taskId)` | `Promise<Record<string, any>>` | 读本应用在 `task.create.field` 里声明的字段值 |
| `setTaskFields(taskId, values)` | `Promise<void>` | 写本应用声明的字段值 |

浏览器里外部打开的页面（`bridge` 通道）调用这一组接口时，会转给用户已打开的 DooTask 桌面客户端执行；用户没有打开桌面客户端时，在新标签打开 DooTask 并执行（`openAI` 只填进输入框，不直接发送）。这个通道里没有 `openApp`、`getHostPreference`、`setHostPreference`。

### 窗口

只在窗口页面里可用（`frame`、`electron`、`app` 通道），扩展和浏览器里外部打开的页面调用时抛 `unsupported`。

| 接口 | 返回 | 说明 |
| --- | --- | --- |
| `close()` | `Promise<void>` | 关闭当前窗口 |
| `onBeforeClose(cb)` | `() => void` | 关闭前询问：`cb` 返回 `false` 或一段提示文字时，DooTask 弹确认框「(提示文字)，确定关闭？」，用户选「关闭」才关；返回其他值、抛错、2 秒内没回复都直接关闭。自己调用 `close()` 不会问 |
| `onNativeBack(cb)` | `() => void` | 手机上的原生返回（Android 返回键、iOS 侧滑、手机浏览器返回）交给应用处理。默认原生返回会把应用收成浮窗；登记期间改为调用最后登记的 `cb`（后进先出），全部取消后恢复默认。电脑上不会触发 |
| `popout()` | `Promise<void>` | 弹出到浏览器新标签或桌面客户端的独立窗口，要求入口配置了 `popout: true` |
| `openWindow(url, options?)` | `Promise<void>` | 打开一个外部窗口：桌面客户端里是独立窗口，其他端是浏览器新标签 |
| `getPresentation()` | `Promise<"page" \| "window">` | 当前的打开方式（手机上总是 `page`） |
| `onPresentationChange(cb)` | `() => void` | 用户在胶囊菜单里改了打开方式 |
| `getCapsuleRect()` | `Promise<Rect>` | 胶囊的位置和大小（`top` 已含安全距离，`right` 是到窗口右边的距离），沉浸式页面用来给胶囊留位置；非沉浸式和外部页面全为 0 |
| `onCapsuleChange(cb)` | `() => void` | 胶囊位置变了（换打开方式、转屏、安全距离变化），`cb(rect)` |
| `setCapsuleMenus(menus)` | `Promise<void>` | 替换胶囊「···」里应用自己的菜单项 |
| `onCapsuleMenu(cb)` | `() => void` | 胶囊菜单项被点击，`cb(value)` |

```typescript
interface Rect {
  top: number
  right: number         // 胶囊右边到窗口右边的距离
  width: number
  height: number
}

interface CapsuleMenu {
  label: string
  value: string | number
}

interface OpenWindowOptions {
  title?: string
  width?: number
  height?: number
}
```

### 扩展

只在扩展里可用（`inline` 通道）。

| 接口 | 返回 | 说明 |
| --- | --- | --- |
| `defineExtension({activate, deactivate})` | 扩展对象 | 扩展入口的默认导出 |
| `registerHandler(name, fn)` | `() => void` | 登记 `action: {call: name}` 的处理函数，`fn(context)` |
| `defineSlot(point, key, slot)` | `() => void` | 渲染型挂载点，`slot` 见下方类型 |
| `registerSearch(key, fn)` | `() => void` | 搜索源：`fn(keyword)` 返回搜索结果，DooTask 会设超时 |
| `getOverlayRoot()` | `ShadowRoot` | 本扩展的浮层根节点，全屏层、弹窗、下拉菜单都挂到这里 |
| `openLayer(layer)` | `Promise<() => void>` | 在浮层根节点里打开一个全屏层，resolve 为关闭函数 |
| `showBanner(banner)` | `Promise<() => void>` | 显示顶部横幅，resolve 为关闭函数 |
| `nextZIndex()` | `Promise<number>` | 在浮层根节点里叠多层时取层级 |
| `onHostEvent(name, cb)` | `() => void` | 监听 DooTask 的事件，给后台常驻扩展用。第一批事件：`message.new`（新消息）、`task.update`（任务变更） |

```typescript
interface Extension {
  activate(): void | Promise<void>      // 超过 10 秒没完成或抛错，这个扩展的挂载点都不显示
  deactivate?(): void | Promise<void>   // 应用卸载或更新时调用
}

interface Slot {
  mount(el: HTMLElement, ctx: SlotContext): void | Promise<void>  // el 在一个 Shadow Root 里
  update?(ctx: SlotContext): void       // 上下文变化，比如切换到另一个任务
  unmount?(): void
}

type SlotContext = Record<string, any>  // 比如 {task_id, project_id}，字段见开发指南的挂载点表

interface SearchResult {
  title: string
  desc?: string
  icon?: string
  action: {open: string} | {call: string}
}

interface Layer {
  mount(el: HTMLElement): void | Promise<void>
  unmount?(): void
}

interface Banner {
  text: string
  tone?: "success" | "danger" | "warning" | "info" | "muted"
  action?: {label: string} & ({open: string} | {call: string})  // 横幅上的按钮
}
```

扩展的几条规则：

- `mount(el, ctx)` 里的 `el` 在 Shadow Root 里：样式只在里面生效，要注入到 Shadow Root，不要插到 `document.head`。DooTask 的主题 CSS 变量能穿进去。
- 组件库的弹窗、下拉要挂到 `getOverlayRoot()`，比如 Vue 的 `<Teleport :to="overlayRoot">`。
- 不能修改 `window` 上的属性、全局样式和原型链；框架和依赖打包进自己的入口。
- 定时器、监听在 `deactivate` 里清掉。

### 错误

所有错误都是 `DooTaskError`，用 `code` 区分：

| `code` | 含义 |
| --- | --- |
| `unsupported` | 当前环境不支持这个接口 |
| `cancelled` | 用户取消，比如关掉了选人 |
| `permission_denied` | 应用没有在 `permissions` 里声明这个主程序接口 |
| `api_error` | 主程序接口返回失败，`data` 是 `{ret, msg, data}` |
| `invalid_args` | 参数不对 |
| `timeout` | 超时 |
| `version_mismatch` | 工具库和 DooTask 的通信协议版本不一致，需要升级其中一方 |
| `auth_failed` | 启动码无效或者应用令牌已失效 |

```typescript
import { selectUsers, DooTaskError } from "@dootask/tools"

try {
  const { users } = await selectUsers({ max: 5 })
} catch (e) {
  if (e instanceof DooTaskError && e.code === "cancelled") return
  throw e
}
```

## 使用示例

### 判断能不能用

```typescript
import { isHosted, isSupported, popout } from "@dootask/tools"

if (!(await isHosted())) {
  // 不在 DooTask 里，比如开发时直接打开了页面
}

if (await isSupported("popout")) {
  await popout()
}
```

### 关闭前确认

```typescript
import { onBeforeClose } from "@dootask/tools"

let dirty = false

const off = onBeforeClose(() => {
  // 返回提示文字：DooTask 弹「有未保存的内容，确定关闭？」，用户选「关闭」才关
  if (dirty) return "有未保存的内容"
})

// 不再需要时取消登记
// off()
```

### 手机返回先关自己的弹层

```typescript
import { onNativeBack } from "@dootask/tools"

// 打开弹层时登记，关掉时取消：登记期间手机返回只关弹层，不会把应用收起来
let offBack: (() => void) | null = null

function openSheet() {
  showSheet.value = true
  offBack = onNativeBack(closeSheet)
}

function closeSheet() {
  showSheet.value = false
  offBack?.()
  offBack = null
}
```

### 沉浸式页面给胶囊让位

胶囊的位置由 DooTask 决定，手机、电脑整页、窗口各不相同，应用的顶栏按它对齐：右边留出胶囊的宽度，顶栏里的按钮和胶囊上下居中对齐。

```typescript
import { getCapsuleRect, getSafeArea, onCapsuleChange, onSafeAreaChange } from "@dootask/tools"

async function layoutHeader() {
  const [capsule, safe] = await Promise.all([getCapsuleRect(), getSafeArea()])
  const header = document.querySelector<HTMLElement>(".app-header")!
  header.style.paddingTop = `${safe.top}px`
  if (!capsule.width) {
    return // 非沉浸式，胶囊不在页面里
  }
  // 胶囊上下留出同样的空隙，顶栏内容上下居中就和胶囊对齐
  header.style.height = `${capsule.top * 2 - safe.top + capsule.height}px`
  header.style.paddingRight = `${capsule.right + capsule.width + 12}px`
}

layoutHeader()
onCapsuleChange(layoutHeader)
onSafeAreaChange(layoutHeader)
```

### 胶囊菜单

```typescript
import { setCapsuleMenus, onCapsuleMenu } from "@dootask/tools"

await setCapsuleMenus([{ label: "使用帮助", value: "help" }])

onCapsuleMenu(value => {
  if (value === "help") showHelp()
})
```

### 调用应用自己的后端

```typescript
import { getAppToken } from "@dootask/tools"

const res = await fetch("/apps/example/api/list", {
  headers: { Authorization: `Bearer ${await getAppToken()}` },
})
```

后端拿到令牌后，用后端工具库（或者直接请求主程序的 `users/info`，同样放在 `Authorization: Bearer` 里）反查当前用户。

### 外部打开的页面

页面照常调用 `ready()` 即可。浏览器里，DooTask 打开页面时会在地址上带 `dootask_launch`（启动码，60 秒内有效、只能用一次）和 `dootask_origin`（DooTask 的地址），工具库自动取走并从地址栏删掉，再插入一个隐藏的 DooTask 桥接页换取应用令牌：

- 页面所在的域名要登记在应用 `config.yml` 的 `origins`（或 `menu_items` 的地址）里，否则 DooTask 不签发启动码，桥接页也不允许被嵌入；
- 页面内跳转、刷新后不需要启动码，桥接页沿用上次换到的凭证（7 天内用过就一直有效，用户退出 DooTask 后失效）；
- 换不到凭证时 `ready()` 抛 `auth_failed`，请提示用户从 DooTask 重新打开应用；
- 弹框、选人、图片预览由桥接页临时铺满页面显示，只有轻提示时只占页面顶部，铺满时键盘焦点在桥接页里（Esc 可关闭弹框），关闭后还给页面原来的元素；平时隐藏、不挡点击；
- 页面不能用 `Content-Security-Policy` 的 `frame-src` / `child-src` 拦掉 DooTask 的地址。

## 从 1.x 升级

2.x 不兼容 1.x，接口对照如下：

| 1.x | 2.x |
| --- | --- |
| `appReady()` | `ready()` |
| `isMicroApp()`、`isEEUIApp()`、`isElectron()`、`isMainElectron()`、`isSubElectron()`、`isFullScreen()`、`isIframe()`、`getWindowType()` | `isHosted()`、`getEnv()`；具体能力用 `isSupported(name)` 判断 |
| `getThemeName()` | `getTheme()`、`onThemeChange(cb)` |
| `getLanguageName()` | `getLanguage()`、`onLanguageChange(cb)` |
| `getUserId()` | `getUserInfo()` 的 `userid` |
| `getUserToken()` | `getAppToken()`（应用令牌，不再提供用户的登录 token） |
| `getBaseUrl()` | `getSystemInfo()` 的 `baseUrl` |
| `requestAPI({url, ...})` | `request(path, options)`，只能调用 `permissions` 里声明的接口 |
| `fetchUserBasic(ids)` | `fetchUsers(ids)` |
| `selectUsers({value, uncancelable, disabledChoice, multipleMax, showBot, showDisable, showDialog, onlyGroup})` | `selectUsers({selected, locked, disabled, max, includeBots, includeDisabled, includeGroups, onlyGroups})`，返回 `{users, groups}`，取消时抛 `cancelled` |
| `modalSuccess()` 等 `modalXxx` | `dialog.success()` 等，`modalConfirm` 改为 `dialog.confirm` |
| `messageSuccess()` 等 `messageXxx` | `toast.success()` 等 |
| `closeApp()`、`backApp()` | `close()`，返回由 DooTask 处理 |
| `interceptBack(cb)`（返回 `true` 阻止） | 关闭前确认用 `onBeforeClose(cb)`（返回 `false` 或提示文字）；拦截手机返回关掉自己的弹层用 `onNativeBack(cb)` |
| `popoutWindow(params)` | `popout()` |
| `openWindow(params)`、`openTabWindow(url)` | `openWindow(url, options)` |
| `openAppPage(params)` | `openApp(appid, key, query)` |
| `openDialogUserid(userid)` | `openDialogWithUser(userid)` |
| `openDialogNewWindow(id)` | 删除，用 `openDialog(id)` |
| `downloadUrl(url)` | `download(url)` |
| `setCapsuleConfig(config)`、`addMenuClickListener(cb)` | `setCapsuleMenus(menus)`、`onCapsuleMenu(cb)`；胶囊的显隐和位置由 DooTask 决定 |
| `nextZIndex()` | 窗口里删除；扩展里保留 |
| `callExtraA()`、`callExtraStore()`、`callExtraEmitter()` | 删除，改用上面的显式接口（比如 `previewImage`、`openTask`） |
| `addDataListener()`、`removeDataListener()` | 删除，用各个 `on*Change` |
| `initSlideBack()` | 删除，侧滑返回由 DooTask 处理 |
| `UnsupportedError`、`ApiError` | `DooTaskError`，用 `code` 区分 |

配置文件也要一起改：`menu_items` 的 `type`、`immersive`、`width`、`{user_token}` 这类旧字段和占位符都已删除，见开发指南。

## 后端服务

DooTask 提供多语言后端工具库，方便在服务端调用主程序接口、用应用令牌校验用户：

- [Go](server/go)
- [Python](server/python)
- [Node](server/node)

## 许可证

MIT
