package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// 风险级别。
const (
	RiskRead   = "read"
	RiskWrite  = "write"
	RiskDanger = "danger"
)

// Invocation 本次执行的命令信息（CLI 单次执行，进程级单例）。
type Invocation struct {
	Path    string            // 命令路径，如 "task delete"
	Risk    string            // 风险级别
	Args    []string          // 位置参数
	Flags   map[string]string // 显式给出的局部 flag
	Started bool              // 已通过参数校验进入执行

	confirmed bool
	data      any
	hasData   bool
	messages  []string
}

// Current 本次执行的命令信息。
var Current Invocation

// Begin 由根命令在参数校验通过后调用。
func Begin(path, risk string, args []string, flags map[string]string) {
	Current = Invocation{Path: path, Risk: risk, Args: args, Flags: flags, Started: true}
}

// guard 发起任何远端请求前的统一护栏：
// 非只读命令在 --dry-run 时不执行；danger 命令未经确认时走 Confirm（兜底，防止命令漏调）。
func guard() error {
	if Current.Risk == RiskRead || !Current.Started {
		return nil
	}
	if Opts.DryRun {
		return dryRunResult()
	}
	if Current.Risk == RiskDanger && !Current.confirmed {
		return Confirm(strings.TrimSpace(Current.Path + " " + strings.Join(Current.Args, " ")))
	}
	return nil
}

// dryRunError 预演结果：以成功形式输出，命令不执行。
type dryRunError struct{}

func (dryRunError) Error() string { return "dry-run" }

func dryRunResult() error { return dryRunError{} }

// Finish 统一收尾：--json 时输出信封 {"ok":...}，否则人类可读；返回进程退出码。
func Finish(err error) int {
	if _, ok := err.(dryRunError); ok {
		err = nil
		Current.data = map[string]any{
			"dry_run": true,
			"command": Current.Path,
			"risk":    Current.Risk,
			"args":    Current.Args,
			"flags":   Current.Flags,
		}
		Current.hasData = true
		if !Opts.JSON {
			fmt.Printf("[dry-run] %s（%s）未执行\n", Current.Path, Current.Risk)
			return 0
		}
	}
	if err != nil {
		e := Classify(err)
		if !Current.Started && e.Code == CodeError {
			e.Code = CodeInvalidArgs // cobra 参数/命令校验阶段的错误
		}
		if Opts.JSON {
			writeJSON(errEnvelope{OK: false, Code: e.Code, Message: e.Message, Confirm: e.Confirm})
		} else {
			fmt.Fprintln(os.Stderr, "doo: "+e.Message)
		}
		return exitCodeOf(e.Code)
	}
	if Opts.JSON && Current.Started {
		out := okEnvelope{OK: true, Message: strings.Join(Current.messages, "\n")}
		if Current.hasData {
			data, ferr := shape(Current.data)
			if ferr != nil {
				return Finish(ferr)
			}
			out.Data = data
		}
		writeJSON(out)
	}
	return 0
}

type okEnvelope struct {
	OK      bool   `json:"ok"`
	Data    any    `json:"data"`
	Message string `json:"message,omitempty"`
}

type errEnvelope struct {
	OK      bool         `json:"ok"`
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Confirm *ConfirmInfo `json:"confirm,omitempty"`
}

func writeJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}
