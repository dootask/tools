package commands

import (
	"fmt"

	dootask "github.com/dootask/tools/server/go"
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

func newProjectCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "project", Short: "项目"}
	cmd.AddCommand(
		newProjectListCmd(),
		newProjectViewCmd(),
		newProjectCreateCmd(),
		newProjectUpdateCmd(),
		newProjectExitCmd(),
		newProjectDeleteCmd(),
	)
	return cmd
}

func newProjectListCmd() *cobra.Command {
	var typ, archived string
	var page, pageSize int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出项目",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			res, err := c.GetProjectList(dootask.GetProjectListRequest{
				Type:     typ,
				Archived: archived,
				Page:     page,
				PageSize: pageSize,
			})
			if err != nil {
				return err
			}
			return cli.Output(res, []string{"id", "name", "desc", "owner_userid", "dialog_id"})
		},
	}
	f := cmd.Flags()
	f.StringVar(&typ, "type", "", "类型，如 personal")
	f.StringVar(&archived, "archived", "", "归档过滤 all|yes|no")
	f.IntVar(&page, "page", 0, "页码")
	f.IntVar(&pageSize, "page-size", 0, "每页数量")
	return cmd
}

func newProjectViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view <项目ID>",
		Short: "查看项目详情",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "项目ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			p, err := c.GetProject(dootask.GetProjectRequest{ProjectID: id})
			if err != nil {
				return err
			}
			return cli.Output(p, nil)
		},
	}
}

func newProjectCreateCmd() *cobra.Command {
	var name, desc, columns, flow string
	var personal bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建项目",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			req := dootask.CreateProjectRequest{Name: name, Desc: desc, Columns: columns, Flow: flow}
			if personal {
				req.Personal = 1
			}
			p, err := c.CreateProject(req)
			if err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(p, nil)
			}
			cli.OK("✓ 已创建项目 #%d：%s", p.ID, p.Name)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&name, "name", "", "项目名称（必填）")
	f.StringVar(&desc, "desc", "", "项目描述")
	f.StringVar(&columns, "columns", "", "初始列名，逗号分隔")
	f.StringVar(&flow, "flow", "", "工作流 open|close")
	f.BoolVar(&personal, "personal", false, "创建为个人项目")
	return cmd
}

func newProjectUpdateCmd() *cobra.Command {
	var name, desc, archiveMethod string
	var archiveDays int
	cmd := &cobra.Command{
		Use:   "update <项目ID>",
		Short: "更新项目（--name 必填）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "项目ID")
			if err != nil {
				return err
			}
			if name == "" {
				return fmt.Errorf("--name 必填（接口要求）")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			f := cmd.Flags()
			params := map[string]any{"project_id": id, "name": name}
			if f.Changed("desc") {
				params["desc"] = desc
			}
			if f.Changed("archive-method") {
				params["archive_method"] = archiveMethod
			}
			if f.Changed("archive-days") {
				params["archive_days"] = archiveDays
			}
			var out map[string]any
			if err := c.NewGetRequest("/api/project/update", params, &out); err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(out, nil)
			}
			cli.OK("✓ 已更新项目 #%d", id)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&name, "name", "", "项目名称（必填）")
	f.StringVar(&desc, "desc", "", "项目描述")
	f.StringVar(&archiveMethod, "archive-method", "", "归档方式")
	f.IntVar(&archiveDays, "archive-days", 0, "自动归档天数")
	return cmd
}

func newProjectExitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "exit <项目ID>",
		Short: "退出项目",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "项目ID")
			if err != nil {
				return err
			}
			if err := cli.Confirm(fmt.Sprintf("确认退出项目 #%d?", id)); err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.ExitProject(id); err != nil {
				return err
			}
			cli.OK("✓ 已退出项目 #%d", id)
			return nil
		},
	}
}

func newProjectDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <项目ID>",
		Short: "删除项目",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "项目ID")
			if err != nil {
				return err
			}
			if err := cli.Confirm(fmt.Sprintf("确认删除项目 #%d（不可逆）?", id)); err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.DeleteProject(id); err != nil {
				return err
			}
			cli.OK("✓ 已删除项目 #%d", id)
			return nil
		},
	}
}
