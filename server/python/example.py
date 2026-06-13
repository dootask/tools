#!/usr/bin/env python3
"""
DooTask Tools 使用示例
"""

from dootask import (
    DooTaskClient,
    SendMessageToUserRequest,
    CreateProjectRequest,
    CreateTaskRequest,
    SearchMessageRequest,
    GetMessageRequest,
    ToggleMessageTodoRequest,
    MarkMessageDoneRequest,
    DooTaskException
)

def main():
    """主函数"""
    
    # 创建客户端
    # version 作为 Version 头随请求发送，供后端 checkClientVersion 判定接口兼容性；
    # 部分接口（如消息待办 msg/todo）有最低版本要求，需传当前适配的主程序版本。
    client = DooTaskClient(
        token="YIG8ANC8q2ROQF91r8Pe6-53rIG3oCxcqQN-mMdZpQKe7mKwNqIHenDNqbDDdyQIdo9w2KdveEpF1NaH-5Nfmv0dBr9TkjJ7KFMkfEUL11wOjyId0nuoSJaAliRz8d5z",
        server="http://127.0.0.1:2222",
        version="1.7.91"
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

        # 5. 搜索消息（走 /api/search/message；不传 dialog_id 为全局搜索）
        results = client.search_message(SearchMessageRequest(key="测试", take=5))
        print(f"搜索到 {len(results)} 条消息")
        for item in results:
            print(f"  msg_id={item.msg_id} dialog={item.dialog_id} type={item.type}")

        # 6. 待办流：把消息设为待办 → 取待办记录 → 完成
        if results:
            msg_id = results[0].msg_id
            # 6.1 设为待办（type=all 表示对话内全体成员）
            client.toggle_message_todo(ToggleMessageTodoRequest(msg_id=msg_id))
            # 6.2 取该消息的待办记录（每条记录的 id 即“待办数据ID”）
            todos = client.get_message_todo_list(GetMessageRequest(msg_id=msg_id))
            print(f"消息 {msg_id} 有 {len(todos)} 条待办记录")
            # 6.3 完成属于自己的那条待办（注意：传待办数据ID，不是消息ID）
            mine = next((t for t in todos if t.userid == user.userid), None)
            if mine:
                client.mark_message_done(MarkMessageDoneRequest(id=mine.id))
                print(f"已完成待办 #{mine.id}")

    except DooTaskException as e:
        print(f"错误: {e}")

if __name__ == "__main__":
    main() 