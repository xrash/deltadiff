package deltadiff

import (
	"encoding/binary"
	"fmt"
	"github.com/xrash/deltadiff/readseeker"
	"io"
)

func Patch(base, delta io.Reader, out io.Writer) error {
	basers := readseeker.NewBasicReadSeeker(base)

	for {
		opcodeBytes := make([]byte, 2)
		bytesRead, err := delta.Read(opcodeBytes)
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		if bytesRead != 2 {
			return fmt.Errorf("Didn't read expected 2 bytes, read %d instead", bytesRead)
		}

		opcode := binary.BigEndian.Uint16(opcodeBytes)

		switch opcode {
		case OP_WRITE:
			if err := doPatchWrite(delta, out); err != nil {
				return fmt.Errorf("doPatchWrite: %v", err)
			}

		case OP_READ:
			if err := doPatchRead(basers, delta, out); err != nil {
				return fmt.Errorf("doPatchRead: %v", err)
			}

		default:
			return fmt.Errorf("Unknown operation %v", opcode)

		}
	}

	return nil

}

func doPatchRead(base io.ReadSeeker, delta io.Reader, out io.Writer) error {
	fromBytes := make([]byte, 4)
	fromRead, err := delta.Read(fromBytes)

	if err != nil {
		return err
	}

	if fromRead != 4 {
		return fmt.Errorf("Didn't read expected 4 bytes, read %d instead", fromRead)
	}

	toBytes := make([]byte, 4)
	toRead, err := delta.Read(toBytes)

	if err != nil {
		return err
	}

	if toRead != 4 {
		return fmt.Errorf("Didn't read expected 4 bytes, read %d instead", toRead)
	}

	from := binary.BigEndian.Uint32(fromBytes)
	to := binary.BigEndian.Uint32(toBytes)

	seekd, err := base.Seek(int64(from), 0)
	if err != nil {
		return err
	}

	if seekd != int64(from) {
		return err
	}

	buffer := make([]byte, to-from)
	bufferRead, err := base.Read(buffer)
	if err != nil {
		return err
	}

	if bufferRead != len(buffer) {
		return fmt.Errorf("Didn't read expected %d bytes, read %d instead", len(buffer), bufferRead)
	}

	written, err := out.Write(buffer)
	if written != len(buffer) {
		return fmt.Errorf("Didnt write expected %d, wrote %d instead", len(buffer), written)
	}

	return nil
}

func doPatchWrite(delta io.Reader, out io.Writer) error {
	datalenBytes := make([]byte, 4)
	datalenRead, err := delta.Read(datalenBytes)

	if err != nil {
		return err
	}

	if datalenRead != 4 {
		return fmt.Errorf("Didn't read expected 4 bytes, read %d instead", datalenRead)
	}

	datalen := binary.BigEndian.Uint32(datalenBytes)

	data := make([]byte, datalen)
	dataRead, err := delta.Read(data)
	if err != nil {
		return err
	}

	if uint32(dataRead) != datalen {
		return fmt.Errorf("Didn't read expected %d bytes, read %d instead", datalen, datalenRead)
	}

	written, err := out.Write(data)
	if written != len(data) {
		return fmt.Errorf("Didnt write expected %d, wrote %d instead", len(data), written)
	}

	return nil
}
