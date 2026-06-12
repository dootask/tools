package commands

import (
	"fmt"
	"os"
	"time"

	dootask "github.com/dootask/tools/server/go"
	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

var taskListCols = []string{"id", "name", "complete_at", "end_at", "project_id", "column_name"}

func newTaskCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "task", Short: "任务"}
	cmd.AddCommand(
		newTaskListCmd(),
		newTaskViewCmd(),
		newTaskFilesCmd(),
		newTaskCreateCmd(),
		newTaskSubtaskCmd(),
		newTaskUpdateCmd(),
		newTaskDoneCmd(false),
		newTaskDoneCmd(true),
		newTaskDialogCmd(),
		newTaskNotifyCmd(),
		newTaskArchiveCmd(),
		newTaskDeleteCmd(),
	)
	return cmd
}

func newTaskListCmd() *cobra.Command {
	var project, parent, page, pageSize int
	var status, search, tag, archived, withExtend string
	var allProject bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出任务",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{}
			if project > 0 {
				params["project_id"] = project
			}
			if parent != 0 {
				params["parent_id"] = parent
			}
			if status != "" {
				params["keys[status]"] = status
			}
			if search != "" {
				params["keys[name]"] = search
			}
			if tag != "" {
				params["keys[tag]"] = tag
			}
			if archived != "" {
				params["archived"] = archived
			}
			if withExtend != "" {
				params["with_extend"] = withExtend
			}
			if allProject {
				params["scope"] = "all_project"
			}
			if page > 0 {
				params["page"] = page
			}
			if pageSize > 0 {
				params["pagesize"] = pageSize
			}
			var out any
			if err := c.NewGetRequest("/api/project/task/lists", params, &out); err != nil {
				return err
			}
			return cli.Output(out, taskListCols)
		},
	}
	f := cmd.Flags()
	f.IntVar(&project, "project", 0, "项目 ID")
	f.IntVar(&parent, "parent", 0, "父任务 ID（>0 取子任务，-1 仅主任务）")
	f.StringVar(&status, "status", "", "状态过滤 completed|uncompleted|flow-<x>")
	f.StringVar(&search, "search", "", "按名称/描述搜索")
	f.StringVar(&tag, "tag", "", "按标签过滤")
	f.StringVar(&archived, "archived", "", "归档过滤 all|yes|no")
	f.StringVar(&withExtend, "with", "", "附带扩展字段，如 project_name,column_name")
	f.BoolVar(&allProject, "all-project", false, "跨全部项目")
	f.IntVar(&page, "page", 0, "页码")
	f.IntVar(&pageSize, "page-size", 0, "每页数量")
	return cmd
}

func newTaskViewCmd() *cobra.Command {
	var content, files bool
	cmd := &cobra.Command{
		Use:   "view <任务ID>",
		Short: "查看任务详情",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "任务ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			t, err := c.GetTask(dootask.GetTaskRequest{TaskID: id})
			if err != nil {
				return err
			}
			if err := cli.Output(t, nil); err != nil {
				return err
			}
			if content {
				tc, err := c.GetTaskContent(dootask.GetTaskContentRequest{TaskID: id})
				if err != nil {
					return err
				}
				if !cli.Opts.JSON {
					fmt.Println("\n--- 内容 ---")
				}
				if err := cli.Output(tc, nil); err != nil {
					return err
				}
			}
			if files {
				ff, err := c.GetTaskFiles(dootask.GetTaskFilesRequest{TaskID: id})
				if err != nil {
					return err
				}
				if !cli.Opts.JSON {
					fmt.Println("\n--- 附件 ---")
				}
				if err := cli.Output(ff, []string{"id", "name", "ext", "size"}); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&content, "content", false, "附带任务正文")
	cmd.Flags().BoolVar(&files, "files", false, "附带附件列表")
	return cmd
}

func newTaskFilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "files <任务ID>",
		Short: "查看任务附件",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "任务ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			ff, err := c.GetTaskFiles(dootask.GetTaskFilesRequest{TaskID: id})
			if err != nil {
				return err
			}
			return cli.Output(ff, []string{"id", "name", "ext", "size"})
		},
	}
}

func newTaskCreateCmd() *cobra.Command {
	var project int
	var name, content, contentFile, column, start, end, owner string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建任务",
		RunE: func(cmd *cobra.Command, args []string) error {
			if project <= 0 || name == "" {
				return fmt.Errorf("--project 与 --name 必填")
			}
			if contentFile != "" {
				b, err := os.ReadFile(contentFile)
				if err != nil {
					return err
				}
				content = string(b)
			}
			owners, err := cli.ParseIDList(owner)
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			req := dootask.CreateTaskRequest{
				ProjectID: project,
				Name:      name,
				Content:   content,
				Times:     cli.BuildTimes(start, end),
				Owner:     owners,
			}
			if column != "" {
				if n, err := cli.ParseInt(column, "列"); err == nil {
					req.ColumnID = n
				} else {
					req.ColumnID = column
				}
			}
			t, err := c.CreateTask(req)
			if err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(t, nil)
			}
			cli.OK("✓ 已创建任务 #%d：%s", t.ID, t.Name)
			return nil
		},
	}
	f := cmd.Flags()
	f.IntVar(&project, "project", 0, "项目 ID（必填）")
	f.StringVar(&name, "name", "", "任务名称（必填）")
	f.StringVar(&content, "content", "", "任务内容")
	f.StringVar(&contentFile, "content-file", "", "从文件读取任务内容")
	f.StringVar(&column, "column", "", "看板列 ID 或列名")
	f.StringVar(&start, "start", "", "开始时间 YYYY-MM-DD HH:MM:SS")
	f.StringVar(&end, "end", "", "截止时间 YYYY-MM-DD HH:MM:SS")
	f.StringVar(&owner, "owner", "", "负责人 ID 列表，逗号分隔")
	return cmd
}

func newTaskSubtaskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "subtask <父任务ID> <名称>",
		Short: "创建子任务",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "父任务ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			t, err := c.CreateSubTask(dootask.CreateSubTaskRequest{TaskID: id, Name: args[1]})
			if err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(t, nil)
			}
			cli.OK("✓ 已创建子任务 #%d：%s", t.ID, t.Name)
			return nil
		},
	}
}

func newTaskUpdateCmd() *cobra.Command {
	var name, content, start, end, owner, assist, color string
	cmd := &cobra.Command{
		Use:   "update <任务ID>",
		Short: "更新任务（仅提交指定字段）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "任务ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			// 仅提交用户实际修改的字段：SDK 结构体无 omitempty，整体提交会误清空字段。
			f := cmd.Flags()
			params := map[string]any{"task_id": id}
			if f.Changed("name") {
				params["name"] = name
			}
			if f.Changed("content") {
				params["content"] = content
			}
			if f.Changed("start") || f.Changed("end") {
				params["times"] = cli.BuildTimes(start, end)
			}
			if f.Changed("owner") {
				ids, err := cli.ParseIDList(owner)
				if err != nil {
					return err
				}
				params["owner"] = ids
			}
			if f.Changed("assist") {
				ids, err := cli.ParseIDList(assist)
				if err != nil {
					return err
				}
				params["assist"] = ids
			}
			if f.Changed("color") {
				params["color"] = color
			}
			if len(params) == 1 {
				return fmt.Errorf("没有要更新的字段")
			}
			var out map[string]any
			if err := c.NewPostRequest("/api/project/task/update", params, &out); err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(out, nil)
			}
			cli.OK("✓ 已更新任务 #%d", id)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&name, "name", "", "任务名称")
	f.StringVar(&content, "content", "", "任务内容")
	f.StringVar(&start, "start", "", "开始时间")
	f.StringVar(&end, "end", "", "截止时间")
	f.StringVar(&owner, "owner", "", "负责人 ID 列表")
	f.StringVar(&assist, "assist", "", "协助者 ID 列表")
	f.StringVar(&color, "color", "", "颜色")
	return cmd
}

// newTaskDoneCmd 复用一个构造器实现 done / undone。
func newTaskDoneCmd(undone bool) *cobra.Command {
	use, short := "done <任务ID>", "标记任务完成"
	if undone {
		use, short = "undone <任务ID>", "取消任务完成（实验性）"
	}
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "任务ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			// 服务端要求 complete_at 为日期字符串才标记完成（实际取 now()）；
			// 传 false（非日期）则标记未完成。仅提交 task_id+complete_at，避免清空其它字段。
			var completeAt any = false
			if !undone {
				completeAt = time.Now().Format("2006-01-02 15:04:05")
			}
			params := map[string]any{"task_id": id, "complete_at": completeAt}
			if err := c.NewPostRequest("/api/project/task/update", params, nil); err != nil {
				return err
			}
			if undone {
				cli.OK("✓ 已取消完成 #%d", id)
			} else {
				cli.OK("✓ 已完成 #%d", id)
			}
			return nil
		},
	}
}

func newTaskDialogCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dialog <任务ID>",
		Short: "打开/创建任务对话",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "任务ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			d, err := c.CreateTaskDialog(dootask.CreateTaskDialogRequest{TaskID: id})
			if err != nil {
				return err
			}
			return cli.Output(d, nil)
		},
	}
}

func newTaskNotifyCmd() *cobra.Command {
	var text, textFile, nickname string
	var silence bool
	cmd := &cobra.Command{
		Use:   "notify <任务ID>",
		Short: "向任务对话发送一条 AI 助手消息",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "任务ID")
			if err != nil {
				return err
			}
			if textFile != "" {
				b, err := os.ReadFile(textFile)
				if err != nil {
					return err
				}
				text = string(b)
			}
			if text == "" {
				return fmt.Errorf("--text 或 --text-file 必填")
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			params := map[string]any{
				"task_id":   id,
				"text":      text,
				"text_type": "md",
			}
			if nickname != "" {
				params["nickname"] = nickname
			}
			if silence {
				params["silence"] = "yes"
			} else {
				params["silence"] = "no"
			}
			var out map[string]any
			if err := c.NewPostRequest("/api/dialog/msg/send_ai_assistant", params, &out); err != nil {
				return err
			}
			if cli.Opts.JSON {
				return cli.Output(out, nil)
			}
			cli.OK("✓ 已向任务 #%d 发送 AI 助手消息", id)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&text, "text", "", "消息内容（Markdown）")
	f.StringVar(&textFile, "text-file", "", "从文件读取消息内容")
	f.StringVar(&nickname, "nickname", "", "AI 助手昵称")
	f.BoolVar(&silence, "silence", false, "静默发送（不触发通知）")
	return cmd
}

func newTaskArchiveCmd() *cobra.Command {
	var recover bool
	cmd := &cobra.Command{
		Use:   "archive <任务ID>",
		Short: "归档任务（--recover 恢复）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "任务ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			action := "add"
			if recover {
				action = "recovery"
			}
			if err := c.ArchiveTask(id, action); err != nil {
				return err
			}
			cli.OK("✓ 已%s任务 #%d", map[bool]string{true: "恢复", false: "归档"}[recover], id)
			return nil
		},
	}
	cmd.Flags().BoolVar(&recover, "recover", false, "恢复已归档任务")
	return cmd
}

func newTaskDeleteCmd() *cobra.Command {
	var recover bool
	cmd := &cobra.Command{
		Use:   "delete <任务ID>",
		Short: "删除任务（--recover 恢复）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cli.ParseInt(args[0], "任务ID")
			if err != nil {
				return err
			}
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			action := "delete"
			if recover {
				action = "recovery"
			} else if err := cli.Confirm(fmt.Sprintf("确认删除任务 #%d?", id)); err != nil {
				return err
			}
			if err := c.DeleteTask(id, action); err != nil {
				return err
			}
			cli.OK("✓ 已%s任务 #%d", map[bool]string{true: "恢复", false: "删除"}[recover], id)
			return nil
		},
	}
	cmd.Flags().BoolVar(&recover, "recover", false, "恢复已删除任务")
	return cmd
}
