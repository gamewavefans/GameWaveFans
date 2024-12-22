package zbc

import (
	"fmt"
	"io"
	"reflect"

	"github.com/namgo/GameWaveFans/pkg/common"
)

// IsPacked return bool if the file is zlib-packed
func IsPacked(r io.ReadSeeker) (bool, error) {
	header, err := common.ReadBytes(r, 0, len(packedHeader))
	if err != nil {
		return false, err
	}
	_, err = r.Seek(0, 0)
	if err != nil {
		return false, err
	}
	if reflect.DeepEqual(header, []byte(packedHeader)) {
		return true, nil
	}
	return false, nil
}

// Unpack return unzlibbed ZBC bytecode
func Unpack(r io.ReadSeeker) ([]byte, error) {
	packed, err := IsPacked(r)
	if err != nil {
		return []byte{}, err
	}
	if !packed {
		return []byte{}, fmt.Errorf("file is not packed")
	}

	// unpackedSize, err := common.ReadUint32(r, 8)
	// if err != nil {
	// 	return []byte{}, err
	// }

	// packedSize, err := common.ReadUint32(r, 0xC)
	// if err != nil {
	// 	return []byte{}, err
	// }

	_, err = r.Seek(0x10, 0)
	if err != nil {
		return []byte{}, err
	}
	unpacked, err := common.ReadZlib(r)
	if err != nil {
		return nil, err
	}
	return unpacked, nil
}
