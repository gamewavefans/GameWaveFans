package zbm

import (
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
)

type config struct {
	r            io.Reader
	unknown1     uint32
	unknown2     uint32
	unknown3     uint32
	unknown4     uint32
	width        uint32
	height       uint32
	unknown5     uint32
	unknown6     uint32
	unknown7     uint32
	sizePacked   uint32
	sizeUnpacked uint32
	unknown8     uint32
}

//	func getPixelValue(value uint16) (r, g, b, a uint8) {
//		r = uint8(math.Round(float64(value>>10&31) * (255.0 / 31.0)))
//		g = uint8(math.Round(float64(value>>5&31) * (255.0 / 31.0)))
//		b = uint8(math.Round(float64(value&31) * (255.0 / 31.0)))
//		// TODO parse transparency
//		a = 0xFF
//		return r, g, b, a
//	}
func getPixelValue(value uint16) (r, g, b, a uint8) {
	// TODO colors are broken
	rRaw := value >> 10 & 31
	r = uint8(math.Round(float64(rRaw)) * (255.0 / 31.0))
	gRaw := value >> 5 & 31
	g = uint8(math.Round(float64(gRaw) * (255.0 / 31.0)))
	bRaw := value >> 0 & 31
	b = uint8(math.Round(float64(bRaw) * (255.0 / 31.0)))

	a = uint8(math.Round(float64(value>>15&1) * (255.0)))
	return r, g, b, a
}

// Decode reads zbm file and returns image.Image
func Decode(r io.Reader) (image.Image, error) {
	var c config
	c.r = r
	if err := c.decodeConfig(); err != nil {
		return nil, err
	}

	zlibDecoder, err := zlib.NewReader(c.r)
	if err != nil {
		return nil, err
	}
	buffer, err := io.ReadAll(zlibDecoder)
	if err != nil {
		return nil, err
	}
	if len(buffer) != int(c.sizeUnpacked) {
		return nil, FormatError(fmt.Sprintf("unpacked size mismatch: got %d, expected %d\n", len(buffer), c.sizeUnpacked))
	}
	err = zlibDecoder.Close()
	if err != nil {
		return nil, err
	}

	pixelBuffer := make([]uint16, c.width*c.height)
	for i := range pixelBuffer {
		pixelBuffer[i] = binary.BigEndian.Uint16(buffer[i*2 : (i*2)+2])
	}

	img := image.NewRGBA(image.Rect(0, 0, int(c.width), int(c.height)))
	//pixels := img.Pix[0 : int(c.width)*int(c.height)*4]
	for i := 0; i < int(c.width*c.height); i++ {
		r, g, b, a := getPixelValue(pixelBuffer[i])
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

func init() {
	image.RegisterFormat(FormatName, textureHeader, Decode, DecodeConfig)
}
