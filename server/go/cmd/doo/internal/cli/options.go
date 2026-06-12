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

// Options 是合并后的全局运行参数。
type Options struct {
	Server string
	Token  string
	JSON   bool
	Yes    bool
	Quiet  bool
}

// Opts 是本次调用生效的全局参数（CLI 单次执行，进程级单例）。
var Opts Options

// Resolve 按 flag > env > 配置文件 > 默认 的优先级合并参数。
func Resolve(flagServer, flagToken string, jsonOut, yes, quiet bool) {
	cfg, _ := config.Load()
	Opts = Options{
		Server: first(flagServer, os.Getenv("DOO_SERVER"), cfg.Server, "http://nginx"),
		Token:  first(flagToken, os.Getenv("DOO_TOKEN"), cfg.Token),
		JSON:   jsonOut,
		Yes:    yes,
		Quiet:  quiet,
	}
}

// Client 用当前 token/server 构造 SDK 客户端；缺 token 时返回 ErrNoAuth。
func (o Options) Client() (*dootask.Client, error) {
	if o.Token == "" {
		return nil, ErrNoAuth
	}
	return dootask.NewClient(o.Token, dootask.WithServer(o.Server), dootask.WithTimeout(30*time.Second)), nil
}

// AnonClient 构造无 token 的客户端（仅用于登录换 token）。
func (o Options) AnonClient() *dootask.Client {
	return dootask.NewClient("", dootask.WithServer(o.Server), dootask.WithTimeout(30*time.Second))
}

// ExitCode 把错误映射为进程退出码。
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	if errors.Is(err, ErrNoAuth) {
		return 3
	}
	return 1
}

// OK 打印一条成功/提示信息（--quiet 时抑制）。
func OK(format string, a ...any) {
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
