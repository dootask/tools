package dootask

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
)

// ------------------------------------------------------------------------------------------
// 客户端相关方法
// ------------------------------------------------------------------------------------------

// WithServer 设置服务器地址
func WithServer(server string) ClientOption {
	return func(c *Client) {
		c.server = server
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// NewClient 创建客户端实例
func NewClient(token string, opts ...ClientOption) *Client {
	client := &Client{
		token:     token,
		server:    "http://nginx",
		cache:     make(map[string]UserCache),
		cacheTime: 10 * time.Minute,
		timeout:   10 * time.Second,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// buildURL 构建带查询参数的URL
func buildURL(baseURL string, params map[string]any) string {
	if len(params) == 0 {
		return baseURL
	}

	var queryParams []string
	for key, value := range params {
		if value == nil {
			continue
		}

		switch v := value.(type) {
		case string:
			if v != "" {
				queryParams = append(queryParams, key+"="+url.QueryEscape(v))
			}
		case int:
			queryParams = append(queryParams, key+"="+strconv.Itoa(v))
		case int64:
			queryParams = append(queryParams, key+"="+strconv.FormatInt(v, 10))
		case float64:
			queryParams = append(queryParams, key+"="+strconv.FormatFloat(v, 'f', -1, 64))
		case []int:
			for _, id := range v {
				queryParams = append(queryParams, key+"[]="+strconv.Itoa(id))
			}
		case []string:
			for _, str := range v {
				if str != "" {
					queryParams = append(queryParams, key+"[]="+url.QueryEscape(str))
				}
			}
		case []any:
			for _, id := range v {
				switch id := id.(type) {
				case int:
					queryParams = append(queryParams, key+"[]="+strconv.Itoa(id))
				case string:
					queryParams = append(queryParams, key+"[]="+url.QueryEscape(id))
				default:
					queryParams = append(queryParams, key+"[]="+url.QueryEscape(fmt.Sprintf("%v", id)))
				}
			}
		case bool:
			if v {
				queryParams = append(queryParams, key+"=1")
			} else {
				queryParams = append(queryParams, key+"=0")
			}
		default:
			// 对于其他类型，转换为字符串
			queryParams = append(queryParams, key+"="+url.QueryEscape(fmt.Sprintf("%v", v)))
		}
	}

	if len(queryParams) == 0 {
		return baseURL
	}

	separator := "?"
	if strings.Contains(baseURL, "?") {
		separator = "&"
	}

	return baseURL + separator + strings.Join(queryParams, "&")
}

// structToMap 将结构体转换为 map[string]any
func structToMap(data any) (map[string]any, error) {
	if data == nil {
		return nil, nil
	}

	// 如果已经是 map，直接返回
	if m, ok := data.(map[string]any); ok {
		return m, nil
	}

	// 使用 JSON 序列化然后反序列化
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal to json failed: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("unmarshal from json failed: %w", err)
	}

	return result, nil
}

// NewRequest 创建请求
func (c *Client) NewRequest(method, api string, requestData any, responseData any) error {
	// 验证 responseData 必须是指针（如果不为 nil）
	if responseData != nil {
		rv := reflect.ValueOf(responseData)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return errors.New("responseData must be a non-nil pointer")
		}
	}

	var req *http.Request
	var err error
	fullURL := c.server + api

	switch strings.ToUpper(method) {
	case "GET":
		// GET 请求：将 requestData 作为查询参数
		if requestData != nil {
			// 如果 requestData 是结构体，则将 requestData 转换为 map[string]any
			if params, err := structToMap(requestData); err == nil && len(params) > 0 {
				fullURL = buildURL(fullURL, params)
			}
		}
		req, err = http.NewRequest("GET", fullURL, nil)

	case "POST", "PUT", "PATCH":
		// POST/PUT/PATCH 请求：将 requestData 作为 JSON body
		var body io.Reader
		if requestData != nil {
			jsonData, jsonErr := json.Marshal(requestData)
			if jsonErr != nil {
				return fmt.Errorf("marshal request data failed: %w", jsonErr)
			}
			body = bytes.NewBuffer(jsonData)
		}
		req, err = http.NewRequest(method, fullURL, body)
		if err == nil && requestData != nil {
			req.Header.Set("Content-Type", "application/json")
		}

	case "DELETE":
		// DELETE 请求：支持查询参数
		if requestData != nil {
			if params, err := structToMap(requestData); err == nil && len(params) > 0 {
				fullURL = buildURL(fullURL, params)
			}
		}
		req, err = http.NewRequest("DELETE", fullURL, nil)

	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	// 设置通用请求头
	req.Header.Set("Token", c.token)
	req.Header.Set("User-Agent", "DooTask-Go-Client/1.0")

	// 发送请求
	client := &http.Client{Timeout: c.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s, body: %s", resp.StatusCode, resp.Status, string(bodyBytes))
	}

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %w", err)
	}

	// 解析响应
	var apiResp Response[json.RawMessage]
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return fmt.Errorf("parse response failed: %w", err)
	}

	// 检查业务状态
	if apiResp.Ret != 1 {
		if apiResp.Msg != "" {
			return fmt.Errorf("%s", apiResp.Msg)
		}
		return fmt.Errorf("API error: %d", apiResp.Ret)
	}

	// 如果不需要响应数据，直接返回
	if responseData == nil {
		return nil
	}

	// 解析数据到目标结构
	if apiResp.Data != nil {
		if err := json.Unmarshal(apiResp.Data, responseData); err != nil {
			return fmt.Errorf("unmarshal data failed: %w", err)
		}
	}

	return nil
}

// NewGetRequest 创建GET请求
func (c *Client) NewGetRequest(api string, requestData any, responseData any) error {
	return c.NewRequest("GET", api, requestData, responseData)
}

// NewPostRequest 创建POST请求
func (c *Client) NewPostRequest(api string, requestData any, responseData any) error {
	return c.NewRequest("POST", api, requestData, responseData)
}

// ------------------------------------------------------------------------------------------
// 用户相关接口
// ------------------------------------------------------------------------------------------

// GetUserInfo 获取用户信息
func (c *Client) GetUserInfo(noCache ...bool) (*UserInfo, error) {
	// 检查缓存
	if cache, ok := c.cache[c.token]; ok {
		if time.Now().Before(cache.ExpiresAt) && !slices.Contains(noCache, true) {
			return &cache.User, nil
		}
		delete(c.cache, c.token)
	}

	// 验证 token
	var response UserInfo
	err := c.NewGetRequest("/api/users/info", nil, &response)
	if err != nil {
		return nil, err
	}

	// 更新缓存
	c.cache[c.token] = UserCache{
		User:      response,
		ExpiresAt: time.Now().Add(c.cacheTime),
	}

	// 返回用户信息
	return &response, nil
}

// CheckUserIdentity 检查用户是否具有指定身份
func (c *Client) CheckUserIdentity(identity string) (*UserInfo, error) {
	user, err := c.GetUserInfo()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(user.Identity, identity) {
		return nil, errors.New("insufficient permissions")
	}

	return user, nil
}

// GetUserDepartments 获取用户部门信息
func (c *Client) GetUserDepartments() ([]Department, error) {
	var response []Department
	err := c.NewGetRequest("/api/users/info/departments", nil, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetUsersBasic 获取指定用户基础信息（支持多个用户）
func (c *Client) GetUsersBasic(userids []int) ([]UserBasic, error) {
	var response []UserBasic
	err := c.NewGetRequest("/api/users/basic", map[string]any{
		"userid": userids,
	}, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetUserBasic 获取指定用户基础信息（单个用户）
func (c *Client) GetUserBasic(userid int) (*UserBasic, error) {
	users, err := c.GetUsersBasic([]int{userid})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("用户不存在")
	}

	return &users[0], nil
}

// ------------------------------------------------------------------------------------------
// 消息相关接口
// ------------------------------------------------------------------------------------------

// SendMessage 发送消息
func (c *Client) SendMessage(message SendMessageRequest) error {
	if message.TextType == "" {
		message.TextType = "md"
	}

	return c.NewPostRequest("/api/dialog/msg/sendtext", message, nil)
}

// SendMessageToUser 发送消息到用户
func (c *Client) SendMessageToUser(message SendMessageToUserRequest) error {
	// 获取用户对话ID
	queryParams := map[string]any{
		"userid": message.UserID,
	}

	var response DialogOpenUserResponse
	err := c.NewGetRequest("/api/dialog/open/user", queryParams, &response)
	if err != nil {
		return err
	}

	// 发送消息
	return c.SendMessage(SendMessageRequest{
		DialogID: response.DialogUser.DialogID,
		Text:     message.Text,
		TextType: message.TextType,
		Silence:  message.Silence,
	})
}

// SendBotMessage 发送机器人消息
func (c *Client) SendBotMessage(message SendBotMessageRequest) error {
	if message.BotType == "" {
		message.BotType = "system-msg"
	}

	return c.NewPostRequest("/api/dialog/msg/sendbot", message, nil)
}

// SendAnonymousMessage 发送匿名消息
func (c *Client) SendAnonymousMessage(message SendAnonymousMessageRequest) error {
	return c.NewPostRequest("/api/dialog/msg/sendanon", message, nil)
}

// ------------------------------------------------------------------------------------------
// 对话相关接口
// ------------------------------------------------------------------------------------------

// GetDialogList 获取对话列表
func (c *Client) GetDialogList(params TimeRangeRequest) (*ResponsePaginate[DialogInfo], error) {
	var response ResponsePaginate[DialogInfo]
	err := c.NewGetRequest("/api/dialog/lists", params, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SearchDialog 搜索会话
func (c *Client) SearchDialog(params SearchDialogRequest) ([]DialogInfo, error) {
	var response []DialogInfo
	err := c.NewGetRequest("/api/dialog/search", params, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetDialogOne 获取单个会话信息
func (c *Client) GetDialogOne(params GetDialogRequest) (*DialogInfo, error) {
	var response DialogInfo
	err := c.NewGetRequest("/api/dialog/one", params, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetDialogUser 获取会话成员
func (c *Client) GetDialogUser(params GetDialogUserRequest) ([]DialogMember, error) {
	var response []DialogMember
	err := c.NewGetRequest("/api/dialog/user", params, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ------------------------------------------------------------------------------------------
// 群组相关接口
// ------------------------------------------------------------------------------------------

// CreateGroup 新增群组
func (c *Client) CreateGroup(params CreateGroupRequest) (*DialogInfo, error) {
	var response DialogInfo
	err := c.NewGetRequest("/api/dialog/group/add", params, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// EditGroup 修改群组
func (c *Client) EditGroup(params EditGroupRequest) error {
	return c.NewGetRequest("/api/dialog/group/edit", params, nil)
}

// AddGroupUser 添加群成员
func (c *Client) AddGroupUser(params AddGroupUserRequest) error {
	return c.NewGetRequest("/api/dialog/group/adduser", params, nil)
}

// RemoveGroupUser 移除群成员
func (c *Client) RemoveGroupUser(params RemoveGroupUserRequest) error {
	return c.NewGetRequest("/api/dialog/group/deluser", params, nil)
}

// ExitGroup 退出群组
func (c *Client) ExitGroup(dialogID int) error {
	return c.RemoveGroupUser(RemoveGroupUserRequest{
		DialogID: dialogID,
		UserIDs:  []int{}, // 空数组表示自己退出
	})
}

// TransferGroup 转让群组
func (c *Client) TransferGroup(params TransferGroupRequest) error {
	return c.NewGetRequest("/api/dialog/group/transfer", params, nil)
}

// DisbandGroup 解散群组
func (c *Client) DisbandGroup(params DisbandGroupRequest) error {
	return c.NewGetRequest("/api/dialog/group/disband", params, nil)
}
