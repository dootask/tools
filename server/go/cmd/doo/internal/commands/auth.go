package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/dootask/tools/server/go/cmd/doo/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "auth", Short: "登录与凭证管理"}
	cmd.AddCommand(newAuthLoginCmd(), newAuthStatusCmd(), newAuthLogoutCmd())
	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	var email, password string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "用邮箱密码登录并保存 token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if email == "" {
				fmt.Fprint(os.Stderr, "邮箱: ")
				line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
				email = strings.TrimSpace(line)
			}
			if password == "" {
				if term.IsTerminal(int(os.Stdin.Fd())) {
					fmt.Fprint(os.Stderr, "密码: ")
					b, err := term.ReadPassword(int(os.Stdin.Fd()))
					fmt.Fprintln(os.Stderr)
					if err != nil {
						return err
					}
					password = string(b)
				} else {
					return fmt.Errorf("缺少 --password")
				}
			}

			c := cli.Opts.AnonClient()
			var resp map[string]any
			err := c.NewGetRequest("/api/users/login", map[string]any{
				"type":     "login",
				"email":    email,
				"password": password,
			}, &resp)
			if err != nil {
				return err
			}
			token, _ := resp["token"].(string)
			if token == "" {
				if code, _ := resp["code"].(string); code == "need" {
					return fmt.Errorf("登录需要验证码，请在浏览器登录后用 --token / DOO_TOKEN 直接传入")
				}
				return fmt.Errorf("登录响应未包含 token")
			}
			if err := config.Save(config.Config{Server: cli.Opts.Server, Token: token}); err != nil {
				return err
			}
			nickname, _ := resp["nickname"].(string)
			cli.OK("✓ 已登录：%s（%s）\n  配置已写入 %s", nickname, cli.Opts.Server, config.Path())
			return nil
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "登录邮箱")
	cmd.Flags().StringVar(&password, "password", "", "登录密码")
	return cmd
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "查看当前登录状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			u, err := c.GetUserInfo()
			if err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(u, nil)
			}
			cli.OK("服务器: %s\n用户:   #%d %s <%s>\n身份:   %s",
				cli.Opts.Server, u.UserID, u.Nickname, u.Email, strings.Join(u.Identity, ","))
			return nil
		},
	}
}

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "清除已保存的 token",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.Load()
			cfg.Token = ""
			if err := config.Save(cfg); err != nil {
				return err
			}
			cli.OK("✓ 已登出")
			return nil
		},
	}
}
