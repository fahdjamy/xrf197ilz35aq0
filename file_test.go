package xrf197ilz35aq0

import (
	"io"
)

// Mock for os.Open
type mockFileDataCopier struct {
	content []byte
	err     error
	readPos int // Keep track of the current read position
	closed  bool
}

// since this is a mock, close should do nothing and just return nil error
func (m *mockFileDataCopier) Close() error {
	m.closed = true
	return nil
}

// mockFileDataCopier's Read needs to fill the p buffer with data from its internal content,
// but starting from the readPos which tracks where we are in the simulated file
func (m *mockFileDataCopier) Read(p []byte) (n int, err error) {
	// In real file, when you call Read,
	// it fills the provided buffer (p) with data from the file (content),
	// starting from the current file position.
	if m.readPos >= len(m.content) {
		return 0, io.EOF // Simulate end-of-file when all content is read
	}

	// The copy is used to efficiently transfer data from one slice to another.
	// copies a portion of the m.content slice (starting from m.readPos) into the p slice
	// n captures the number of bytes actually copied from m.content to p.
	n = copy(p, m.content[m.readPos:])

	m.readPos += n
	return n, m.err
}

type mockFileDataTracker struct {
	closed      bool
	opened      bool
	openedTimes int
	closedTimes int
}

func (m *mockFileDataTracker) Close() error {
	m.closed = true
	m.closedTimes++
	return nil
}

func (m *mockFileDataTracker) Open() error {
	m.opened = true
	m.openedTimes++
	return nil
}
