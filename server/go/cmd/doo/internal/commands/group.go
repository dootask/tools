package commands

import (
	"fmt"

	dootask "github.com/dootask/tools/server/go"
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

func newGroupCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "group", Short: "群组"}
	cmd.AddCommand(
		newGroupCreateCmd(),
		newGroupEditCmd(),
		newGroupAddUserCmd(),
		newGroupRemoveUserCmd(),
		newGroupExitCmd(),
		newGroupTransferCmd(),
		newGroupDisbandCmd(),
	)
	return cmd
}

func newGroupCreateCmd() *cobra.Command {
	var users, name, avatar string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建群组",
		RunE: func(cmd *cobra.Command, args []string) error {
			ids, err := cli.ParseIDList(users)
			if err != nil {
				return err
			}
			if len(ids) == 0 {
				return fmt.Errorf("--users 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			d, err := c.CreateGroup(dootask.CreateGroupRequest{ChatName: name, Avatar: avatar, UserIDs: ids})
			if err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(d, nil)
			}
			cli.OK("✓ 已创建群组 #%d", d.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&users, "users", "", "成员 ID 列表（必填）")
	cmd.Flags().StringVar(&name, "name", "", "群名称")
	cmd.Flags().StringVar(&avatar, "avatar", "", "群头像 URL")
	return cmd
}

func newGroupEditCmd() *cobra.Command {
	var name, avatar string
	cmd := &cobra.Command{
		Use:   "edit <对话ID>",
		Short: "编辑群组",
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
			f := cmd.Flags()
			params := map[string]any{"dialog_id": id}
			if f.Changed("name") {
				params["chat_name"] = name
			}
			if f.Changed("avatar") {
				params["avatar"] = avatar
			}
			if len(params) == 1 {
				return fmt.Errorf("没有要更新的字段")
			}
			if err := c.NewGetRequest("/api/dialog/group/edit", params, nil); err != nil {
				return err
			}
			cli.OK("✓ 已更新群组 #%d", id)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "群名称")
	cmd.Flags().StringVar(&avatar, "avatar", "", "群头像 URL")
	return cmd
}

func newGroupAddUserCmd() *cobra.Command {
	var users string
	cmd := &cobra.Command{
		Use:   "add-user <对话ID>",
		Short: "添加群成员",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "对话ID")
			if err != nil {
				return err
			}
			ids, err := cli.ParseIDList(users)
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.AddGroupUser(dootask.AddGroupUserRequest{DialogID: id, UserIDs: ids}); err != nil {
				return err
			}
			cli.OK("✓ 已添加成员到群 #%d", id)
			return nil
		},
	}
	cmd.Flags().StringVar(&users, "users", "", "成员 ID 列表（必填）")
	return cmd
}

func newGroupRemoveUserCmd() *cobra.Command {
	var users string
	cmd := &cobra.Command{
		Use:   "remove-user <对话ID>",
		Short: "移除群成员",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "对话ID")
			if err != nil {
				return err
			}
			ids, err := cli.ParseIDList(users)
			if err != nil {
				return err
			}
			if err := cli.Confirm(fmt.Sprintf("确认从群 #%d 移除成员 %v?", id, ids)); err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.RemoveGroupUser(dootask.RemoveGroupUserRequest{DialogID: id, UserIDs: ids}); err != nil {
				return err
			}
			cli.OK("✓ 已从群 #%d 移除成员", id)
			return nil
		},
	}
	cmd.Flags().StringVar(&users, "users", "", "成员 ID 列表")
	return cmd
}

func newGroupExitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "exit <对话ID>",
		Short: "退出群组",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "对话ID")
			if err != nil {
				return err
			}
			if err := cli.Confirm(fmt.Sprintf("确认退出群 #%d?", id)); err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.ExitGroup(id); err != nil {
				return err
			}
			cli.OK("✓ 已退出群 #%d", id)
			return nil
		},
	}
}

func newGroupTransferCmd() *cobra.Command {
	var user int
	cmd := &cobra.Command{
		Use:   "transfer <对话ID>",
		Short: "转让群主",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "对话ID")
			if err != nil {
				return err
			}
			if user <= 0 {
				return fmt.Errorf("--user 必填")
			}
			if err := cli.Confirm(fmt.Sprintf("确认把群 #%d 转让给用户 #%d?", id, user)); err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.TransferGroup(dootask.TransferGroupRequest{DialogID: id, UserID: user}); err != nil {
				return err
			}
			cli.OK("✓ 已转让群 #%d 给 #%d", id, user)
			return nil
		},
	}
	cmd.Flags().IntVar(&user, "user", 0, "新群主用户 ID（必填）")
	return cmd
}

func newGroupDisbandCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disband <对话ID>",
		Short: "解散群组",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "对话ID")
			if err != nil {
				return err
			}
			if err := cli.Confirm(fmt.Sprintf("确认解散群 #%d（不可逆）?", id)); err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.DisbandGroup(dootask.DisbandGroupRequest{DialogID: id}); err != nil {
				return err
			}
			cli.OK("✓ 已解散群 #%d", id)
			return nil
		},
	}
}
