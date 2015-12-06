package packer

import "fmt"

// Canvas contains location information for all the sprites
type Canvas struct {
	Root    *Sprite
	Sprites Sprites
	layout  Layout
}

// Sprite holds describes the location of each and every image within the canvas.
type Sprite struct {
	Name   string
	X      int
	Y      int
	Width  int
	Height int
	used   bool
	fit    *Sprite
	right  *Sprite
	down   *Sprite
}

// Sprites is a slice of Sprite
type Sprites []*Sprite

// NewSprite records the size of an image to be laid out.
func NewSprite(name string, w, h int) *Sprite {
	return &Sprite{Name: name, Width: w, Height: h}
}

func newPosSprite(name string, x, y, w, h int) *Sprite {
	return &Sprite{Name: name, X: x, Y: y, Width: w, Height: h}
}

func (s *Sprite) String() string {
	return fmt.Sprintf("<%q @(%d, %d) %dx%d>", s.Name, s.X, s.Y, s.Width, s.Height)
}
