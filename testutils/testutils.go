package testutils

import "testing"

// ErrorfHandler エラーメッセージ出力処理
func ErrorfHandler(tb testing.TB, want interface{}, got interface{}) {
	tb.Helper()
	tb.Errorf("want = %d, got = %d", want, got)
}
