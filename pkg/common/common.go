// Package common contains commonly used functions
package common

import (
	"bufio"
	"bytes"
	"compress/zlib"
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

// ReadUint32Big reads 4 bytes from the input as Big Endian int
func ReadUint32Big(r io.ReadSeeker, offset int64) (uint32, error) {
	dataBytes, err := ReadBytes(r, offset, 4)
	if err != nil {
		return 0, err
	}

	headerInt := binary.BigEndian.Uint32(dataBytes)

	return headerInt, nil
}

// WriteUint32 writes a Little Endian number to an io.Writer
func WriteUint32(w io.Writer, data uint32) (int, error) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, data)
	return w.Write(buf)
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

// ReadZlib reads Zlib packed data from reader until EOF
func ReadZlib(r io.Reader) ([]byte, error) {
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
	return buffer, nil
}

// ReadZlibFromBuffer unpakcs Zlib data from a slice
func ReadZlibFromBuffer(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	r := bufio.NewReader(buf)
	return ReadZlib(r)
}

// WriteZlib writes Zlib-packed data to a writer, and return
func WriteZlib(data []byte, w io.Writer) (int, error) {
	zlibEncoder, err := zlib.NewWriterLevel(w, zlib.BestCompression)
	if err != nil {
		return 0, err
	}

	n, err := zlibEncoder.Write(data)
	if err != nil {
		return n, err
	}

	err = zlibEncoder.Close()
	if err != nil {
		return 0, err
	}
	return n, err
}

// WriteZlibToBuffer packs data with Zlib to a slice
func WriteZlibToBuffer(data []byte) ([]byte, error) {
	buf := bytes.Buffer{}
	w := bufio.NewWriter(&buf)
	_, err := WriteZlib(data, w)
	if err != nil {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
