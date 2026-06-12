package commands

import (
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

func newUserCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "user", Short: "用户"}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "info",
			Short: "当前登录用户信息",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := cli.Opts.Client()
				if err != nil {
					return err
				}
				u, err := c.GetUserInfo()
				if err != nil {
					return err
				}
				return cli.Output(u, nil)
			},
		},
		&cobra.Command{
			Use:   "departments",
			Short: "我的部门",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := cli.Opts.Client()
				if err != nil {
					return err
				}
				d, err := c.GetUserDepartments()
				if err != nil {
					return err
				}
				return cli.Output(d, []string{"id", "name", "parent_id", "owner_userid"})
			},
		},
		&cobra.Command{
			Use:   "basic <id> [id...]",
			Short: "按 ID 批量获取用户基础信息",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := cli.Opts.Client()
				if err != nil {
					return err
				}
				ids := make([]int, 0, len(args))
				for _, a := range args {
					id, err := cli.ParseInt(a, "用户ID")
					if err != nil {
						return err
					}
					ids = append(ids, id)
				}
				u, err := c.GetUsersBasic(ids)
				if err != nil {
					return err
				}
				return cli.Output(u, []string{"userid", "nickname", "email", "profession", "department_name"})
			},
		},
		newUserSearchCmd(),
	)
	return cmd
}

func newUserSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <关键词>",
		Short: "搜索用户（兜底端点）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			var out any
			if err := c.NewGetRequest("/api/users/search", map[string]any{"keyword": args[0]}, &out); err != nil {
				return err
			}
			return cli.Output(out, []string{"userid", "nickname", "email"})
		},
	}
}
