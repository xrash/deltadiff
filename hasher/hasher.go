package hasher

import (
	"encoding/binary"
	"fmt"
)

const (
	HASHER_CODE_POLYROLL uint16 = 0
	HASHER_CODE_MD5      uint16 = 1
	HASHER_CODE_CRC32    uint16 = 2

	POLYROLL_BASE = 257
	//POLYROLL_MOD = 8509909
	POLYROLL_MOD = 15485863
	//POLYROLL_MOD = 32452843
	//POLYROLL_MOD = 256203161
)

type Hasher interface {
	Hash([]byte) ([]byte, error)
	HashSize() int
	Code() []byte
	Reset()
}

func GetHasherByName(name string) (Hasher, error) {
	switch name {
	case "polyroll":
		return &PolyrollHasher{
			Base: POLYROLL_BASE,
			Mod:  POLYROLL_MOD,
		}, nil
	case "md5":
		return &MD5Hasher{}, nil
	case "crc32":
		return &CRC32Hasher{}, nil
	}

	return nil, fmt.Errorf("Unknown hasher %s", name)
}

func GetHasherByCode(codebytes []byte) (Hasher, error) {
	code := binary.BigEndian.Uint16(codebytes)

	switch code {
	case HASHER_CODE_POLYROLL:
		return &PolyrollHasher{
			Base: POLYROLL_BASE,
			Mod:  POLYROLL_MOD,
		}, nil
	case HASHER_CODE_MD5:
		return &MD5Hasher{}, nil
	case HASHER_CODE_CRC32:
		return &CRC32Hasher{}, nil
	}

	return nil, fmt.Errorf("Unknown hasher [%v] %v", codebytes, code)
}
