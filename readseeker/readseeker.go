package readseeker

import (
	"io"
)

type BasicReadSeeker struct {
	reader io.Reader
	cursor int64
}

func NewBasicReadSeeker(r io.Reader) io.ReadSeeker {
	return &BasicReadSeeker{
		reader: r,
	}
}

// Whence doesn't work, this is forward seek from current
// position only.
func (rs *BasicReadSeeker) Seek(offset int64, whence int) (int64, error) {
	if int64(rs.cursor) >= offset {
		return rs.cursor, nil
	}

	todiscard := offset - rs.cursor
	buffer := make([]byte, todiscard)
	n, err := rs.reader.Read(buffer)

	// Just in case it returns something < 0
	if n > 0 {
		rs.cursor += int64(n)
	}

	return rs.cursor, err
}

func (rs *BasicReadSeeker) Read(p []byte) (int, error) {
	n, err := rs.reader.Read(p)

	// Just in case it returns something < 0
	if n > 0 {
		rs.cursor += int64(n)
	}

	return n, err
}
