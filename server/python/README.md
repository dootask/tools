# DooTask Tools

一个用于与 DooTask 系统交互的 Python 客户端库，提供了完整的 API 封装和类型支持。

## 安装

```bash
pip install dootask-tools
```

## 快速开始

### 初始化客户端

```python
from dootask import DooTaskClient

# 创建客户端
client = DooTaskClient(
    token="your_token_here",
    server="https://your-dootask-server.com"
)

# 获取用户信息
user = client.get_user_info()
print(f"用户: {user.nickname}")
```

### 发送消息示例

```python
from dootask import SendMessageRequest, SendMessageToUserRequest

# 发送消息到指定对话
client.send_message(SendMessageRequest(
    dialog_id=123,
    text="Hello, World!",
    text_type="md"
))

# 发送消息到用户
client.send_message_to_user(SendMessageToUserRequest(
    userid=456,
    text="私信内容",
    text_type="md"
))
```

### 项目管理示例

```python
from dootask import GetProjectListRequest, CreateProjectRequest, CreateTaskRequest

# 获取项目列表
projects = client.get_project_list(GetProjectListRequest(
    page=1,
    pagesize=20
))

# 创建项目
project = client.create_project(CreateProjectRequest(
    name="新项目",
    desc="项目描述"
))

# 创建任务
task = client.create_task(CreateTaskRequest(
    project_id=project.id,
    name="新任务",
    content="任务内容"
))
```

## API 方法列表

### 客户端配置

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `DooTaskClient` | 创建客户端实例 | `token, server, timeout` | `DooTaskClient` |

### 用户相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `get_user_info` | 获取用户信息 | `no_cache=False` | `UserInfo` |
| `check_user_identity` | 检查用户身份 | `identity` | `UserInfo` |
| `get_user_departments` | 获取用户部门信息 | - | `List[Department]` |
| `get_users_basic` | 获取多个用户基础信息 | `userids: List[int]` | `List[UserBasic]` |
| `get_user_basic` | 获取单个用户基础信息 | `userid: int` | `UserBasic` |

### 消息相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `send_message` | 发送消息 | `SendMessageRequest` | `None` |
| `send_message_to_user` | 发送消息到用户 | `SendMessageToUserRequest` | `None` |
| `send_bot_message` | 发送机器人消息 | `SendBotMessageRequest` | `None` |
| `send_anonymous_message` | 发送匿名消息 | `SendAnonymousMessageRequest` | `None` |

### 对话相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `get_dialog_list` | 获取对话列表 | `TimeRangeRequest` | `ResponsePaginate[DialogInfo]` |
| `search_dialog` | 搜索会话 | `SearchDialogRequest` | `List[DialogInfo]` |
| `get_dialog_one` | 获取单个会话信息 | `GetDialogRequest` | `DialogInfo` |
| `get_dialog_user` | 获取会话成员 | `GetDialogUserRequest` | `List[DialogMember]` |

### 群组相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `create_group` | 创建群组 | `CreateGroupRequest` | `DialogInfo` |
| `edit_group` | 修改群组 | `EditGroupRequest` | `None` |
| `add_group_user` | 添加群成员 | `AddGroupUserRequest` | `None` |
| `remove_group_user` | 移除群成员 | `RemoveGroupUserRequest` | `None` |
| `exit_group` | 退出群组 | `dialog_id: int` | `None` |
| `transfer_group` | 转让群组 | `TransferGroupRequest` | `None` |
| `disband_group` | 解散群组 | `DisbandGroupRequest` | `None` |

### 项目管理相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `get_project_list` | 获取项目列表 | `GetProjectListRequest` | `ResponsePaginate[Project]` |
| `get_project` | 获取项目信息 | `GetProjectRequest` | `Project` |
| `create_project` | 创建项目 | `CreateProjectRequest` | `Project` |
| `update_project` | 更新项目 | `UpdateProjectRequest` | `Project` |
| `exit_project` | 退出项目 | `project_id: int` | `None` |
| `delete_project` | 删除项目 | `project_id: int` | `None` |

### 任务列表相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `get_column_list` | 获取任务列表 | `GetColumnListRequest` | `ResponsePaginate[ProjectColumn]` |
| `create_column` | 创建任务列表 | `CreateColumnRequest` | `ProjectColumn` |
| `update_column` | 更新任务列表 | `UpdateColumnRequest` | `ProjectColumn` |
| `delete_column` | 删除任务列表 | `column_id: int` | `None` |

### 任务相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `get_task_list` | 获取任务列表 | `GetTaskListRequest` | `ResponsePaginate[ProjectTask]` |
| `get_task` | 获取任务信息 | `GetTaskRequest` | `ProjectTask` |
| `get_task_content` | 获取任务内容 | `GetTaskContentRequest` | `TaskContent` |
| `get_task_files` | 获取任务文件列表 | `GetTaskFilesRequest` | `List[TaskFile]` |
| `create_task` | 创建任务 | `CreateTaskRequest` | `ProjectTask` |
| `create_sub_task` | 创建子任务 | `CreateSubTaskRequest` | `ProjectTask` |
| `update_task` | 更新任务 | `UpdateTaskRequest` | `ProjectTask` |
| `create_task_dialog` | 创建任务对话 | `CreateTaskDialogRequest` | `CreateTaskDialogResponse` |
| `archive_task` | 归档任务 | `task_id: int, archive_type: str` | `None` |
| `delete_task` | 删除任务 | `task_id: int, delete_type: str` | `None` |

### 系统设置相关接口

| 方法 | 描述 | 参数 | 返回值 |
|------|------|------|--------|
| `get_system_settings` | 获取系统设置 | - | `SystemSettings` |
| `get_version` | 获取版本信息 | - | `VersionInfo` |

## 主要数据类型

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

## 详细示例

```python
from dootask import (
    DooTaskClient,
    CreateProjectRequest,
    CreateTaskRequest,
    DooTaskException
)

client = DooTaskClient(
    token="your_token",
    server="https://your-server.com"
)

try:
    # 创建项目
    project = client.create_project(CreateProjectRequest(
        name="我的项目",
        desc="项目描述"
    ))
    
    # 创建任务
    task = client.create_task(CreateTaskRequest(
        project_id=project.id,
        name="任务名称",
        content="任务内容",
        owner=[user.userid]
    ))
    
    print(f"项目创建成功: {project.name}")
    print(f"任务创建成功: {task.name}")
    
except DooTaskException as e:
    print(f"错误: {e}")
```

## 异常处理

所有方法都可能抛出异常，包含详细的错误信息：

```python
from dootask import (
    DooTaskException,
    DooTaskAPIException,
    DooTaskHTTPException,
    DooTaskAuthException,
    DooTaskPermissionException
)

try:
    user = client.get_user_info()
except DooTaskAuthException:
    print("认证失败，请检查 token")
except DooTaskPermissionException:
    print("权限不足")
except DooTaskAPIException as e:
    print(f"API 错误: {e}")
except DooTaskHTTPException as e:
    print(f"网络错误: {e}")
except DooTaskException as e:
    print(f"其他错误: {e}")
```

## 许可证

MIT License 