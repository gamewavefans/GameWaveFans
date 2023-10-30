// Package zwf contains functions for handling Gamewave .zwf audio files
// This package is partially compatible with go-audio interface
package zwf

import (
	"compress/zlib"
	"encoding/binary"
	"io"

	"github.com/go-audio/audio"
	"github.com/namgo/GameWaveFans/pkg/common"
)

// Decode reads .zwf file and returns raw audio data
func Decode(r io.ReadSeeker) (*audio.IntBuffer, error) {
	format := &audio.Format{
		NumChannels: 2,
		SampleRate:  22050,
	}

	// read samples count
	samplesCount, err := common.ReadUint32(r, 4)
	if err != nil {
		return nil, err
	}

	// read packed data  size
	// packedSize, err := common.ReadUint32(r, 0xC)
	// if err != nil {
	// 	return nil, err
	// }

	// seek to zlib data
	if _, err := r.Seek(0x14, 0); err != nil {
		return nil, err
	}

	zlibDecoder, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	buffer, err := io.ReadAll(zlibDecoder)
	if err != nil {
		return nil, err
	}
	err = zlibDecoder.Close()
	if err != nil {
		return nil, err
	}

	samples := make([]int, samplesCount)
	for i := range samples {

		sample := binary.BigEndian.Uint16(buffer[i*2 : (i+1)*2])
		samples[i] = int(sample)
	}

	buf := &audio.IntBuffer{Format: format, SourceBitDepth: 16, Data: samples}

	return buf, nil
}
