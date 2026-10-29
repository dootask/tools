package cli

import (
	"reflect"
	"testing"

	dootask "github.com/dootask/tools/server/go"
)

func TestShapeFieldsLimit(t *testing.T) {
	Opts = Options{Fields: []string{"id", "name"}, Limit: 2}
	defer func() { Opts = Options{} }()
	in := map[string]any{"total": 3, "data": []any{
		map[string]any{"id": 1, "name": "a", "desc": "x"},
		map[string]any{"id": 2, "name": "b", "desc": "y"},
		map[string]any{"id": 3, "name": "c", "desc": "z"},
	}}
	got, err := shape(in)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"total": float64(3), "data": []any{
		map[string]any{"id": float64(1), "name": "a"},
		map[string]any{"id": float64(2), "name": "b"},
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("shape=%v", got)
	}
}

func TestShapeJQ(t *testing.T) {
	Opts = Options{JQ: "[.[] | .id]"}
	defer func() { Opts = Options{} }()
	got, err := shape([]any{map[string]any{"id": 1}, map[string]any{"id": 2}})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []any{float64(1), float64(2)}) {
		t.Errorf("jq=%v", got)
	}
	Opts = Options{JQ: ".["}
	if _, err := shape(1); err == nil || Classify(err).Code != CodeInvalidArgs {
		t.Errorf("非法 jq 应返回 invalid_args，得到 %v", err)
	}
}

func TestClassifyAPIError(t *testing.T) {
	cases := map[string]error{
		CodeAuthRequired:     &dootask.APIError{Ret: -1, Msg: "身份已失效,请重新登录"},
		CodeNotFound:         &dootask.APIError{Ret: -4002, Msg: "任务不存在"},
		CodePermissionDenied: &dootask.APIError{Ret: 0, Msg: "权限不足"},
		CodeAPIError:         &dootask.APIError{Ret: 0, Msg: "请输入任务名称"},
		CodeServerError:      &dootask.HTTPError{StatusCode: 502},
	}
	for want, err := range cases {
		if got := Classify(err).Code; got != want {
			t.Errorf("Classify(%v)=%s，期望 %s", err, got, want)
		}
	}
}
