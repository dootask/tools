package commands

import (
	"strings"
	"sync"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

var searchAllTypes = []string{"task", "project", "file", "contact", "message"}

// search 跨域统一搜索：对每个 type 并行打 /api/search/<type>（兜底端点，实验性）。
func newSearchCmd() *cobra.Command {
	var types string
	var take int
	cmd := &cobra.Command{
		Use:   "search <关键词>",
		Short: "跨任务/项目/文件/联系人/消息统一搜索（实验性）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.Opts.Client()
			if err != nil {
				return err
			}
			selected := searchAllTypes
			if strings.TrimSpace(types) != "" {
				selected = nil
				for _, t := range strings.Split(types, ",") {
					if t = strings.TrimSpace(t); t != "" {
						selected = append(selected, t)
					}
				}
			}

			results := make(map[string]any, len(selected))
			var mu sync.Mutex
			var wg sync.WaitGroup
			for _, t := range selected {
				wg.Add(1)
				go func(t string) {
					defer wg.Done()
					params := map[string]any{"key": args[0]}
					if take > 0 {
						params["take"] = take
					}
					var out any
					if err := c.NewGetRequest("/api/search/"+t, params, &out); err != nil {
						out = map[string]any{"error": err.Error()}
					}
					mu.Lock()
					results[t] = out
					mu.Unlock()
				}(t)
			}
			wg.Wait()
			return cli.Output(results, nil)
		},
	}
	cmd.Flags().StringVar(&types, "types", "", "限定类型，逗号分隔：task,project,file,contact,message")
	cmd.Flags().IntVar(&take, "take", 0, "每类结果数量")
	return cmd
}
