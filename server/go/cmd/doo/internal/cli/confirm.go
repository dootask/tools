package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Target 危险操作涉及的对象（供确认卡片按 type+id 查名称展示）。
type Target struct {
	Type string `json:"type"`
	ID   any    `json:"id"`
}

// ConfirmInfo needs_confirm 时附带的确认信息。
type ConfirmInfo struct {
	Summary string   `json:"summary"`           // 操作摘要，如「删除任务 #123」
	Command string   `json:"command"`           // 命令路径，如「task delete」
	Risk    string   `json:"risk"`              // 风险级别
	Targets []Target `json:"targets,omitempty"` // 涉及对象
}

// Confirm 危险操作护栏，action 为操作摘要（如「删除任务 #123」）：
//   - --dry-run 返回预演结果，不执行；
//   - --yes 直接放行；
//   - --json（机器调用）返回 needs_confirm，由调用方确认后带 --yes 重跑；
//   - 交互终端弹 y/N 提问；
//   - 其余非交互环境同样返回 needs_confirm（防脚本误删/卡死）。
func Confirm(action string, targets ...Target) error {
	if Opts.DryRun {
		return dryRunResult()
	}
	if Opts.Yes {
		Current.confirmed = true
		return nil
	}
	if Opts.JSON || !term.IsTerminal(int(os.Stdin.Fd())) {
		return &Error{
			Code:    CodeNeedsConfirm,
			Message: fmt.Sprintf("%s（危险操作需确认；确认后加 --yes 重新执行）", action),
			Confirm: &ConfirmInfo{Summary: action, Command: Current.Path, Risk: Current.Risk, Targets: targets},
		}
	}
	fmt.Fprintf(os.Stderr, "确认%s? [y/N]: ", action)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		Current.confirmed = true
		return nil
	default:
		return Errorf(CodeCancelled, "已取消")
	}
}
