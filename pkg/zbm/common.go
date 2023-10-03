// Package zbm helps interfacing with zbm files, un unpack and repack them
package zbm

// A FormatError reports that the input is not a valid Gamewave texture.
type FormatError string

func (e FormatError) Error() string { return "gamewave zbm error:" + string(e) }

// FormatName is the name of the registered texture format
const FormatName = "Gamewave texture"
const textureHeader = ""
