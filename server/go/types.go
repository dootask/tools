package dootask

import "time"

// Client DooTask客户端类
type Client struct {
	token     string
	server    string
	cache     map[string]UserCache
	cacheTime time.Duration
	timeout   time.Duration
}

// ClientOption 客户端选项
type ClientOption func(*Client)

// Response 基础响应结构
type Response[T any] struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// ResponsePaginate 分页数据
type ResponsePaginate[T any] struct {
	CurrentPage int     `json:"current_page"`
	Data        []T     `json:"data"`
	NextPageUrl *string `json:"next_page_url"`
	Path        string  `json:"path"`
	PerPage     int     `json:"per_page"`
	PrevPageUrl *string `json:"prev_page_url"`
	To          int     `json:"to"`
	Total       int     `json:"total"`
}

// UserCache 用户缓存
type UserCache struct {
	User      UserInfo
	ExpiresAt time.Time
}

// ------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------

// UserInfo 用户信息
type UserInfo struct {
	UserID         uint     `json:"userid"`          // 用户ID
	Identity       []string `json:"identity"`        // 身份
	Email          string   `json:"email"`           // 邮箱
	Nickname       string   `json:"nickname"`        // 昵称
	Profession     string   `json:"profession"`      // 职业
	UserImg        string   `json:"userimg"`         // 头像
	Bot            int      `json:"bot"`             // 是否机器人
	Department     []int    `json:"department"`      // 部门
	DepartmentName string   `json:"department_name"` // 部门名称
}

// UserBasic 用户基础信息
type UserBasic struct {
	UserID         uint   `json:"userid"`          // 用户ID
	Email          string `json:"email"`           // 邮箱
	Nickname       string `json:"nickname"`        // 昵称
	Profession     string `json:"profession"`      // 职业
	UserImg        string `json:"userimg"`         // 头像
	Bot            int    `json:"bot"`             // 是否机器人
	Online         bool   `json:"online"`          // 是否在线
	Department     []int  `json:"department"`      // 部门
	DepartmentName string `json:"department_name"` // 部门名称
}

// Department 部门信息
type Department struct {
	ID          int    `json:"id"`           // 部门ID
	Name        string `json:"name"`         // 部门名称
	ParentID    int    `json:"parent_id"`    // 父部门ID
	OwnerUserID uint   `json:"owner_userid"` // 负责人ID
}

// DialogOpenUserResponse 打开用户对话响应数据
type DialogOpenUserResponse struct {
	DialogUser DialogUserResponse `json:"dialog_user"`
}

// DialogUserResponse 对话用户信息
type DialogUserResponse struct {
	DialogID int `json:"dialog_id"`
	UserID   int `json:"userid"`
	Bot      int `json:"bot"`
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

// DialogInfo 对话信息
type DialogInfo struct {
	ID         int    `json:"id"`          // 会话ID
	Type       string `json:"type"`        // 会话类型
	GroupType  string `json:"group_type"`  // 群组类型
	Name       string `json:"name"`        // 会话名称
	Avatar     string `json:"avatar"`      // 会话头像
	OwnerID    int    `json:"owner_id"`    // 群主ID
	CreatedAt  string `json:"created_at"`  // 创建时间
	UpdatedAt  string `json:"updated_at"`  // 更新时间
	LastAt     string `json:"last_at"`     // 最后活跃时间
	MarkUnread int    `json:"mark_unread"` // 标记未读
	Silence    int    `json:"silence"`     // 静默
	Hide       int    `json:"hide"`        // 是否隐藏
	Color      string `json:"color"`       // 颜色
	Unread     int    `json:"unread"`      // 未读数
	UnreadOne  int    `json:"unread_one"`  // 单聊未读数
	Mention    int    `json:"mention"`     // @消息数
	MentionIDs []int  `json:"mention_ids"` // @用户ID列表
	People     int    `json:"people"`      // 总人数
	PeopleUser int    `json:"people_user"` // 用户人数
	PeopleBot  int    `json:"people_bot"`  // 机器人数量
	TodoNum    int    `json:"todo_num"`    // 待办数量
	LastMsg    any    `json:"last_msg"`    // 最后一条消息（结构体或map，具体类型视接口返回）
	Pinyin     string `json:"pinyin"`      // 拼音
	Bot        int    `json:"bot"`         // 机器人所有者
	TopAt      string `json:"top_at"`      // 置顶时间
}

// DialogMember 会话成员信息
type DialogMember struct {
	ID       int    `json:"id"`        // 会话成员ID
	DialogID int    `json:"dialog_id"` // 会话ID
	UserID   int    `json:"userid"`    // 用户ID
	Nickname string `json:"nickname"`  // 昵称
	Email    string `json:"email"`     // 邮箱
	UserImg  string `json:"userimg"`   // 头像
	Bot      int    `json:"bot"`       // 是否机器人
	Online   bool   `json:"online"`    // 是否在线
}

// CreateGroupRequest 创建群组请求
type CreateGroupRequest struct {
	Avatar   string `json:"avatar"`    // 可选：群头像
	ChatName string `json:"chat_name"` // 可选：群名称
	UserIDs  []int  `json:"userids"`   // 必填：群成员ID列表
}

// EditGroupRequest 修改群组请求
type EditGroupRequest struct {
	DialogID int    `json:"dialog_id"` // 必填：会话ID
	Avatar   string `json:"avatar"`    // 可选：群头像
	ChatName string `json:"chat_name"` // 可选：群名称
	Admin    int    `json:"admin"`     // 可选：系统管理员操作（1：只判断是不是系统管理员，否则判断是否群管理员）
}

// AddGroupUserRequest 添加群成员请求
type AddGroupUserRequest struct {
	DialogID int   `json:"dialog_id"` // 必填：会话ID
	UserIDs  []int `json:"userids"`   // 必填：新增的群成员ID列表
}

// RemoveGroupUserRequest 移除群成员请求
type RemoveGroupUserRequest struct {
	DialogID int   `json:"dialog_id"` // 必填：会话ID
	UserIDs  []int `json:"userids"`   // 可选：移出的群成员ID列表，留空表示自己退出
}

// TransferGroupRequest 转让群组请求
type TransferGroupRequest struct {
	DialogID   int    `json:"dialog_id"`   // 必填：会话ID
	UserID     int    `json:"userid"`      // 必填：新的群主ID
	CheckOwner string `json:"check_owner"` // 可选：转让验证，yes-需要验证，no-不需要验证
	Key        string `json:"key"`         // 可选：密钥（APP_KEY）
}

// DisbandGroupRequest 解散群组请求
type DisbandGroupRequest struct {
	DialogID int `json:"dialog_id"` // 必填：会话ID
}

// TimeRangeRequest 时间范围请求参数
type TimeRangeRequest struct {
	TimeRange string `json:"timerange"` // 可选：时间范围，例如：1752711205,1751776557，表示时间戳范围
	Page      int    `json:"page"`      // 可选：当前页，默认1
	PageSize  int    `json:"pagesize"`  // 可选：每页显示数量，默认50，最大100
}

// SearchDialogRequest 搜索会话请求
type SearchDialogRequest struct {
	Key string `json:"key"` // 必填：搜索关键词
}

// GetDialogRequest 获取单个会话请求
type GetDialogRequest struct {
	DialogID int `json:"dialog_id"` // 必填：对话ID
}

// GetDialogUserRequest 获取会话成员请求
type GetDialogUserRequest struct {
	DialogID int `json:"dialog_id"` // 必填：会话ID
	GetUser  int `json:"getuser"`   // 可选：获取会员详情（1: 返回会员昵称、邮箱等基本信息，0: 默认不返回）
}
