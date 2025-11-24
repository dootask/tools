# DooTask Go 客户端

DooTask Go 客户端库，用于与 DooTask 系统进行交互。提供完整的 API 封装，支持用户管理、消息通信、项目管理、任务协作等功能。

## 安装

```bash
go get github.com/dootask/tools/server/go
```

## 快速开始

```go
package main

import (
    "fmt"
    "time"
    
    dootask "github.com/dootask/tools/server/go"
)

func main() {
    // 创建客户端
    client := dootask.NewClient(
        "your-token-here",
        dootask.WithServer("http://your-server.com"),
        dootask.WithTimeout(30*time.Second),
    )
    
    // 获取用户信息
    user, err := client.GetUserInfo()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("用户: %s\n", user.Nickname)
}
```

> 说明：如果不调用 `WithServer`，客户端默认使用 `http://nginx` 作为服务地址，适用于在 Docker / K8s 等环境中通过服务名 `nginx` 访问 DooTask 主程序。

## 主要功能示例

### 发送消息

```go
// 发送消息到指定对话
err := client.SendMessage(dootask.SendMessageRequest{
    DialogID: 123,
    Text:     "Hello, World!",
    TextType: "md",
})

// 发送消息到用户
err := client.SendMessageToUser(dootask.SendMessageToUserRequest{
    UserID:   456,
    Text:     "私信内容",
    TextType: "md",
})
```

### 项目和任务管理

```go
// 创建项目
project, err := client.CreateProject(dootask.CreateProjectRequest{
    Name:     "新项目",
    Desc:     "项目描述",
    Columns:  "待办,进行中,已完成",
})

// 创建任务
task, err := client.CreateTask(dootask.CreateTaskRequest{
    ProjectID: project.ID,
    Name:      "新任务",
    Content:   "任务内容",
    Owner:     []int{123},
})

// 更新任务
updatedTask, err := client.UpdateTask(dootask.UpdateTaskRequest{
    TaskID:  task.ID,
    Name:    "更新后的任务名",
    Content: "更新后的内容",
})
```

### 群组管理

```go
// 创建群组
group, err := client.CreateGroup(dootask.CreateGroupRequest{
    ChatName: "新群组",
    UserIDs:  []int{123, 456, 789},
})

// 添加群成员
err := client.AddGroupUser(dootask.AddGroupUserRequest{
    DialogID: group.ID,
    UserIDs:  []int{999},
})
```

## API 方法列表

### 基础方法

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `NewClient` | 创建客户端实例 | `token`, `...ClientOption` | `*Client` |
| `WithServer` | 设置服务器地址 | `server` | `ClientOption` |
| `WithTimeout` | 设置超时时间 | `timeout` | `ClientOption` |

### 用户相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `GetUserInfo` | 获取用户信息 | `noCache ...bool` | `*UserInfo, error` |
| `CheckUserIdentity` | 检查用户身份 | `identity string` | `*UserInfo, error` |
| `GetUserDepartments` | 获取用户部门信息 | - | `[]Department, error` |
| `GetUsersBasic` | 获取多个用户基础信息 | `userids []int` | `[]UserBasic, error` |
| `GetUserBasic` | 获取单个用户基础信息 | `userid int` | `*UserBasic, error` |

### 机器人相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `GetBotList` | 获取机器人列表 | - | `*BotListResponse, error` |
| `GetBot` | 获取机器人信息 | `GetBotRequest` | `*Bot, error` |
| `CreateBot` | 创建机器人 | `CreateBotRequest` | `*Bot, error` |
| `UpdateBot` | 更新机器人 | `EditBotRequest` | `*Bot, error` |
| `DeleteBot` | 删除机器人 | `DeleteBotRequest` | `error` |

### 消息相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `SendMessage` | 发送消息 | `SendMessageRequest` | `error` |
| `SendMessageToUser` | 发送消息到用户 | `SendMessageToUserRequest` | `error` |
| `SendBotMessage` | 发送机器人消息 | `SendBotMessageRequest` | `error` |
| `SendAnonymousMessage` | 发送匿名消息 | `SendAnonymousMessageRequest` | `error` |
| `SendStreamMessage` | 通知成员监听消息 | `SendStreamMessageRequest` | `error` |
| `SendNoticeMessage` | 发送通知 | `SendNoticeMessageRequest` | `error` |
| `SendTemplateMessage` | 发送模板消息 | `SendTemplateMessageRequest` | `error` |
| `GetMessageList` | 获取消息列表 | `GetMessageListRequest` | `*DialogMessageListResponse, error` |
| `SearchMessage` | 搜索消息 | `SearchMessageRequest` | `*DialogMessageSearchResponse, error` |
| `GetMessage` | 获取单个消息详情 | `GetMessageRequest` | `*DialogMessage, error` |
| `GetMessageDetail` | 获取消息详情（兼容性） | `GetMessageRequest` | `*DialogMessage, error` |
| `WithdrawMessage` | 撤回消息 | `WithdrawMessageRequest` | `error` |
| `ForwardMessage` | 转发消息 | `ForwardMessageRequest` | `error` |
| `ToggleMessageTodo` | 切换消息待办状态 | `ToggleMessageTodoRequest` | `error` |
| `GetMessageTodoList` | 获取消息待办列表 | `GetMessageRequest` | `*TodoListResponse, error` |
| `MarkMessageDone` | 标记消息完成 | `MarkMessageDoneRequest` | `error` |
| `ConvertWebhookMessageToAI` | 转换webhook消息为AI对话格式 | `ConvertWebhookMessageRequest` | `*ConvertWebhookMessageResponse, error` |

### 对话相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `GetDialogList` | 获取对话列表 | `TimeRangeRequest` | `*ResponsePaginate[DialogInfo], error` |
| `SearchDialog` | 搜索会话 | `SearchDialogRequest` | `[]DialogInfo, error` |
| `GetDialogOne` | 获取单个会话信息 | `GetDialogRequest` | `*DialogInfo, error` |
| `GetDialogUser` | 获取会话成员 | `GetDialogUserRequest` | `[]DialogMember, error` |

### 群组相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `CreateGroup` | 创建群组 | `CreateGroupRequest` | `*DialogInfo, error` |
| `EditGroup` | 修改群组 | `EditGroupRequest` | `error` |
| `AddGroupUser` | 添加群成员 | `AddGroupUserRequest` | `error` |
| `RemoveGroupUser` | 移除群成员 | `RemoveGroupUserRequest` | `error` |
| `ExitGroup` | 退出群组 | `dialogID int` | `error` |
| `TransferGroup` | 转让群组 | `TransferGroupRequest` | `error` |
| `DisbandGroup` | 解散群组 | `DisbandGroupRequest` | `error` |

### 项目管理相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `GetProjectList` | 获取项目列表 | `GetProjectListRequest` | `*ResponsePaginate[Project], error` |
| `GetProject` | 获取项目信息 | `GetProjectRequest` | `*Project, error` |
| `CreateProject` | 创建项目 | `CreateProjectRequest` | `*Project, error` |
| `UpdateProject` | 更新项目 | `UpdateProjectRequest` | `*Project, error` |
| `ExitProject` | 退出项目 | `projectID int` | `error` |
| `DeleteProject` | 删除项目 | `projectID int` | `error` |

### 任务列表相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `GetColumnList` | 获取任务列表 | `GetColumnListRequest` | `*ResponsePaginate[ProjectColumn], error` |
| `CreateColumn` | 创建任务列表 | `CreateColumnRequest` | `*ProjectColumn, error` |
| `UpdateColumn` | 更新任务列表 | `UpdateColumnRequest` | `*ProjectColumn, error` |
| `DeleteColumn` | 删除任务列表 | `columnID int` | `error` |

### 任务相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `GetTaskList` | 获取任务列表 | `GetTaskListRequest` | `*ResponsePaginate[ProjectTask], error` |
| `GetTask` | 获取任务信息 | `GetTaskRequest` | `*ProjectTask, error` |
| `GetTaskContent` | 获取任务内容 | `GetTaskContentRequest` | `*TaskContent, error` |
| `GetTaskFiles` | 获取任务文件列表 | `GetTaskFilesRequest` | `[]TaskFile, error` |
| `CreateTask` | 创建任务 | `CreateTaskRequest` | `*ProjectTask, error` |
| `CreateSubTask` | 创建子任务 | `CreateSubTaskRequest` | `*ProjectTask, error` |
| `UpdateTask` | 更新任务 | `UpdateTaskRequest` | `*ProjectTask, error` |
| `CreateTaskDialog` | 创建任务对话 | `CreateTaskDialogRequest` | `*CreateTaskDialogResponse, error` |
| `ArchiveTask` | 归档任务 | `taskID int, archiveType string` | `error` |
| `DeleteTask` | 删除任务 | `taskID int, deleteType string` | `error` |

### 系统相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `GetSystemSettings` | 获取系统设置 | - | `*SystemSettings, error` |
| `GetVersion` | 获取版本信息 | - | `*VersionInfo, error` |

## 主要数据类型

### 基础类型
- `Client` - 客户端实例
- `Response[T]` - 基础响应结构
- `ResponsePaginate[T]` - 分页响应结构

### 用户相关
- `UserInfo` - 用户信息
- `UserBasic` - 用户基础信息
- `Department` - 部门信息

### 消息相关
- `SendMessageRequest` - 发送消息请求
- `DialogMessage` - 对话消息
- `DialogMessageListResponse` - 消息列表响应
- `TodoListResponse` - 待办列表响应

### 对话相关
- `DialogInfo` - 对话信息
- `DialogMember` - 对话成员

### 项目和任务相关
- `Project` - 项目信息
- `ProjectColumn` - 项目列表
- `ProjectTask` - 项目任务
- `TaskFile` - 任务文件
- `TaskContent` - 任务内容

### 系统相关
- `SystemSettings` - 系统设置
- `VersionInfo` - 版本信息

## 错误处理

```go
user, err := client.GetUserInfo()
if err != nil {
    fmt.Printf("获取用户信息失败: %v\n", err)
    return
}
```

## 缓存机制

客户端内置用户信息缓存机制，默认缓存时间为10分钟：

```go
// 强制刷新缓存
user, err := client.GetUserInfo(true)

// 使用缓存（默认）
user, err := client.GetUserInfo()
```

## 测试

```bash
cd server/go/test
go test -v
```

## 许可证

MIT License
