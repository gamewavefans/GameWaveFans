// Package common contains commonly used functions
package common

import (
	"encoding/binary"
	"io"
)

// ReadUint32 reads 4 bytes from the input as Little Endian int
func ReadUint32(r io.ReadSeeker, offset int64) (uint32, error) {
	dataBytes, err := ReadBytes(r, offset, 4)
	if err != nil {
		return 0, err
	}

	headerInt := binary.LittleEndian.Uint32(dataBytes)

	return headerInt, nil
}

// ReadBytes reads length bytes from a specified location in a stream
func ReadBytes(r io.ReadSeeker, offset int64, length int) ([]byte, error) {
	if _, err := r.Seek(offset, 0); err != nil {
		return []byte{}, err
	}

	dataBytes := make([]byte, length)
	if _, err := r.Read(dataBytes); err != nil {
		return []byte{}, err
	}
	return dataBytes, nil
}
