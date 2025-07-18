"""
DooTask Tools 客户端
"""

import json
import time
from typing import Optional, List, Dict, Any, Union, TypeVar, Type
from datetime import datetime, timedelta
from dataclasses import asdict, is_dataclass
from urllib.parse import urlencode, quote

import requests

from .models import *
from .exceptions import *

T = TypeVar('T')

class DooTaskClient:
    """DooTask 客户端"""
    
    def __init__(self, token: str, server: str = "http://nginx", timeout: int = 10):
        """
        初始化客户端
        
        Args:
            token: 认证令牌
            server: 服务器地址
            timeout: 请求超时时间（秒）
        """
        self.token = token
        self.server = server.rstrip('/')
        self.timeout = timeout
        self._cache: Dict[str, Dict[str, Any]] = {}
        self._cache_time = 600  # 10分钟缓存
        
        # 创建会话
        self.session = requests.Session()
        self.session.headers.update({
            'Token': token,
            'User-Agent': 'DooTask-Tools/1.0',
            'Content-Type': 'application/json'
        })
    
    def _build_url(self, base_url: str, params: Optional[Dict[str, Any]] = None) -> str:
        """构建带查询参数的URL"""
        if not params:
            return base_url
        
        query_params = []
        for key, value in params.items():
            if value is None:
                continue
            
            if isinstance(value, bool):
                query_params.append(f"{key}={1 if value else 0}")
            elif isinstance(value, str):
                if value:
                    query_params.append(f"{key}={quote(value)}")
            elif isinstance(value, (int, float)):
                query_params.append(f"{key}={value}")
            elif isinstance(value, list):
                for item in value:
                    if isinstance(item, str):
                        query_params.append(f"{key}[]={quote(item)}")
                    else:
                        query_params.append(f"{key}[]={item}")
            else:
                query_params.append(f"{key}={quote(str(value))}")
        
        if not query_params:
            return base_url
        
        separator = "&" if "?" in base_url else "?"
        return base_url + separator + "&".join(query_params)
    
    def _dataclass_to_dict(self, obj: Any) -> Dict[str, Any]:
        """将数据类转换为字典"""
        if obj is None:
            return {}
        
        if is_dataclass(obj):
            return asdict(obj)
        
        if isinstance(obj, dict):
            return obj
        
        # 如果是普通对象，尝试转换为字典
        if hasattr(obj, '__dict__'):
            return obj.__dict__
        
        return {}
    
    def _make_request(self, method: str, api: str, request_data: Any = None, 
                      response_type: Optional[Type[T]] = None, 
                      headers: Optional[Dict[str, Any]] = None) -> Optional[T]:
        """发送请求"""
        url = self.server + api
        
        # 设置请求头
        request_headers = self.session.headers.copy()
        if headers:
            request_headers.update(headers)
        
        # 处理请求数据
        if method.upper() == 'GET':
            # GET 请求：将数据作为查询参数
            if request_data:
                params = self._dataclass_to_dict(request_data)
                url = self._build_url(url, params)
            
            response = self.session.get(url, timeout=self.timeout)
        
        elif method.upper() in ['POST', 'PUT', 'PATCH']:
            # POST/PUT/PATCH 请求：将数据作为 JSON body
            json_data = None
            if request_data:
                json_data = self._dataclass_to_dict(request_data)
            
            response = self.session.request(
                method, url, json=json_data, 
                headers=request_headers, timeout=self.timeout
            )
        
        elif method.upper() == 'DELETE':
            # DELETE 请求：支持查询参数
            if request_data:
                params = self._dataclass_to_dict(request_data)
                url = self._build_url(url, params)
            
            response = self.session.delete(url, timeout=self.timeout)
        
        else:
            raise DooTaskException(f"不支持的 HTTP 方法: {method}")
        
        # 检查 HTTP 状态码
        if response.status_code != 200:
            raise DooTaskHTTPException(
                f"HTTP {response.status_code}: {response.reason}, body: {response.text}",
                response.status_code
            )
        
        # 解析响应
        try:
            api_response = response.json()
        except (json.JSONDecodeError, ValueError):
            raise DooTaskException("响应不是有效的 JSON 格式")
        
        # 检查业务状态
        if api_response.get('ret') != 1:
            error_msg = api_response.get('msg', f"API 错误: {api_response.get('ret', 'unknown')}")
            raise DooTaskAPIException(error_msg, api_response.get('ret', 0))
        
        # 如果不需要响应数据，直接返回 None
        if response_type is None:
            return None
        
        # 解析数据到目标类型
        data = api_response.get('data')
        if data is None:
            return None
        
        # 如果是原始类型，直接返回
        if response_type in [str, int, float, bool, dict, list]:
            return data
        
        # 如果是数据类，进行转换
        if is_dataclass(response_type):
            if isinstance(data, dict):
                # 递归处理嵌套数据类
                return self._convert_to_dataclass(data, response_type)
            elif isinstance(data, list):
                return [self._convert_to_dataclass(item, response_type) 
                       if isinstance(item, dict) else item for item in data]
        
        return data
    
    def _convert_to_dataclass(self, data: Dict[str, Any], dataclass_type: Type[T]) -> T:
        """递归转换字典为数据类，只使用数据类中定义的字段"""
        from dataclasses import fields
        
        # 获取数据类的字段定义
        field_definitions = {field.name: field for field in fields(dataclass_type)}
        converted_data = {}
        
        for field_name, field_def in field_definitions.items():
            if field_name in data:
                field_value = data[field_name]
                field_type = field_def.type
                
                # 如果字段类型是数据类，递归转换
                if is_dataclass(field_type) and isinstance(field_value, dict):
                    converted_data[field_name] = self._convert_to_dataclass(field_value, field_type)
                else:
                    converted_data[field_name] = field_value
        
        return dataclass_type(**converted_data)
    
    def _get_request(self, api: str, request_data: Any = None, 
                     response_type: Optional[Type[T]] = None,
                     headers: Optional[Dict[str, Any]] = None) -> Optional[T]:
        """发送 GET 请求"""
        return self._make_request('GET', api, request_data, response_type, headers)
    
    def _post_request(self, api: str, request_data: Any = None, 
                      response_type: Optional[Type[T]] = None) -> Optional[T]:
        """发送 POST 请求"""
        return self._make_request('POST', api, request_data, response_type)
    
    # ------------------------------------------------------------------------------------------
    # 用户相关接口
    # ------------------------------------------------------------------------------------------
    
    def get_user_info(self, no_cache: bool = False) -> UserInfo:
        """获取用户信息"""
        cache_key = f"user_info_{self.token}"
        
        # 检查缓存
        if not no_cache and cache_key in self._cache:
            cache_data = self._cache[cache_key]
            if time.time() < cache_data['expires_at']:
                return cache_data['data']
        
        # 获取用户信息
        user_info = self._get_request('/api/users/info', response_type=UserInfo)
        
        # 更新缓存
        self._cache[cache_key] = {
            'data': user_info,
            'expires_at': time.time() + self._cache_time
        }
        
        return user_info
    
    def check_user_identity(self, identity: str) -> UserInfo:
        """检查用户是否具有指定身份"""
        user = self.get_user_info()
        
        if identity not in user.identity:
            raise DooTaskPermissionException("权限不足")
        
        return user
    
    def get_user_departments(self) -> List[Department]:
        """获取用户部门信息"""
        return self._get_request('/api/users/info/departments', response_type=List[Department])
    
    def get_users_basic(self, userids: List[int]) -> List[UserBasic]:
        """获取指定用户基础信息（支持多个用户）"""
        params = {'userid': userids}
        return self._get_request('/api/users/basic', params, response_type=List[UserBasic])
    
    def get_user_basic(self, userid: int) -> UserBasic:
        """获取指定用户基础信息（单个用户）"""
        users = self.get_users_basic([userid])
        if not users:
            raise DooTaskException("用户不存在")
        return users[0]
    
    # ------------------------------------------------------------------------------------------
    # 消息相关接口
    # ------------------------------------------------------------------------------------------
    
    def send_message(self, message: SendMessageRequest) -> None:
        """发送消息"""
        if not message.text_type:
            message.text_type = "md"
        
        self._post_request('/api/dialog/msg/sendtext', message)
    
    def send_message_to_user(self, message: SendMessageToUserRequest) -> None:
        """发送消息到用户"""
        # 获取用户对话ID
        query_params = {'userid': message.userid}
        response = self._get_request('/api/dialog/open/user', query_params, 
                                   response_type=DialogOpenUserResponse)
        
        # 发送消息
        self.send_message(SendMessageRequest(
            dialog_id=response.dialog_user.dialog_id,
            text=message.text,
            text_type=message.text_type,
            silence=message.silence
        ))
    
    def send_bot_message(self, message: SendBotMessageRequest) -> None:
        """发送机器人消息"""
        if not message.bot_type:
            message.bot_type = "system-msg"
        
        self._post_request('/api/dialog/msg/sendbot', message)
    
    def send_anonymous_message(self, message: SendAnonymousMessageRequest) -> None:
        """发送匿名消息"""
        self._post_request('/api/dialog/msg/sendanon', message)
    
    # ------------------------------------------------------------------------------------------
    # 对话相关接口
    # ------------------------------------------------------------------------------------------
    
    def get_dialog_list(self, params: Optional[TimeRangeRequest] = None) -> ResponsePaginate[DialogInfo]:
        """获取对话列表"""
        if params is None:
            params = TimeRangeRequest()
        return self._get_request('/api/dialog/lists', params, response_type=ResponsePaginate[DialogInfo])
    
    def search_dialog(self, params: SearchDialogRequest) -> List[DialogInfo]:
        """搜索会话"""
        return self._get_request('/api/dialog/search', params, response_type=List[DialogInfo])
    
    def get_dialog_one(self, params: GetDialogRequest) -> DialogInfo:
        """获取单个会话信息"""
        return self._get_request('/api/dialog/one', params, response_type=DialogInfo)
    
    def get_dialog_user(self, params: GetDialogUserRequest) -> List[DialogMember]:
        """获取会话成员"""
        return self._get_request('/api/dialog/user', params, response_type=List[DialogMember])
    
    # ------------------------------------------------------------------------------------------
    # 群组相关接口
    # ------------------------------------------------------------------------------------------
    
    def create_group(self, params: CreateGroupRequest) -> DialogInfo:
        """新增群组"""
        return self._get_request('/api/dialog/group/add', params, response_type=DialogInfo)
    
    def edit_group(self, params: EditGroupRequest) -> None:
        """修改群组"""
        self._get_request('/api/dialog/group/edit', params)
    
    def add_group_user(self, params: AddGroupUserRequest) -> None:
        """添加群成员"""
        self._get_request('/api/dialog/group/adduser', params)
    
    def remove_group_user(self, params: RemoveGroupUserRequest) -> None:
        """移除群成员"""
        self._get_request('/api/dialog/group/deluser', params)
    
    def exit_group(self, dialog_id: int) -> None:
        """退出群组"""
        self.remove_group_user(RemoveGroupUserRequest(dialog_id=dialog_id))
    
    def transfer_group(self, params: TransferGroupRequest) -> None:
        """转让群组"""
        self._get_request('/api/dialog/group/transfer', params)
    
    def disband_group(self, params: DisbandGroupRequest) -> None:
        """解散群组"""
        self._get_request('/api/dialog/group/disband', params)
    
    # ------------------------------------------------------------------------------------------
    # 项目管理相关接口
    # ------------------------------------------------------------------------------------------
    
    def get_project_list(self, params: Optional[GetProjectListRequest] = None) -> ResponsePaginate[Project]:
        """获取项目列表"""
        if params is None:
            params = GetProjectListRequest()
        return self._get_request('/api/project/lists', params, response_type=ResponsePaginate[Project])
    
    def get_project(self, params: GetProjectRequest) -> Project:
        """获取项目信息"""
        return self._get_request('/api/project/one', params, response_type=Project)
    
    def create_project(self, params: CreateProjectRequest) -> Project:
        """创建项目"""
        return self._get_request('/api/project/add', params, response_type=Project)
    
    def update_project(self, params: UpdateProjectRequest) -> Project:
        """更新项目"""
        return self._get_request('/api/project/update', params, response_type=Project)
    
    def exit_project(self, project_id: int) -> None:
        """退出项目"""
        params = ProjectActionRequest(project_id=project_id)
        self._get_request('/api/project/exit', params)
    
    def delete_project(self, project_id: int) -> None:
        """删除项目"""
        params = ProjectActionRequest(project_id=project_id)
        self._get_request('/api/project/remove', params)
    
    # ------------------------------------------------------------------------------------------
    # 任务列表相关接口
    # ------------------------------------------------------------------------------------------
    
    def get_column_list(self, params: GetColumnListRequest) -> ResponsePaginate[ProjectColumn]:
        """获取任务列表"""
        return self._get_request('/api/project/column/lists', params, 
                               response_type=ResponsePaginate[ProjectColumn])
    
    def create_column(self, params: CreateColumnRequest) -> ProjectColumn:
        """创建任务列表"""
        return self._get_request('/api/project/column/add', params, response_type=ProjectColumn)
    
    def update_column(self, params: UpdateColumnRequest) -> ProjectColumn:
        """更新任务列表"""
        return self._get_request('/api/project/column/update', params, response_type=ProjectColumn)
    
    def delete_column(self, column_id: int) -> None:
        """删除任务列表"""
        params = ColumnActionRequest(column_id=column_id)
        self._get_request('/api/project/column/remove', params)
    
    # ------------------------------------------------------------------------------------------
    # 任务相关接口
    # ------------------------------------------------------------------------------------------
    
    def get_task_list(self, params: Optional[GetTaskListRequest] = None) -> ResponsePaginate[ProjectTask]:
        """获取任务列表"""
        if params is None:
            params = GetTaskListRequest()
        return self._get_request('/api/project/task/lists', params, 
                               response_type=ResponsePaginate[ProjectTask])
    
    def get_task(self, params: GetTaskRequest) -> ProjectTask:
        """获取任务信息"""
        return self._get_request('/api/project/task/one', params, response_type=ProjectTask)
    
    def get_task_content(self, params: GetTaskContentRequest) -> TaskContent:
        """获取任务内容"""
        return self._get_request('/api/project/task/content', params, response_type=TaskContent)
    
    def get_task_files(self, params: GetTaskFilesRequest) -> List[TaskFile]:
        """获取任务文件列表"""
        return self._get_request('/api/project/task/files', params, response_type=List[TaskFile])
    
    def create_task(self, params: CreateTaskRequest) -> ProjectTask:
        """创建任务"""
        return self._post_request('/api/project/task/add', params, response_type=ProjectTask)
    
    def create_sub_task(self, params: CreateSubTaskRequest) -> ProjectTask:
        """创建子任务"""
        return self._get_request('/api/project/task/addsub', params, response_type=ProjectTask)
    
    def update_task(self, params: UpdateTaskRequest) -> ProjectTask:
        """更新任务"""
        return self._post_request('/api/project/task/update', params, response_type=ProjectTask)
    
    def create_task_dialog(self, params: CreateTaskDialogRequest) -> CreateTaskDialogResponse:
        """创建任务对话"""
        return self._get_request('/api/project/task/dialog', params, 
                               response_type=CreateTaskDialogResponse)
    
    def archive_task(self, task_id: int, archive_type: str = "add") -> None:
        """归档任务"""
        params = TaskActionRequest(task_id=task_id, type=archive_type)
        self._get_request('/api/project/task/archived', params)
    
    def delete_task(self, task_id: int, delete_type: str = "delete") -> None:
        """删除任务"""
        params = TaskActionRequest(task_id=task_id, type=delete_type)
        self._get_request('/api/project/task/remove', params)
    
    # ------------------------------------------------------------------------------------------
    # 系统设置相关方法
    # ------------------------------------------------------------------------------------------
    
    def get_system_settings(self) -> SystemSettings:
        """获取系统设置"""
        return self._get_request('/api/system/setting', response_type=SystemSettings)
    
    def get_version(self) -> VersionInfo:
        """获取版本信息"""
        headers = {'version': 'true'}
        return self._get_request('/api/system/version', headers=headers, response_type=VersionInfo) 