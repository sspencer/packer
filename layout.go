package packer

// Layout is used to sort sprites in different ways to achieve differnt packs.
//go:generate stringer -type=Layout
type Layout int

const (
	// LayoutByWidth and others describe sort to use when laying out image
	LayoutByWidth Layout = iota
	LayoutByHeight
	LayoutByArea
	LayoutByMax
)
