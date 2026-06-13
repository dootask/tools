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
		newDialogInboxCmd(),
		newDialogMyTodoCmd(),
		newDialogUnreadCmd(),
		newDialogReadCmd(),
	)
	return cmd
}

func newDialogInboxCmd() *cobra.Command {
	var unreadAt, todoAt string
	cmd := &cobra.Command{
		Use:   "inbox",
		Short: "列出（列表外的）有未读或待办的对话",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{}
			if unreadAt != "" {
				params["unread_at"] = unreadAt
			}
			if todoAt != "" {
				params["todo_at"] = todoAt
			}
			var out any
			if err := c.NewGetRequest("/api/dialog/beyond", params, &out); err != nil {
				return err
			}
			return cli.Output(out, nil)
		},
	}
	f := cmd.Flags()
	f.StringVar(&unreadAt, "unread-at", "", "只取该时间之后有未读的对话")
	f.StringVar(&todoAt, "todo-at", "", "只取该时间之后有待办的对话")
	return cmd
}

func newDialogMyTodoCmd() *cobra.Command {
	var dialog int
	cmd := &cobra.Command{
		Use:   "mytodo",
		Short: "列出我未完成的待办（可按对话过滤；其 id 用于 message done）",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{}
			if dialog > 0 {
				params["dialog_id"] = dialog
			}
			var out any
			if err := c.NewGetRequest("/api/dialog/todo", params, &out); err != nil {
				return err
			}
			return cli.Output(out, []string{"id", "dialog_id", "msg_id", "userid", "done_at"})
		},
	}
	cmd.Flags().IntVar(&dialog, "dialog", 0, "限定对话 ID（默认全部）")
	return cmd
}

func newDialogUnreadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unread <对话ID>",
		Short: "查看对话的未读/提及统计",
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
			var out any
			if err := c.NewGetRequest("/api/dialog/msg/unread", map[string]any{"dialog_id": id}, &out); err != nil {
				return err
			}
			return cli.Output(out, nil)
		},
	}
}

func newDialogReadCmd() *cobra.Command {
	var after int
	cmd := &cobra.Command{
		Use:   "read <对话ID>",
		Short: "把对话标记为已读（可选只标记某消息ID及之后）",
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
			params := map[string]any{"dialog_id": id, "type": "read"}
			if after > 0 {
				params["after_msg_id"] = after
			}
			if err := c.NewGetRequest("/api/dialog/msg/mark", params, nil); err != nil {
				return err
			}
			cli.OK("✓ 对话 #%d 已标记为已读", id)
			return nil
		},
	}
	cmd.Flags().IntVar(&after, "after", 0, "只标记该消息 ID 及之后为已读")
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
