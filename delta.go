package deltadiff

import (
	"encoding/binary"
	"fmt"
	"github.com/xrash/deltadiff/hasher"
	"io"
	"io/ioutil"
	"os"
)

type match struct {
	block        int
	segmentBegin int
	segmentEnd   int
}

type operation struct {
	// `kind` can be either "read" or "write"
	kind string
	from int
	to   int
	data []byte
}

type DeltaConfig struct {
	Debug       bool
	DebugWriter io.Writer
}

func Delta(signature, target io.Reader, result io.Writer, c *DeltaConfig) error {

	if c.Debug && c.DebugWriter == nil {
		c.DebugWriter = os.Stderr
	}

	hashcode, err := readHashCode(signature)
	if err != nil {
		return err
	}

	blockSize, err := readBlockSize(signature)
	if err != nil {
		return err
	}

	baseSize, err := readBaseSize(signature)
	if err != nil {
		return err
	}

	h, err := hasher.GetHasherByCode(hashcode)
	if err != nil {
		return err
	}

	blocks, err := readBlocks(signature, h)
	if err != nil {
		return err
	}

	buffer, err := readTarget(target)
	if err != nil {
		return err
	}

	matches, err := collectMatches(
		blocks,
		buffer,
		h,
		blockSize,
	)

	if err != nil {
		return err
	}

	operations := calculateOperations(
		matches,
		buffer,
		blockSize,
		baseSize,
	)

	operations = mergeConsecutiveReads(
		operations,
	)

	if c.Debug {
		for _, m := range matches {
			fmt.Printf("match\t%d:%d-%d\n", m.block, m.segmentBegin, m.segmentEnd)
		}

		for _, o := range operations {
			fmt.Printf("op\t%s:%d-%d\n", o.kind, o.from, o.to)
		}
	}

	if err := writeDelta(operations, result); err != nil {
		return fmt.Errorf("Error writing delta: %v", err)
	}

	return nil
}

func writeDelta(operations []*operation, out io.Writer) error {
	for i := 0; i < len(operations); i++ {
		op := operations[i]

		var opbytes []byte
		switch op.kind {
		case "write":
			opbytes = opWrite(op.data)
		case "read":
			opbytes = opRead(op.from, op.to)
		default:
			return fmt.Errorf("Unexpected op.kind %s", op.kind)
		}

		written, err := out.Write(opbytes)
		if err != nil {
			return err
		}

		if written != len(opbytes) {
			return fmt.Errorf("Couldn't write everything, wrote only %d", written)
		}
	}

	return nil
}

func mergeConsecutiveReads(operations []*operation) []*operation {
	ops := make([]*operation, 0)

	for i := 0; i < len(operations); i++ {
		op := operations[i]

		switch op.kind {
		case "write":
			ops = append(ops, op)
		case "read":
			updatePrevInsteadOfAppending := i > 0 &&
				ops[len(ops)-1].kind == "read" &&
				ops[len(ops)-1].to == op.from

			if updatePrevInsteadOfAppending {
				ops[len(ops)-1].to = op.to
			} else {
				ops = append(ops, op)
			}
		}
	}

	return ops
}

func calculateOperations(matches []*match, target []byte, blockSize, baseSize int) []*operation {
	ops := make([]*operation, 0)

	if len(matches) == 0 {
		from := 0
		to := len(target)
		op := &operation{
			kind: "write",
			from: from,
			to:   to,
			data: target[from:to],
		}
		ops = append(ops, op)
		return ops
	}

	for i := 0; i < len(matches); i++ {
		match := matches[i]

		if i == 0 && match.block > 0 {
			from := 0
			to := match.segmentBegin
			op := &operation{
				kind: "write",
				from: from,
				to:   to,
				data: target[from:to],
			}
			ops = append(ops, op)
		}

		if i > 0 && match.segmentBegin != matches[i-1].segmentEnd {
			from := matches[i-1].segmentEnd
			to := match.segmentBegin
			op := &operation{
				kind: "write",
				from: from,
				to:   to,
				data: target[from:to],
			}
			ops = append(ops, op)
		}

		from := match.block * blockSize
		to := match.block*blockSize + blockSize

		if to > baseSize {
			to = baseSize
		}

		op := &operation{
			kind: "read",
			from: from,
			to:   to,
		}

		ops = append(ops, op)
	}

	if matches[len(matches)-1].segmentEnd != len(target) {
		from := matches[len(matches)-1].segmentEnd
		to := len(target)
		op := &operation{
			kind: "write",
			from: from,
			to:   to,
			data: target[from:to],
		}
		ops = append(ops, op)
	}

	return ops
}

func collectMatches(signature [][]byte, target []byte, h hasher.Hasher, blockSize int) ([]*match, error) {

	matches := make([]*match, 0)

	anchor := 0

	for block := 0; block < len(signature); block++ {
		if block == len(signature) {
			return matches, nil
		}

		segmentBegin := anchor
		segmentEnd := segmentBegin + blockSize

		if segmentEnd > len(target) {
			segmentEnd = len(target)
		}

		h.Reset()

		for {
			if (segmentEnd) > len(target) {
				break
			}

			segment := target[segmentBegin:segmentEnd]
			hashedSegment, err := h.Hash(segment)
			if err != nil {
				return nil, err
			}

			if hashesAreEqual(signature[block], hashedSegment) {
				match := &match{
					block:        block,
					segmentBegin: segmentBegin,
					segmentEnd:   segmentEnd,
				}
				matches = append(matches, match)
				anchor = segmentEnd
				break
			}

			segmentBegin++
			segmentEnd++
		}
	}

	return matches, nil
}

func readHashCode(signature io.Reader) ([]byte, error) {
	buffer := make([]byte, 2)
	read, err := signature.Read(buffer)
	if err != nil {
		return nil, err
	}

	if read != len(buffer) {
		return nil, fmt.Errorf("Read %d instead of expected hash code len %d", read, len(buffer))
	}

	return buffer, nil
}

func readBlockSize(signature io.Reader) (int, error) {
	buffer := make([]byte, 4)
	read, err := signature.Read(buffer)
	if err != nil {
		return -1, err
	}

	if read != len(buffer) {
		return -1, fmt.Errorf("Read %d instead of expected block size len %d", read, len(buffer))
	}

	blockSize := binary.BigEndian.Uint32(buffer)

	return int(blockSize), nil
}

func readBaseSize(signature io.Reader) (int, error) {
	buffer := make([]byte, 4)
	read, err := signature.Read(buffer)
	if err != nil {
		return -1, err
	}

	if read != len(buffer) {
		return -1, fmt.Errorf("Read %d instead of expected base size len %d", read, len(buffer))
	}

	baseSize := binary.BigEndian.Uint32(buffer)

	return int(baseSize), nil
}

func readBlocks(signature io.Reader, h hasher.Hasher) ([][]byte, error) {
	blocks := make([][]byte, 0)

	for {
		buffer := make([]byte, h.HashSize())
		read, err := signature.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if read > 0 {
			blocks = append(blocks, buffer)
		}

		if err == io.EOF {
			break
		}
	}

	return blocks, nil
}

func readTarget(in io.Reader) ([]byte, error) {
	content, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func hashesAreEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
