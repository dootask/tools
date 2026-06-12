package commands

import (
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

func newSystemCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "system", Short: "系统信息"}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "version",
			Short: "查看 DooTask 版本",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := cli.Opts.Client()
				if err != nil {
					return err
				}
				v, err := c.GetVersion()
				if err != nil {
					return err
				}
				return cli.Output(v, nil)
			},
		},
		&cobra.Command{
			Use:   "settings",
			Short: "查看系统设置",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := cli.Opts.Client()
				if err != nil {
					return err
				}
				s, err := c.GetSystemSettings()
				if err != nil {
					return err
				}
				return cli.Output(s, nil)
			},
		},
	)
	return cmd
}
