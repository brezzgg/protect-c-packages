package lg

import (
	"path/filepath"
	"runtime"
	"strings"
)

func GetCallerInfo(skip int, splitPrefix bool) Caller {
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

		if strings.Contains(file, "/") && splitPrefix {
			if strings.HasPrefix(file, "github.com/") {
				file = splitFilePath(file, 3)
			} else {
				file = splitFilePath(file, 1)
			}
		}
	}

	return Caller{
		Method: method,
		File:   file,
		Line:   line,
	}
}

func isBadOffset(skip int) bool {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return false
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return false
	}
	return strings.HasPrefix(fn.Name(), "github.com/brezzgg/go-packages/lg")
}

func splitFilePath(path string, count int) string {
	buf := path
	for range count {
		index := strings.Index(buf, "/")
		if index == -1 {
			return path
		}
		buf = buf[index+1:]
	}
	return buf
}
