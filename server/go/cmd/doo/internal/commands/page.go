package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

// 页面操作经主程序常驻 WebSocket（/ws）派发：
// POST assistant/operation/dispatch 拿 requestId → 轮询 GET assistant/operation/result。
// 对调用者表现为一条同步命令。
const (
	pagePollInterval = 200 * time.Millisecond
	pagePollTimeout  = 30 * time.Second
)

func newPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "page",
		Short: "页面操作（驱动当前用户浏览器，依赖 AI 助手会话）",
		Long:  "通过主程序 WebSocket 向用户浏览器派发页面操作。需 --session <fd>（或环境变量 DOO_SESSION）指定目标会话。",
	}
	cmd.PersistentFlags().Int("session", 0, "目标会话 fd（默认环境变量 DOO_SESSION）")
	cmd.AddCommand(
		newPageContextCmd(),
		newPageActionCmd(),
		newPageElementCmd(),
	)
	return cmd
}

// resolveSession 取 --session > env DOO_SESSION。
func resolveSession(cmd *cobra.Command) (int, error) {
	fd, _ := cmd.Flags().GetInt("session")
	if fd <= 0 {
		if env := os.Getenv("DOO_SESSION"); env != "" {
			if n, err := cli.ParseInt(env, "DOO_SESSION"); err == nil {
				fd = n
			}
		}
	}
	if fd <= 0 {
		return 0, fmt.Errorf("缺少会话：请用 --session <fd> 或设置环境变量 DOO_SESSION")
	}
	return fd, nil
}

// dispatchPageOp 派发一次操作并轮询结果，返回 result 字段。
func dispatchPageOp(cmd *cobra.Command, action string, payload map[string]any) (any, error) {
	fd, err := resolveSession(cmd)
	if err != nil {
		return nil, err
	}
	c, err := cli.Opts.Client()
	if err != nil {
		return nil, err
	}

	var disp struct {
		RequestID string `json:"requestId"`
	}
	dispParams := map[string]any{"fd": fd, "action": action}
	if payload != nil {
		dispParams["payload"] = payload
	}
	if err := c.NewPostRequest("/api/assistant/operation/dispatch", dispParams, &disp); err != nil {
		return nil, err
	}
	if disp.RequestID == "" {
		return nil, fmt.Errorf("派发失败：未返回 requestId")
	}

	deadline := time.Now().Add(pagePollTimeout)
	for {
		var res struct {
			Status  string `json:"status"`
			Success bool   `json:"success"`
			Result  any    `json:"result"`
			Error   string `json:"error"`
		}
		err := c.NewGetRequest("/api/assistant/operation/result", map[string]any{"request_id": disp.RequestID}, &res)
		if err != nil {
			return nil, err
		}
		if res.Status == "ready" {
			if !res.Success {
				if res.Error != "" {
					return nil, fmt.Errorf("页面操作失败：%s", res.Error)
				}
				return nil, fmt.Errorf("页面操作失败")
			}
			return res.Result, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("页面操作超时（%s）：浏览器未在限定时间内回包", pagePollTimeout)
		}
		time.Sleep(pagePollInterval)
	}
}

func newPageContextCmd() *cobra.Command {
	var query, container string
	var maxElements, offset int
	var interactiveOnly, noElements bool
	cmd := &cobra.Command{
		Use:   "context",
		Short: "获取当前页面上下文与可交互元素",
		RunE: func(cmd *cobra.Command, args []string) error {
			payload := map[string]any{}
			if query != "" {
				payload["query"] = query
			}
			if container != "" {
				payload["container"] = container
			}
			if maxElements > 0 {
				payload["max_elements"] = maxElements
			}
			if offset > 0 {
				payload["offset"] = offset
			}
			if interactiveOnly {
				payload["interactive_only"] = true
			}
			if noElements {
				payload["include_elements"] = false
			}
			result, err := dispatchPageOp(cmd, "get_page_context", payload)
			if err != nil {
				return err
			}
			return cli.Output(result, nil)
		},
	}
	f := cmd.Flags()
	f.StringVar(&query, "query", "", "按语义查找相关元素")
	f.StringVar(&container, "container", "", "限定容器")
	f.IntVar(&maxElements, "max-elements", 0, "返回元素上限")
	f.IntVar(&offset, "offset", 0, "元素分页偏移")
	f.BoolVar(&interactiveOnly, "interactive-only", false, "仅可交互元素")
	f.BoolVar(&noElements, "no-elements", false, "不返回元素列表")
	return cmd
}

func newPageActionCmd() *cobra.Command {
	var paramsJSON string
	cmd := &cobra.Command{
		Use:   "action <名称>",
		Short: "执行一个业务页面操作",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload := map[string]any{"name": args[0]}
			if paramsJSON != "" {
				var params map[string]any
				if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
					return fmt.Errorf("--params 不是合法 JSON：%w", err)
				}
				payload["params"] = params
			}
			result, err := dispatchPageOp(cmd, "execute_action", payload)
			if err != nil {
				return err
			}
			return cli.Output(result, nil)
		},
	}
	cmd.Flags().StringVar(&paramsJSON, "params", "", "操作参数（JSON 对象）")
	return cmd
}

func newPageElementCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "element <元素UID> <动作> [值]",
		Short: "对页面元素执行动作（如 click、fill）",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload := map[string]any{
				"element_uid": args[0],
				"action":      args[1],
			}
			if len(args) == 3 {
				payload["value"] = args[2]
			}
			result, err := dispatchPageOp(cmd, "execute_element_action", payload)
			if err != nil {
				return err
			}
			return cli.Output(result, nil)
		},
	}
	return cmd
}
