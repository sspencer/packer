package packer

import "fmt"

// Block holds describes the location of each and every image within the canvas.
type Block struct {
	Name   string
	X      int
	Y      int
	Width  int
	Height int
	used   bool
	fit    *Block
	right  *Block
	down   *Block
}

// Blocks is a slice of Blocks
type Blocks []*Block

func (s *Block) String() string {
	return fmt.Sprintf("<%q @(%d, %d) %dx%d>", s.Name, s.X, s.Y, s.Width, s.Height)
}
