// Package commands 实现 doo 的命令树。
package commands

import (
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Version 由构建期 -ldflags 注入。
var Version = "dev"

// NewRootCmd 构造根命令并挂载全部子命令与全局 flag。
func NewRootCmd() *cobra.Command {
	var o cli.Options

	root := &cobra.Command{
		Use:           "doo",
		Short:         "DooTask 命令行工具",
		Long:          "doo —— DooTask 的命令行工具，覆盖任务/项目/对话/消息/群组/文件/报告等操作。",
		Version:       Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cli.Resolve(o)
			flags := map[string]string{}
			cmd.LocalNonPersistentFlags().Visit(func(f *pflag.Flag) {
				flags[f.Name] = f.Value.String()
			})
			cli.Begin(cmdPath(cmd), riskOf(cmd), args, flags)
			return nil
		},
	}

	pf := root.PersistentFlags()
	pf.StringVar(&o.Server, "server", "", "DooTask 服务器地址（默认 env DOO_SERVER 或配置）")
	pf.StringVar(&o.Token, "token", "", "API token（默认 env DOO_TOKEN 或配置）")
	pf.BoolVar(&o.JSON, "json", false, `以 JSON 输出：{"ok":true,"data":...} / {"ok":false,"code":...,"message":...}`)
	pf.BoolVarP(&o.Yes, "yes", "y", false, "跳过危险操作确认")
	pf.BoolVarP(&o.Quiet, "quiet", "q", false, "精简输出")
	pf.BoolVar(&o.DryRun, "dry-run", false, "预演：校验参数并返回将执行的操作，不实际执行")
	pf.StringSliceVar(&o.Fields, "fields", nil, "只输出指定字段（逗号分隔）")
	pf.IntVar(&o.Limit, "limit", 0, "列表最多输出条数")
	pf.StringVar(&o.JQ, "jq", "", "用 jq 表达式处理输出")

	root.AddCommand(
		newAuthCmd(),
		newTaskCmd(),
		newProjectCmd(),
		newColumnCmd(),
		newFlowCmd(),
		newTagCmd(),
		newDialogCmd(),
		newMessageCmd(),
		newGroupCmd(),
		newUserCmd(),
		newBotCmd(),
		newFileCmd(),
		newReportCmd(),
		newSearchCmd(),
		newPageCmd(),
		newAppCmd(),
		newSystemCmd(),
		newMetaCmd(),
	)
	return root
}
