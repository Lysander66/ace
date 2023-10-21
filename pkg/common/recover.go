package common

import (
	"log/slog"
	"runtime/debug"
)

func PanicRecover() {
	if e := recover(); e != nil {
		slog.Error("recover", "err", e, "stack", string(debug.Stack()))
	}
}
