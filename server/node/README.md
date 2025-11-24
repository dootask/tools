# DooTask Node 客户端

一个用于与 DooTask 系统交互的 Node.js 客户端。

## 安装

```bash
npm install @dootask/tools --save
```

> 建议使用 Node 18+（内置 `fetch`）。较低版本可在创建客户端时通过 `fetch` 选项传入自定义实现（例如 `node-fetch`）。

## 快速开始

```typescript
import { DooTaskClient } from "@dootask/tools"

// 创建客户端
const client = new DooTaskClient({
  token: "your-token-here",
  // server: "https://your-dootask-server.com", // 不传时默认 http://nginx
  timeoutMs: 10_000,
  // fetch: customFetch, // 旧版 Node 可自定义 fetch
})

// 获取用户信息
const user = await client.getUserInfo()
console.log("用户:", user.nickname)

// 发送一条消息
await client.sendMessage({
  dialog_id: 123,
  text: "Hello, World!",
  text_type: "md",
})
```

## 主要功能示例

### 发送消息

```typescript
import { DooTaskClient } from "@dootask/tools"

const client = new DooTaskClient({
  token: "YOUR_TOKEN",
})

// 发送消息到指定对话
await client.sendMessage({
  dialog_id: 123,
  text: "Hello, World!",
  text_type: "md",
})

// 发送消息到用户
await client.sendMessageToUser({
  userid: 456,
  text: "私信内容",
  text_type: "md",
})
```

### 项目和任务管理

```typescript
import { DooTaskClient } from "@dootask/tools"

const client = new DooTaskClient({
  token: "YOUR_TOKEN",
})

// 创建项目
const project = await client.createProject({
  name: "新项目",
  desc: "项目描述",
  columns: "待办,进行中,已完成",
})

// 创建任务
const task = await client.createTask({
  project_id: project.id,
  name: "新任务",
  content: "任务内容",
  owner: [123],
})

// 更新任务
const updatedTask = await client.updateTask({
  task_id: task.id,
  name: "更新后的任务名",
  content: "更新后的内容",
})
```

### 群组管理

```typescript
import { DooTaskClient } from "@dootask/tools"

const client = new DooTaskClient({
  token: "YOUR_TOKEN",
  server: "https://your-dootask-server.com",
})

// 创建群组
const group = await client.createGroup({
  chat_name: "新群组",
  userids: [123, 456, 789],
})

// 添加群成员
await client.addGroupUser({
  dialog_id: group.id,
  userids: [999],
})
```

## API 方法列表

### 基础方法

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `DooTaskClient` | 创建客户端实例 | `options: { token, server?, timeoutMs?, fetch? }` | `DooTaskClient` |

> 说明：`server` 默认值为 `http://nginx`，适用于在 Docker / K8s 等环境中通过服务名 `nginx` 访问 DooTask 主程序；在这种部署方式下生产环境通常只需传入 `token`。
>
> 所有方法返回 `Promise`。请求失败或业务 `ret != 1` 时，会抛出 `DooTaskApiError`。

### 用户相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `getUserInfo` | 获取用户信息 | `noCache?: boolean` | `Promise<UserInfo>` |
| `checkUserIdentity` | 检查用户身份 | `identity: string` | `Promise<UserInfo>` |
| `getUserDepartments` | 获取用户部门信息 | - | `Promise<Department[]>` |
| `getUsersBasic` | 获取多个用户基础信息 | `userids: number[]` | `Promise<UserBasic[]>` |
| `getUserBasic` | 获取单个用户基础信息 | `userid: number` | `Promise<UserBasic>` |

### 机器人相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `getBotList` | 获取机器人列表 | - | `Promise<BotListResponse>` |
| `getBot` | 获取机器人信息 | `params: GetBotRequest` | `Promise<Bot>` |
| `createBot` | 创建机器人 | `params: CreateBotRequest` | `Promise<Bot>` |
| `updateBot` | 更新机器人 | `params: EditBotRequest` | `Promise<Bot>` |
| `deleteBot` | 删除机器人 | `params: DeleteBotRequest` | `Promise<void>` |

### 消息相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `sendMessage` | 发送消息 | `message: SendMessageRequest` | `Promise<any>` |
| `sendMessageToUser` | 发送消息到用户 | `message: SendMessageToUserRequest` | `Promise<any>` |
| `sendBotMessage` | 发送机器人消息 | `message: SendBotMessageRequest` | `Promise<any>` |
| `sendAnonymousMessage` | 发送匿名消息 | `message: SendAnonymousMessageRequest` | `Promise<any>` |
| `sendStreamMessage` | 通知成员监听消息 | `message: SendStreamMessageRequest` | `Promise<any>` |
| `sendNoticeMessage` | 发送通知 | `message: SendNoticeMessageRequest` | `Promise<any>` |
| `sendTemplateMessage` | 发送模板消息 | `message: SendTemplateMessageRequest` | `Promise<any>` |
| `getMessageList` | 获取消息列表 | `params: GetMessageListRequest` | `Promise<DialogMessageListResponse>` |
| `searchMessage` | 搜索消息 | `params: SearchMessageRequest` | `Promise<DialogMessageSearchResponse>` |
| `getMessage` | 获取单个消息详情 | `params: GetMessageRequest` | `Promise<DialogMessage>` |
| `getMessageDetail` | 获取消息详情（兼容性） | `params: GetMessageRequest` | `Promise<DialogMessage>` |
| `withdrawMessage` | 撤回消息 | `params: WithdrawMessageRequest` | `Promise<void>` |
| `forwardMessage` | 转发消息 | `params: ForwardMessageRequest` | `Promise<void>` |
| `toggleMessageTodo` | 切换消息待办状态 | `params: ToggleMessageTodoRequest` | `Promise<void>` |
| `getMessageTodoList` | 获取消息待办列表 | `params: GetMessageRequest` | `Promise<TodoListResponse>` |
| `markMessageDone` | 标记消息完成 | `params: MarkMessageDoneRequest` | `Promise<void>` |
| `convertWebhookMessageToAi` | 转换 webhook 消息为 AI 对话格式 | `params: ConvertWebhookMessageRequest` | `Promise<ConvertWebhookMessageResponse>` |

### 对话相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `getDialogList` | 获取对话列表 | `params: TimeRangeRequest` | `Promise<ResponsePaginate<DialogInfo>>` |
| `searchDialog` | 搜索会话 | `params: SearchDialogRequest` | `Promise<DialogInfo[]>` |
| `getDialogOne` | 获取单个会话信息 | `params: GetDialogRequest` | `Promise<DialogInfo>` |
| `getDialogUser` | 获取会话成员 | `params: GetDialogUserRequest` | `Promise<DialogMember[]>` |

### 群组相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `createGroup` | 创建群组 | `params: CreateGroupRequest` | `Promise<DialogInfo>` |
| `editGroup` | 修改群组 | `params: EditGroupRequest` | `Promise<void>` |
| `addGroupUser` | 添加群成员 | `params: AddGroupUserRequest` | `Promise<void>` |
| `removeGroupUser` | 移除群成员 | `params: RemoveGroupUserRequest` | `Promise<void>` |
| `exitGroup` | 退出群组 | `dialog_id: number` | `Promise<void>` |
| `transferGroup` | 转让群组 | `params: TransferGroupRequest` | `Promise<void>` |
| `disbandGroup` | 解散群组 | `params: DisbandGroupRequest` | `Promise<void>` |

### 项目管理相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `getProjectList` | 获取项目列表 | `params: GetProjectListRequest` | `Promise<ResponsePaginate<Project>>` |
| `getProject` | 获取项目信息 | `params: GetProjectRequest` | `Promise<Project>` |
| `createProject` | 创建项目 | `params: CreateProjectRequest` | `Promise<Project>` |
| `updateProject` | 更新项目 | `params: UpdateProjectRequest` | `Promise<Project>` |
| `exitProject` | 退出项目 | `project_id: number` | `Promise<void>` |
| `deleteProject` | 删除项目 | `project_id: number` | `Promise<void>` |

### 任务列表相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `getColumnList` | 获取任务列表 | `params: GetColumnListRequest` | `Promise<ResponsePaginate<ProjectColumn>>` |
| `createColumn` | 创建任务列表 | `params: CreateColumnRequest` | `Promise<ProjectColumn>` |
| `updateColumn` | 更新任务列表 | `params: UpdateColumnRequest` | `Promise<ProjectColumn>` |
| `deleteColumn` | 删除任务列表 | `column_id: number` | `Promise<void>` |

### 任务相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `getTaskList` | 获取任务列表 | `params: GetTaskListRequest` | `Promise<ResponsePaginate<ProjectTask>>` |
| `getTask` | 获取任务信息 | `params: GetTaskRequest` | `Promise<ProjectTask>` |
| `getTaskContent` | 获取任务内容 | `params: GetTaskContentRequest` | `Promise<TaskContent>` |
| `getTaskFiles` | 获取任务文件列表 | `params: GetTaskFilesRequest` | `Promise<TaskFile[]>` |
| `createTask` | 创建任务 | `params: CreateTaskRequest` | `Promise<ProjectTask>` |
| `createSubTask` | 创建子任务 | `params: CreateSubTaskRequest` | `Promise<ProjectTask>` |
| `updateTask` | 更新任务 | `params: UpdateTaskRequest` | `Promise<ProjectTask>` |
| `createTaskDialog` | 创建任务对话 | `params: CreateTaskDialogRequest` | `Promise<CreateTaskDialogResponse>` |
| `archiveTask` | 归档任务 | `task_id: number, archiveType: string` | `Promise<void>` |
| `deleteTask` | 删除任务 | `task_id: number, deleteType?: string` | `Promise<void>` |

### 系统相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `getSystemSettings` | 获取系统设置 | - | `Promise<SystemSettings>` |
| `getVersion` | 获取版本信息 | - | `Promise<VersionInfo>` |
