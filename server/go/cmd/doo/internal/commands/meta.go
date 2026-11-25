package commands

import (
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type metaFlag struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Usage   string `json:"usage"`
	Default string `json:"default,omitempty"`
}

type metaCmd struct {
	Command string     `json:"command"`
	Use     string     `json:"use"`
	Short   string     `json:"short"`
	Risk    string     `json:"risk"`
	Flags   []metaFlag `json:"flags,omitempty"`
}

func newMetaCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "meta",
		Short: "输出机器可读的命令清单（名称、用法、参数、风险级别）",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var out []metaCmd
			var walk func(c *cobra.Command)
			walk = func(c *cobra.Command) {
				for _, s := range c.Commands() {
					if s.Hidden || s.Name() == "help" || s.Name() == "completion" {
						continue
					}
					if s.HasSubCommands() {
						walk(s)
						continue
					}
					m := metaCmd{Command: cmdPath(s), Use: s.UseLine(), Short: s.Short, Risk: riskOf(s)}
					s.LocalNonPersistentFlags().VisitAll(func(f *pflag.Flag) {
						m.Flags = append(m.Flags, metaFlag{Name: f.Name, Type: f.Value.Type(), Usage: f.Usage, Default: f.DefValue})
					})
					out = append(out, m)
				}
			}
			walk(cmd.Root())
			return cli.Output(map[string]any{"version": Version, "compat_version": cli.CompatVersion, "commands": out}, nil)
		},
	}
}
