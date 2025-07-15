# DooTask Go SDK

DooTask API的Go语言客户端SDK。

## 功能特性

- 用户身份验证和信息获取
- 用户权限检查
- 部门信息获取
- 获取用户基础信息（支持批量获取）
- 发送消息到对话或用户
- 发送匿名消息
- 发送机器人消息

## 基本用法

```go
package main

import (
    "fmt"
    "github.com/dootask/tools/go/utils"
)

func main() {
    // 创建客户端
    client := utils.NewClient("your-token", "https://your-server.com")  // your-server 生产环境可以不填
    
    // 获取用户信息
    user, err := client.GetUserInfo()
    if err != nil {
        panic(err)
    }
    fmt.Printf("用户: %s\n", user.Nickname)
    
    // 发送消息到用户
    response, err := client.SendMessageToUser(123, "Hello World!")
    if err != nil {
        panic(err)
    }
    fmt.Printf("消息发送成功: %+v\n", response)
}
```

## API 接口

### 用户相关

#### GetUserInfo() 获取用户信息
```go
user, err := client.GetUserInfo()
```

#### CheckUserIdentity(identity string) 检查用户权限
```go
user, err := client.CheckUserIdentity("admin")
```

#### GetUserDepartments() 获取用户部门
```go
departments, err := client.GetUserDepartments()
```

#### GetUserBasic(userid int) 获取单个用户基础信息
```go
userBasic, err := client.GetUserBasic(123)
```

#### GetUsersBasic(userids []int) 获取多个用户基础信息
```go
// 获取多个用户基础信息（最多50个）
usersBasic, err := client.GetUsersBasic([]int{123, 456, 789})
```

### 消息相关

#### SendMessage(dialogId int, text string, args ...interface{}) 发送消息到指定对话
```go
// 基本用法（默认md格式）
response, err := client.SendMessage(100, "这是一条**markdown**消息")

// 指定消息类型
response, err := client.SendMessage(100, "纯文本消息", "text")

// 静默发送
response, err := client.SendMessage(100, "静默消息", "md", true)
```

支持的消息类型：
- `md`: Markdown格式（默认）
- `text`: 纯文本
- `template`: 模板消息

#### SendMessageToUser(userid int, text string, args ...interface{}) 发送消息到用户
```go
// 基本用法（自动获取对话ID并发送）
response, err := client.SendMessageToUser(123, "这是一条**markdown**消息")

// 指定消息类型
response, err := client.SendMessageToUser(123, "纯文本消息", "text")

// 静默发送
response, err := client.SendMessageToUser(123, "静默消息", "md", true)
```

#### SendAnonymousMessage(userid int, text string) 发送匿名消息
```go
response, err := client.SendAnonymousMessage(123, "这是一条匿名消息")
```

#### SendBotMessage(userid int, text string, args ...interface{}) 发送机器人消息
```go
// 基本用法（使用默认system-msg机器人）
response, err := client.SendBotMessage(123, "消息内容")

// 指定机器人类型
response, err := client.SendBotMessage(123, "任务提醒", "task-alert")

// 静默发送
response, err := client.SendBotMessage(123, "静默消息", "system-msg", true)
```

支持的机器人类型：
- `system-msg`: 系统消息（默认）
- `task-alert`: 任务提醒
- `check-in`: 签到打卡
- `approval-alert`: 审批
- `meeting-alert`: 会议通知

### 配置方法

#### SetCacheTime(duration time.Duration) 设置缓存时间
```go
client.SetCacheTime(5 * time.Minute)
```

#### SetTimeout(duration time.Duration) 设置请求超时
```go
client.SetTimeout(30 * time.Second)
```

#### ClearCache() 清空缓存
```go
client.ClearCache()
```

## 消息发送方式对比

| 方法 | 适用场景 | 发送方式 | 特点 |
|------|----------|----------|------|
| `SendMessage` | 已知对话ID | 直接发送到对话 | 最高效，需要预先获取dialog_id |
| `SendMessageToUser` | 知道用户ID | 自动获取对话ID后发送 | 便捷，自动处理对话获取 |
| `SendAnonymousMessage` | 匿名消息 | 通过匿名机器人发送 | 保护发送者身份 |
| `SendBotMessage` | 系统通知 | 通过指定机器人发送 | 支持多种机器人类型 |

## 完整示例

```go
package main

import (
    "fmt"
    "log"
    "github.com/dootask/tools/go/utils"
)

func main() {
    // 创建客户端
    client := utils.NewClient("your-token", "https://your-server.com")  // your-server 生产环境可以不填
    
    // 获取用户信息
    user, err := client.GetUserInfo()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("当前用户: %s\n", user.Nickname)
    
    // 获取用户基础信息
    userBasic, err := client.GetUserBasic(123)
    if err != nil {
        log.Printf("获取用户基础信息失败: %v", err)
    } else {
        fmt.Printf("目标用户: %s\n", userBasic.Nickname)
    }
    
    // 发送消息到用户（推荐方式）
    response, err := client.SendMessageToUser(123, "## 消息标题\n这是一条**重要**消息！")
    if err != nil {
        log.Printf("发送消息失败: %v", err)
    } else {
        fmt.Printf("消息发送成功: %+v\n", response)
    }
    
    // 发送任务提醒
    _, err = client.SendBotMessage(123, "您有新的任务需要处理", "task-alert")
    if err != nil {
        log.Printf("发送任务提醒失败: %v", err)
    } else {
        fmt.Println("任务提醒发送成功")
    }
}
```
