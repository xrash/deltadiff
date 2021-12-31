package deltadiff

import (
	"encoding/binary"
	"fmt"
	"github.com/xrash/deltadiff/hasher"
	"io"
)

type SignatureConfig struct {
	Hasher    string
	BlockSize int
	BaseSize  int
}

func Signature(base io.Reader, out io.Writer, c *SignatureConfig) error {
	if c.BaseSize < 0 {
		return fmt.Errorf("Must provide valid BaseSize in config")
	}

	h, err := hasher.GetHasherByName(c.Hasher)
	if err != nil {
		return fmt.Errorf("Didn't find hasher %s", c.Hasher)
	}

	if err := writeHasherCode(out, h); err != nil {
		return fmt.Errorf("Couldn't write hasher code %s", err)
	}

	if err := writeBlockSize(out, c.BlockSize); err != nil {
		return fmt.Errorf("Couldn't write hasher code %s", err)
	}

	if err := writeBaseSize(out, c.BaseSize); err != nil {
		return fmt.Errorf("Couldn't write hasher code %s", err)
	}

	for {
		b := make([]byte, c.BlockSize)
		_, err := base.Read(b)

		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		h.Reset()
		hashed, err := h.Hash(b)
		if err != nil {
			return err
		}

		written, err := out.Write(hashed)
		if err != nil {
			return err
		}

		if written != h.HashSize() {
			return fmt.Errorf("Wrote %d instead of expected hash size %d", written, h.HashSize())
		}
	}

	return nil
}

func writeHasherCode(out io.Writer, h hasher.Hasher) error {
	written, err := out.Write(h.Code())
	if err != nil {
		return err
	}

	if len(h.Code()) != written {
		return fmt.Errorf(
			"Wrote %d instead of expected hash code len %d",
			written,
			len(h.Code()),
		)
	}

	return nil
}

func writeBlockSize(out io.Writer, blockSize int) error {

	blocksizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(blocksizeBytes, uint32(blockSize))

	written, err := out.Write(blocksizeBytes)
	if err != nil {
		return err
	}

	if len(blocksizeBytes) != written {
		return fmt.Errorf(
			"Wrote %d instead of expected block size len %d",
			written,
			len(blocksizeBytes),
		)
	}

	return nil
}

func writeBaseSize(out io.Writer, baseSize int) error {

	basesizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(basesizeBytes, uint32(baseSize))

	written, err := out.Write(basesizeBytes)
	if err != nil {
		return err
	}

	if len(basesizeBytes) != written {
		return fmt.Errorf(
			"Wrote %d instead of expected base size len %d",
			written,
			len(basesizeBytes),
		)
	}

	return nil
}
