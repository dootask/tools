#!/usr/bin/env python3
"""
DooTask Tools 使用示例
"""

from dootask import (
    DooTaskClient,
    SendMessageToUserRequest,
    CreateProjectRequest,
    CreateTaskRequest,
    DooTaskException
)

def main():
    """主函数"""
    
    # 创建客户端
    client = DooTaskClient(
        token="YIG8ANC8q2ROQF91r8Pe6-53rIG3oCxcqQN-mMdZpQKe7mKwNqIHenDNqbDDdyQIdo9w2KdveEpF1NaH-5Nfmv0dBr9TkjJ7KFMkfEUL11wOjyId0nuoSJaAliRz8d5z",
        server="http://127.0.0.1:2222"
    )
    
    try:
        # 1. 获取用户信息
        user = client.get_user_info()
        print(f"用户: {user.nickname} ({user.email})")
        
        # 2. 发送消息
        message = SendMessageToUserRequest(
            userid=user.userid,
            text="Hello from Python! 🐍"
        )
        client.send_message_to_user(message)
        print("消息发送成功！")
        
        # 3. 创建项目
        project = client.create_project(CreateProjectRequest(
            name="Python 测试项目",
            desc="这是一个测试项目"
        ))
        print(f"项目创建成功: {project.name}")
        
        # 4. 创建任务
        task = client.create_task(CreateTaskRequest(
            project_id=project.id,
            name="测试任务",
            content="这是一个测试任务",
            owner=[user.userid]
        ))
        print(f"任务创建成功: {task.name}")
        
    except DooTaskException as e:
        print(f"错误: {e}")

if __name__ == "__main__":
    main() 