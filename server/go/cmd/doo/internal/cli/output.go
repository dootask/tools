package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/itchyny/gojq"
)

// Output 统一输出：--json 时暂存，由 Finish 包成 {"ok":true,"data":...} 输出；否则人类可读表格/键值。
// cols 是表格首选列（json key），为空时自动挑选标量字段。--fields/--limit/--jq 两种模式都生效。
func Output(v any, cols []string) error {
	if Opts.JSON {
		Current.data = v
		Current.hasData = true
		return nil
	}
	g, err := shape(v)
	if err != nil {
		return err
	}
	if len(Opts.Fields) > 0 {
		cols = Opts.Fields
	}
	renderHuman(g, cols)
	return nil
}

// shape 按 --fields / --limit / --jq 裁剪输出，大结果不进模型上下文。
// 列表（数组或分页对象的 data）逐项取字段、截断条数；单个对象只取字段。--jq 最后执行。
func shape(v any) (any, error) {
	g, err := toGeneric(v)
	if err != nil {
		return nil, err
	}
	if len(Opts.Fields) > 0 || Opts.Limit > 0 {
		g = trimValue(g, Opts.Fields, Opts.Limit)
	}
	if Opts.JQ != "" {
		return runJQ(g, Opts.JQ)
	}
	return g, nil
}

func trimValue(v any, fields []string, limit int) any {
	switch t := v.(type) {
	case []any:
		if limit > 0 && len(t) > limit {
			t = t[:limit]
		}
		out := make([]any, len(t))
		for i, item := range t {
			out[i] = pickFields(item, fields)
		}
		return out
	case map[string]any:
		if arr, ok := t["data"].([]any); ok {
			cp := make(map[string]any, len(t))
			for k, val := range t {
				cp[k] = val
			}
			cp["data"] = trimValue(arr, fields, limit)
			return cp
		}
		return pickFields(t, fields)
	}
	return v
}

func pickFields(v any, fields []string) any {
	m, ok := v.(map[string]any)
	if !ok || len(fields) == 0 {
		return v
	}
	out := make(map[string]any, len(fields))
	for _, f := range fields {
		if val, ok := m[f]; ok {
			out[f] = val
		}
	}
	return out
}

// runJQ 执行 jq 表达式：单个结果直接返回，多个结果收成数组。
func runJQ(v any, expr string) (any, error) {
	q, err := gojq.Parse(expr)
	if err != nil {
		return nil, Errorf(CodeInvalidArgs, "无效的 --jq 表达式: "+err.Error())
	}
	var results []any
	iter := q.Run(v)
	for {
		r, ok := iter.Next()
		if !ok {
			break
		}
		if e, isErr := r.(error); isErr {
			return nil, Errorf(CodeInvalidArgs, "--jq 执行失败: "+e.Error())
		}
		results = append(results, r)
	}
	if len(results) == 1 {
		return results[0], nil
	}
	return results, nil
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
