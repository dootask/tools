package commands

import (
	"bytes"
	"testing"
)

// 这些用例都在 RunE 之前（参数/命令校验阶段）失败，因此不触网。
func TestRootArgValidation(t *testing.T) {
	cases := [][]string{
		{"task", "view"},           // 缺少必填位置参数
		{"task", "view", "a", "b"}, // 位置参数过多
		{"nonexistent-command"},    // 未知命令
	}
	for _, args := range cases {
		root := NewRootCmd()
		root.SetArgs(args)
		root.SetOut(&bytes.Buffer{})
		root.SetErr(&bytes.Buffer{})
		if err := root.Execute(); err == nil {
			t.Errorf("args %v 期望报错，实际成功", args)
		}
	}
}

func TestRootHasAllNouns(t *testing.T) {
	root := NewRootCmd()
	want := []string{"auth", "task", "project", "column", "dialog", "message", "group", "user", "bot", "file", "report", "search", "system"}
	have := map[string]bool{}
	for _, c := range root.Commands() {
		have[c.Name()] = true
	}
	for _, w := range want {
		if !have[w] {
			t.Errorf("缺少子命令: %s", w)
		}
	}
}
