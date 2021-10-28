package slog

import (
	"runtime"
	"strings"
	"sync"
)

var (
	LogIDs = map[string]string{}
	locker = sync.RWMutex{}
)

func goid() string {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine ")
	)
	idField := strings.Fields(stk)[0]
	return idField
}

func SetLogID(xid string) {
	locker.Lock()
	defer locker.Unlock()
	LogIDs[goid()] = xid
}

func GetLogID() string {
	locker.RLock()
	defer locker.RUnlock()
	if logID, ok := LogIDs[goid()]; ok {
		return logID
	}
	return ""
}

func Close() {
	locker.Lock()
	defer locker.Unlock()
	goId := goid()
	if _, ok := LogIDs[goId]; ok {
		delete(LogIDs, goId)
	}
}
