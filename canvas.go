package packer

import "fmt"

// Canvas contains location information for all the sprites
type Canvas struct {
	Root   *Block
	Blocks Blocks
	layout Layout
}

func (c *Canvas) String() string {
	return fmt.Sprintf("<Canvas %dx%d>", c.Root.Width, c.Root.Height)
}
