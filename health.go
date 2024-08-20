package xrf197ilz35aq0

import (
	"runtime"
	"time"
)

// Runtime application stats
type Runtime struct {
	Arch            string // running program's architecture target: one of 386, amd64, arm, s390x, e.t.c.
	OS              string // running program's OS target: e.g darwin, freebsd, linux, e.t.c
	Alloc           int    // Total amount of memory allocated on the heap (includes both in-use and freed memory)
	CollectedAt     time.Time
	HeapAlloc       int    // Memory allocated and currently in use on the heap
	GoroutinesCount int    // num of goroutines that are currently running.
	Version         string // either the commit hash & date at the time of the build or, a release tag e.g "go1.3"
}

const version = "0.0.0"

type Health struct {
	// version contains current app version.
	version string

	Runtime Runtime
}

func (h *Health) Version() string {
	return h.version
}

func NewHealth() Health {
	// MemStats is a struct that contains statistics about memory allocation and garbage collection
	ms := runtime.MemStats{}
	// call runtime.ReadMemStats(&memStats) to populate a MemStats struct with the current memory statistics
	// TODO. better to do periodically
	runtime.ReadMemStats(&ms)

	return Health{
		version: version,
		Runtime: Runtime{
			CollectedAt:     time.Now(),
			OS:              runtime.GOOS,
			Alloc:           int(ms.Alloc),
			Arch:            runtime.GOARCH,
			HeapAlloc:       int(ms.HeapAlloc),
			Version:         runtime.Version(),
			GoroutinesCount: runtime.NumGoroutine(),
		},
	}
}
