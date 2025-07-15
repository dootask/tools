package utils

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
