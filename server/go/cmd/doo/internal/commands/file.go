package commands

import (
	"encoding/json"
	"fmt"
	"os"

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
		newFileHistoryCmd(),
		newFileRestoreCmd(),
		newFileCreateCmd(),
		newFileSaveCmd(),
		newFileRenameCmd(),
		newFileDeleteCmd(),
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
			params := map[string]any{"key": args[0]}
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
			params := map[string]any{"id": args[0]}
			if content {
				params["with_text"] = "yes"
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

func newFileHistoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history <文件ID>",
		Short: "查看文件内容修改历史",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "文件ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			var out any
			if err := c.NewGetRequest("/api/file/content/history", map[string]any{"id": id}, &out); err != nil {
				return err
			}
			return cli.Output(out, []string{"id", "userid", "size", "created_at"})
		},
	}
	return cmd
}

func newFileRestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore <文件ID> <历史ID>",
		Short: "恢复文件到指定历史版本",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "文件ID")
			if err != nil {
				return err
			}
			historyID, err := cli.ParseInt(args[1], "历史ID")
			if err != nil {
				return err
			}
			if err := cli.Confirm(fmt.Sprintf("将文件 #%d 恢复到历史版本 #%d", id, historyID), cli.Target{Type: "file", ID: id}); err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{"id": id, "history_id": historyID}
			if err := c.NewGetRequest("/api/file/content/restore", params, nil); err != nil {
				return err
			}
			cli.OK("✓ 已将文件 #%d 恢复到历史版本 #%d", id, historyID)
			return nil
		},
	}
	return cmd
}

var fileCreateTypes = []string{"folder", "document", "mind", "drawio", "word", "excel", "ppt"}

func newFileCreateCmd() *cobra.Command {
	var ftype string
	var pid int
	cmd := &cobra.Command{
		Use:   "create <名称>",
		Short: "新建文件（夹）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if !containsStr(fileCreateTypes, ftype) {
				return fmt.Errorf("--type 只能是：%v", fileCreateTypes)
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{"name": name, "type": ftype}
			if pid > 0 {
				params["pid"] = pid
			}
			var out any
			if err := c.NewGetRequest("/api/file/add", params, &out); err != nil {
				return err
			}
			return cli.Output(out, []string{"id", "name", "type", "ext"})
		},
	}
	cmd.Flags().StringVar(&ftype, "type", "document", "类型：folder/document/mind/drawio/word/excel/ppt")
	cmd.Flags().IntVar(&pid, "pid", 0, "父文件夹 ID（默认根目录）")
	return cmd
}

func newFileSaveCmd() *cobra.Command {
	var content, file, format string
	cmd := &cobra.Command{
		Use:   "save <文件ID>",
		Short: "保存文件内容（仅 document 类型）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "文件ID")
			if err != nil {
				return err
			}
			if (content == "") == (file == "") {
				return fmt.Errorf("请二选一提供内容：--content 或 --file")
			}
			text := content
			if file != "" {
				b, err := os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("读取文件失败：%w", err)
				}
				text = string(b)
			}
			if format != "md" && format != "text" {
				return fmt.Errorf("--format 只能是 md 或 text")
			}
			// document 类型要求 content 为 {"type","content"} 结构
			payload, err := json.Marshal(map[string]string{"type": format, "content": text})
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			var out any
			params := map[string]any{"id": id, "content": string(payload)}
			if err := c.NewGetRequest("/api/file/content/save", params, &out); err != nil {
				return err
			}
			cli.OK("✓ 已保存文件 #%d", id)
			return cli.Output(out, []string{"id", "fid", "size", "userid", "created_at"})
		},
	}
	cmd.Flags().StringVar(&content, "content", "", "文本内容")
	cmd.Flags().StringVar(&file, "file", "", "从本地文件读取内容")
	cmd.Flags().StringVar(&format, "format", "md", "内容格式：md 或 text")
	return cmd
}

func newFileRenameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename <文件ID> <新名称>",
		Short: "重命名文件（夹）",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "文件ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			var out any
			params := map[string]any{"id": id, "name": args[1]}
			if err := c.NewGetRequest("/api/file/add", params, &out); err != nil {
				return err
			}
			cli.OK("✓ 已重命名文件 #%d", id)
			return cli.Output(out, []string{"id", "name"})
		},
	}
	return cmd
}

func newFileDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <文件ID> [文件ID...]",
		Short: "删除文件（夹）",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids := make([]int, 0, len(args))
			for _, a := range args {
				id, err := cli.ParseInt(a, "文件ID")
				if err != nil {
					return err
				}
				ids = append(ids, id)
			}
			if err := cli.Confirm(fmt.Sprintf("删除文件 %v", ids), fileTargets(ids)...); err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			if err := c.NewGetRequest("/api/file/remove", map[string]any{"ids": ids}, nil); err != nil {
				return err
			}
			cli.OK("✓ 已删除文件 %v", ids)
			return nil
		},
	}
	return cmd
}

func containsStr(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func fileTargets(ids []int) []cli.Target {
	out := make([]cli.Target, len(ids))
	for i, id := range ids {
		out[i] = cli.Target{Type: "file", ID: id}
	}
	return out
}
