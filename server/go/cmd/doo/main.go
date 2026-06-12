// Command doo 是 DooTask 的命令行工具，建在官方 Go SDK 之上。
package main

import (
	"fmt"
	"os"

	"github.com/dootask/tools/server/go/cmd/doo/internal/cli"
	"github.com/dootask/tools/server/go/cmd/doo/internal/commands"
)

func main() {
	root := commands.NewRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "doo: "+err.Error())
		os.Exit(cli.ExitCode(err))
	}
}
