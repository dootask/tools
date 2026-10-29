package commands

import (
	"strings"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

// dangerCmds 危险命令：不可逆或影响他人，执行前必须确认（--yes / 交互 y / 调用方确认卡片）。
var dangerCmds = map[string]bool{
	"task delete":       true,
	"project delete":    true,
	"project exit":      true,
	"column delete":     true,
	"tag delete":        true,
	"group remove-user": true,
	"group exit":        true,
	"group transfer":    true,
	"group disband":     true,
	"file delete":       true,
	"file restore":      true,
	"bot delete":        true,
	"message withdraw":  true,
	"app uninstall":     true,
	"app remove":        true,
	"app reinstall":     true,
}

// readCmds 只读命令，未列入的叶子命令按 write 处理。
var readCmds = map[string]bool{
	"app catalog": true, "app container-logs": true, "app containers": true, "app fields": true,
	"app list": true, "app logs": true, "app updates": true,
	"auth status": true,
	"bot list":    true, "bot view": true,
	"column list":  true,
	"dialog inbox": true, "dialog list": true, "dialog mytodo": true, "dialog search": true,
	"dialog unread": true, "dialog users": true, "dialog view": true,
	"file fetch": true, "file history": true, "file list": true, "file search": true, "file view": true,
	"flow list":    true,
	"message list": true, "message search": true, "message todolist": true, "message view": true,
	"page context": true,
	"project list": true, "project view": true,
	"report my": true, "report received": true, "report template": true, "report unread": true, "report view": true,
	"search":          true,
	"system settings": true, "system version": true,
	"tag list":   true,
	"task files": true, "task list": true, "task view": true,
	"user basic": true, "user departments": true, "user info": true, "user search": true,
	"meta": true,
}

// cmdPath 命令路径（去掉根命令名），如 "task delete"。
func cmdPath(cmd *cobra.Command) string {
	return strings.TrimSpace(strings.TrimPrefix(cmd.CommandPath(), cmd.Root().Name()))
}

// riskOf 命令风险级别；--recover 恢复类操作降为 write。
func riskOf(cmd *cobra.Command) string {
	p := cmdPath(cmd)
	switch {
	case readCmds[p]:
		return cli.RiskRead
	case dangerCmds[p]:
		if f := cmd.Flags().Lookup("recover"); f != nil && f.Changed {
			return cli.RiskWrite
		}
		return cli.RiskDanger
	}
	return cli.RiskWrite
}
