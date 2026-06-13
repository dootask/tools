package commands

import (
	"fmt"

	dootask "github.com/dootask/tools/server/go"
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

func newMessageCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "message", Aliases: []string{"msg"}, Short: "消息"}
	cmd.AddCommand(
		newMessageSendCmd(),
		newMessageSendUserCmd(),
		newMessageListCmd(),
		newMessageSearchCmd(),
		newMessageViewCmd(),
		newMessageWithdrawCmd(),
		newMessageForwardCmd(),
		newMessageTodoCmd(),
		newMessageTodoListCmd(),
		newMessageTodoRemindCmd(),
		newMessageDoneCmd(),
	)
	return cmd
}

func newMessageSendCmd() *cobra.Command {
	var dialog int
	var text, textType string
	var silence bool
	var reply int
	cmd := &cobra.Command{
		Use:   "send",
		Short: "向对话发送消息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dialog <= 0 || text == "" {
				return fmt.Errorf("--dialog 与 --text 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			req := dootask.SendMessageRequest{DialogID: dialog, Text: text, TextType: textType, Silence: silence, ReplyID: reply}
			var out map[string]any
			if err := c.SendMessage(req, &out); err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(out, nil)
			}
			cli.OK("✓ 已发送到对话 #%d", dialog)
			return nil
		},
	}
	f := cmd.Flags()
	f.IntVar(&dialog, "dialog", 0, "对话 ID（必填）")
	f.StringVar(&text, "text", "", "消息内容（必填）")
	f.StringVar(&textType, "type", "md", "内容类型 md|text|html")
	f.BoolVar(&silence, "silence", false, "静默发送（不通知）")
	f.IntVar(&reply, "reply", 0, "回复的消息 ID")
	return cmd
}

func newMessageSendUserCmd() *cobra.Command {
	var user int
	var text, textType string
	var silence bool
	cmd := &cobra.Command{
		Use:   "send-user",
		Short: "向用户发送私聊消息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if user <= 0 || text == "" {
				return fmt.Errorf("--user 与 --text 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			req := dootask.SendMessageToUserRequest{UserID: user, Text: text, TextType: textType, Silence: silence}
			var out map[string]any
			if err := c.SendMessageToUser(req, &out); err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(out, nil)
			}
			cli.OK("✓ 已发送给用户 #%d", user)
			return nil
		},
	}
	f := cmd.Flags()
	f.IntVar(&user, "user", 0, "用户 ID（必填）")
	f.StringVar(&text, "text", "", "消息内容（必填）")
	f.StringVar(&textType, "type", "md", "内容类型 md|text|html")
	f.BoolVar(&silence, "silence", false, "静默发送")
	return cmd
}

func newMessageListCmd() *cobra.Command {
	var dialog, take int
	var msgType string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出对话消息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dialog <= 0 {
				return fmt.Errorf("--dialog 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			res, err := c.GetMessageList(dootask.GetMessageListRequest{DialogID: dialog, MsgType: msgType, Take: take})
			if err != nil {
				return err
			}
			return cli.Output(res, []string{"id", "userid", "type", "created_at"})
		},
	}
	f := cmd.Flags()
	f.IntVar(&dialog, "dialog", 0, "对话 ID（必填）")
	f.IntVar(&take, "take", 0, "数量（最大 100）")
	f.StringVar(&msgType, "type", "", "类型 tag|todo|link|text|image|file|record|meeting")
	return cmd
}

func newMessageSearchCmd() *cobra.Command {
	var dialog, take int
	var key string
	cmd := &cobra.Command{
		Use:   "search",
		Short: "搜索消息（可选 --dialog 限定对话）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if key == "" {
				return fmt.Errorf("--key 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			res, err := c.SearchMessage(dootask.SearchMessageRequest{Key: key, DialogID: dialog, Take: take})
			if err != nil {
				return err
			}
			return cli.Output(res, []string{"msg_id", "dialog_id", "userid", "type", "created_at"})
		},
	}
	cmd.Flags().IntVar(&dialog, "dialog", 0, "限定对话 ID（默认全局搜索）")
	cmd.Flags().StringVar(&key, "key", "", "关键词（必填）")
	cmd.Flags().IntVar(&take, "take", 0, "返回数量（默认 20，最大 50）")
	return cmd
}

func newMessageViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view <消息ID>",
		Short: "查看消息详情",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "消息ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			m, err := c.GetMessage(dootask.GetMessageRequest{MsgID: id})
			if err != nil {
				return err
			}
			return cli.Output(m, nil)
		},
	}
}

func newMessageWithdrawCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw <消息ID>",
		Short: "撤回消息",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "消息ID")
			if err != nil {
				return err
			}
			if err := cli.Confirm(fmt.Sprintf("确认撤回消息 #%d?", id)); err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.WithdrawMessage(dootask.WithdrawMessageRequest{MsgID: id}); err != nil {
				return err
			}
			cli.OK("✓ 已撤回消息 #%d", id)
			return nil
		},
	}
}

func newMessageForwardCmd() *cobra.Command {
	var dialogs, users string
	cmd := &cobra.Command{
		Use:   "forward <消息ID>",
		Short: "转发消息",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "消息ID")
			if err != nil {
				return err
			}
			dialogIDs, err := cli.ParseIDList(dialogs)
			if err != nil {
				return err
			}
			userIDs, err := cli.ParseIDList(users)
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.ForwardMessage(dootask.ForwardMessageRequest{MsgID: id, DialogIDs: dialogIDs, UserIDs: userIDs}); err != nil {
				return err
			}
			cli.OK("✓ 已转发消息 #%d", id)
			return nil
		},
	}
	cmd.Flags().StringVar(&dialogs, "dialogs", "", "目标对话 ID 列表")
	cmd.Flags().StringVar(&users, "users", "", "目标用户 ID 列表")
	return cmd
}

func newMessageTodoCmd() *cobra.Command {
	var users string
	cmd := &cobra.Command{
		Use:   "todo <消息ID>",
		Short: "把消息标记为待办",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "消息ID")
			if err != nil {
				return err
			}
			userIDs, err := cli.ParseIDList(users)
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.ToggleMessageTodo(dootask.ToggleMessageTodoRequest{MsgID: id, UserIDs: userIDs}); err != nil {
				return err
			}
			cli.OK("✓ 已切换待办 #%d", id)
			return nil
		},
	}
	cmd.Flags().StringVar(&users, "users", "", "指定用户 ID 列表（默认全部）")
	return cmd
}

func newMessageTodoListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "todolist <消息ID>",
		Short: "列出消息的待办记录（其 id 用于 done）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "消息ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			res, err := c.GetMessageTodoList(dootask.GetMessageRequest{MsgID: id})
			if err != nil {
				return err
			}
			return cli.Output(res, []string{"id", "userid", "done_at", "remind_at"})
		},
	}
}

func newMessageTodoRemindCmd() *cobra.Command {
	var users, at string
	cmd := &cobra.Command{
		Use:   "todo-remind <消息ID>",
		Short: "设置/取消待办提醒时间（--at 为空表示取消提醒）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "消息ID")
			if err != nil {
				return err
			}
			userIDs, err := cli.ParseIDList(users)
			if err != nil {
				return err
			}
			if len(userIDs) == 0 {
				return fmt.Errorf("--users 必填（指定要设提醒的成员）")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{"msg_id": id, "userids": userIDs, "remind_at": at}
			if err := c.NewGetRequest("/api/dialog/msg/todoremind", params, nil); err != nil {
				return err
			}
			if at == "" {
				cli.OK("✓ 已取消消息 #%d 的待办提醒", id)
			} else {
				cli.OK("✓ 已设置消息 #%d 的待办提醒：%s", id, at)
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&users, "users", "", "目标成员 ID 列表（必填）")
	f.StringVar(&at, "at", "", "提醒时间（如 2026-06-20 09:00:00，空=取消）")
	return cmd
}

func newMessageDoneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "done <待办ID>",
		Short: "完成待办（待办ID 来自 message todolist，非消息ID）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "待办ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.MarkMessageDone(dootask.MarkMessageDoneRequest{ID: id}); err != nil {
				return err
			}
			cli.OK("✓ 已完成待办 #%d", id)
			return nil
		},
	}
}
