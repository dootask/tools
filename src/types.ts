/**
 * 微应用数据类型定义
 */

/** eslint-disable-next-line @typescript-eslint/no-explicit-any */
export type Any = any

/** eslint-disable-next-line @typescript-eslint/no-unsafe-function-type */
export type Func = Function

/** 扩展Window全局接口 */
declare global {
  interface Window {
    microApp?: {
      getData: () => MicroAppData
      addDataListener?: (callback: Func, autoTrigger?: boolean) => void
      removeDataListener?: (callback: Func) => void
    }
    modalTransferIndex?: number
    systemInfo?: Any
  }
}

/** 微应用属性接口 */
export interface MicroAppProps {
  name: string
  url: string
  urlType: string

  userId: number
  userToken: string
  userInfo: Any

  baseUrl: string
  systemInfo: Any
  windowType: "popout" | "embed"

  isEEUIApp: boolean
  isElectron: boolean
  isMainElectron: boolean
  isSubElectron: boolean

  languageList: Any[]
  languageName: string
  themeName: string

  [key: string]: Any
}

/** 微应用方法接口 */
export interface MicroAppMethods {
  close: (destroy?: boolean) => void
  back: () => void
  openWindow: (params: OpenWindowParams) => void
  openTabWindow: (url: string) => void
  openAppPage: (params: OpenAppPageParams) => void
  requestAPI: (params: requestParams) => Promise<responseSuccess | responseError>
  selectUsers: (params: SelectUsersParams) => Promise<Any>
  nextZIndex: () => number
  extraCallA: (...args: Any[]) => Any

  [key: string]: Any
}

/** 微应用实例类型 */
export interface MicroAppInstance {
  Vue: Any
  store: Any
  components: {
    DialogWrapper: Any
    UserSelect: Any
    DatePicker: Any
    [key: string]: Any
  }
}

/** 完整微应用数据接口 */
export interface MicroAppData {
  type: string
  props: MicroAppProps
  methods?: MicroAppMethods
  instance?: MicroAppInstance
}

/** 窗口配置接口 */
export interface WindowConfig {
  /** 窗口标题 */
  title?: string
  /** 窗口标题是否固定 */
  titleFixed?: boolean
  /** 窗口宽度 */
  width?: number
  /** 窗口高度 */
  height?: number
  /** 窗口最小宽度 */
  minWidth?: number
  [key: string]: Any

  /** 更多配置项参考 https://www.electronjs.org/docs/latest/api/structures/base-window-options */
}

/** 打开独立窗口参数接口 */
export interface PopoutWindowParams extends WindowConfig {
  /** 自定义访问地址，如果为空则打开当前页面 */
  url?: string
  [key: string]: Any
}

/** 打开窗口参数接口 */
export interface OpenWindowParams {
  /** 窗口唯一标识 */
  name?: string
  /** 访问地址 */
  url?: string
  /** 是否强制创建新窗口，而不是重用已有窗口 */
  force?: boolean
  config?: WindowConfig

  [key: string]: Any
}

/** 打开应用页面参数接口 */
export interface OpenAppPageParams {
  /** 页面标题 */
  title?: string
  /** 窗口标题是否固定 */
  titleFixed?: boolean
  /** 访问地址 */
  url?: string
  [key: string]: Any
}

/** 选择用户参数接口 */
export interface SelectUsersParams {
  /** 已选择的值，默认值: [] */
  value?: string | number | Array<Any>
  /** 不允许取消的列表，默认值: [] */
  uncancelable?: Array<Any>
  /** 禁止选择的列表，默认值: [] */
  disabledChoice?: Array<Any>
  /** 指定项目ID，默认值: 0 */
  projectId?: number
  /** 指定非项目ID，默认值: 0 */
  noProjectId?: number
  /** 指定会话ID，默认值: 0 */
  dialogId?: number
  /** 是否显示机器人，默认值: false */
  showBot?: boolean
  /** 是否显示禁用的，默认值: false */
  showDisable?: boolean
  /** 最大选择数量 */
  multipleMax?: number
  /** 弹窗标题 */
  title?: string
  /** 搜索提示 */
  placeholder?: string
  /** 显示全选项，默认值: true */
  showSelectAll?: boolean
  /** 是否显示会话，默认值: false */
  showDialog?: boolean
  /** 仅显示群组，默认值: false */
  onlyGroup?: boolean
  /** 提交前的回调 */
  beforeSubmit?: Func
  [key: string]: Any
}

/** 弹出提示框参数接口 */
export interface ModalParams {
  title: string
  content?: string
  width?: number
  okText?: string
  cancelText?: string
  scrollable?: boolean
  closable?: boolean
  [key: string]: Any
}

/** 请求服务器API参数接口 */
export interface requestParams {
  /** 请求地址 */
  url: string
  /** 请求方式 */
  method?: string
  /** 请求数据 */
  data?: Any
  /** 请求超时时间 */
  timeout?: number
  /** 请求头 */
  header?: Any
  /** 是否显示加载动画 */
  spinner?: boolean
  [key: string]: Any
}

/** 请求服务器API返回接口（成功） */
export interface responseSuccess {
  /** 返回消息 */
  msg: string
  /** 返回数据 */
  data: Any
  /** XMLHttpRequest对象 */
  xhr: Any
  [key: string]: Any
}

/** 请求服务器API返回接口（错误） */
export interface responseError {
  /** 返回状态 */
  ret: number
  /** 返回消息 */
  msg: string
  /** 返回数据 */
  data: Any
  [key: string]: Any
}

/** 用户信息接口 */
export interface DooTaskUserInfo {
  /** 用户ID */
  userid: number
  /** 身份列表 */
  identity: string[]
  /** 部门ID列表 */
  department: number[]
  /** 邮箱 */
  email: string
  /** 电话 */
  tel: string
  /** 昵称 */
  nickname: string
  /** 职业 */
  profession: string
  /** 用户头像 */
  userimg: string
  /** 是否机器人 */
  bot: number
  /** 创建时间 */
  created_at: string
  [key: string]: Any
}

/** 系统信息接口 */
export interface DooTaskSystemInfo {
  /** 应用标题 */
  title?: string
  /** 系统版本 */
  version: string
  /** API地址 */
  apiUrl?: string
  [key: string]: Any
}

/** 语言列表接口 */
export type DooTaskLanguage = "zh" | "zh-CHT" | "en" | "ko" | "ja" | "de" | "fr" | "id" | "ru"
