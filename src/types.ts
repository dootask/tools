/**
 * 微应用数据类型定义
 */

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type Any = any

// eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
export type Func = Function

// 扩展Window全局接口
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

// 微应用属性接口
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

// 微应用方法接口
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

// 微应用实例类型
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

// 完整微应用数据接口
export interface MicroAppData {
  type: string
  props: MicroAppProps
  methods?: MicroAppMethods
  instance?: MicroAppInstance
}

// 窗口配置接口
export interface WindowConfig {
  title?: string // 窗口标题
  titleFixed?: boolean // 窗口标题是否固定
  width?: number // 窗口宽度
  height?: number // 窗口高度
  minWidth?: number // 窗口最小宽度
  [key: string]: Any

  // 更多配置项参考 https://www.electronjs.org/docs/latest/api/structures/base-window-options
}

// 打开独立窗口参数接口
export interface PopoutWindowParams extends WindowConfig {
  url?: string // 自定义访问地址，如果为空则打开当前页面
  [key: string]: Any
}

// 打开窗口参数接口
export interface OpenWindowParams {
  name?: string // 窗口唯一标识
  url?: string // 访问地址
  force?: boolean // 是否强制创建新窗口，而不是重用已有窗口
  config?: WindowConfig

  [key: string]: Any
}

// 打开应用页面参数接口
export interface OpenAppPageParams {
  title?: string // 页面标题
  titleFixed?: boolean // 窗口标题是否固定
  url?: string // 访问地址
  [key: string]: Any
}

// 选择用户参数接口
export interface SelectUsersParams {
  value?: string | number | Array<Any> // 已选择的值，默认值: []
  uncancelable?: Array<Any> // 不允许取消的列表，默认值: []
  disabledChoice?: Array<Any> // 禁止选择的列表，默认值: []
  projectId?: number // 指定项目ID，默认值: 0
  noProjectId?: number // 指定非项目ID，默认值: 0
  dialogId?: number // 指定会话ID，默认值: 0
  showBot?: boolean // 是否显示机器人，默认值: false
  showDisable?: boolean // 是否显示禁用的，默认值: false
  multipleMax?: number // 最大选择数量
  title?: string // 弹窗标题
  placeholder?: string // 搜索提示
  showSelectAll?: boolean // 显示全选项，默认值: true
  showDialog?: boolean // 是否显示会话，默认值: false
  onlyGroup?: boolean // 仅显示群组，默认值: false
  beforeSubmit?: Func // 提交前的回调
  [key: string]: Any
}

// 弹出提示框参数接口
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

// 请求服务器API参数接口
export interface requestParams {
  url: string // 请求地址
  method?: string // 请求方式
  data?: Any // 请求数据
  timeout?: number // 请求超时时间
  header?: Any // 请求头
  spinner?: boolean // 是否显示加载动画
  [key: string]: Any
}

// 请求服务器API返回接口（成功）
export interface responseSuccess {
  msg: string // 返回消息
  data: Any // 返回数据
  xhr: Any // XMLHttpRequest对象
  [key: string]: Any
}

// 请求服务器API返回接口（错误）
export interface responseError {
  ret: number // 返回状态
  msg: string // 返回消息
  data: Any // 返回数据
  [key: string]: Any
}

// 用户信息接口
export interface DooTaskUserInfo {
  userid: number // 用户ID
  identity: string[] // 身份列表
  department: number[] // 部门ID列表
  email: string // 邮箱
  tel: string // 电话
  nickname: string // 昵称
  profession: string // 职业
  userimg: string // 用户头像
  bot: number // 是否机器人
  created_at: string // 创建时间
  [key: string]: Any
}

// 系统信息接口
export interface DooTaskSystemInfo {
  title?: string // 应用标题
  version: string // 系统版本
  apiUrl?: string // API地址
  [key: string]: Any
}

// 语言列表接口
export type DooTaskLanguage = "zh" | "zh-CHT" | "en" | "ko" | "ja" | "de" | "fr" | "id" | "ru"
