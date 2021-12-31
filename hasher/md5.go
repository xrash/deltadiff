package hasher

import (
	"crypto/md5"
	"encoding/binary"
)

type MD5Hasher struct{}

func (h *MD5Hasher) Hash(data []byte) ([]byte, error) {
	hasher := md5.New()

	_, err := hasher.Write(data)
	if err != nil {
		return nil, err
	}

	hash := hasher.Sum(nil)

	return hash, nil
}

func (h *MD5Hasher) HashSize() int {
	return 16
}

func (h *MD5Hasher) Code() []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, HASHER_CODE_MD5)
	return b
}

func (h *MD5Hasher) Reset() {}
