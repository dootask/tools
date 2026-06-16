package commands

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

// 指令目录（与 AppStore /internal/cli/:appId/spec 返回结构对齐）。
type cliParam struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Default     any    `json:"default,omitempty"`
	Description string `json:"description,omitempty"`
}

type cliCommand struct {
	Name        string     `json:"name"`
	Summary     string     `json:"summary,omitempty"`
	Description string     `json:"description,omitempty"`
	Method      string     `json:"method"`
	Path        string     `json:"path"`
	HasBody     bool       `json:"has_body"`
	Params      []cliParam `json:"params"`
	Tags        []string   `json:"tags,omitempty"`
}

type cliCatalog struct {
	AppID    string       `json:"app_id"`
	Version  string       `json:"version"`
	Source   string       `json:"source"`
	Title    string       `json:"title,omitempty"`
	Commands []cliCommand `json:"commands"`
}

func fetchCliCatalog(appid string, refresh bool) (*cliCatalog, error) {
	var q map[string]string
	if refresh {
		q = map[string]string{"refresh": "1"}
	}
	var cat cliCatalog
	if err := cli.AppStoreRequest("GET", "/internal/cli/"+url.PathEscape(appid)+"/spec", q, nil, &cat); err != nil {
		return nil, err
	}
	return &cat, nil
}

// 调用请求体（与 AppStore models.CliCallRequest 对齐）。
type cliRawRequest struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   any    `json:"body,omitempty"`
}

type cliCallRequest struct {
	Command string            `json:"command,omitempty"`
	Args    map[string]string `json:"args,omitempty"`
	Data    json.RawMessage   `json:"data,omitempty"`
	Raw     *cliRawRequest    `json:"raw,omitempty"`
}

func newAppCallCmd() *cobra.Command {
	var argPairs []string
	var data string
	var raw, describe, refresh, list bool

	cmd := &cobra.Command{
		Use:   "call <应用ID> [指令] [K=V ...]",
		Short: "调用应用自定义指令（基于应用的 OpenAPI/Swagger）",
		Long: "根据应用暴露的 OpenAPI/Swagger 规范调用其后端接口。\n\n" +
			"  doo app call <应用ID>                 列出应用可用指令\n" +
			"  doo app call <应用ID> <指令> --describe  查看某指令的参数\n" +
			"  doo app call <应用ID> <指令> [K=V ...]   调用指令（也可用 --arg K=V / --data JSON）\n" +
			"  doo app call <应用ID> --raw GET /api/x   直接调用未声明的端点",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appid := args[0]
			rest := args[1:]

			// 解析 body：--data 优先，否则 K=V（含位置参数与 --arg）。
			argMap, err := collectArgs(argPairs, rest, raw)
			if err != nil {
				return err
			}
			var dataRaw json.RawMessage
			if data != "" {
				if !json.Valid([]byte(data)) {
					return fmt.Errorf("--data 不是合法 JSON")
				}
				dataRaw = json.RawMessage(data)
			}

			// raw 模式：call <id> --raw <METHOD> <path>
			if raw {
				if len(rest) < 2 {
					return fmt.Errorf("--raw 需要 <METHOD> <path>，如 --raw GET /api/list")
				}
				req := cliCallRequest{
					Raw:  &cliRawRequest{Method: strings.ToUpper(rest[0]), Path: rest[1]},
					Data: dataRaw,
				}
				return doCliCall(appid, req)
			}

			// 无指令：列出指令目录。
			if len(rest) == 0 || list {
				cat, err := fetchCliCatalog(appid, refresh)
				if err != nil {
					return err
				}
				if cli.Opts.JSON {
					return cli.Output(cat, nil)
				}
				if len(cat.Commands) == 0 {
					fmt.Println("（该应用未暴露任何指令）")
					return nil
				}
				return cli.Output(cat.Commands, []string{"name", "method", "path", "summary"})
			}

			command := rest[0]

			// describe：查看某指令参数。
			if describe {
				cat, err := fetchCliCatalog(appid, refresh)
				if err != nil {
					return err
				}
				var found *cliCommand
				for i := range cat.Commands {
					if strings.EqualFold(cat.Commands[i].Name, command) {
						found = &cat.Commands[i]
						break
					}
				}
				if found == nil {
					return fmt.Errorf("应用没有指令 %q（用 doo app call %s 查看可用指令）", command, appid)
				}
				if cli.Opts.JSON {
					return cli.Output(found, nil)
				}
				fmt.Printf("%s  %s %s\n", found.Name, found.Method, found.Path)
				if found.Summary != "" {
					fmt.Printf("  %s\n", found.Summary)
				}
				if len(found.Params) == 0 {
					fmt.Println("（无参数）")
					return nil
				}
				return cli.Output(found.Params, []string{"name", "in", "type", "required", "description"})
			}

			// 执行调用。
			req := cliCallRequest{Command: command, Args: argMap, Data: dataRaw}
			return doCliCall(appid, req)
		},
	}

	f := cmd.Flags()
	f.StringArrayVar(&argPairs, "arg", nil, "指令参数 K=V（可重复）")
	f.StringVar(&data, "data", "", "请求体 JSON（覆盖 K=V body 参数）")
	f.BoolVar(&raw, "raw", false, "直接调用端点：--raw <METHOD> <path>")
	f.BoolVar(&describe, "describe", false, "查看指令的参数说明")
	f.BoolVar(&refresh, "refresh", false, "重新解析应用接口规范（跳过缓存）")
	f.BoolVar(&list, "list", false, "列出应用可用指令")
	return cmd
}

func doCliCall(appid string, req cliCallRequest) error {
	var out any
	if err := cli.AppStoreRequest("POST", "/internal/cli/"+url.PathEscape(appid)+"/call", nil, req, &out); err != nil {
		return err
	}
	return cli.Output(out, nil)
}

// collectArgs 合并 --arg K=V 与位置 K=V。raw 模式下位置参数是 METHOD/path，不当作 K=V。
func collectArgs(argPairs, rest []string, raw bool) (map[string]string, error) {
	m := map[string]string{}
	add := func(p string) error {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 || kv[0] == "" {
			return fmt.Errorf("非法参数 %q（需 K=V 形式）", p)
		}
		m[kv[0]] = kv[1]
		return nil
	}
	for _, p := range argPairs {
		if err := add(p); err != nil {
			return nil, err
		}
	}
	if !raw {
		// rest[0] 是指令名，rest[1:] 视为位置 K=V。
		for _, p := range rest[min(1, len(rest)):] {
			if err := add(p); err != nil {
				return nil, err
			}
		}
	}
	if len(m) == 0 {
		return nil, nil
	}
	return m, nil
}
