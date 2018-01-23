package packer

// Layout is used to sort sprites in different ways to achieve different packs.
//go:generate stringer -type=Layout
type Layout int

const (
	// LayoutByWidth sorts objects by width
	LayoutByWidth Layout = iota
	// LayoutByHeight sorts objects by height
	LayoutByHeight
	// LayoutByArea sorts objects by total area
	LayoutByArea
	// LayoutByMax sorts objects by max dimension, width or height
	LayoutByMax
)
