package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AppStore 由独立 Go 微服务提供，主程序反代在 /appstore/api/v1/...，用 Token 头鉴权，
// 响应格式是 {code,message,data}（code==200 成功）——与主程序 {ret,msg,data} 不同，
// 故不能复用 SDK 的 NewGetRequest/NewPostRequest，这里单独走原始 HTTP。
const appStoreTimeout = 300 * time.Second // 安装/卸载会跑 docker compose，留足时间

// AppStore 用 Version 头判定主程序版本（缺省按 1.0.0，会让 require_version 校验失败）。
// 缓存一次主程序版本，避免每条 app 命令重复请求。
var cachedMainVersion string

func mainAppVersion() string {
	if cachedMainVersion != "" {
		return cachedMainVersion
	}
	c, err := Opts.Client()
	if err != nil {
		return ""
	}
	v, err := c.GetVersion()
	if err != nil || v == nil {
		return ""
	}
	cachedMainVersion = v.Version
	return cachedMainVersion
}

// AppStoreRequest 调用 AppStore 接口。path 以 / 开头、相对 /appstore/api/v1。
// body 非 nil 时以 JSON POST；out 非 nil 时把 data 反序列化进去。
func AppStoreRequest(method, path string, query map[string]string, body any, out any) error {
	if Opts.Token == "" {
		return ErrNoAuth
	}
	u := strings.TrimRight(Opts.Server, "/") + "/appstore/api/v1" + path
	if len(query) > 0 {
		vals := url.Values{}
		for k, v := range query {
			if v != "" {
				vals.Set(k, v)
			}
		}
		if enc := vals.Encode(); enc != "" {
			u += "?" + enc
		}
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("编码请求失败: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, u, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Token", Opts.Token)
	req.Header.Set("User-Agent", "doo-cli")
	if v := mainAppVersion(); v != "" {
		req.Header.Set("Version", v) // 供 AppStore 校验 require_version，缺省会被当 1.0.0
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := (&http.Client{Timeout: appStoreTimeout}).Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	return parseAppStoreResponse(resp, out)
}

// AppStoreUpload 以 multipart/form-data 上传一个文件到 AppStore（如上传本地应用压缩包）。
// fileField 为文件表单字段名，filePath 为本地文件路径；fields 为附带的普通文本字段。
// 鉴权、Version 头、响应解析与 AppStoreRequest 一致。
func AppStoreUpload(path, fileField, filePath string, fields map[string]string, out any) error {
	if Opts.Token == "" {
		return ErrNoAuth
	}
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range fields {
		if v != "" {
			if err := w.WriteField(k, v); err != nil {
				return fmt.Errorf("写入表单字段失败: %w", err)
			}
		}
	}
	fw, err := w.CreateFormFile(fileField, filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("创建上传字段失败: %w", err)
	}
	if _, err := io.Copy(fw, f); err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("封装上传数据失败: %w", err)
	}

	u := strings.TrimRight(Opts.Server, "/") + "/appstore/api/v1" + path
	req, err := http.NewRequest("POST", u, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Token", Opts.Token)
	req.Header.Set("User-Agent", "doo-cli")
	if v := mainAppVersion(); v != "" {
		req.Header.Set("Version", v)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := (&http.Client{Timeout: appStoreTimeout}).Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	return parseAppStoreResponse(resp, out)
}

// parseAppStoreResponse 解析 AppStore 统一响应 {code,message,data}：code!=200 转为 error；
// out 非 nil 时把 data 反序列化进去。
func parseAppStoreResponse(resp *http.Response, out any) error {
	data, _ := io.ReadAll(resp.Body)
	var r struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if r.Code != 200 {
		msg := strings.TrimSpace(r.Message)
		if msg == "" {
			msg = fmt.Sprintf("AppStore 错误 %d", r.Code)
		}
		return fmt.Errorf("%s", msg)
	}
	if out != nil && len(r.Data) > 0 {
		if err := json.Unmarshal(r.Data, out); err != nil {
			return fmt.Errorf("解析响应失败: %w", err)
		}
	}
	return nil
}
