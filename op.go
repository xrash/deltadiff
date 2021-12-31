package deltadiff

import (
	"encoding/binary"
)

const (
	OP_WRITE uint16 = 0
	OP_READ  uint16 = 1
)

func opRead(from, to int) []byte {
	opcodeBytes := make([]byte, 2)
	fromBytes := make([]byte, 4)
	toBytes := make([]byte, 4)

	binary.BigEndian.PutUint16(opcodeBytes, OP_READ)
	binary.BigEndian.PutUint32(fromBytes, uint32(from))
	binary.BigEndian.PutUint32(toBytes, uint32(to))

	op := make([]byte, 0)
	op = append(op, opcodeBytes...)
	op = append(op, fromBytes...)
	op = append(op, toBytes...)

	return op
}

func opWrite(data []byte) []byte {
	opcodeBytes := make([]byte, 2)
	datalenBytes := make([]byte, 4)

	binary.BigEndian.PutUint16(opcodeBytes, OP_WRITE)
	binary.BigEndian.PutUint32(datalenBytes, uint32(len(data)))

	op := make([]byte, 0)
	op = append(op, opcodeBytes...)
	op = append(op, datalenBytes...)
	op = append(op, data...)

	return op
}
