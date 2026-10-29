// Package cli 提供 doo 各命令共享的运行时：配置解析、客户端构造、输出与确认。
package cli

import (
	"errors"
	"fmt"
	"os"
	"time"

	dootask "github.com/dootask/tools/server/go"
	"github.com/dootask/tools/server/go/cmd/doo/internal/config"
)

// ErrNoAuth 表示未提供 token（未登录）。
var ErrNoAuth = errors.New("未登录：请先执行 `doo auth login`，或设置环境变量 DOO_TOKEN")

// CompatVersion 是 doo 当前适配到的 DooTask 主程序版本，作为 Version 头随请求发送，
// 供后端 checkClientVersion 判定接口兼容性。当 doo 适配更新的主程序 API 时，同步上调。
const CompatVersion = "1.7.91"

// Options 是合并后的全局运行参数。
type Options struct {
	Server string
	Token  string
	JSON   bool
	Yes    bool
	Quiet  bool
	DryRun bool
	Fields []string
	Limit  int
	JQ     string
}

// Opts 是本次调用生效的全局参数（CLI 单次执行，进程级单例）。
var Opts Options

// Resolve 按 flag > env > 配置文件 > 默认 的优先级补齐 Server/Token，其余字段取自 flag。
func Resolve(o Options) {
	cfg, _ := config.Load()
	o.Server = first(o.Server, os.Getenv("DOO_SERVER"), cfg.Server, "http://nginx")
	o.Token = first(o.Token, os.Getenv("DOO_TOKEN"), cfg.Token)
	Opts = o
}

// Client 用当前 token/server 构造 SDK 客户端；缺 token 时返回 ErrNoAuth。
func (o Options) Client() (*dootask.Client, error) {
	if o.Token == "" {
		return nil, ErrNoAuth
	}
	if err := guard(); err != nil {
		return nil, err
	}
	return dootask.NewClient(o.Token, dootask.WithServer(o.Server), dootask.WithTimeout(30*time.Second), dootask.WithVersion(CompatVersion)), nil
}

// AnonClient 构造无 token 的客户端（仅用于登录换 token）。
func (o Options) AnonClient() *dootask.Client {
	return dootask.NewClient("", dootask.WithServer(o.Server), dootask.WithTimeout(30*time.Second))
}

// OK 打印一条成功/提示信息（--quiet 时抑制）；--json 时并入信封的 message 字段。
func OK(format string, a ...any) {
	if Opts.JSON {
		Current.messages = append(Current.messages, fmt.Sprintf(format, a...))
		return
	}
	if !Opts.Quiet {
		fmt.Printf(format+"\n", a...)
	}
}

func first(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
