package hasher

import (
	"encoding/binary"
	"hash/crc32"
)

type CRC32Hasher struct{}

func (h *CRC32Hasher) Hash(data []byte) ([]byte, error) {
	hasher := crc32.NewIEEE()

	_, err := hasher.Write(data)
	if err != nil {
		return nil, err
	}

	hash := hasher.Sum(nil)

	return hash, nil
}

func (h *CRC32Hasher) HashSize() int {
	return 4
}

func (h *CRC32Hasher) Code() []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, HASHER_CODE_CRC32)
	return b
}

func (h *CRC32Hasher) Reset() {}
