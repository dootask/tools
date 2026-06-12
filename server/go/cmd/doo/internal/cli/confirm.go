package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Confirm 危险操作护栏：
//   - --yes 直接放行；
//   - 交互终端弹 y/N 提问；
//   - 非交互且未给 --yes 时报错（防脚本误删/卡死）。
func Confirm(prompt string) error {
	if Opts.Yes {
		return nil
	}
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("%s（危险操作需确认；非交互环境请加 --yes）", prompt)
	}
	fmt.Fprintf(os.Stderr, "%s [y/N]: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return nil
	default:
		return fmt.Errorf("已取消")
	}
}
