package commands

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/spf13/cobra"
)

func leafPaths() map[string]bool {
	out := map[string]bool{}
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		for _, s := range c.Commands() {
			if s.HasSubCommands() {
				walk(s)
			} else {
				out[cmdPath(s)] = true
			}
		}
	}
	walk(NewRootCmd())
	return out
}

// 风险表里的命令名必须真实存在，防止改名后表项失效。
func TestRiskTablesMatchCommands(t *testing.T) {
	leaves := leafPaths()
	for p := range dangerCmds {
		if !leaves[p] {
			t.Errorf("dangerCmds 中的命令不存在: %s", p)
		}
		if readCmds[p] {
			t.Errorf("命令同时出现在 read 与 danger: %s", p)
		}
	}
	for p := range readCmds {
		if !leaves[p] {
			t.Errorf("readCmds 中的命令不存在: %s", p)
		}
	}
}

// runJSON 执行命令并捕获 Finish 输出的 JSON 信封。
func runJSON(t *testing.T, args ...string) (map[string]any, int) {
	t.Helper()
	t.Setenv("DOO_TOKEN", "test-token")
	t.Setenv("DOO_SERVER", "http://127.0.0.1:1") // 不可达；危险命令应在发请求前就返回
	t.Setenv("HOME", t.TempDir())
	cli.Opts = cli.Options{}
	cli.Current = cli.Invocation{}

	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w
	root := NewRootCmd()
	root.SetArgs(append(args, "--json"))
	root.SetErr(&bytes.Buffer{})
	err := root.Execute()
	if !cli.Current.Started {
		cli.Opts.JSON = true
	}
	code := cli.Finish(err)
	w.Close()
	os.Stdout = stdout
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	var out map[string]any
	if jerr := json.Unmarshal(buf.Bytes(), &out); jerr != nil {
		t.Fatalf("输出不是 JSON: %q", buf.String())
	}
	return out, code
}

func TestDangerNeedsConfirm(t *testing.T) {
	out, code := runJSON(t, "task", "delete", "123")
	if out["ok"] != false || out["code"] != cli.CodeNeedsConfirm || code != 4 {
		t.Fatalf("期望 needs_confirm，得到 %v (exit %d)", out, code)
	}
	c, _ := out["confirm"].(map[string]any)
	if c["command"] != "task delete" || c["risk"] != "danger" || !strings.Contains(c["summary"].(string), "#123") {
		t.Errorf("confirm 字段异常: %v", c)
	}
	targets, _ := c["targets"].([]any)
	if len(targets) != 1 {
		t.Errorf("targets 异常: %v", c["targets"])
	}
}

// 未显式调 Confirm 的危险命令，由 guard 兜底拦截。
func TestDangerGuardFallback(t *testing.T) {
	saved := dangerCmds["dialog read"]
	dangerCmds["dialog read"] = true
	defer func() { dangerCmds["dialog read"] = saved }()
	out, _ := runJSON(t, "dialog", "read", "1")
	if out["code"] != cli.CodeNeedsConfirm {
		t.Fatalf("期望 guard 兜底 needs_confirm，得到 %v", out)
	}
}

func TestRecoverIsNotDanger(t *testing.T) {
	out, _ := runJSON(t, "task", "delete", "123", "--recover", "--dry-run")
	data, _ := out["data"].(map[string]any)
	if out["ok"] != true || data["risk"] != "write" {
		t.Fatalf("--recover 应为 write 且 dry-run 成功，得到 %v", out)
	}
}

func TestDryRun(t *testing.T) {
	out, code := runJSON(t, "task", "delete", "123", "--dry-run")
	data, _ := out["data"].(map[string]any)
	if out["ok"] != true || code != 0 || data["dry_run"] != true || data["command"] != "task delete" {
		t.Fatalf("dry-run 输出异常: %v", out)
	}
}

func TestInvalidArgsEnvelope(t *testing.T) {
	for _, args := range [][]string{{"task", "view"}, {"nonexistent"}, {"task", "view", "abc"}} {
		out, code := runJSON(t, args...)
		if out["ok"] != false || out["code"] != cli.CodeInvalidArgs || code != 2 {
			t.Errorf("args %v 期望 invalid_args，得到 %v (exit %d)", args, out, code)
		}
	}
}

func TestNetworkErrorEnvelope(t *testing.T) {
	out, _ := runJSON(t, "task", "view", "1")
	if out["ok"] != false || out["code"] != cli.CodeNetworkError {
		t.Fatalf("期望 network_error，得到 %v", out)
	}
}

func TestMeta(t *testing.T) {
	out, _ := runJSON(t, "meta", "--jq", `.commands[] | select(.command=="task delete") | .risk`)
	if out["ok"] != true || out["data"] != "danger" {
		t.Fatalf("meta 输出异常: %v", out)
	}
}
