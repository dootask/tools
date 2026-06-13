// Package commands 实现 doo 的命令树。
package commands

import (
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

// Version 由构建期 -ldflags 注入。
var Version = "dev"

// NewRootCmd 构造根命令并挂载全部子命令与全局 flag。
func NewRootCmd() *cobra.Command {
	var fServer, fToken string
	var fJSON, fYes, fQuiet bool

	root := &cobra.Command{
		Use:               "doo",
		Short:             "DooTask 命令行工具",
		Long:              "doo —— DooTask 的命令行工具，覆盖任务/项目/对话/消息/群组/文件/报告等操作。",
		Version:           Version,
		SilenceUsage:      true,
		SilenceErrors:     true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cli.Resolve(fServer, fToken, fJSON, fYes, fQuiet)
			return nil
		},
	}

	pf := root.PersistentFlags()
	pf.StringVar(&fServer, "server", "", "DooTask 服务器地址（默认 env DOO_SERVER 或配置）")
	pf.StringVar(&fToken, "token", "", "API token（默认 env DOO_TOKEN 或配置）")
	pf.BoolVar(&fJSON, "json", false, "以 JSON 输出")
	pf.BoolVarP(&fYes, "yes", "y", false, "跳过危险操作确认")
	pf.BoolVarP(&fQuiet, "quiet", "q", false, "精简输出")

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
	)
	return root
}
