package packer

import (
	"image"
	"image/color"
	"strconv"
	"strings"
)

// hexToRGBA converts a hex rgba string (either #abcd or #aabbccdd) to 4 bytes.
func hexToRGBA(h string) (uint8, uint8, uint8, uint8) {
	if len(h) > 0 && h[0] == '#' {
		h = h[1:]
	}

	if len(h) == 3 {
		h = h + "f"
	}

	if len(h) == 6 {
		h = h + "ff"
	}

	if len(h) == 4 {
		h = h[:1] + h[:1] + h[1:2] + h[1:2] + h[2:3] + h[2:3] + h[3:4] + h[3:4]
	}

	if len(h) == 8 {
		if c, err := strconv.ParseUint(string(h), 16, 32); err == nil {
			return uint8(c >> 24), uint8(c >> 16), uint8(c >> 8), uint8(c & 0xFF)
		}
	}

	return 0, 0, 0, 0
}

// 	colorToUniform converts a string color name (only Black, White, Transparent) or
// hex number (#rgba or #rrggbbaa) to a uniform color.
func colorToUniform(hex string) *image.Uniform {
	lc := strings.ToLower(hex)
	switch lc {
	case "transparent":
		return image.Transparent
	case "black":
		return image.Black
	case "white":
		return image.White
	}

	r, g, b, a := hexToRGBA(hex)

	if r+g+b+a == 0 || a == 0 {
		return image.Transparent
	}

	return &image.Uniform{color.RGBA{r, g, b, a}}
}
