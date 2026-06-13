package commands

import (
	"fmt"
	"os"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

// report 全部走兜底端点（SDK 暂无报告域），标记为实验性。
func newReportCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "report", Short: "工作报告（实验性）"}
	cmd.AddCommand(
		newReportReceivedCmd(),
		newReportMyCmd(),
		newReportViewCmd(),
		newReportTemplateCmd(),
		newReportSubmitCmd(),
		newReportMarkCmd(),
	)
	return cmd
}

var reportCols = []string{"id", "title", "type", "username", "created_at"}

func newReportReceivedCmd() *cobra.Command {
	var typ, status, search string
	cmd := &cobra.Command{
		Use:   "received",
		Short: "收到的报告",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{}
			if typ != "" {
				params["keys[type]"] = typ
			}
			if status != "" {
				params["keys[status]"] = status
			}
			if search != "" {
				params["keys[key]"] = search
			}
			var out any
			if err := c.NewGetRequest("/api/report/receive", params, &out); err != nil {
				return err
			}
			return cli.Output(out, reportCols)
		},
	}
	f := cmd.Flags()
	f.StringVar(&typ, "type", "", "类型 weekly|daily|all")
	f.StringVar(&status, "status", "", "状态 read|unread|all")
	f.StringVar(&search, "search", "", "关键词")
	return cmd
}

func newReportMyCmd() *cobra.Command {
	var typ, search string
	cmd := &cobra.Command{
		Use:   "my",
		Short: "我发出的报告",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{}
			if typ != "" {
				params["keys[type]"] = typ
			}
			if search != "" {
				params["keys[key]"] = search
			}
			var out any
			if err := c.NewGetRequest("/api/report/my", params, &out); err != nil {
				return err
			}
			return cli.Output(out, reportCols)
		},
	}
	cmd.Flags().StringVar(&typ, "type", "", "类型 weekly|daily|all")
	cmd.Flags().StringVar(&search, "search", "", "关键词")
	return cmd
}

func newReportViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view <报告ID>",
		Short: "查看报告详情",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "报告ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			var out any
			if err := c.NewGetRequest("/api/report/detail", map[string]any{"id": id}, &out); err != nil {
				return err
			}
			return cli.Output(out, nil)
		},
	}
}

func newReportTemplateCmd() *cobra.Command {
	var typ string
	var offset int
	cmd := &cobra.Command{
		Use:   "template",
		Short: "生成报告模板",
		RunE: func(cmd *cobra.Command, args []string) error {
			if typ == "" {
				return fmt.Errorf("--type 必填 weekly|daily")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			var out any
			if err := c.NewGetRequest("/api/report/template", map[string]any{"type": typ, "offset": offset}, &out); err != nil {
				return err
			}
			return cli.Output(out, nil)
		},
	}
	cmd.Flags().StringVar(&typ, "type", "", "类型 weekly|daily（必填）")
	cmd.Flags().IntVar(&offset, "offset", 0, "周期偏移（0 当前）")
	return cmd
}

func newReportSubmitCmd() *cobra.Command {
	var typ, title, content, contentFile, sign string
	var receive string
	var offset int
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "提交工作报告",
		RunE: func(cmd *cobra.Command, args []string) error {
			if typ == "" || title == "" {
				return fmt.Errorf("--type 与 --title 必填")
			}
			if contentFile != "" {
				b, err := os.ReadFile(contentFile)
				if err != nil {
					return err
				}
				content = string(b)
			}
			if content == "" {
				return fmt.Errorf("--content 或 --content-file 必填")
			}
			receiveIDs, err := cli.ParseIDList(receive)
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			body := map[string]any{"type": typ, "title": title, "content": content, "offset": offset}
			if len(receiveIDs) > 0 {
				body["receive"] = receiveIDs
			}
			if sign != "" {
				body["sign"] = sign
			}
			var out any
			if err := c.NewPostRequest("/api/report/store", body, &out); err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(out, nil)
			}
			cli.OK("✓ 已提交报告：%s", title)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&typ, "type", "", "类型 weekly|daily（必填）")
	f.StringVar(&title, "title", "", "标题（必填）")
	f.StringVar(&content, "content", "", "内容")
	f.StringVar(&contentFile, "content-file", "", "从文件读取内容")
	f.StringVar(&receive, "receive", "", "接收人 ID 列表")
	f.StringVar(&sign, "sign", "", "模板签名")
	f.IntVar(&offset, "offset", 0, "周期偏移")
	return cmd
}

func newReportMarkCmd() *cobra.Command {
	var action string
	cmd := &cobra.Command{
		Use:   "mark <报告ID> [报告ID...]",
		Short: "标记报告已读/未读",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids := make([]int, 0, len(args))
			for _, a := range args {
				id, err := cli.ParseInt(a, "报告ID")
				if err != nil {
					return err
				}
				ids = append(ids, id)
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			var out any
			if err := c.NewGetRequest("/api/report/mark", map[string]any{"id": ids, "action": action}, &out); err != nil {
				return err
			}
			cli.OK("✓ 已标记 %d 份报告为 %s", len(ids), action)
			return nil
		},
	}
	cmd.Flags().StringVar(&action, "action", "read", "read|unread")
	return cmd
}
