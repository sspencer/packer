package packer

import "fmt"

// Rect name and size
type Rect struct {
	Name   string
	X      int
	Y      int
	Width  int
	Height int
}

// Sprite holds describes the location of images within a sprite
type Sprite struct {
	Rect
	Used  bool
	Fit   *Sprite
	Right *Sprite
	Down  *Sprite
}

// Sprites is a slice of Sprite
type Sprites []*Sprite

// NewSprite creates sprite with additional positional parameters.
func NewSprite(name string, x, y, w, h int) *Sprite {
	s := Sprite{}
	s.Name = name
	s.X = x
	s.Y = y
	s.Width = w
	s.Height = h
	return &s
}

func (s *Sprite) String() string {
	return fmt.Sprintf("<Sprite %q @(%d, %d) %dx%d>", s.Name, s.X, s.Y, s.Width, s.Height)
}
