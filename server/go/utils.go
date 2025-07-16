package dootask

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"
)

// NewClient 创建DooTask客户端实例
func NewClient(token string, args ...string) *Client {
	server := "http://nginx"
	if len(args) > 0 {
		server = args[0]
	}
	return &Client{
		token:     token,
		server:    server,
		cache:     make(map[string]UserCache),
		cacheTime: 10 * time.Minute,
		timeout:   10 * time.Second,
	}
}

// SetCacheTime 设置缓存时间
func (c *Client) SetCacheTime(duration time.Duration) {
	c.cacheTime = duration
}

// SetTimeout 设置请求超时时间
func (c *Client) SetTimeout(duration time.Duration) {
	c.timeout = duration
}

// ClearCache 清空缓存
func (c *Client) ClearCache() {
	c.cache = make(map[string]UserCache)
}

// GetCacheSize 获取缓存大小
func (c *Client) GetCacheSize() int {
	return len(c.cache)
}

// NewGetRequest 创建GET请求
func (c *Client) NewGetRequest(api string, response any) (*http.Response, error) {
	client := &http.Client{
		Timeout: c.timeout,
	}
	req, err := http.NewRequest("GET", c.server+api, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Token", c.token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return resp, nil
}

// NewPostRequest 创建POST请求
func (c *Client) NewPostRequest(api string, body interface{}, response any) (*http.Response, error) {
	client := &http.Client{
		Timeout: c.timeout,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.server+api, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", c.token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return resp, nil
}

// ------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------

// GetUserInfo 获取用户信息
func (c *Client) GetUserInfo() (*UserInfo, error) {
	// 检查缓存
	if cache, ok := c.cache[c.token]; ok {
		if time.Now().Before(cache.ExpiresAt) {
			return &cache.User, nil
		}
		delete(c.cache, c.token)
	}

	// 验证 token
	var response UserResponse
	_, err := c.NewGetRequest("/api/users/info", &response)
	if err != nil {
		return nil, err
	}

	if response.Ret != 1 {
		return nil, errors.New(response.Msg)
	}

	// 更新缓存
	c.cache[c.token] = UserCache{
		User:      response.Data,
		ExpiresAt: time.Now().Add(c.cacheTime),
	}

	// 返回用户信息
	return &response.Data, nil
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
	var response UserInfoDepartmentsResponse
	_, err := c.NewGetRequest("/api/users/info/departments", &response)
	if err != nil {
		return nil, err
	}

	if response.Ret != 1 {
		return nil, errors.New(response.Msg)
	}

	return response.Data, nil
}

// GetUsersBasic 获取指定用户基础信息（支持多个用户）
func (c *Client) GetUsersBasic(userids []int) ([]UserBasic, error) {
	q := ""
	for i, id := range userids {
		if i > 0 {
			q += "&"
		}
		q += fmt.Sprintf("userid[]=%d", id)
	}

	var response UserBasicResponse
	_, err := c.NewGetRequest("/api/users/basic?"+q, &response)
	if err != nil {
		return nil, err
	}

	if response.Ret != 1 {
		return nil, errors.New(response.Msg)
	}

	return response.Data, nil
}

// GetUserBasic 获取单个用户基础信息（便利方法）
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

// SendMessage 发送消息
func (c *Client) SendMessage(message SendMessageRequest) (*Response, error) {
	if message.TextType == "" {
		message.TextType = "md"
	}

	var response Response
	_, err := c.NewPostRequest("/api/dialog/msg/sendtext", message, &response)
	if err != nil {
		return nil, err
	}

	if response.Ret != 1 {
		return nil, errors.New(response.Msg)
	}

	return &response, nil
}

// SendMessageToUser 发送消息到用户
func (c *Client) SendMessageToUser(message SendMessageToUserRequest) (*Response, error) {
	// 获取用户对话ID
	var response DialogOpenUserResponse
	_, err := c.NewGetRequest("/api/dialog/open/user?userid="+strconv.Itoa(message.UserID), &response)
	if err != nil {
		return nil, err
	} else if response.Ret != 1 {
		return nil, errors.New(response.Msg)
	}

	// 发送消息
	return c.SendMessage(SendMessageRequest{
		DialogID: response.Data.DialogUser.DialogID,
		Text:     message.Text,
		TextType: message.TextType,
		Silence:  message.Silence,
	})
}

// SendBotMessage 发送机器人消息
func (c *Client) SendBotMessage(message SendBotMessageRequest) (*Response, error) {
	var response Response
	_, err := c.NewPostRequest("/api/dialog/msg/sendbot", message, &response)
	if err != nil {
		return nil, err
	}

	if response.Ret != 1 {
		return nil, errors.New(response.Msg)
	}

	return &response, nil
}

// SendAnonymousMessage 发送匿名消息
func (c *Client) SendAnonymousMessage(message SendAnonymousMessageRequest) (*Response, error) {
	var response Response
	_, err := c.NewPostRequest("/api/dialog/msg/sendanon", message, &response)
	if err != nil {
		return nil, err
	}

	if response.Ret != 1 {
		return nil, errors.New(response.Msg)
	}

	return &response, nil
}
