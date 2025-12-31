package lg

import (
	"path/filepath"
	"runtime"
	"strings"
)

var appRoot string

func GetCallerInfo(skip int) Caller {
	var file, method string
	pc, _, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return Caller{}
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		method = "<unk>"
		file = "<unk>"
	} else {
		fnName := filepath.FromSlash(fn.Name())
		dotNum := strings.LastIndexByte(fnName, '.')

		method = fnName[dotNum+1:]
		file = strings.ReplaceAll(fnName[:dotNum], "\\", "/")
	}

	return Caller{
		Method: method,
		File:   file,
		Line:   line,
	}
}
