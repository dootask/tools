import type {
  AddGroupUserRequest,
  ApiResponse,
  Bot,
  BotListResponse,
  ColumnActionRequest,
  DeleteBotRequest,
  ConvertWebhookMessageRequest,
  ConvertWebhookMessageResponse,
  CreateBotRequest,
  CreateColumnRequest,
  CreateGroupRequest,
  CreateProjectRequest,
  CreateSubTaskRequest,
  CreateTaskDialogRequest,
  CreateTaskDialogResponse,
  CreateTaskRequest,
  Department,
  DialogInfo,
  DialogMember,
  DialogMessage,
  DialogMessageListResponse,
  DialogMessageSearchResponse,
  DialogOpenUserResponse,
  DialogUserResponse,
  DisbandGroupRequest,
  EditBotRequest,
  EditGroupRequest,
  ForwardMessageRequest,
  GetBotRequest,
  GetColumnListRequest,
  GetDialogRequest,
  GetDialogUserRequest,
  GetMessageListRequest,
  GetMessageRequest,
  GetProjectListRequest,
  GetProjectRequest,
  GetTaskContentRequest,
  GetTaskFilesRequest,
  GetTaskListRequest,
  GetTaskRequest,
  MarkMessageDoneRequest,
  Project,
  ProjectActionRequest,
  ProjectColumn,
  ProjectTask,
  RemoveGroupUserRequest,
  ResponsePaginate,
  SearchDialogRequest,
  SearchMessageRequest,
  SendAnonymousMessageRequest,
  SendBotMessageRequest,
  SendMessageRequest,
  SendMessageToUserRequest,
  SendNoticeMessageRequest,
  SendStreamMessageRequest,
  SendTemplateMessageRequest,
  SystemSettings,
  TaskActionRequest,
  TaskContent,
  TaskFile,
  TimeRangeRequest,
  TransferGroupRequest,
  ToggleMessageTodoRequest,
  TodoListResponse,
  UpdateColumnRequest,
  UpdateProjectRequest,
  UpdateTaskRequest,
  UserBasic,
  UserInfo,
  VersionInfo,
  WithdrawMessageRequest,
} from "./types"

type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE"

export class DooTaskApiError extends Error {
  status?: number
  statusText?: string
  rawBody?: string
  ret?: number
  apiData?: unknown

  constructor(message: string, extras: Partial<DooTaskApiError> = {}) {
    super(message)
    Object.assign(this, extras)
  }
}

export interface DooTaskClientOptions {
  token: string
  server?: string
  timeoutMs?: number
  fetch?: typeof fetch
}

interface UserCache {
  user: UserInfo
  expiresAt: number
}

export class DooTaskClient {
  private readonly token: string
  private readonly server: string
  private readonly timeoutMs: number
  private readonly fetchFn: typeof fetch
  private userCache?: UserCache
  private cacheTimeMs = 10 * 60 * 1000

  constructor(options: DooTaskClientOptions) {
    this.token = options.token
    this.server = options.server ?? "http://nginx"
    this.timeoutMs = options.timeoutMs ?? 10_000
    this.fetchFn = options.fetch ?? (globalThis.fetch as typeof fetch)
    if (!this.fetchFn) {
      throw new Error("fetch is not available; provide a fetch implementation in options.fetch")
    }
  }

  // ------------------------------------------------------------------------------------------
  // 基础 HTTP
  // ------------------------------------------------------------------------------------------
  private buildURL(baseURL: string, params?: Record<string, unknown>): string {
    if (!params || Object.keys(params).length === 0) return baseURL

    const queryParts: string[] = []
    const append = (key: string, value: unknown) => {
      if (value === undefined || value === null) return
      if (Array.isArray(value)) {
        value.forEach(v => append(key + "[]", v))
        return
      }
      if (typeof value === "boolean") {
        queryParts.push(`${encodeURIComponent(key)}=${value ? "1" : "0"}`)
        return
      }
      queryParts.push(`${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
    }

    for (const [key, value] of Object.entries(params)) {
      append(key, value)
    }

    if (queryParts.length === 0) return baseURL
    const separator = baseURL.includes("?") ? "&" : "?"
    return `${baseURL}${separator}${queryParts.join("&")}`
  }

  private async request<T>(method: HttpMethod, api: string, payload?: unknown): Promise<T> {
    const isGetLike = method === "GET" || method === "DELETE"
    const url = isGetLike ? this.buildURL(this.server + api, payload as Record<string, unknown>) : this.server + api

    const headers: Record<string, string> = {
      Token: this.token,
      "User-Agent": "DooTask-Node-Client/1.0",
    }

    const controller = new AbortController()
    const timeout = setTimeout(() => controller.abort(), this.timeoutMs)

    try {
      const init: RequestInit = { method, headers, signal: controller.signal }
      if (!isGetLike && payload !== undefined) {
        init.body = JSON.stringify(payload)
        headers["Content-Type"] = "application/json"
      }

      const res = await this.fetchFn(url, init)
      if (!res.ok) {
        const rawBody = await res.text().catch(() => "")
        throw new DooTaskApiError(`HTTP ${res.status} ${res.statusText}`, {
          status: res.status,
          statusText: res.statusText,
          rawBody,
        })
      }

      const json = (await res.json()) as ApiResponse<T>
      if (typeof json?.ret !== "number") {
        throw new DooTaskApiError("Invalid response format", { rawBody: JSON.stringify(json ?? null) })
      }
      if (json.ret !== 1) {
        throw new DooTaskApiError(json.msg || "API error", { ret: json.ret, apiData: json.data })
      }
      return json.data
    } catch (error: any) {
      if (error?.name === "AbortError") {
        throw new DooTaskApiError("Request timeout", { statusText: "timeout" })
      }
      if (error instanceof DooTaskApiError) {
        throw error
      }
      throw new DooTaskApiError(error?.message ?? "Request failed", { rawBody: error?.stack })
    } finally {
      clearTimeout(timeout)
    }
  }

  private async get<T>(api: string, params?: unknown): Promise<T> {
    return this.request<T>("GET", api, params)
  }

  private async post<T>(api: string, body?: unknown): Promise<T> {
    return this.request<T>("POST", api, body)
  }

  // ------------------------------------------------------------------------------------------
  // 用户
  // ------------------------------------------------------------------------------------------
  async getUserInfo(noCache = false): Promise<UserInfo> {
    if (!noCache && this.userCache && Date.now() < this.userCache.expiresAt) {
      return this.userCache.user
    }
    const data = await this.get<UserInfo>("/api/users/info")
    this.userCache = { user: data, expiresAt: Date.now() + this.cacheTimeMs }
    return data
  }

  async checkUserIdentity(identity: string): Promise<UserInfo> {
    const user = await this.getUserInfo()
    if (!user.identity?.includes(identity)) {
      throw new DooTaskApiError("insufficient permissions")
    }
    return user
  }

  async getUserDepartments(): Promise<Department[]> {
    return this.get<Department[]>("/api/users/info/departments")
  }

  async getUsersBasic(userids: number[]): Promise<UserBasic[]> {
    return this.get<UserBasic[]>("/api/users/basic", { userid: userids })
  }

  async getUserBasic(userid: number): Promise<UserBasic> {
    const list = await this.getUsersBasic([userid])
    if (!list.length) {
      throw new DooTaskApiError("用户不存在")
    }
    return list[0]
  }

  // ------------------------------------------------------------------------------------------
  // 机器人
  // ------------------------------------------------------------------------------------------
  async getBotList(): Promise<BotListResponse> {
    return this.get<BotListResponse>("/api/users/bot/list")
  }

  async getBot(params: GetBotRequest): Promise<Bot> {
    return this.get<Bot>("/api/users/bot/info", params)
  }

  async createBot(params: CreateBotRequest): Promise<Bot> {
    return this.post<Bot>("/api/users/bot/edit", params)
  }

  async updateBot(params: EditBotRequest): Promise<Bot> {
    return this.post<Bot>("/api/users/bot/edit", params)
  }

  async deleteBot(params: DeleteBotRequest): Promise<void> {
    await this.get<void>("/api/users/bot/delete", params)
  }

  // ------------------------------------------------------------------------------------------
  // 消息
  // ------------------------------------------------------------------------------------------
  async sendMessage<T = any>(message: SendMessageRequest): Promise<T> {
    if (!message.text_type) message.text_type = "md"
    return this.post<T>("/api/dialog/msg/sendtext", message)
  }

  async sendMessageToUser<T = any>(message: SendMessageToUserRequest): Promise<T> {
    const dialog = await this.get<DialogOpenUserResponse>("/api/dialog/open/user", { userid: message.userid })
    return this.sendMessage<T>({
      dialog_id: dialog.dialog_user.dialog_id,
      text: message.text,
      text_type: message.text_type,
      silence: message.silence,
    })
  }

  async sendBotMessage<T = any>(message: SendBotMessageRequest): Promise<T> {
    if (!message.bot_type) message.bot_type = "system-msg"
    return this.post<T>("/api/dialog/msg/sendbot", message)
  }

  async sendAnonymousMessage<T = any>(message: SendAnonymousMessageRequest): Promise<T> {
    return this.post<T>("/api/dialog/msg/sendanon", message)
  }

  async sendStreamMessage<T = any>(message: SendStreamMessageRequest): Promise<T> {
    if (!message.source) message.source = "api"
    return this.post<T>("/api/dialog/msg/stream", message)
  }

  async sendNoticeMessage<T = any>(message: SendNoticeMessageRequest): Promise<T> {
    if (!message.source) message.source = "api"
    return this.post<T>("/api/dialog/msg/sendnotice", message)
  }

  async sendTemplateMessage<T = any>(message: SendTemplateMessageRequest): Promise<T> {
    if (!message.source) message.source = "api"
    return this.post<T>("/api/dialog/msg/sendtemplate", message)
  }

  async getMessageList(params: GetMessageListRequest): Promise<DialogMessageListResponse> {
    return this.get<DialogMessageListResponse>("/api/dialog/msg/list", params)
  }

  async searchMessage(params: SearchMessageRequest): Promise<DialogMessageSearchResponse> {
    return this.get<DialogMessageSearchResponse>("/api/dialog/msg/search", params)
  }

  async getMessage(params: GetMessageRequest): Promise<DialogMessage> {
    return this.get<DialogMessage>("/api/dialog/msg/one", params)
  }

  async getMessageDetail(params: GetMessageRequest): Promise<DialogMessage> {
    return this.get<DialogMessage>("/api/dialog/msg/detail", params)
  }

  async withdrawMessage(params: WithdrawMessageRequest): Promise<void> {
    await this.get<void>("/api/dialog/msg/withdraw", params)
  }

  async forwardMessage(params: ForwardMessageRequest): Promise<void> {
    await this.get<void>("/api/dialog/msg/forward", params)
  }

  async toggleMessageTodo(params: ToggleMessageTodoRequest): Promise<void> {
    if (!params.type) params.type = "all"
    await this.get<void>("/api/dialog/msg/todo", params)
  }

  async getMessageTodoList(params: GetMessageRequest): Promise<TodoListResponse> {
    return this.get<TodoListResponse>("/api/dialog/msg/todolist", params)
  }

  async markMessageDone(params: MarkMessageDoneRequest): Promise<void> {
    await this.get<void>("/api/dialog/msg/done", params)
  }

  async convertWebhookMessageToAi(params: ConvertWebhookMessageRequest): Promise<ConvertWebhookMessageResponse> {
    return this.post<ConvertWebhookMessageResponse>("/api/dialog/msg/webhookmsg2ai", params)
  }

  // ------------------------------------------------------------------------------------------
  // 对话
  // ------------------------------------------------------------------------------------------
  async getDialogList(params: TimeRangeRequest): Promise<ResponsePaginate<DialogInfo>> {
    return this.get<ResponsePaginate<DialogInfo>>("/api/dialog/lists", params)
  }

  async searchDialog(params: SearchDialogRequest): Promise<DialogInfo[]> {
    return this.get<DialogInfo[]>("/api/dialog/search", params)
  }

  async getDialogOne(params: GetDialogRequest): Promise<DialogInfo> {
    return this.get<DialogInfo>("/api/dialog/one", params)
  }

  async getDialogUser(params: GetDialogUserRequest): Promise<DialogMember[]> {
    return this.get<DialogMember[]>("/api/dialog/user", params)
  }

  // ------------------------------------------------------------------------------------------
  // 群组
  // ------------------------------------------------------------------------------------------
  async createGroup(params: CreateGroupRequest): Promise<DialogInfo> {
    return this.get<DialogInfo>("/api/dialog/group/add", params)
  }

  async editGroup(params: EditGroupRequest): Promise<void> {
    await this.get<void>("/api/dialog/group/edit", params)
  }

  async addGroupUser(params: AddGroupUserRequest): Promise<void> {
    await this.get<void>("/api/dialog/group/adduser", params)
  }

  async removeGroupUser(params: RemoveGroupUserRequest): Promise<void> {
    await this.get<void>("/api/dialog/group/deluser", params)
  }

  async exitGroup(dialog_id: number): Promise<void> {
    await this.removeGroupUser({ dialog_id, userids: [] })
  }

  async transferGroup(params: TransferGroupRequest): Promise<void> {
    await this.get<void>("/api/dialog/group/transfer", params)
  }

  async disbandGroup(params: DisbandGroupRequest): Promise<void> {
    await this.get<void>("/api/dialog/group/disband", params)
  }

  // ------------------------------------------------------------------------------------------
  // 项目
  // ------------------------------------------------------------------------------------------
  async getProjectList(params: GetProjectListRequest): Promise<ResponsePaginate<Project>> {
    return this.get<ResponsePaginate<Project>>("/api/project/lists", params)
  }

  async getProject(params: GetProjectRequest): Promise<Project> {
    return this.get<Project>("/api/project/one", params)
  }

  async createProject(params: CreateProjectRequest): Promise<Project> {
    return this.get<Project>("/api/project/add", params)
  }

  async updateProject(params: UpdateProjectRequest): Promise<Project> {
    return this.get<Project>("/api/project/update", params)
  }

  async exitProject(project_id: number): Promise<void> {
    const payload: ProjectActionRequest = { project_id }
    await this.get<void>("/api/project/exit", payload)
  }

  async deleteProject(project_id: number): Promise<void> {
    const payload: ProjectActionRequest = { project_id }
    await this.get<void>("/api/project/remove", payload)
  }

  // ------------------------------------------------------------------------------------------
  // 列表
  // ------------------------------------------------------------------------------------------
  async getColumnList(params: GetColumnListRequest): Promise<ResponsePaginate<ProjectColumn>> {
    return this.get<ResponsePaginate<ProjectColumn>>("/api/project/column/lists", params)
  }

  async createColumn(params: CreateColumnRequest): Promise<ProjectColumn> {
    return this.get<ProjectColumn>("/api/project/column/add", params)
  }

  async updateColumn(params: UpdateColumnRequest): Promise<ProjectColumn> {
    return this.get<ProjectColumn>("/api/project/column/update", params)
  }

  async deleteColumn(column_id: number): Promise<void> {
    const payload: ColumnActionRequest = { column_id }
    await this.get<void>("/api/project/column/remove", payload)
  }

  // ------------------------------------------------------------------------------------------
  // 任务
  // ------------------------------------------------------------------------------------------
  async getTaskList(params: GetTaskListRequest): Promise<ResponsePaginate<ProjectTask>> {
    return this.get<ResponsePaginate<ProjectTask>>("/api/project/task/lists", params)
  }

  async getTask(params: GetTaskRequest): Promise<ProjectTask> {
    return this.get<ProjectTask>("/api/project/task/one", params)
  }

  async getTaskContent(params: GetTaskContentRequest): Promise<TaskContent> {
    return this.get<TaskContent>("/api/project/task/content", params)
  }

  async getTaskFiles(params: GetTaskFilesRequest): Promise<TaskFile[]> {
    return this.get<TaskFile[]>("/api/project/task/files", params)
  }

  async createTask(params: CreateTaskRequest): Promise<ProjectTask> {
    return this.post<ProjectTask>("/api/project/task/add", params)
  }

  async createSubTask(params: CreateSubTaskRequest): Promise<ProjectTask> {
    return this.get<ProjectTask>("/api/project/task/addsub", params)
  }

  async updateTask(params: UpdateTaskRequest): Promise<ProjectTask> {
    return this.post<ProjectTask>("/api/project/task/update", params)
  }

  async createTaskDialog(params: CreateTaskDialogRequest): Promise<CreateTaskDialogResponse> {
    return this.get<CreateTaskDialogResponse>("/api/project/task/dialog", params)
  }

  async archiveTask(task_id: number, archiveType: string): Promise<void> {
    const payload: TaskActionRequest = { task_id, type: archiveType }
    await this.get<void>("/api/project/task/archived", payload)
  }

  async deleteTask(task_id: number, deleteType?: string): Promise<void> {
    const payload: TaskActionRequest = { task_id, type: deleteType }
    await this.get<void>("/api/project/task/remove", payload)
  }

  // ------------------------------------------------------------------------------------------
  // 系统
  // ------------------------------------------------------------------------------------------
  async getSystemSettings(): Promise<SystemSettings> {
    return this.get<SystemSettings>("/api/system/setting")
  }

  async getVersion(): Promise<VersionInfo> {
    return this.get<VersionInfo>("/api/system/version", { version: true })
  }
}
