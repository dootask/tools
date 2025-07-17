# DooTask Go 客户端

DooTask Go 客户端库，用于与 DooTask 系统进行交互。

## 安装

```bash
go get github.com/dootask/tools/server/go
```

## 快速开始

### 初始化客户端

```go
package main

import (
    "fmt"
    
    dootask "github.com/dootask/tools/server/go"
)

func main() {
    // 创建客户端
    client := dootask.NewClient(
        "your-token-here",
        dootask.WithServer("http://your-server.com"),
    )
    
    // 获取用户信息
    user, err := client.GetUserInfo()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("用户: %s\n", user.Nickname)
}
```

### 发送消息示例

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

### 项目管理示例

```go
// 获取项目列表
projects, err := client.GetProjectList(dootask.GetProjectListRequest{
    Page: 1,
    Size: 20,
})

// 创建项目
project, err := client.CreateProject(dootask.CreateProjectRequest{
    Name:        "新项目",
    Description: "项目描述",
})

// 创建任务
task, err := client.CreateTask(dootask.CreateTaskRequest{
    ProjectID: project.ID,
    ColumnID:  1,
    Name:      "新任务",
    Content:   "任务内容",
})
```

## API 方法列表

### 客户端配置

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `NewClient` | 创建客户端实例 | `token`, `...ClientOption` | `*Client` |
| `WithServer` | 设置服务器地址 | `server` | `ClientOption` |

### 用户相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `GetUserInfo` | 获取用户信息 | `noCache` | `*UserInfo, error` |
| `CheckUserIdentity` | 检查用户身份 | `identity` | `*UserInfo, error` |
| `GetUserDepartments` | 获取用户部门信息 | - | `[]Department, error` |
| `GetUsersBasic` | 获取多个用户基础信息 | `userids []int` | `[]UserBasic, error` |
| `GetUserBasic` | 获取单个用户基础信息 | `userid` | `*UserBasic, error` |

### 消息相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `SendMessage` | 发送消息 | `SendMessageRequest` | `error` |
| `SendMessageToUser` | 发送消息到用户 | `SendMessageToUserRequest` | `error` |
| `SendBotMessage` | 发送机器人消息 | `SendBotMessageRequest` | `error` |
| `SendAnonymousMessage` | 发送匿名消息 | `SendAnonymousMessageRequest` | `error` |

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
| `ExitGroup` | 退出群组 | `dialogID` | `error` |
| `TransferGroup` | 转让群组 | `TransferGroupRequest` | `error` |
| `DisbandGroup` | 解散群组 | `DisbandGroupRequest` | `error` |

### 项目管理相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `GetProjectList` | 获取项目列表 | `GetProjectListRequest` | `*ResponsePaginate[Project], error` |
| `GetProject` | 获取项目信息 | `GetProjectRequest` | `*Project, error` |
| `CreateProject` | 创建项目 | `CreateProjectRequest` | `*Project, error` |
| `UpdateProject` | 更新项目 | `UpdateProjectRequest` | `*Project, error` |
| `ExitProject` | 退出项目 | `projectID` | `error` |
| `DeleteProject` | 删除项目 | `projectID` | `error` |

### 任务列表相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `GetColumnList` | 获取任务列表 | `GetColumnListRequest` | `*ResponsePaginate[ProjectColumn], error` |
| `CreateColumn` | 创建任务列表 | `CreateColumnRequest` | `*ProjectColumn, error` |
| `UpdateColumn` | 更新任务列表 | `UpdateColumnRequest` | `*ProjectColumn, error` |
| `DeleteColumn` | 删除任务列表 | `columnID` | `error` |

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
| `ArchiveTask` | 归档任务 | `taskID, archiveType` | `error` |
| `DeleteTask` | 删除任务 | `taskID, deleteType` | `error` |

### 系统设置相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `GetSystemSettings` | 获取系统设置 | - | `*SystemSettings, error` |
| `GetVersion` | 获取版本信息 | - | `*VersionInfo, error` |

### 主要数据类型

- `UserInfo` - 用户信息
- `UserBasic` - 用户基础信息
- `Department` - 部门信息
- `DialogInfo` - 对话信息
- `DialogMember` - 对话成员
- `Project` - 项目信息
- `ProjectColumn` - 项目列表
- `ProjectTask` - 项目任务
- `TaskFile` - 任务文件
- `TaskContent` - 任务内容
- `SystemSettings` - 系统设置
- `VersionInfo` - 版本信息

## 错误处理

所有方法都返回 `error` 类型，包含详细的错误信息：

```go
user, err := client.GetUserInfo()
if err != nil {
    // 处理错误
    fmt.Printf("获取用户信息失败: %v\n", err)
    return
}
```

## 测试

```bash
cd server/go/test
go test -v
```

## 许可证

MIT License
