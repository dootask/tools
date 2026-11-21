/**
 * 对外接口，见 README「API 文档」
 * - 调用类接口返回 Promise；登记类接口（on*）同步返回取消函数
 * - 文字原样显示，由应用自己翻译
 */
import { appToken, call, listen, listenerCount, ready as connect } from "./client"
import { DooTaskError } from "./errors"
import type {
  AddTaskParams,
  CapsuleMenu,
  DialogOptions,
  Env,
  Language,
  MeetingParams,
  OpenAIParams,
  OpenWindowOptions,
  Presentation,
  Rect,
  RequestOptions,
  SafeArea,
  SelectedIds,
  SelectUsersOptions,
  SystemInfo,
  Theme,
  UserBasic,
  UserInfo,
} from "./types"

async function state<T>(name: string): Promise<T> {
  return (await connect()).state[name] as T
}

// ---------- 环境 ----------

/** 等待握手完成。不在 DooTask 里时抛 unsupported；协议版本不一致时抛 version_mismatch */
export async function ready(): Promise<Env> {
  return { ...(await connect()).env }
}

/** 是否运行在 DooTask 里，不抛错 */
export async function isHosted(): Promise<boolean> {
  try {
    await connect()
    return true
  } catch {
    return false
  }
}

export async function getEnv(): Promise<Env> {
  return ready()
}

/** 当前环境支持的接口名列表 */
export async function capabilities(): Promise<string[]> {
  return [...(await connect()).capabilities]
}

/** 某个接口能不能用，name 是接口名，比如 "popout"、"toast.success" */
export async function isSupported(name: string): Promise<boolean> {
  try {
    return (await connect()).capabilities.includes(name)
  } catch {
    return false
  }
}

export async function getSystemInfo(): Promise<SystemInfo> {
  return { ...(await state<SystemInfo>("system")) }
}

export function getTheme(): Promise<Theme> {
  return state<Theme>("theme")
}

export function onThemeChange(cb: (theme: Theme) => void): () => void {
  return listen("theme", cb)
}

export function getLanguage(): Promise<Language> {
  return state<Language>("language")
}

export async function getLanguageList(): Promise<Record<Language, string>> {
  return { ...(await state<Record<Language, string>>("languageList")) }
}

export function onLanguageChange(cb: (language: Language) => void): () => void {
  return listen("language", cb)
}

export async function getSafeArea(): Promise<SafeArea> {
  return { top: 0, bottom: 0, ...(await state<SafeArea | null>("safe_area")) }
}

export function onSafeAreaChange(cb: (area: SafeArea) => void): () => void {
  return listen("safe_area", cb)
}

// ---------- 身份和请求 ----------

/** 应用令牌，调用应用自己的后端时放在 Authorization: Bearer 里。每次请求时取，不要长期保存 */
export function getAppToken(): Promise<string> {
  return appToken()
}

export function getUserInfo(): Promise<UserInfo> {
  return call("getUserInfo")
}

export function fetchUsers(ids: number[]): Promise<UserBasic[]> {
  return call("fetchUsers", ids)
}

/** 调用主程序接口，返回接口的 data。只能调用 permissions 里声明过的接口 */
export function request<T = any>(path: string, options: RequestOptions = {}): Promise<T> {
  return call("request", path, { ...options, method: options.method === "POST" ? "post" : "get" })
}

// ---------- 界面 ----------

export const toast = {
  success: (text: string): Promise<void> => call("toast.success", text),
  error: (text: string): Promise<void> => call("toast.error", text),
  warning: (text: string): Promise<void> => call("toast.warning", text),
  info: (text: string): Promise<void> => call("toast.info", text),
}

export const dialog = {
  alert: (options: DialogOptions): Promise<void> => call("dialog.alert", options),
  info: (options: DialogOptions): Promise<void> => call("dialog.info", options),
  success: (options: DialogOptions): Promise<void> => call("dialog.success", options),
  warning: (options: DialogOptions): Promise<void> => call("dialog.warning", options),
  error: (options: DialogOptions): Promise<void> => call("dialog.error", options),
  /** 确定返回 true，取消返回 false */
  confirm: async (options: DialogOptions): Promise<boolean> => (await call("dialog.confirm", options)) === true,
}

/** 选人，用户关掉选择器时抛 cancelled */
export function selectUsers(options: SelectUsersOptions = {}): Promise<SelectedIds> {
  return call("selectUsers", options)
}

export function previewImage(urls: string | string[], index = 0): Promise<void> {
  return call("previewImage", Array.isArray(urls) ? urls : [urls], index)
}

export function download(url: string): Promise<void> {
  return call("download", url)
}

// ---------- 操作主程序 ----------

export function openTask(id: number): Promise<void> {
  return call("openTask", id)
}

export function addTask(params: AddTaskParams = {}): Promise<void> {
  return call("addTask", params)
}

export function openDialog(id: number): Promise<void> {
  return call("openDialog", id)
}

export function openDialogWithUser(userid: number): Promise<void> {
  return call("openDialogWithUser", userid)
}

export function createGroup(userids: number[]): Promise<void> {
  return call("createGroup", userids)
}

export function startMeeting(params: MeetingParams = {}): Promise<void> {
  return call("startMeeting", params)
}

export function joinMeeting(): Promise<void> {
  return call("joinMeeting")
}

export function openSearch(keyword?: string): Promise<void> {
  return call("openSearch", keyword)
}

export function openAI(params: OpenAIParams = {}): Promise<void> {
  return call("openAI", params)
}

/** 打开某个应用的窗口入口，key 是 menu_items 的 key，query 拼到地址上 */
export function openApp(appid: string, key: string, query?: Record<string, string | number | boolean>): Promise<void> {
  return call("openApp", appid, key, query)
}

/** 读 DooTask 的偏好设置，目前只有 ai.float_button */
export function getHostPreference<T = any>(key: string): Promise<T> {
  return call("getHostPreference", key)
}

export function setHostPreference(key: string, value: any): Promise<void> {
  return call("setHostPreference", key, value)
}

export function getTaskFields(taskId: number): Promise<Record<string, any>> {
  return call("getTaskFields", taskId)
}

export function setTaskFields(taskId: number, values: Record<string, any>): Promise<void> {
  return call("setTaskFields", taskId, values)
}

// ---------- 窗口 ----------

export function close(): Promise<void> {
  return call("close")
}

/**
 * 关闭前询问：cb 返回 false（或 resolve 为 false）时阻止关闭
 * 第一个登记时告诉 DooTask 关闭前要问，全部取消后告诉它不用问了
 */
export function onBeforeClose(cb: () => boolean | void | Promise<boolean | void>): () => void {
  const first = listenerCount("before_close") === 0
  const off = listen("before_close", cb)
  if (first) {
    call("sdk.onBeforeClose", true).catch(error => console.error("[dootask]", error))
  }
  return () => {
    off()
    if (listenerCount("before_close") === 0) {
      call("sdk.onBeforeClose", false).catch(error => console.error("[dootask]", error))
    }
  }
}

/** 弹出到浏览器新标签或桌面客户端的独立窗口，要求入口配置了 popout: true */
export function popout(): Promise<void> {
  return call("popout")
}

/** 打开一个外部窗口：桌面客户端里是独立窗口，其他端是浏览器新标签 */
export function openWindow(url: string, options: OpenWindowOptions = {}): Promise<void> {
  if (!/^https?:\/\//i.test(url)) {
    return Promise.reject(new DooTaskError("invalid_args", "url must be an absolute http(s) address"))
  }
  return call("openWindow", url, options)
}

export function getPresentation(): Promise<Presentation> {
  return call("getPresentation")
}

export function onPresentationChange(cb: (presentation: Presentation) => void): () => void {
  return listen("presentation", cb)
}

/** 胶囊的位置和大小，沉浸式页面用来给胶囊留位置 */
export function getCapsuleRect(): Promise<Rect> {
  return call("getCapsuleRect")
}

/** 替换胶囊「···」里应用自己的菜单项 */
export function setCapsuleMenus(menus: CapsuleMenu[]): Promise<void> {
  return call("setCapsuleMenus", menus)
}

export function onCapsuleMenu(cb: (value: string) => void): () => void {
  return listen("capsule.menu", data => cb(String(data?.value ?? "")))
}
