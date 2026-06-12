package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseInt 解析单个整型位置参数。
func ParseInt(s, name string) (int, error) {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0, fmt.Errorf("无效的%s: %q", name, s)
	}
	return n, nil
}

// ParseIDList 把 "1,2,3" 解析为 []int；空串返回 nil。
func ParseIDList(s string) ([]int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("无效的 ID 列表项: %q", p)
		}
		out = append(out, n)
	}
	return out, nil
}

// BuildTimes 把开始/结束时间拼成 SDK 的 Times 切片。
// 都为空返回 nil；只给其一时另一位置留空字符串。
func BuildTimes(start, end string) []string {
	start, end = strings.TrimSpace(start), strings.TrimSpace(end)
	if start == "" && end == "" {
		return nil
	}
	return []string{start, end}
}
