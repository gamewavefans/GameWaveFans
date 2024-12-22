package zbm

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"

	"github.com/namgo/GameWaveFans/pkg/common"
)

func getPixelValue(value uint16) (uint8, uint8, uint8, uint8) {
	// most images store pixels in big endian, 4633

	// values are shifted and multiplied for better color mapping:
	// https://reshax.com/topic/265-gamewave-texture-zbm-format/#comment-856
	cr := uint8(value&0x7) << 5          // 3 bits
	cb := uint8((value>>3)&0x7) << 5     // 3 bits
	y := uint8((value>>6)&0x3F) << 2     // 6 bits
	alpha := uint8((value>>12)&0xF) * 17 // 4 bits

	return cr, cb, y, alpha
}

func clampUint8(value int32) uint8 {
	return uint8(math.Max(math.Min(float64(value), 255), 0))
}

// Decode reads zbm file and returns image.Image
func Decode(r io.Reader) (image.Image, error) {
	var c config
	c.r = r
	if err := c.decodeConfig(); err != nil {
		return nil, err
	}

	buffer, err := common.ReadZlib(c.r)
	if err != nil {
		return nil, err
	}
	if len(buffer) != int(c.sizeUnpacked) {
		return nil, FormatError(fmt.Sprintf("unpacked size mismatch: got %d, expected %d\n", len(buffer), c.sizeUnpacked))
	}

	pixelBuffer := make([]uint16, c.width*c.height)

	// swap every two pixels, endianness changes a bit
	for i := 0; i < len(pixelBuffer)-1; i += 2 {
		pixelBuffer[i+1] = binary.BigEndian.Uint16(buffer[i*2 : (i*2)+2])
		pixelBuffer[i] = binary.BigEndian.Uint16(buffer[(i+1)*2 : ((i+1)*2)+2])
	}

	img := image.NewNRGBA(image.Rect(0, 0, int(c.width), int(c.height)))

	for i := 0; i < int(c.width*c.height); i++ {
		cr, cb, y, a := getPixelValue(pixelBuffer[i])

		// convert YCrCb colors to RGB with clamping to avoid overflows when converting to uint8
		cb1 := int32(cb) - 128
		cr1 := int32(cr) - 128
		r := clampUint8(int32(y) + (int32(45*cr1) / 32))
		g := clampUint8(int32(y) - (int32(11*cb1+23*cr1) / 32))
		b := clampUint8(int32(y) + (int32(113*cb1) / 64))

		img.Pix[4*i] = r
		img.Pix[(4*i)+1] = g
		img.Pix[(4*i)+2] = b
		img.Pix[(4*i)+3] = a
	}
	return img, nil
}

func (c *config) decodeConfig() error {
	var err error
	buf := make([]byte, 12*4)
	if _, err = io.ReadFull(c.r, buf); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return err
	}

	c.unknown1 = binary.LittleEndian.Uint32(buf[0x0:0x4])
	c.unknown2 = binary.LittleEndian.Uint32(buf[0x4:0x8])
	c.unknown3 = binary.LittleEndian.Uint32(buf[0x8:0xC])
	c.unknown4 = binary.LittleEndian.Uint32(buf[0xC:0x10])
	c.width = binary.LittleEndian.Uint32(buf[0x10:0x14])
	c.height = binary.LittleEndian.Uint32(buf[0x14:0x18])
	c.unknown5 = binary.LittleEndian.Uint32(buf[0x18:0x1C])
	c.unknown6 = binary.LittleEndian.Uint32(buf[0x1C:0x20])
	c.unknown7 = binary.LittleEndian.Uint32(buf[0x20:0x24])
	c.sizePacked = binary.LittleEndian.Uint32(buf[0x24:0x28])
	c.sizeUnpacked = binary.LittleEndian.Uint32(buf[0x28:0x2C])
	c.unknown8 = binary.LittleEndian.Uint32(buf[0x2C:0x30])

	if c.width == 0 || c.height == 0 {
		return FormatError(fmt.Sprintf("unsupported size: %dx%d\n", c.width, c.height))
	}

	return nil
}

// DecodeConfig returns the color model and dimensions of an image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	var c config
	c.r = r
	if err := c.decodeConfig(); err != nil {
		return image.Config{}, err
	}
	colorModel := color.RGBAModel

	return image.Config{
		ColorModel: colorModel,
		Width:      int(c.width),
		Height:     int(c.height),
	}, nil
}

// RegisterFormat registers format to be used by image.DecodeConfig
// in normal circumstances we'd call image.RegisterFormat in init, but since there's no magic bytes to detect this file format, Go thinks all files are in zbm format
func RegisterFormat() {
	image.RegisterFormat(FormatName, textureHeader, Decode, DecodeConfig)
}
