// Package zbm helps interfacing with zbm files, un unpack and repack them
package zbm

import "io"

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

// A FormatError reports that the input is not a valid Gamewave texture.
type FormatError string

func (e FormatError) Error() string { return "gamewave zbm error:" + string(e) }

// FormatName is the name of the registered texture format
const FormatName = "zbm"
const textureHeader = ""
