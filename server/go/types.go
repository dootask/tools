package dootask

import "time"

// Response 基础响应结构
type Response struct {
	Ret  int         `json:"ret"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// UserInfo 用户信息
type UserInfo struct {
	UserID     uint     `json:"userid"`
	Identity   []string `json:"identity"`
	Email      string   `json:"email"`
	Nickname   string   `json:"nickname"`
	UserImg    string   `json:"userimg"`
	Profession string   `json:"profession"`
}

// UserResponse 用户响应
type UserResponse struct {
	Response
	Data UserInfo `json:"data"`
}

// UserCache 用户缓存
type UserCache struct {
	User      UserInfo
	ExpiresAt time.Time
}

// Department 部门信息
type Department struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ParentID    int    `json:"parent_id"`
	OwnerUserID uint   `json:"owner_userid"`
}

// UserBasic 用户基础信息
type UserBasic struct {
	UserID   uint   `json:"userid"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	UserImg  string `json:"userimg"`
	Bot      int    `json:"bot"`
}

// UserBasicResponse 用户基础信息响应
type UserBasicResponse struct {
	Response
	Data []UserBasic `json:"data"`
}

// UserInfoDepartmentsResponse 部门响应
type UserInfoDepartmentsResponse struct {
	Response
	Data []Department `json:"data"`
}

// DialogOpenUserResponseData 打开用户对话响应数据
type DialogOpenUserResponseData struct {
	DialogUser DialogUserInfo `json:"dialog_user"`
}

// DialogUserInfo 对话用户信息
type DialogUserInfo struct {
	DialogID int `json:"dialog_id"`
	UserID   int `json:"userid"`
	Bot      int `json:"bot"`
}

// DialogOpenUserResponse 打开用户对话响应
type DialogOpenUserResponse struct {
	Response
	Data DialogOpenUserResponseData `json:"data"`
}

// Client DooTask客户端类
type Client struct {
	token     string
	server    string
	cache     map[string]UserCache
	cacheTime time.Duration
	timeout   time.Duration
}

// SendMessageRequest 消息
type SendMessageRequest struct {
	DialogID int    `json:"dialog_id"` // 必填：对话ID
	Text     string `json:"text"`      // 必填：消息内容
	TextType string `json:"text_type"` // 可选：消息类型，可选值：md、text
	Silence  bool   `json:"silence"`   // 可选：是否静默，可选值：true、false，默认 false
}

// SendMessageToUserRequest 发送消息到用户
type SendMessageToUserRequest struct {
	UserID   int    `json:"userid"`    // 必填：用户ID，发送给指定用户
	Text     string `json:"text"`      // 必填：消息内容
	TextType string `json:"text_type"` // 可选：消息类型，可选值：md、text
	Silence  bool   `json:"silence"`   // 可选：是否静默，可选值：true、false，默认 false
}

// SendBotMessageRequest 发送机器人消息
type SendBotMessageRequest struct {
	UserID  int    `json:"userid"`   // 必填：用户ID，发送给指定用户
	Text    string `json:"text"`     // 必填：消息内容，支持 markdown 格式
	BotType string `json:"bot_type"` // 可选：机器人类型，可选值：system-msg、task-alert、check-in、approval-alert、meeting-alert、xxxxxx（自定义机器人，6-20个字符）
	BotName string `json:"bot_name"` // 可选：机器人名称，当 bot_type 为 xxxxxx 时有效
	Silence bool   `json:"silence"`  // 可选：是否静默，可选值：true、false，默认 false
}

// SendAnonymousMessageRequest 发送匿名消息
type SendAnonymousMessageRequest struct {
	UserID int    `json:"userid"` // 必填：用户ID，发送给指定用户
	Text   string `json:"text"`   // 必填：消息内容
}
