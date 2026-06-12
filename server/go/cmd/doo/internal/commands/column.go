package commands

import (
	"fmt"

	dootask "github.com/dootask/tools/server/go"
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

func newColumnCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "column", Short: "看板列"}
	cmd.AddCommand(
		newColumnListCmd(),
		newColumnCreateCmd(),
		newColumnUpdateCmd(),
		newColumnDeleteCmd(),
	)
	return cmd
}

func newColumnListCmd() *cobra.Command {
	var project, page, pageSize int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出项目看板列",
		RunE: func(cmd *cobra.Command, args []string) error {
			if project <= 0 {
				return fmt.Errorf("--project 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			res, err := c.GetColumnList(dootask.GetColumnListRequest{ProjectID: project, Page: page, PageSize: pageSize})
			if err != nil {
				return err
			}
			return cli.Output(res, []string{"id", "name", "color", "project_id"})
		},
	}
	f := cmd.Flags()
	f.IntVar(&project, "project", 0, "项目 ID（必填）")
	f.IntVar(&page, "page", 0, "页码")
	f.IntVar(&pageSize, "page-size", 0, "每页数量")
	return cmd
}

func newColumnCreateCmd() *cobra.Command {
	var project int
	var name string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建看板列",
		RunE: func(cmd *cobra.Command, args []string) error {
			if project <= 0 || name == "" {
				return fmt.Errorf("--project 与 --name 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			col, err := c.CreateColumn(dootask.CreateColumnRequest{ProjectID: project, Name: name})
			if err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(col, nil)
			}
			cli.OK("✓ 已创建看板列 #%d：%s", col.ID, col.Name)
			return nil
		},
	}
	cmd.Flags().IntVar(&project, "project", 0, "项目 ID（必填）")
	cmd.Flags().StringVar(&name, "name", "", "列名（必填）")
	return cmd
}

func newColumnUpdateCmd() *cobra.Command {
	var name, color string
	cmd := &cobra.Command{
		Use:   "update <列ID>",
		Short: "更新看板列",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "列ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			f := cmd.Flags()
			params := map[string]any{"column_id": id}
			if f.Changed("name") {
				params["name"] = name
			}
			if f.Changed("color") {
				params["color"] = color
			}
			if len(params) == 1 {
				return fmt.Errorf("没有要更新的字段")
			}
			var out map[string]any
			if err := c.NewGetRequest("/api/project/column/update", params, &out); err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(out, nil)
			}
			cli.OK("✓ 已更新看板列 #%d", id)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "列名")
	cmd.Flags().StringVar(&color, "color", "", "颜色")
	return cmd
}

func newColumnDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <列ID>",
		Short: "删除看板列",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "列ID")
			if err != nil {
				return err
			}
			if err := cli.Confirm(fmt.Sprintf("确认删除看板列 #%d?", id)); err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.DeleteColumn(id); err != nil {
				return err
			}
			cli.OK("✓ 已删除看板列 #%d", id)
			return nil
		},
	}
}
