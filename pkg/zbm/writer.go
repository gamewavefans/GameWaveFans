package zbm

import (
	"encoding/binary"
	"image"
	"image/color"
	"io"
	"math"

	"github.com/namgo/GameWaveFans/pkg/common"
)

/*
Currently, writer only outputs images in 3364CrCbYA image format, which all official games use
*/

func convertColorToCrCbYA(c color.Color) uint16 {
	r, g, b, a := c.RGBA()
	col := uint16(0)

	r8 := uint8((float64(r) / 65536.0) * 255.0)
	g8 := uint8((float64(g) / 65536.0) * 255.0)
	b8 := uint8((float64(b) / 65536.0) * 255.0)

	y, cb, cr := color.RGBToYCbCr(r8, g8, b8)

	y6 := uint16(math.Round((float64(y) / 255.0) * 63.0))
	cr3 := uint16(math.Round((float64(cr) / 255.0) * 7.0))
	cb3 := uint16(math.Round((float64(cb) / 255.0) * 7.0))
	a4 := uint16(math.Round((float64(a) / 65536.0) * 15.0))

	col = a4<<12 | y6<<6 | cb3<<3 | cr3
	return col
}

func convertImage(m image.Image) []byte {
	pixelBuffer := make([]uint16, m.Bounds().Dx()*m.Bounds().Dy())

	for y := 0; y < m.Bounds().Dy(); y++ {
		for x := 0; x < m.Bounds().Dx(); x++ {
			pixelBuffer[x+(y*m.Bounds().Dx())] = convertColorToCrCbYA(m.At(x, y))
		}
	}

	data := make([]byte, len(pixelBuffer)*2)

	// swap every two pixels, endianness changes a bit
	for i := 0; i < len(pixelBuffer)-1; i += 2 {
		pixelData := make([]byte, 2)
		pixelData2 := make([]byte, 2)
		binary.BigEndian.PutUint16(pixelData, pixelBuffer[i])
		binary.BigEndian.PutUint16(pixelData2, pixelBuffer[i+1])

		data[i*2] = pixelData2[0]
		data[i*2+1] = pixelData2[1]
		data[i*2+2] = pixelData[0]
		data[i*2+3] = pixelData[1]
	}
	return data
}

// write header of zbm file, before packed data is written
func writeHeader(w io.Writer, im image.Image, unpackedSize, packedSize uint32) error {
	if _, err := common.WriteUint32(w, 1); err != nil {
		return err
	}

	// TEXTURE_OSD
	if _, err := common.WriteUint32(w, 1); err != nil {
		return err
	}
	// 3364 format
	if _, err := common.WriteUint32(w, 4); err != nil {
		return err
	}
	// BPP
	if _, err := common.WriteUint32(w, 2); err != nil {
		return err
	}
	// width
	if _, err := common.WriteUint32(w, uint32(im.Bounds().Dx())); err != nil {
		return err
	}
	// height
	if _, err := common.WriteUint32(w, uint32(im.Bounds().Dy())); err != nil {
		return err
	}

	if _, err := common.WriteUint32(w, 0); err != nil {
		return err
	}
	if _, err := common.WriteUint32(w, 0); err != nil {
		return err
	}
	if _, err := common.WriteUint32(w, 1); err != nil {
		return err
	}

	if _, err := common.WriteUint32(w, packedSize); err != nil {
		return err
	}
	if _, err := common.WriteUint32(w, unpackedSize); err != nil {
		return err
	}

	_, err := common.WriteUint32(w, 0)
	return err
}

// Encode encodes image.Image to .zbm file
func Encode(w io.Writer, m image.Image) error {
	//convert data
	convertedData := convertImage(m)

	// pack data
	packedData, err := common.WriteZlibToBuffer(convertedData)
	if err != nil {
		return err
	}

	// write header
	err = writeHeader(w, m, uint32(len(convertedData)), uint32(len(packedData)))
	if err != nil {
		return err
	}

	// write data
	_, err = w.Write(packedData)
	return err
}
