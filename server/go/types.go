package dootask

import "time"

// ------------------------------------------------------------------------------------------
// 基础结构定义
// ------------------------------------------------------------------------------------------

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
	PerPage     any     `json:"per_page"` // 可能是 int 或 string
	PrevPageUrl *string `json:"prev_page_url"`
	To          any     `json:"to"` // 可能是 int 或 string
	Total       int     `json:"total"`
}

// UserCache 用户缓存
type UserCache struct {
	User      UserInfo
	ExpiresAt time.Time
}

// ------------------------------------------------------------------------------------------
// 用户相关结构体
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

// ------------------------------------------------------------------------------------------
// 消息相关结构体
// ------------------------------------------------------------------------------------------

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

// DialogMessage 对话消息
type DialogMessage struct {
	ID         int           `json:"id"`          // 消息ID
	DialogID   int           `json:"dialog_id"`   // 对话ID
	UserID     int           `json:"userid"`      // 用户ID
	Bot        int           `json:"bot"`         // 是否机器人
	CreatedAt  string        `json:"created_at"`  // 创建时间
	Type       string        `json:"type"`        // 消息类型
	MType      string        `json:"mtype"`       // 消息媒体类型
	Msg        any           `json:"msg"`         // 消息内容
	ReplyID    int           `json:"reply_id"`    // 回复消息ID
	ReplyNum   int           `json:"reply_num"`   // 回复数量
	ForwardID  int           `json:"forward_id"`  // 转发消息ID
	ForwardNum int           `json:"forward_num"` // 转发数量
	Tag        int           `json:"tag"`         // 标签
	Todo       int           `json:"todo"`        // 待办
	Read       int           `json:"read"`        // 已读人数
	Send       int           `json:"send"`        // 发送人数
	ReadAt     *string       `json:"read_at"`     // 已读时间
	Mention    int           `json:"mention"`     // 提及
	Dot        int           `json:"dot"`         // 点标记
	Emoji      []interface{} `json:"emoji"`       // 表情回应
	Link       int           `json:"link"`        // 链接
	Modify     int           `json:"modify"`      // 修改标记
	Percentage int           `json:"percentage"`  // 百分比
}

// DialogMessageListResponse 消息列表响应
type DialogMessageListResponse struct {
	List   []DialogMessage `json:"list"`   // 消息列表
	Time   int64           `json:"time"`   // 时间戳
	Dialog DialogInfo      `json:"dialog"` // 对话信息
	Todo   []interface{}   `json:"todo"`   // 待办列表
	Top    *interface{}    `json:"top"`    // 置顶消息ID
}

// DialogMessageSearchResponse 搜索消息响应
type DialogMessageSearchResponse struct {
	Data []int `json:"data"`
}

// TodoUser 待办用户
type TodoUser struct {
	UserID   int    `json:"userid"`   // 用户ID
	Nickname string `json:"nickname"` // 昵称
	UserImg  string `json:"userimg"`  // 头像
	Done     bool   `json:"done"`     // 是否完成
	DoneAt   string `json:"done_at"`  // 完成时间
}

// TodoListResponse 待办列表响应
type TodoListResponse struct {
	Users []TodoUser `json:"users"`
}

// GetMessageListRequest 获取消息列表请求
type GetMessageListRequest struct {
	DialogID   int    `json:"dialog_id"`   // 必填：对话ID
	MsgID      int    `json:"msg_id"`      // 可选：消息ID
	PositionID int    `json:"position_id"` // 可选：位置ID
	PrevID     int    `json:"prev_id"`     // 可选：前一个消息ID
	NextID     int    `json:"next_id"`     // 可选：下一个消息ID
	MsgType    string `json:"msg_type"`    // 可选：消息类型(tag/todo/link/text/image/file/record/meeting)
	Take       int    `json:"take"`        // 可选：获取数量，默认50，最大100
}

// SearchMessageRequest 搜索消息请求
type SearchMessageRequest struct {
	DialogID int    `json:"dialog_id"` // 必填：对话ID
	Key      string `json:"key"`       // 必填：搜索关键词
}

// GetMessageRequest 获取单个消息请求
type GetMessageRequest struct {
	MsgID int `json:"msg_id"` // 必填：消息ID
}

// WithdrawMessageRequest 撤回消息请求
type WithdrawMessageRequest struct {
	MsgID int `json:"msg_id"` // 必填：消息ID
}

// ForwardMessageRequest 转发消息请求
type ForwardMessageRequest struct {
	MsgID        int    `json:"msg_id"`        // 必填：消息ID
	DialogIDs    []int  `json:"dialogids"`     // 可选：目标对话ID列表
	UserIDs      []int  `json:"userids"`       // 可选：目标用户ID列表
	ShowSource   int    `json:"show_source"`   // 可选：是否显示来源，1显示，0不显示
	LeaveMessage string `json:"leave_message"` // 可选：留言
}

// ToggleMessageTodoRequest 切换消息待办请求
type ToggleMessageTodoRequest struct {
	MsgID   int    `json:"msg_id"`  // 必填：消息ID
	Type    string `json:"type"`    // 可选：类型(all/指定用户)，默认all
	UserIDs []int  `json:"userids"` // 可选：用户ID列表，当type不为all时使用
}

// MarkMessageDoneRequest 标记消息完成请求
type MarkMessageDoneRequest struct {
	MsgID int `json:"msg_id"` // 必填：消息ID
}

// ------------------------------------------------------------------------------------------
// 对话相关结构体
// ------------------------------------------------------------------------------------------

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

// ------------------------------------------------------------------------------------------
// 群组相关结构体
// ------------------------------------------------------------------------------------------

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

// ------------------------------------------------------------------------------------------
// 项目管理相关结构体
// ------------------------------------------------------------------------------------------

// Project 项目信息
type Project struct {
	ID          int    `json:"id"`           // 项目ID
	Name        string `json:"name"`         // 项目名称
	Desc        string `json:"desc"`         // 项目描述
	UserID      int    `json:"userid"`       // 创建者ID
	DialogID    int    `json:"dialog_id"`    // 对话ID
	ArchivedAt  string `json:"archived_at"`  // 归档时间
	CreatedAt   string `json:"created_at"`   // 创建时间
	UpdatedAt   string `json:"updated_at"`   // 更新时间
	Owner       int    `json:"owner"`        // 是否项目负责人
	OwnerUserID int    `json:"owner_userid"` // 项目负责人ID
	Personal    int    `json:"personal"`     // 是否个人项目
	// 任务统计
	TaskNum        int `json:"task_num"`         // 任务总数
	TaskComplete   int `json:"task_complete"`    // 已完成任务数
	TaskPercent    int `json:"task_percent"`     // 任务完成百分比
	TaskMyNum      int `json:"task_my_num"`      // 我的任务数
	TaskMyComplete int `json:"task_my_complete"` // 我的完成任务数
	TaskMyPercent  int `json:"task_my_percent"`  // 我的任务完成百分比
}

// GetProjectListRequest 获取项目列表请求
type GetProjectListRequest struct {
	Type          string `json:"type"`          // 可选：项目类型，all、team、personal
	Archived      string `json:"archived"`      // 可选：归档状态，all、yes、no
	GetColumn     string `json:"getcolumn"`     // 可选：同时取列表，yes、no
	GetUserID     string `json:"getuserid"`     // 可选：同时取成员ID，yes、no
	GetStatistics string `json:"getstatistics"` // 可选：同时取任务统计，yes、no
	TimeRange     string `json:"timerange"`     // 可选：时间范围
	Page          int    `json:"page"`          // 可选：当前页，默认1
	PageSize      int    `json:"pagesize"`      // 可选：每页数量，默认50
}

// GetProjectRequest 获取项目信息请求
type GetProjectRequest struct {
	ProjectID int `json:"project_id"` // 必填：项目ID
}

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Name     string `json:"name"`     // 必填：项目名称
	Desc     string `json:"desc"`     // 可选：项目描述
	Columns  string `json:"columns"`  // 可选：列表，格式：列表名称1,列表名称2
	Flow     string `json:"flow"`     // 可选：开启流程，open、close
	Personal int    `json:"personal"` // 可选：是否个人项目
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	ProjectID     int    `json:"project_id"`     // 必填：项目ID
	Name          string `json:"name"`           // 必填：项目名称
	Desc          string `json:"desc"`           // 可选：项目描述
	ArchiveMethod string `json:"archive_method"` // 可选：归档方式
	ArchiveDays   int    `json:"archive_days"`   // 可选：自动归档天数
}

// ProjectActionRequest 项目操作请求
type ProjectActionRequest struct {
	ProjectID int    `json:"project_id"` // 必填：项目ID
	Type      string `json:"type"`       // 可选：操作类型，如 add、recovery
}

// ------------------------------------------------------------------------------------------
// 任务列表相关结构体
// ------------------------------------------------------------------------------------------

// ProjectColumn 项目列表
type ProjectColumn struct {
	ID        int    `json:"id"`         // 列表ID
	ProjectID int    `json:"project_id"` // 项目ID
	Name      string `json:"name"`       // 列表名称
	Color     string `json:"color"`      // 颜色
	Sort      int    `json:"sort"`       // 排序
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}

// GetColumnListRequest 获取列表请求
type GetColumnListRequest struct {
	ProjectID int `json:"project_id"` // 必填：项目ID
	Page      int `json:"page"`       // 可选：当前页，默认1
	PageSize  int `json:"pagesize"`   // 可选：每页数量，默认100
}

// CreateColumnRequest 创建列表请求
type CreateColumnRequest struct {
	ProjectID int    `json:"project_id"` // 必填：项目ID
	Name      string `json:"name"`       // 必填：列表名称
}

// UpdateColumnRequest 更新列表请求
type UpdateColumnRequest struct {
	ColumnID int    `json:"column_id"` // 必填：列表ID
	Name     string `json:"name"`      // 可选：列表名称
	Color    string `json:"color"`     // 可选：颜色
}

// ColumnActionRequest 列表操作请求
type ColumnActionRequest struct {
	ColumnID int `json:"column_id"` // 必填：列表ID
}

// ------------------------------------------------------------------------------------------
// 任务相关结构体
// ------------------------------------------------------------------------------------------

// ProjectTask 项目任务
type ProjectTask struct {
	ID           int    `json:"id"`             // 任务ID
	ProjectID    int    `json:"project_id"`     // 项目ID
	ColumnID     int    `json:"column_id"`      // 列表ID
	ParentID     int    `json:"parent_id"`      // 父任务ID
	Name         string `json:"name"`           // 任务名称
	Desc         string `json:"desc"`           // 任务描述
	StartAt      string `json:"start_at"`       // 开始时间
	EndAt        string `json:"end_at"`         // 结束时间
	CompleteAt   string `json:"complete_at"`    // 完成时间
	ArchivedAt   string `json:"archived_at"`    // 归档时间
	CreatedAt    string `json:"created_at"`     // 创建时间
	UpdatedAt    string `json:"updated_at"`     // 更新时间
	UserID       int    `json:"userid"`         // 创建者ID
	DialogID     int    `json:"dialog_id"`      // 对话ID
	FlowItemID   int    `json:"flow_item_id"`   // 流程状态ID
	FlowItemName string `json:"flow_item_name"` // 流程状态名称
	Visibility   int    `json:"visibility"`     // 可见性
	Color        string `json:"color"`          // 颜色
	// 统计信息
	FileNum     int `json:"file_num"`     // 文件数量
	MsgNum      int `json:"msg_num"`      // 消息数量
	SubNum      int `json:"sub_num"`      // 子任务数量
	SubComplete int `json:"sub_complete"` // 子任务完成数量
	Percent     int `json:"percent"`      // 完成百分比
	// 关联数据
	ProjectName string `json:"project_name"` // 项目名称
	ColumnName  string `json:"column_name"`  // 列表名称
}

// TaskFile 任务文件
type TaskFile struct {
	ID        int    `json:"id"`         // 文件ID
	TaskID    int    `json:"task_id"`    // 任务ID
	Name      string `json:"name"`       // 文件名
	Ext       string `json:"ext"`        // 文件扩展名
	Size      int    `json:"size"`       // 文件大小
	Path      string `json:"path"`       // 文件路径
	Thumb     string `json:"thumb"`      // 缩略图
	UserID    int    `json:"userid"`     // 上传者ID
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}

// TaskContent 任务内容
type TaskContent struct {
	Content string `json:"content"` // 任务内容
	Type    string `json:"type"`    // 内容类型
}

// GetTaskListRequest 获取任务列表请求
type GetTaskListRequest struct {
	ProjectID int    `json:"project_id"` // 可选：项目ID
	ParentID  int    `json:"parent_id"`  // 可选：主任务ID
	Archived  string `json:"archived"`   // 可选：归档状态，all、yes、no
	Deleted   string `json:"deleted"`    // 可选：删除状态，all、yes、no
	TimeRange string `json:"timerange"`  // 可选：时间范围
	Page      int    `json:"page"`       // 可选：当前页，默认1
	PageSize  int    `json:"pagesize"`   // 可选：每页数量，默认100
}

// GetTaskRequest 获取任务信息请求
type GetTaskRequest struct {
	TaskID   int    `json:"task_id"`  // 必填：任务ID
	Archived string `json:"archived"` // 可选：归档状态，all、yes、no
}

// GetTaskContentRequest 获取任务内容请求
type GetTaskContentRequest struct {
	TaskID    int `json:"task_id"`    // 必填：任务ID
	HistoryID int `json:"history_id"` // 可选：历史ID
}

// GetTaskFilesRequest 获取任务文件请求
type GetTaskFilesRequest struct {
	TaskID int `json:"task_id"` // 必填：任务ID
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	ProjectID int      `json:"project_id"` // 必填：项目ID
	ColumnID  any      `json:"column_id"`  // 可选：列表ID
	Name      string   `json:"name"`       // 必填：任务名称
	Content   string   `json:"content"`    // 可选：任务内容
	Times     []string `json:"times"`      // 可选：计划时间
	Owner     []int    `json:"owner"`      // 可选：负责人
	Top       int      `json:"top"`        // 可选：置顶
}

// CreateSubTaskRequest 创建子任务请求
type CreateSubTaskRequest struct {
	TaskID int    `json:"task_id"` // 必填：任务ID
	Name   string `json:"name"`    // 必填：任务名称
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	TaskID     int      `json:"task_id"`     // 必填：任务ID
	Name       string   `json:"name"`        // 可选：任务名称
	Content    string   `json:"content"`     // 可选：任务内容
	Times      []string `json:"times"`       // 可选：计划时间
	Owner      []int    `json:"owner"`       // 可选：负责人
	Assist     []int    `json:"assist"`      // 可选：协助人
	Color      string   `json:"color"`       // 可选：颜色
	Visibility int      `json:"visibility"`  // 可选：可见性
	CompleteAt any      `json:"complete_at"` // 可选：完成时间
}

// TaskActionRequest 任务操作请求
type TaskActionRequest struct {
	TaskID int    `json:"task_id"` // 必填：任务ID
	Type   string `json:"type"`    // 可选：操作类型，如 add、recovery、delete
}

// CreateTaskDialogRequest 创建任务对话请求
type CreateTaskDialogRequest struct {
	TaskID int `json:"task_id"` // 必填：任务ID
}

// CreateTaskDialogResponse 创建任务对话响应
type CreateTaskDialogResponse struct {
	ID         int `json:"id"`          // 任务ID
	DialogID   int `json:"dialog_id"`   // 对话ID
	DialogData any `json:"dialog_data"` // 对话数据
}

// ------------------------------------------------------------------------------------------
// 系统相关结构体
// ------------------------------------------------------------------------------------------

// SystemSettings 系统设置
type SystemSettings struct {
	// 用户注册设置
	Reg *string `json:"reg"` // 用户注册开关：open/close

	// 任务设置
	TaskDefaultTime *[]string `json:"task_default_time"` // 任务默认时间

	// 系统信息
	SystemAlias   *string `json:"system_alias"`   // 系统别名
	SystemWelcome string  `json:"system_welcome"` // 系统欢迎语

	// 服务器信息
	ServerTimezone *string `json:"server_timezone"` // 服务器时区
	ServerVersion  *string `json:"server_version"`  // 服务器版本
}

// VersionInfo 版本信息
type VersionInfo struct {
	DeviceCount int    `json:"device_count"` // 设备数量
	Version     string `json:"version"`      // 版本号
}
