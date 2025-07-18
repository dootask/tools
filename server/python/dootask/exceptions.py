"""
DooTask Tools 异常定义
"""

class DooTaskException(Exception):
    """DooTask 基础异常类"""
    pass

class DooTaskAPIException(DooTaskException):
    """DooTask API 异常"""
    def __init__(self, message: str, ret_code: int = 0):
        super().__init__(message)
        self.ret_code = ret_code

class DooTaskHTTPException(DooTaskException):
    """DooTask HTTP 异常"""
    def __init__(self, message: str, status_code: int = 0):
        super().__init__(message)
        self.status_code = status_code

class DooTaskAuthException(DooTaskException):
    """DooTask 认证异常"""
    pass

class DooTaskPermissionException(DooTaskException):
    """DooTask 权限异常"""
    pass 