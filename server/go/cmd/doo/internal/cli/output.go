package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

// Output 统一输出：--json 输出紧凑 JSON；否则人类可读表格/键值。
// cols 是表格首选列（json key），为空时自动挑选标量字段。
func Output(v any, cols []string) error {
	if Opts.JSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(v)
	}
	g, err := toGeneric(v)
	if err != nil {
		return err
	}
	renderHuman(g, cols)
	return nil
}

// toGeneric 把任意值经 JSON 归一成 map/slice/标量，渲染与 SDK 结构体字段名解耦。
func toGeneric(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func renderHuman(v any, cols []string) {
	switch t := v.(type) {
	case map[string]any:
		// 分页对象：{data:[...], total, current_page}
		if d, ok := t["data"]; ok {
			if arr, ok := d.([]any); ok {
				renderTable(arr, cols)
				renderPaginateFooter(t)
				return
			}
		}
		renderObject(t)
	case []any:
		renderTable(t, cols)
	case nil:
		fmt.Println("(空)")
	default:
		fmt.Println(t)
	}
}

func renderTable(rows []any, cols []string) {
	if len(rows) == 0 {
		fmt.Println("(无数据)")
		return
	}
	maps := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		if m, ok := r.(map[string]any); ok {
			maps = append(maps, m)
		}
	}
	// 非对象数组：逐行打印
	if len(maps) != len(rows) {
		for _, r := range rows {
			fmt.Println(cellString(r))
		}
		return
	}
	columns := pickColumns(maps, cols)
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(columns, "\t"))
	for _, m := range maps {
		cellsRow := make([]string, len(columns))
		for i, c := range columns {
			cellsRow[i] = tableCell(m[c])
		}
		fmt.Fprintln(w, strings.Join(cellsRow, "\t"))
	}
	w.Flush()
}

func renderObject(m map[string]any) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	for _, k := range keys {
		fmt.Fprintf(w, "%s\t%s\n", k, cellString(m[k]))
	}
	w.Flush()
}

func renderPaginateFooter(m map[string]any) {
	total := m["total"]
	page := m["current_page"]
	if total != nil || page != nil {
		fmt.Printf("\n— 共 %s 条，当前第 %s 页 —\n", cellString(total), cellString(page))
	}
}

// pickColumns 选定表格列：优先用 cols 中实际存在的；否则取所有行里出现的标量字段（排序）。
func pickColumns(maps []map[string]any, cols []string) []string {
	if len(cols) > 0 {
		out := make([]string, 0, len(cols))
		for _, c := range cols {
			for _, m := range maps {
				if _, ok := m[c]; ok {
					out = append(out, c)
					break
				}
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	seen := map[string]bool{}
	var out []string
	for _, m := range maps {
		for k, v := range m {
			if seen[k] || !isScalar(v) {
				continue
			}
			seen[k] = true
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

func isScalar(v any) bool {
	switch v.(type) {
	case nil, bool, float64, string:
		return true
	default:
		return false
	}
}

// tableCell 把单元格压成单行并截断，保证表格对齐；完整值用 --json 查看。
func tableCell(v any) string {
	s := cellString(v)
	s = strings.Join(strings.Fields(s), " ") // 折叠所有空白（含换行/制表）为单空格
	const max = 50
	if r := []rune(s); len(r) > max {
		s = string(r[:max-1]) + "…"
	}
	return s
}

func cellString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float64:
		// JSON 数字统一为 float64；整数去掉小数。
		if t == float64(int64(t)) {
			return fmt.Sprintf("%d", int64(t))
		}
		return fmt.Sprintf("%g", t)
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Sprintf("%v", t)
		}
		return string(b)
	}
}
