package commands

import (
	"fmt"

	dootask "github.com/dootask/tools/server/go"
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

func newBotCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "bot", Short: "机器人"}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "列出机器人",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := cli.Opts.Client()
				if err != nil {
					return err
				}
				res, err := c.GetBotList()
				if err != nil {
					return err
				}
				return cli.Output(res, nil)
			},
		},
		&cobra.Command{
			Use:   "view <机器人ID>",
			Short: "查看机器人",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				id, err := cli.ParseInt(args[0], "机器人ID")
				if err != nil {
					return err
				}
				c, err := cli.Opts.Client()
				if err != nil {
					return err
				}
				b, err := c.GetBot(dootask.GetBotRequest{ID: id})
				if err != nil {
					return err
				}
				return cli.Output(b, nil)
			},
		},
		newBotCreateCmd(),
		newBotUpdateCmd(),
		newBotDeleteCmd(),
	)
	return cmd
}

func newBotCreateCmd() *cobra.Command {
	var name, avatar, webhook string
	var clearDay int
	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建机器人",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			b, err := c.CreateBot(dootask.CreateBotRequest{Name: name, Avatar: avatar, WebhookURL: webhook, ClearDay: clearDay})
			if err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(b, nil)
			}
			cli.OK("✓ 已创建机器人 #%d：%s", b.ID, b.Name)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&name, "name", "", "名称（必填）")
	f.StringVar(&avatar, "avatar", "", "头像 URL")
	f.StringVar(&webhook, "webhook", "", "Webhook 地址")
	f.IntVar(&clearDay, "clear-day", 0, "消息清理天数")
	return cmd
}

func newBotUpdateCmd() *cobra.Command {
	var name, avatar, webhook string
	var clearDay int
	cmd := &cobra.Command{
		Use:   "update <机器人ID>",
		Short: "更新机器人",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "机器人ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			f := cmd.Flags()
			params := map[string]any{"id": id}
			if f.Changed("name") {
				params["name"] = name
			}
			if f.Changed("avatar") {
				params["avatar"] = avatar
			}
			if f.Changed("webhook") {
				params["webhook_url"] = webhook
			}
			if f.Changed("clear-day") {
				params["clear_day"] = clearDay
			}
			var out map[string]any
			if err := c.NewPostRequest("/api/users/bot/edit", params, &out); err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(out, nil)
			}
			cli.OK("✓ 已更新机器人 #%d", id)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&name, "name", "", "名称")
	f.StringVar(&avatar, "avatar", "", "头像 URL")
	f.StringVar(&webhook, "webhook", "", "Webhook 地址")
	f.IntVar(&clearDay, "clear-day", 0, "消息清理天数")
	return cmd
}

func newBotDeleteCmd() *cobra.Command {
	var remark string
	cmd := &cobra.Command{
		Use:   "delete <机器人ID>",
		Short: "删除机器人",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "机器人ID")
			if err != nil {
				return err
			}
			if remark == "" {
				return fmt.Errorf("--remark 必填")
			}
			if err := cli.Confirm(fmt.Sprintf("确认删除机器人 #%d?", id)); err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.DeleteBot(dootask.DeleteBotRequest{ID: id, Remark: remark}); err != nil {
				return err
			}
			cli.OK("✓ 已删除机器人 #%d", id)
			return nil
		},
	}
	cmd.Flags().StringVar(&remark, "remark", "", "删除备注（必填）")
	return cmd
}
