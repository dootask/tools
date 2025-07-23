"""
DooTask Tools 数据类型定义
"""

from dataclasses import dataclass, field
from typing import Optional, List, Dict, Any, Union, Generic, TypeVar
from datetime import datetime

T = TypeVar('T')

# ------------------------------------------------------------------------------------------
# 基础结构定义
# ------------------------------------------------------------------------------------------

@dataclass
class Response(Generic[T]):
    """基础响应结构"""
    ret: int
    msg: str
    data: Optional[T] = None

@dataclass
class ResponsePaginate(Generic[T]):
    """分页数据"""
    current_page: int
    data: List[T]
    next_page_url: Optional[str] = None
    path: str = ""
    per_page: Union[int, str] = 0
    prev_page_url: Optional[str] = None
    to: Union[int, str] = 0
    total: int = 0

# ------------------------------------------------------------------------------------------
# 用户相关结构
# ------------------------------------------------------------------------------------------

@dataclass
class UserInfo:
    """用户信息"""
    userid: int
    identity: List[str] = field(default_factory=list)
    email: str = ""
    nickname: str = ""
    profession: str = ""
    userimg: str = ""
    bot: int = 0
    department: List[int] = field(default_factory=list)
    department_name: str = ""

@dataclass
class UserBasic:
    """用户基础信息"""
    userid: int
    email: str = ""
    nickname: str = ""
    profession: str = ""
    userimg: str = ""
    bot: int = 0
    online: bool = False
    department: List[int] = field(default_factory=list)
    department_name: str = ""

@dataclass
class Department:
    """部门信息"""
    id: int
    name: str
    parent_id: int = 0
    owner_userid: int = 0

# ------------------------------------------------------------------------------------------
# 机器人相关结构
# ------------------------------------------------------------------------------------------

@dataclass
class Bot:
    """机器人信息"""
    id: int
    name: str = ""
    avatar: str = ""
    clear_day: int = 0
    webhook_url: str = ""

@dataclass
class BotListResponse:
    """机器人列表响应"""
    list: List[Bot] = field(default_factory=list)

@dataclass
class GetBotRequest:
    """获取机器人请求"""
    id: int

@dataclass
class CreateBotRequest:
    """创建机器人请求"""
    name: str
    avatar: str = ""
    clear_day: int = 7
    webhook_url: str = ""

@dataclass
class EditBotRequest:
    """编辑机器人请求"""
    id: int
    name: str = ""
    avatar: str = ""
    clear_day: int = 0
    webhook_url: str = ""

@dataclass
class DeleteBotRequest:
    """删除机器人请求"""
    id: int
    remark: str


# ------------------------------------------------------------------------------------------
# 对话相关结构
# ------------------------------------------------------------------------------------------

@dataclass
class DialogUserResponse:
    """对话用户信息"""
    dialog_id: int
    userid: int
    bot: int = 0

@dataclass
class DialogOpenUserResponse:
    """打开用户对话响应数据"""
    dialog_user: DialogUserResponse

@dataclass
class DialogInfo:
    """对话信息"""
    id: int
    type: str = ""
    group_type: str = ""
    name: str = ""
    avatar: str = ""
    owner_id: int = 0
    created_at: str = ""
    updated_at: str = ""
    last_at: str = ""
    mark_unread: int = 0
    silence: int = 0
    hide: int = 0
    color: str = ""
    unread: int = 0
    unread_one: int = 0
    mention: int = 0
    mention_ids: List[int] = field(default_factory=list)
    people: int = 0
    people_user: int = 0
    people_bot: int = 0
    todo_num: int = 0
    last_msg: Any = None
    pinyin: str = ""
    bot: int = 0
    top_at: str = ""

@dataclass
class DialogMember:
    """会话成员信息"""
    id: int
    dialog_id: int
    userid: int
    nickname: str = ""
    email: str = ""
    userimg: str = ""
    bot: int = 0
    online: bool = False

@dataclass
class TimeRangeRequest:
    """时间范围请求参数"""
    timerange: str = ""
    page: int = 1
    pagesize: int = 50

@dataclass
class SearchDialogRequest:
    """搜索会话请求"""
    key: str

@dataclass
class GetDialogRequest:
    """获取单个会话请求"""
    dialog_id: int

@dataclass
class GetDialogUserRequest:
    """获取会话成员请求"""
    dialog_id: int
    getuser: int = 0
    
# ------------------------------------------------------------------------------------------
# 消息相关结构
# ------------------------------------------------------------------------------------------

@dataclass
class SendMessageRequest:
    """发送消息请求"""
    dialog_id: int
    text: str
    text_type: str = "md"
    silence: bool = False

@dataclass
class SendMessageToUserRequest:
    """发送消息到用户请求"""
    userid: int
    text: str
    text_type: str = "md"
    silence: bool = False

@dataclass
class SendBotMessageRequest:
    """发送机器人消息请求"""
    userid: int
    text: str
    bot_type: str = "system-msg"
    bot_name: str = ""
    silence: bool = False

@dataclass
class SendAnonymousMessageRequest:
    """发送匿名消息请求"""
    userid: int
    text: str

@dataclass
class DialogMessage:
    """对话消息"""
    id: int
    dialog_id: int
    userid: int
    bot: int = 0
    created_at: str = ""
    type: str = ""
    mtype: str = ""
    msg: Any = None
    reply_id: int = 0
    reply_num: int = 0
    forward_id: int = 0
    forward_num: int = 0
    tag: int = 0
    todo: int = 0
    read: int = 0
    send: int = 0
    read_at: Optional[str] = None
    mention: int = 0
    dot: int = 0
    emoji: List[Any] = field(default_factory=list)
    link: int = 0
    modify: int = 0
    percentage: int = 100

@dataclass
class DialogMessageListResponse:
    """消息列表响应"""
    list: List[DialogMessage] = field(default_factory=list)
    time: int = 0
    dialog: Optional[DialogInfo] = None
    todo: List[Any] = field(default_factory=list)
    top: Any = None

@dataclass
class DialogMessageSearchResponse:
    """搜索消息响应"""
    data: List[int] = field(default_factory=list)

@dataclass
class TodoUser:
    """待办用户"""
    userid: int
    nickname: str = ""
    userimg: str = ""
    done: bool = False
    done_at: str = ""

@dataclass
class TodoListResponse:
    """待办列表响应"""
    users: List[TodoUser] = field(default_factory=list)

@dataclass
class GetMessageListRequest:
    """获取消息列表请求"""
    dialog_id: int
    msg_id: int = 0
    position_id: int = 0
    prev_id: int = 0
    next_id: int = 0
    msg_type: str = ""
    take: int = 50

@dataclass
class SearchMessageRequest:
    """搜索消息请求"""
    dialog_id: int
    key: str

@dataclass
class GetMessageRequest:
    """获取单个消息请求"""
    msg_id: int

@dataclass
class WithdrawMessageRequest:
    """撤回消息请求"""
    msg_id: int

@dataclass
class ForwardMessageRequest:
    """转发消息请求"""
    msg_id: int
    dialogids: List[int] = field(default_factory=list)
    userids: List[int] = field(default_factory=list)
    show_source: int = 0
    leave_message: str = ""

@dataclass
class ToggleMessageTodoRequest:
    """切换消息待办请求"""
    msg_id: int
    type: str = "all"
    userids: List[int] = field(default_factory=list)

@dataclass
class MarkMessageDoneRequest:
    """标记消息完成请求"""
    msg_id: int

# ------------------------------------------------------------------------------------------
# 群组相关结构
# ------------------------------------------------------------------------------------------

@dataclass
class CreateGroupRequest:
    """创建群组请求"""
    userids: List[int]
    avatar: str = ""
    chat_name: str = ""

@dataclass
class EditGroupRequest:
    """修改群组请求"""
    dialog_id: int
    avatar: str = ""
    chat_name: str = ""
    admin: int = 0

@dataclass
class AddGroupUserRequest:
    """添加群成员请求"""
    dialog_id: int
    userids: List[int]

@dataclass
class RemoveGroupUserRequest:
    """移除群成员请求"""
    dialog_id: int
    userids: List[int] = field(default_factory=list)

@dataclass
class TransferGroupRequest:
    """转让群组请求"""
    dialog_id: int
    userid: int
    check_owner: str = ""
    key: str = ""

@dataclass
class DisbandGroupRequest:
    """解散群组请求"""
    dialog_id: int

# ------------------------------------------------------------------------------------------
# 项目管理相关结构
# ------------------------------------------------------------------------------------------

@dataclass
class Project:
    """项目信息"""
    id: int
    name: str = ""
    desc: str = ""
    userid: int = 0
    dialog_id: int = 0
    archived_at: str = ""
    created_at: str = ""
    updated_at: str = ""
    owner: int = 0
    owner_userid: int = 0
    personal: int = 0
    # 任务统计
    task_num: int = 0
    task_complete: int = 0
    task_percent: int = 0
    task_my_num: int = 0
    task_my_complete: int = 0
    task_my_percent: int = 0

@dataclass
class GetProjectListRequest:
    """获取项目列表请求"""
    type: str = "all"
    archived: str = "no"
    getcolumn: str = "no"
    getuserid: str = "no"
    getstatistics: str = "no"
    timerange: str = ""
    page: int = 1
    pagesize: int = 50

@dataclass
class GetProjectRequest:
    """获取项目信息请求"""
    project_id: int

@dataclass
class CreateProjectRequest:
    """创建项目请求"""
    name: str
    desc: str = ""
    columns: str = ""
    flow: str = ""
    personal: int = 0

@dataclass
class UpdateProjectRequest:
    """更新项目请求"""
    project_id: int
    name: str
    desc: str = ""
    archive_method: str = ""
    archive_days: int = 0

@dataclass
class ProjectActionRequest:
    """项目操作请求"""
    project_id: int
    type: str = ""

# ------------------------------------------------------------------------------------------
# 任务列表相关结构
# ------------------------------------------------------------------------------------------

@dataclass
class ProjectColumn:
    """项目列表"""
    id: int
    project_id: int
    name: str = ""
    color: str = ""
    sort: int = 0
    created_at: str = ""
    updated_at: str = ""

@dataclass
class GetColumnListRequest:
    """获取列表请求"""
    project_id: int
    page: int = 1
    pagesize: int = 100

@dataclass
class CreateColumnRequest:
    """创建列表请求"""
    project_id: int
    name: str

@dataclass
class UpdateColumnRequest:
    """更新列表请求"""
    column_id: int
    name: str = ""
    color: str = ""

@dataclass
class ColumnActionRequest:
    """列表操作请求"""
    column_id: int

# ------------------------------------------------------------------------------------------
# 任务相关结构
# ------------------------------------------------------------------------------------------

@dataclass
class ProjectTask:
    """项目任务"""
    id: int
    project_id: int = 0
    column_id: int = 0
    parent_id: int = 0
    name: str = ""
    desc: str = ""
    start_at: str = ""
    end_at: str = ""
    complete_at: str = ""
    archived_at: str = ""
    created_at: str = ""
    updated_at: str = ""
    userid: int = 0
    dialog_id: int = 0
    flow_item_id: int = 0
    flow_item_name: str = ""
    visibility: int = 0
    color: str = ""
    # 统计信息
    file_num: int = 0
    msg_num: int = 0
    sub_num: int = 0
    sub_complete: int = 0
    percent: int = 0
    # 关联数据
    project_name: str = ""
    column_name: str = ""

@dataclass
class TaskFile:
    """任务文件"""
    id: int
    task_id: int
    name: str = ""
    ext: str = ""
    size: int = 0
    path: str = ""
    thumb: str = ""
    userid: int = 0
    created_at: str = ""
    updated_at: str = ""

@dataclass
class TaskContent:
    """任务内容"""
    content: str = ""
    type: str = ""

@dataclass
class GetTaskListRequest:
    """获取任务列表请求"""
    project_id: int = 0
    parent_id: int = 0
    archived: str = "no"
    deleted: str = "no"
    timerange: str = ""
    page: int = 1
    pagesize: int = 100

@dataclass
class GetTaskRequest:
    """获取任务信息请求"""
    task_id: int
    archived: str = "no"

@dataclass
class GetTaskContentRequest:
    """获取任务内容请求"""
    task_id: int
    history_id: int = 0

@dataclass
class GetTaskFilesRequest:
    """获取任务文件请求"""
    task_id: int

@dataclass
class CreateTaskRequest:
    """创建任务请求"""
    project_id: int
    name: str
    column_id: Union[int, str] = 0
    content: str = ""
    times: List[str] = field(default_factory=list)
    owner: List[int] = field(default_factory=list)
    top: int = 0

@dataclass
class CreateSubTaskRequest:
    """创建子任务请求"""
    task_id: int
    name: str

@dataclass
class UpdateTaskRequest:
    """更新任务请求"""
    task_id: int
    name: str = ""
    content: str = ""
    times: List[str] = field(default_factory=list)
    owner: List[int] = field(default_factory=list)
    assist: List[int] = field(default_factory=list)
    color: str = ""
    visibility: int = 0
    complete_at: Any = None

@dataclass
class TaskActionRequest:
    """任务操作请求"""
    task_id: int
    type: str = ""

@dataclass
class CreateTaskDialogRequest:
    """创建任务对话请求"""
    task_id: int

@dataclass
class CreateTaskDialogResponse:
    """创建任务对话响应"""
    id: int
    dialog_id: int
    dialog_data: Any = None

# ------------------------------------------------------------------------------------------
# 系统相关结构
# ------------------------------------------------------------------------------------------

@dataclass
class SystemSettings:
    """系统设置"""
    reg: Optional[str] = None
    task_default_time: Optional[List[str]] = None
    system_alias: Optional[str] = None
    system_welcome: str = ""
    server_timezone: Optional[str] = None
    server_version: Optional[str] = None

@dataclass
class VersionInfo:
    """版本信息"""
    device_count: int = 0
    version: str = "" 