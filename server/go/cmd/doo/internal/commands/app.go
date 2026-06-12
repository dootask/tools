package commands

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

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
	return &cobra.Command{
		Use:   "catalog",
		Short: "列出应用市场可安装的应用",
		RunE: func(cmd *cobra.Command, args []string) error {
			var out any
			if err := cli.AppStoreRequest("GET", "/list", nil, nil, &out); err != nil {
				return err
			}
			return cli.Output(out, []string{"id", "name", "version", "tags"})
		},
	}
}

// install / update 共用：对已安装应用再 install 即为升级（后端自动判定）。
func newAppInstallCmd(verb string) *cobra.Command {
	var version string
	var params []string
	var pull bool
	short := "安装应用（已安装则升级）"
	defVersion := "latest"
	if verb == "update" {
		short = "更新应用到指定/最新版本"
	}
	cmd := &cobra.Command{
		Use:   verb + " <应用ID>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]any{
				"appid":      args[0],
				"version":    version,
				"pull_image": pull,
			}
			if len(params) > 0 {
				pm := map[string]any{}
				for _, p := range params {
					kv := strings.SplitN(p, "=", 2)
					if len(kv) == 2 {
						pm[kv[0]] = kv[1]
					}
				}
				body["params"] = pm
			}
			if err := cli.AppStoreRequest("POST", "/internal/install", nil, body, nil); err != nil {
				return err
			}
			cli.OK("✓ 已触发%s：%s@%s", map[string]string{"install": "安装", "update": "更新"}[verb], args[0], version)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&version, "version", defVersion, "版本号（默认 latest）")
	f.StringArrayVar(&params, "param", nil, "应用参数 k=v（可重复）")
	f.BoolVar(&pull, "pull", false, "操作前先拉取镜像")
	return cmd
}

func newAppReinstallCmd() *cobra.Command {
	var pull bool
	cmd := &cobra.Command{
		Use:   "reinstall <应用ID>",
		Short: "重新安装应用（按当前已装版本重部署）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appid := args[0]
			var installed []map[string]any
			if err := cli.AppStoreRequest("GET", "/internal/installed", nil, nil, &installed); err != nil {
				return err
			}
			version := ""
			for _, a := range installed {
				if fmt.Sprintf("%v", a["id"]) == appid {
					version = fmt.Sprintf("%v", a["version"])
					break
				}
			}
			if version == "" {
				return fmt.Errorf("应用 %s 未安装，无法重装", appid)
			}
			body := map[string]any{"appid": appid, "version": version, "pull_image": pull}
			if err := cli.AppStoreRequest("POST", "/internal/install", nil, body, nil); err != nil {
				return err
			}
			cli.OK("✓ 已触发重装：%s@%s", appid, version)
			return nil
		},
	}
	cmd.Flags().BoolVar(&pull, "pull", false, "重装前先拉取镜像")
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
