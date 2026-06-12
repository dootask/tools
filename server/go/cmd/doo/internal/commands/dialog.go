package commands

import (
	dootask "github.com/dootask/tools/server/go"
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

var dialogCols = []string{"id", "type", "name", "last_at"}

func newDialogCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "dialog", Short: "对话"}
	cmd.AddCommand(
		newDialogListCmd(),
		newDialogSearchCmd(),
		newDialogViewCmd(),
		newDialogUsersCmd(),
	)
	return cmd
}

func newDialogListCmd() *cobra.Command {
	var timeRange string
	var page, pageSize int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出对话",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			res, err := c.GetDialogList(dootask.TimeRangeRequest{TimeRange: timeRange, Page: page, PageSize: pageSize})
			if err != nil {
				return err
			}
			return cli.Output(res, dialogCols)
		},
	}
	f := cmd.Flags()
	f.StringVar(&timeRange, "time", "", "时间范围")
	f.IntVar(&page, "page", 0, "页码")
	f.IntVar(&pageSize, "page-size", 0, "每页数量")
	return cmd
}

func newDialogSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <关键词>",
		Short: "搜索对话",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			res, err := c.SearchDialog(dootask.SearchDialogRequest{Key: args[0]})
			if err != nil {
				return err
			}
			return cli.Output(res, dialogCols)
		},
	}
}

func newDialogViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view <对话ID>",
		Short: "查看对话详情",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "对话ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			d, err := c.GetDialogOne(dootask.GetDialogRequest{DialogID: id})
			if err != nil {
				return err
			}
			return cli.Output(d, nil)
		},
	}
}

func newDialogUsersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "users <对话ID>",
		Short: "查看对话成员",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "对话ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			m, err := c.GetDialogUser(dootask.GetDialogUserRequest{DialogID: id, GetUser: 1})
			if err != nil {
				return err
			}
			return cli.Output(m, []string{"userid", "nickname", "email"})
		},
	}
}
