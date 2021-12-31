package hasher

import (
	"encoding/binary"
)

// In Go, if we do -3 % 20 we get -3
// but we'd like to get 17, so this
// function does that.
func mod(n, m int) int {
	return ((n % m) + m) % m
}

type PolyrollHasher struct {
	Base int
	Mod  int

	memory    bool
	internal  int
	bootstrap int
	pos       []int
}

func (h *PolyrollHasher) calculatePositionFactors(blockSize int) []int {
	pos := make([]int, blockSize)
	n := 1

	for i := 0; i < blockSize; i++ {
		pos[i] = n
		n = (n * h.Base) % h.Mod
	}

	// Reverse the array as we want the most
	// significant position to be at the left.
	for i := 0; i < blockSize/2; i++ {
		tmp := pos[i]
		pos[i] = pos[len(pos)-1-i]
		pos[len(pos)-1-i] = tmp
	}

	return pos
}

func (h *PolyrollHasher) SingleHash(data []byte) ([]byte, error) {
	pos := h.calculatePositionFactors(len(data))

	var hash int

	for i := 0; i < len(data); i++ {
		value := int(data[i])
		pos := pos[i]
		positioned := value * pos
		hash += positioned
	}

	hash = mod(hash, h.Mod)

	hashBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(hashBytes, uint32(hash))

	return hashBytes, nil
}

func (h *PolyrollHasher) Hash(data []byte) ([]byte, error) {
	if len(h.pos) == 0 {
		pos := h.calculatePositionFactors(len(data))
		h.pos = pos
	}

	var hash int

	if !h.memory {
		var first int

		for i := 0; i < len(data); i++ {
			value := int(data[i])
			pos := h.pos[i]
			positioned := value * pos
			hash += positioned

			if i == 0 {
				first = positioned
			}
		}

		h.bootstrap = hash - mod(first, h.Mod)

		hash = mod(hash, h.Mod)

		h.memory = true

	} else {
		nextValue := int(data[len(data)-1])
		lastPos := h.pos[len(h.pos)-1]
		nextPositioned := nextValue * lastPos

		hash = mod(h.bootstrap*h.Base, h.Mod) + nextPositioned

		// We could have `bootstrap` to be the value we
		// get after applying the mod, however that
		// makes the subtraction to possibly bring up
		// negative numbers and we don't want negative
		// numbers in modulo arithmetic.
		first := int(data[0]) * h.pos[0]

		h.bootstrap = hash - mod(first, h.Mod)

		hash = mod(hash, h.Mod)
	}

	hashBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(hashBytes, uint32(hash))

	return hashBytes, nil
}

func (h *PolyrollHasher) HashSize() int {
	return 4
}

func (h *PolyrollHasher) Code() []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, HASHER_CODE_POLYROLL)
	return b
}

func (h *PolyrollHasher) Reset() {
	h.memory = false
}
