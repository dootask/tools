package cli

import (
	"reflect"
	"testing"
)

func TestParseIDList(t *testing.T) {
	cases := []struct {
		in      string
		want    []int
		wantErr bool
	}{
		{"", nil, false},
		{"  ", nil, false},
		{"1", []int{1}, false},
		{"1,2,3", []int{1, 2, 3}, false},
		{" 1 , 2 ,3 ", []int{1, 2, 3}, false},
		{"1,,2", []int{1, 2}, false},
		{"1,x", nil, true},
	}
	for _, c := range cases {
		got, err := ParseIDList(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("ParseIDList(%q) 期望报错，实际无", c.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseIDList(%q) 意外报错: %v", c.in, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ParseIDList(%q)=%v, 期望 %v", c.in, got, c.want)
		}
	}
}

func TestBuildTimes(t *testing.T) {
	cases := []struct {
		start, end string
		want       []string
	}{
		{"", "", nil},
		{"2026-01-01 00:00:00", "", []string{"2026-01-01 00:00:00", ""}},
		{"", "2026-01-02 00:00:00", []string{"", "2026-01-02 00:00:00"}},
		{"a", "b", []string{"a", "b"}},
	}
	for _, c := range cases {
		got := BuildTimes(c.start, c.end)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("BuildTimes(%q,%q)=%v, 期望 %v", c.start, c.end, got, c.want)
		}
	}
}

func TestParseInt(t *testing.T) {
	if n, err := ParseInt("42", "x"); err != nil || n != 42 {
		t.Errorf("ParseInt(42)=%d,%v", n, err)
	}
	if _, err := ParseInt("abc", "x"); err == nil {
		t.Errorf("ParseInt(abc) 期望报错")
	}
}
