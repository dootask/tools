export type Channel = "frame" | "inline" | "electron" | "app" | "bridge"

export type Platform = "web" | "electron" | "app"

export type Presentation = "page" | "panel" | "modal"

export type Theme = "light" | "dark"

export type Language = "zh" | "zh-CHT" | "en" | "ko" | "ja" | "de" | "fr" | "id" | "ru"

export interface Env {
  channel: Channel
  platform: Platform
  /** 只有窗口有，其他为 null */
  presentation: Presentation | null
  mobile: boolean
  /** 通信协议版本，2.x 为 2 */
  protocol: number
}

export interface SystemInfo {
  /** DooTask 版本 */
  version: string
  /** 系统名称 */
  title: string
  /** DooTask 的访问地址 */
  baseUrl: string
}

export interface SafeArea {
  top: number
  bottom: number
}

export interface UserInfo {
  userid: number
  nickname: string
  email: string
  /** 头像地址 */
  userimg: string
  /** 身份，比如 ["admin"] */
  identity: string[]
  /** 部门 ID */
  department: number[]
  profession: string
  tel: string
  bot: boolean
  created_at: string
  [key: string]: any
}

export interface UserBasic {
  userid: number
  nickname: string
  email: string
  userimg: string
  department: number[]
  profession: string
  bot: boolean
  online: boolean
  /** 最近在线时间 */
  line_at: string
  [key: string]: any
}

export interface RequestOptions {
  /** 默认 GET */
  method?: "GET" | "POST"
  data?: Record<string, any>
  headers?: Record<string, string>
}

export interface DialogOptions {
  title: string
  /** 纯文本 */
  content?: string
  okText?: string
  cancelText?: string
  width?: number
}

export interface SelectedIds {
  users: number[]
  groups: number[]
}

export interface SelectUsersOptions {
  /** 已选的用户 */
  selected?: number[]
  /** 已选、不能取消的用户 */
  locked?: number[]
  /** 不能选的用户 */
  disabled?: number[]
  /** 最多选几个 */
  max?: number
  /** 只在这个项目的成员里选 */
  projectId?: number
  /** 只在这个会话的成员里选 */
  dialogId?: number
  /** 显示机器人，默认 false */
  includeBots?: boolean
  /** 显示已停用的用户，默认 false */
  includeDisabled?: boolean
  /** 同时可以选群组，默认 false */
  includeGroups?: boolean
  /** 只选群组，默认 false */
  onlyGroups?: boolean
  title?: string
  placeholder?: string
  /** 返回 false 阻止提交 */
  beforeSubmit?: (ids: SelectedIds) => boolean | Promise<boolean>
}

export interface AddTaskParams {
  projectId?: number
  columnId?: number
  name?: string
  /** HTML */
  content?: string
  /** 开始、结束时间（YYYY-MM-DD HH:mm） */
  times?: [string, string]
  owner?: number[]
}

export interface MeetingParams {
  /** 预选的成员 */
  userids?: number[]
  /** 会议消息发到哪个会话 */
  dialogId?: number
}

export interface OpenAIParams {
  /** 填进输入框 */
  prompt?: string
  /** 直接发送 */
  send?: boolean
}

export interface Rect {
  top: number
  /** 胶囊右边到窗口右边的距离 */
  right: number
  width: number
  height: number
}

export interface CapsuleMenu {
  label: string
  value: string | number
}

export interface OpenWindowOptions {
  title?: string
  width?: number
  height?: number
}
