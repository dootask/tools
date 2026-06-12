package commands

import (
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

// file 全部走兜底端点（SDK 暂无文件域），标记为实验性。
func newFileCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "file", Short: "文件（实验性）"}
	cmd.AddCommand(
		newFileListCmd(),
		newFileSearchCmd(),
		newFileViewCmd(),
		newFileFetchCmd(),
	)
	return cmd
}

func newFileListCmd() *cobra.Command {
	var pid int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出文件",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{}
			if pid > 0 {
				params["pid"] = pid
			}
			var out any
			if err := c.NewGetRequest("/api/file/lists", params, &out); err != nil {
				return err
			}
			return cli.Output(out, []string{"id", "name", "type", "ext", "size"})
		},
	}
	cmd.Flags().IntVar(&pid, "pid", 0, "父文件夹 ID（默认根目录）")
	return cmd
}

func newFileSearchCmd() *cobra.Command {
	var take int
	cmd := &cobra.Command{
		Use:   "search <关键词>",
		Short: "搜索文件",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{"keyword": args[0]}
			if take > 0 {
				params["take"] = take
			}
			var out any
			if err := c.NewGetRequest("/api/file/search", params, &out); err != nil {
				return err
			}
			return cli.Output(out, []string{"id", "name", "type", "ext", "size"})
		},
	}
	cmd.Flags().IntVar(&take, "take", 0, "返回数量（最大 100）")
	return cmd
}

func newFileViewCmd() *cobra.Command {
	var content bool
	cmd := &cobra.Command{
		Use:   "view <文件ID>",
		Short: "查看文件详情",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{"file_id": args[0]}
			if content {
				params["with_content"] = true
			}
			var out any
			if err := c.NewGetRequest("/api/file/one", params, &out); err != nil {
				return err
			}
			return cli.Output(out, nil)
		},
	}
	cmd.Flags().BoolVar(&content, "content", false, "附带文本内容")
	return cmd
}

func newFileFetchCmd() *cobra.Command {
	var offset, limit int
	cmd := &cobra.Command{
		Use:   "fetch <路径>",
		Short: "按路径获取文件文本内容",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{"path": args[0]}
			if offset > 0 {
				params["offset"] = offset
			}
			if limit > 0 {
				params["limit"] = limit
			}
			var out any
			if err := c.NewGetRequest("/api/file/fetch", params, &out); err != nil {
				return err
			}
			return cli.Output(out, nil)
		},
	}
	cmd.Flags().IntVar(&offset, "offset", 0, "起始偏移")
	cmd.Flags().IntVar(&limit, "limit", 0, "长度上限")
	return cmd
}
