// Command doo 是 DooTask 的命令行工具，建在官方 Go SDK 之上。
package main

import (
	"os"
	"slices"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/dootask/tools/server/go/cmd/doo/internal/commands"
)

func main() {
	root := commands.NewRootCmd()
	err := root.Execute()
	if !cli.Current.Started && slices.Contains(os.Args[1:], "--json") {
		cli.Opts.JSON = true // 参数校验阶段即失败时，全局 flag 尚未解析
	}
	os.Exit(cli.Finish(err))
}
