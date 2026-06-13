package commands

import (
	"fmt"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

func newTagCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "tag", Short: "任务标签"}
	cmd.AddCommand(
		newTagListCmd(),
		newTagSaveCmd(false),
		newTagSaveCmd(true),
		newTagDeleteCmd(),
	)
	return cmd
}

func newTagListCmd() *cobra.Command {
	var project int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出项目的任务标签",
		RunE: func(cmd *cobra.Command, args []string) error {
			if project <= 0 {
				return fmt.Errorf("--project 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			var out any
			if err := c.NewGetRequest("/api/project/tag/list", map[string]any{"project_id": project}, &out); err != nil {
				return err
			}
			return cli.Output(out, []string{"id", "name", "color"})
		},
	}
	cmd.Flags().IntVar(&project, "project", 0, "项目 ID（必填）")
	return cmd
}

// newTagSaveCmd 复用一个构造器实现 create / update（update 需 <标签ID>）。
func newTagSaveCmd(update bool) *cobra.Command {
	var project int
	var name, color, desc string
	use, short, nargs := "create", "新建任务标签", cobra.NoArgs
	if update {
		use, short, nargs = "update <标签ID>", "修改任务标签", cobra.ExactArgs(1)
	}
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  nargs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if project <= 0 {
				return fmt.Errorf("--project 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{"project_id": project, "name": name, "color": color, "desc": desc}
			if update {
				id, err := cli.ParseInt(args[0], "标签ID")
				if err != nil {
					return err
				}
				params["id"] = id
			}
			var out any
			if err := c.NewGetRequest("/api/project/tag/save", params, &out); err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(out, nil)
			}
			cli.OK("✓ 已保存标签：%s", name)
			return nil
		},
	}
	f := cmd.Flags()
	f.IntVar(&project, "project", 0, "项目 ID（必填）")
	f.StringVar(&name, "name", "", "标签名称（必填）")
	f.StringVar(&color, "color", "", "颜色")
	f.StringVar(&desc, "desc", "", "描述")
	return cmd
}

func newTagDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <标签ID>",
		Short: "删除任务标签",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "标签ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.NewGetRequest("/api/project/tag/delete", map[string]any{"id": id}, nil); err != nil {
				return err
			}
			cli.OK("✓ 已删除标签 #%d", id)
			return nil
		},
	}
}
