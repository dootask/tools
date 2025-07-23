"""
DooTask Tools - Python 客户端包

一个用于与 DooTask 系统交互的 Python 客户端库。
"""

from .client import DooTaskClient
from .models import *
from .exceptions import *

__version__ = "0.0.5"
__author__ = "DooTask Team"
__email__ = "support@dootask.com"

# 导出最常用的类
__all__ = [
    # 主要客户端类
    'DooTaskClient',
    
    # 异常类
    'DooTaskException',
    'DooTaskAPIException', 
    'DooTaskHTTPException',
    'DooTaskAuthException',
    'DooTaskPermissionException',
    
    # 常用请求类
    'SendMessageRequest',
    'SendMessageToUserRequest',
    'SendBotMessageRequest',
    'CreateProjectRequest',
    'CreateTaskRequest',
    'CreateGroupRequest',
    
    # 常用响应类
    'UserInfo',
    'UserBasic',
    'Project',
    'ProjectTask',
    'DialogInfo',
] 