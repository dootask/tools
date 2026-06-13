package commands

import (
	"fmt"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

func newFlowCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "flow", Short: "工作流"}
	cmd.AddCommand(newFlowListCmd())
	return cmd
}

func newFlowListCmd() *cobra.Command {
	var project int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出项目工作流及各状态（状态 id 用于 task move/update --flow）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if project <= 0 {
				return fmt.Errorf("--project 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			var out any
			if err := c.NewGetRequest("/api/project/flow/list", map[string]any{"project_id": project}, &out); err != nil {
				return err
			}
			return cli.Output(out, nil)
		},
	}
	cmd.Flags().IntVar(&project, "project", 0, "项目 ID（必填）")
	return cmd
}
