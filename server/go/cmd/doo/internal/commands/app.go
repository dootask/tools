package commands

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

// 字段定义来自 GET /one/:appId 顶层 fields（已按当前语种摊平 label/description）。
type appField struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Type        string `json:"type"` // input/select/password/textarea/...
	Required    bool   `json:"required"`
	Default     any    `json:"default"`
	Description string `json:"description"`
	Options     []struct {
		Label string `json:"label"`
		Value any    `json:"value"`
	} `json:"options"`
}

type appOne struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Fields []appField `json:"fields"`
	Config struct {
		Params         map[string]any `json:"params"`
		Resources      map[string]any `json:"resources"`
		InstallVersion string         `json:"install_version"`
		Status         string         `json:"status"`
	} `json:"config"`
}

func fetchAppOne(appid string) (*appOne, error) {
	var o appOne
	if err := cli.AppStoreRequest("GET", "/one/"+url.PathEscape(appid), nil, nil, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// parseParamPairs 把 ["K=V","K2=V2"] 解析成 map；非法的 K=V 报错。
func parseParamPairs(pairs []string) (map[string]any, error) {
	m := map[string]any{}
	for _, p := range pairs {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 || kv[0] == "" {
			return nil, fmt.Errorf("非法 --param %q（需 K=V 形式）", p)
		}
		m[kv[0]] = kv[1]
	}
	return m, nil
}

// validateRequiredFields 检查 required 字段：若用户未提供，且字段无 default，则报错。
// 同时拒绝 fields 里不存在的多余 key，提示用户拼错。
func validateRequiredFields(fields []appField, given map[string]any) error {
	defined := map[string]appField{}
	for _, f := range fields {
		defined[f.Name] = f
	}
	for k := range given {
		if _, ok := defined[k]; !ok {
			return fmt.Errorf("未知参数 %q（应用 fields 中不存在，可用 doo app fields <appid> 查看）", k)
		}
	}
	var missing []string
	for _, f := range fields {
		if !f.Required {
			continue
		}
		if _, ok := given[f.Name]; ok {
			continue
		}
		if f.Default != nil && fmt.Sprintf("%v", f.Default) != "" {
			continue // 有默认值，后端会兜底
		}
		missing = append(missing, f.Name)
	}
	if len(missing) > 0 {
		return fmt.Errorf("缺少必填参数：%s（用 --param K=V 提供，详见 doo app fields <appid>）", strings.Join(missing, ", "))
	}
	return nil
}

// buildResources 把 --cpu-limit/--memory-limit 拼成 AppStore 期望的形态；
// 任一为空就不覆盖已装值（reinstall 场景下保留 sticky）。
func buildResources(cpu, mem string, fallback map[string]any) map[string]any {
	r := map[string]any{}
	for k, v := range fallback {
		r[k] = v
	}
	if cpu != "" {
		r["cpu_limit"] = cpu
	}
	if mem != "" {
		r["memory_limit"] = mem
	}
	return r
}

// 应用插件（AppStore）管理。安装/卸载/删除/更新列表需管理员；列表/日志/容器普通用户即可。
func newAppCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "应用插件（AppStore）",
		Long:  "管理 DooTask 应用插件：安装/更新/重装/卸载/删除、查看已装清单、应用日志与容器。安装类操作需管理员权限。",
	}
	cmd.AddCommand(
		newAppListCmd(),
		newAppCatalogCmd(),
		newAppFieldsCmd(),
		newAppInstallCmd("install"),
		newAppInstallCmd("update"),
		newAppReinstallCmd(),
		newAppUninstallCmd(),
		newAppRemoveCmd(),
		newAppLogsCmd(),
		newAppContainersCmd(),
		newAppContainerLogsCmd(),
		newAppRefreshCmd(),
	)
	return cmd
}

func newAppFieldsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fields <应用ID>",
		Short: "列出应用的安装参数定义（字段名/类型/是否必填/默认值/可选项）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			one, err := fetchAppOne(args[0])
			if err != nil {
				return err
			}
			return cli.Output(one.Fields, []string{"name", "type", "required", "default", "label", "description"})
		},
	}
}

func newAppListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "列出已安装应用",
		RunE: func(cmd *cobra.Command, args []string) error {
			var out any
			if err := cli.AppStoreRequest("GET", "/internal/installed", nil, nil, &out); err != nil {
				return err
			}
			return cli.Output(out, []string{"id", "version", "status", "install_at"})
		},
	}
}

func newAppCatalogCmd() *cobra.Command {
	var search string
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "列出应用市场可安装的应用（--search 关键词模糊匹配 id/name/description/tags）",
		RunE: func(cmd *cobra.Command, args []string) error {
			var raw []map[string]any
			if err := cli.AppStoreRequest("GET", "/list", nil, nil, &raw); err != nil {
				return err
			}
			items := raw
			if kw := strings.TrimSpace(search); kw != "" {
				kwLower := strings.ToLower(kw)
				matched := make([]map[string]any, 0, len(raw))
				for _, it := range raw {
					if catalogMatch(it, kw, kwLower) {
						matched = append(matched, it)
					}
				}
				items = matched
			}
			return cli.Output(items, []string{"id", "name", "version", "tags"})
		},
	}
	cmd.Flags().StringVar(&search, "search", "", "按关键词模糊匹配 id/name/description/tags（含中文）")
	return cmd
}

// catalogMatch 按关键词在 id/name/description/tags 中做大小写不敏感的子串匹配；
// tags 任一项命中即视为匹配，便于覆盖中文标签如「客户管理」。
func catalogMatch(item map[string]any, kw, kwLower string) bool {
	for _, k := range []string{"id", "name", "description"} {
		if v, ok := item[k].(string); ok && v != "" {
			if strings.Contains(strings.ToLower(v), kwLower) || strings.Contains(v, kw) {
				return true
			}
		}
	}
	if tags, ok := item["tags"].([]any); ok {
		for _, t := range tags {
			if s, ok := t.(string); ok {
				if strings.Contains(strings.ToLower(s), kwLower) || strings.Contains(s, kw) {
					return true
				}
			}
		}
	}
	return false
}

// install / update 共用：对已安装应用再 install 即为升级（后端自动判定）。
// 校验 fields（拒未知 key、缺必填字段在没默认值时报错），暴露 cpu/memory/pull。
func newAppInstallCmd(verb string) *cobra.Command {
	var version, cpuLimit, memLimit string
	var params []string
	var pull bool
	short := "安装应用（已安装则升级）"
	if verb == "update" {
		short = "更新应用到指定/最新版本"
	}
	cmd := &cobra.Command{
		Use:   verb + " <应用ID>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appid := args[0]
			one, err := fetchAppOne(appid)
			if err != nil {
				return err
			}
			given, err := parseParamPairs(params)
			if err != nil {
				return err
			}
			// 已安装则 sticky：以当前已装 params 为底叠加用户传入，避免不传 --param 时
			// 误把令牌/选项清空（与网页表单"初值即当前值"的行为对齐）。
			merged := map[string]any{}
			if one.Config.Status == "installed" {
				for k, v := range one.Config.Params {
					merged[k] = v
				}
			}
			for k, v := range given {
				merged[k] = v
			}
			if err := validateRequiredFields(one.Fields, merged); err != nil {
				return err
			}
			body := map[string]any{
				"appid":      appid,
				"version":    version,
				"pull_image": pull,
			}
			if len(merged) > 0 {
				body["params"] = merged
			}
			if cpuLimit != "" || memLimit != "" {
				body["resources"] = buildResources(cpuLimit, memLimit, one.Config.Resources)
			} else if one.Config.Status == "installed" && len(one.Config.Resources) > 0 {
				body["resources"] = one.Config.Resources
			}
			if err := cli.AppStoreRequest("POST", "/internal/install", nil, body, nil); err != nil {
				return err
			}
			cli.OK("✓ 已触发%s：%s@%s", map[string]string{"install": "安装", "update": "更新"}[verb], appid, version)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&version, "version", "latest", "版本号（默认 latest）")
	f.StringArrayVar(&params, "param", nil, "应用参数 K=V（可重复，详见 doo app fields <appid>；已安装应用未传则沿用当前值）")
	f.StringVar(&cpuLimit, "cpu-limit", "", "CPU 限额（如 1.0；空则沿用当前或应用默认）")
	f.StringVar(&memLimit, "memory-limit", "", "内存限额（如 512M / 2G；空则沿用当前或应用默认）")
	f.BoolVar(&pull, "pull", false, "操作前先拉取镜像")
	return cmd
}

// reinstall：按当前已装版本重部署；sticky 复用当前 params/resources，允许 --param/--cpu-limit 覆盖。
func newAppReinstallCmd() *cobra.Command {
	var pull bool
	var cpuLimit, memLimit string
	var params []string
	cmd := &cobra.Command{
		Use:   "reinstall <应用ID>",
		Short: "重新安装应用（按当前已装版本与参数重部署，可用 --param 覆盖）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appid := args[0]
			one, err := fetchAppOne(appid)
			if err != nil {
				return err
			}
			version := one.Config.InstallVersion
			if version == "" || one.Config.Status != "installed" {
				return fmt.Errorf("应用 %s 未安装，无法重装", appid)
			}
			given, err := parseParamPairs(params)
			if err != nil {
				return err
			}
			// sticky：以当前已装 params 为底，再叠加用户传入的覆盖
			merged := map[string]any{}
			for k, v := range one.Config.Params {
				merged[k] = v
			}
			for k, v := range given {
				merged[k] = v
			}
			if err := validateRequiredFields(one.Fields, merged); err != nil {
				return err
			}
			body := map[string]any{
				"appid":      appid,
				"version":    version,
				"pull_image": pull,
				"params":     merged,
				"resources":  buildResources(cpuLimit, memLimit, one.Config.Resources),
			}
			if err := cli.AppStoreRequest("POST", "/internal/install", nil, body, nil); err != nil {
				return err
			}
			cli.OK("✓ 已触发重装：%s@%s", appid, version)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringArrayVar(&params, "param", nil, "覆盖参数 K=V（可重复；未传则沿用当前已装值）")
	f.StringVar(&cpuLimit, "cpu-limit", "", "CPU 限额（空则沿用当前已装值）")
	f.StringVar(&memLimit, "memory-limit", "", "内存限额（空则沿用当前已装值）")
	f.BoolVar(&pull, "pull", false, "重装前先拉取镜像")
	return cmd
}

func newAppUninstallCmd() *cobra.Command {
	var deleteData bool
	cmd := &cobra.Command{
		Use:   "uninstall <应用ID>",
		Short: "卸载应用（--delete-data 同时删除数据）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg := fmt.Sprintf("确认卸载应用 %s?", args[0])
			if deleteData {
				msg = fmt.Sprintf("确认卸载应用 %s 并删除其数据（不可恢复）?", args[0])
			}
			if err := cli.Confirm(msg); err != nil {
				return err
			}
			q := map[string]string{}
			if deleteData {
				q["delete_data"] = "true"
			}
			if err := cli.AppStoreRequest("GET", "/internal/uninstall/"+url.PathEscape(args[0]), q, nil, nil); err != nil {
				return err
			}
			cli.OK("✓ 已触发卸载：%s", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVar(&deleteData, "delete-data", false, "卸载时同时删除应用数据")
	return cmd
}

func newAppRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <应用ID>",
		Short: "删除社区应用（需先卸载，仅 community_ 应用）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cli.Confirm(fmt.Sprintf("确认删除应用 %s（不可恢复）?", args[0])); err != nil {
				return err
			}
			if err := cli.AppStoreRequest("GET", "/internal/remove/"+url.PathEscape(args[0]), nil, nil, nil); err != nil {
				return err
			}
			cli.OK("✓ 已删除：%s", args[0])
			return nil
		},
	}
}

func newAppLogsCmd() *cobra.Command {
	var lines int
	cmd := &cobra.Command{
		Use:   "logs <应用ID>",
		Short: "查看应用安装/运行日志",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := map[string]string{}
			if lines > 0 {
				q["n"] = fmt.Sprintf("%d", lines)
			}
			var out any
			if err := cli.AppStoreRequest("GET", "/internal/log/"+url.PathEscape(args[0]), q, nil, &out); err != nil {
				return err
			}
			return cli.Output(out, nil)
		},
	}
	cmd.Flags().IntVarP(&lines, "lines", "n", 200, "返回的日志行数")
	return cmd
}

func newAppContainersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "containers <应用ID>",
		Short: "列出应用的容器/服务",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var out any
			if err := cli.AppStoreRequest("GET", "/internal/containers/"+url.PathEscape(args[0]), nil, nil, &out); err != nil {
				return err
			}
			return cli.Output(out, nil)
		},
	}
}

func newAppContainerLogsCmd() *cobra.Command {
	var service string
	var lines int
	cmd := &cobra.Command{
		Use:   "container-logs <应用ID> --service <服务名>",
		Short: "查看应用某个容器的日志",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if service == "" {
				return fmt.Errorf("--service 必填（可先用 app containers 查看服务名）")
			}
			q := map[string]string{"service": service}
			if lines > 0 {
				q["n"] = fmt.Sprintf("%d", lines)
			}
			var out any
			if err := cli.AppStoreRequest("GET", "/internal/containers/"+url.PathEscape(args[0])+"/logs", q, nil, &out); err != nil {
				return err
			}
			return cli.Output(out, nil)
		},
	}
	f := cmd.Flags()
	f.StringVar(&service, "service", "", "容器/服务名（必填）")
	f.IntVarP(&lines, "lines", "n", 200, "返回的日志行数")
	return cmd
}

func newAppRefreshCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "refresh",
		Short: "刷新应用市场可安装列表（拉取远程源）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cli.AppStoreRequest("GET", "/internal/apps/update", nil, nil, nil); err != nil {
				return err
			}
			cli.OK("✓ 已触发刷新应用列表")
			return nil
		},
	}
}
