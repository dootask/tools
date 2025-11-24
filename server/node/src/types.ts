export interface ApiResponse<T> {
  ret: number
  msg: string
  data: T
}

export interface ResponsePaginate<T> {
  current_page: number
  data: T[]
  next_page_url?: string | null
  path: string
  per_page: number | string
  prev_page_url?: string | null
  to: number | string | null
  total: number
}

// 用户
export interface UserInfo {
  userid: number
  identity: string[]
  email: string
  nickname: string
  profession: string
  userimg: string
  bot: number
  department: number[]
  department_name: string
}

export interface UserBasic {
  userid: number
  email: string
  nickname: string
  profession: string
  userimg: string
  bot: number
  online: boolean
  department: number[]
  department_name: string
}

export interface Department {
  id: number
  name: string
  parent_id: number
  owner_userid: number
}

// 机器人
export interface Bot {
  id: number
  name: string
  avatar: string
  clear_day: number
  webhook_url: string
}

export interface BotListResponse {
  list: Bot[]
}

export interface GetBotRequest {
  id: number
}

export interface CreateBotRequest {
  name: string
  avatar?: string
  clear_day?: number
  webhook_url?: string
  session?: number
}

export interface EditBotRequest {
  id: number
  name?: string
  avatar?: string
  clear_day?: number
  webhook_url?: string
}

export interface DeleteBotRequest {
  id: number
  remark: string
}

// 消息
export interface SendMessageRequest {
  dialog_id: number
  text: string
  text_type?: string
  silence?: boolean
  reply_id?: number
  reply_check?: string
  update_id?: number
  update_mark?: string
}

export interface SendMessageToUserRequest {
  userid: number
  text: string
  text_type?: string
  silence?: boolean
}

export interface SendBotMessageRequest {
  userid: number
  text: string
  bot_type?: string
  bot_name?: string
  silence?: boolean
}

export interface SendAnonymousMessageRequest {
  userid: number
  text: string
}

export interface SendStreamMessageRequest {
  userid: number
  stream_url: string
  source?: string
}

export interface SendNoticeMessageRequest {
  dialog_id: number
  dialog_ids?: string
  notice: string
  silence?: boolean
  source?: string
}

export interface TemplateContent {
  content: string
  style: string
}

export interface SendTemplateMessageRequest {
  dialog_id: number
  dialog_ids?: string
  content: TemplateContent[]
  title?: string
  silence?: boolean
  source?: string
}

export interface ConvertWebhookMessageRequest {
  msg: string
}

export interface ConvertWebhookMessageResponse {
  msg: string
}

export interface DialogMessage {
  id: number
  dialog_id: number
  userid: number
  bot: number
  created_at: string
  type: string
  mtype: string
  msg: any
  reply_id: number
  reply_num: number
  forward_id: number
  forward_num: number
  tag: number
  todo: number
  read: number
  send: number
  read_at?: string | null
  mention: number
  dot: number
  emoji: any[]
  link: number
  modify: number
  percentage: number
}

export interface DialogMessageListResponse {
  list: DialogMessage[]
  time: number
  dialog: DialogInfo
  todo: any[]
  top?: any
}

export interface DialogMessageSearchResponse {
  data: number[]
}

export interface TodoUser {
  userid: number
  nickname: string
  userimg: string
  done: boolean
  done_at: string
}

export interface TodoListResponse {
  users: TodoUser[]
}

export interface GetMessageListRequest {
  dialog_id: number
  msg_id?: number
  position_id?: number
  prev_id?: number
  next_id?: number
  msg_type?: string
  take?: number
}

export interface SearchMessageRequest {
  dialog_id: number
  key: string
}

export interface GetMessageRequest {
  msg_id: number
}

export interface WithdrawMessageRequest {
  msg_id: number
}

export interface ForwardMessageRequest {
  msg_id: number
  dialogids?: number[]
  userids?: number[]
  show_source?: number
  leave_message?: string
}

export interface ToggleMessageTodoRequest {
  msg_id: number
  type?: string
  userids?: number[]
}

export interface MarkMessageDoneRequest {
  msg_id: number
}

// 对话
export interface DialogOpenUserResponse {
  dialog_user: DialogUserResponse
}

export interface DialogUserResponse {
  dialog_id: number
  userid: number
  bot: number
}

export interface DialogInfo {
  id: number
  type: string
  group_type: string
  name: string
  avatar: string
  owner_id: number
  created_at: string
  updated_at: string
  last_at: string
  mark_unread: number
  silence: number
  hide: number
  color: string
  unread: number
  unread_one: number
  mention: number
  mention_ids: number[]
  people: number
  people_user: number
  people_bot: number
  todo_num: number
  last_msg: any
  pinyin: string
  bot: number
  top_at: string
}

export interface DialogMember {
  id: number
  dialog_id: number
  userid: number
  nickname: string
  email: string
  userimg: string
  bot: number
  online: boolean
}

export interface TimeRangeRequest {
  timerange?: string
  page?: number
  pagesize?: number
}

export interface SearchDialogRequest {
  key: string
}

export interface GetDialogRequest {
  dialog_id: number
}

export interface GetDialogUserRequest {
  dialog_id: number
  getuser?: number
}

// 群组
export interface CreateGroupRequest {
  avatar?: string
  chat_name?: string
  userids: number[]
}

export interface EditGroupRequest {
  dialog_id: number
  avatar?: string
  chat_name?: string
  admin?: number
}

export interface AddGroupUserRequest {
  dialog_id: number
  userids: number[]
}

export interface RemoveGroupUserRequest {
  dialog_id: number
  userids?: number[]
}

export interface TransferGroupRequest {
  dialog_id: number
  userid: number
  check_owner?: string
  key?: string
}

export interface DisbandGroupRequest {
  dialog_id: number
}

// 项目
export interface Project {
  id: number
  name: string
  desc: string
  userid: number
  dialog_id: number
  archived_at: string
  created_at: string
  updated_at: string
  owner: number
  owner_userid: number
  personal: number
  task_num: number
  task_complete: number
  task_percent: number
  task_my_num: number
  task_my_complete: number
  task_my_percent: number
}

export interface GetProjectListRequest {
  type?: string
  archived?: string
  getcolumn?: string
  getuserid?: string
  getstatistics?: string
  timerange?: string
  page?: number
  pagesize?: number
}

export interface GetProjectRequest {
  project_id: number
}

export interface CreateProjectRequest {
  name: string
  desc?: string
  columns?: string
  flow?: string
  personal?: number
}

export interface UpdateProjectRequest {
  project_id: number
  name: string
  desc?: string
  archive_method?: string
  archive_days?: number
}

export interface ProjectActionRequest {
  project_id: number
  type?: string
}

// 列表
export interface ProjectColumn {
  id: number
  project_id: number
  name: string
  color: string
  sort: number
  created_at: string
  updated_at: string
}

export interface GetColumnListRequest {
  project_id: number
  page?: number
  pagesize?: number
}

export interface CreateColumnRequest {
  project_id: number
  name: string
}

export interface UpdateColumnRequest {
  column_id: number
  name?: string
  color?: string
}

export interface ColumnActionRequest {
  column_id: number
}

// 任务
export interface ProjectTask {
  id: number
  project_id: number
  column_id: number
  parent_id: number
  name: string
  desc: string
  start_at: string
  end_at: string
  complete_at: string
  archived_at: string
  created_at: string
  updated_at: string
  userid: number
  dialog_id: number
  flow_item_id: number
  flow_item_name: string
  visibility: number
  color: string
  file_num: number
  msg_num: number
  sub_num: number
  sub_complete: number
  percent: number
  project_name: string
  column_name: string
}

export interface TaskFile {
  id: number
  task_id: number
  name: string
  ext: string
  size: number
  path: string
  thumb: string
  userid: number
  created_at: string
  updated_at: string
}

export interface TaskContent {
  content: string
  type: string
}

export interface GetTaskListRequest {
  project_id?: number
  parent_id?: number
  archived?: string
  deleted?: string
  timerange?: string
  page?: number
  pagesize?: number
}

export interface GetTaskRequest {
  task_id: number
  archived?: string
}

export interface GetTaskContentRequest {
  task_id: number
  history_id?: number
}

export interface GetTaskFilesRequest {
  task_id: number
}

export interface CreateTaskRequest {
  project_id: number
  column_id?: any
  name: string
  content?: string
  times?: string[]
  owner?: number[]
  top?: number
}

export interface CreateSubTaskRequest {
  task_id: number
  name: string
}

export interface UpdateTaskRequest {
  task_id: number
  name?: string
  content?: string
  times?: string[]
  owner?: number[]
  assist?: number[]
  color?: string
  visibility?: number
  complete_at?: any
}

export interface TaskActionRequest {
  task_id: number
  type?: string
}

export interface CreateTaskDialogRequest {
  task_id: number
}

export interface CreateTaskDialogResponse {
  id: number
  dialog_id: number
  dialog_data: any
}

// 系统
export interface SystemSettings {
  reg?: string | null
  task_default_time?: string[] | null
  system_alias?: string | null
  system_welcome: string
  server_timezone?: string | null
  server_version?: string | null
}

export interface VersionInfo {
  device_count: number
  version: string
}
