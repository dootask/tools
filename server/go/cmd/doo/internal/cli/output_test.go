package cli

import (
	"strings"
	"testing"
)

func TestCellStringNumber(t *testing.T) {
	// JSON 数字解出来是 float64，整数应去掉小数。
	if got := cellString(float64(37999)); got != "37999" {
		t.Errorf("cellString(37999)=%q", got)
	}
	if got := cellString(float64(1.5)); got != "1.5" {
		t.Errorf("cellString(1.5)=%q", got)
	}
	if got := cellString(nil); got != "" {
		t.Errorf("cellString(nil)=%q", got)
	}
	if got := cellString(true); got != "true" {
		t.Errorf("cellString(true)=%q", got)
	}
}

func TestTableCellCollapseAndTruncate(t *testing.T) {
	// 含换行的多行内容应折叠成单行。
	got := tableCell("第一行\n第二行\t带制表")
	if strings.ContainsAny(got, "\n\t") {
		t.Errorf("tableCell 未折叠空白: %q", got)
	}
	if got != "第一行 第二行 带制表" {
		t.Errorf("tableCell 折叠结果=%q", got)
	}
	// 超长内容应截断并加省略号。
	long := strings.Repeat("中", 80)
	tc := tableCell(long)
	if r := []rune(tc); len(r) != 50 || r[len(r)-1] != '…' {
		t.Errorf("tableCell 截断异常: 长度=%d 末字符=%q", len([]rune(tc)), string([]rune(tc)[len([]rune(tc))-1]))
	}
}

func TestToGenericPaginate(t *testing.T) {
	type row struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	type paginate struct {
		Data  []row `json:"data"`
		Total int   `json:"total"`
	}
	g, err := toGeneric(paginate{Data: []row{{1, "a"}}, Total: 1})
	if err != nil {
		t.Fatal(err)
	}
	m, ok := g.(map[string]any)
	if !ok {
		t.Fatalf("toGeneric 未得到 map, 得到 %T", g)
	}
	if _, ok := m["data"].([]any); !ok {
		t.Errorf("data 应为 []any, 得到 %T", m["data"])
	}
}
