package packer

import (
	"sort"
)

// Fit packs blocks into a rectangle using 4 different sorting algorithms,
// modifying the x/y of the block to give the tighest pack in a rectangle.
func Fit(blocks Blocks) *Canvas {

	// compute area of the shapes to determine best layout below
	blockArea := 0

	// Copy each list of blocks so they can be packed most efficiently.
	// Since we're dealing with a slice of pointers, seems like we have
	// to copy each struct.  Is there a better way??
	byWidth := make(Blocks, len(blocks))
	byHeight := make(Blocks, len(blocks))
	byArea := make(Blocks, len(blocks))
	byMax := make(Blocks, len(blocks))

	for i, s := range blocks {
		byWidth[i] = &Block{Name: s.Name, Width: s.Width, Height: s.Height}
		byHeight[i] = &Block{Name: s.Name, Width: s.Width, Height: s.Height}
		byArea[i] = &Block{Name: s.Name, Width: s.Width, Height: s.Height}
		byMax[i] = &Block{Name: s.Name, Width: s.Width, Height: s.Height}

		blockArea += s.Width * s.Height
	}

	// Try to layout Blocks 4 different ways.  What we have here
	// is an "embarrassingly parallel" problem, the easiest kind
	// to perform concurrently
	ch := make(chan *Canvas)

	// Canvi ... canvases
	numCanvi := 4

	// TODO DANGER what if we're laying huge, int64 range area here
	// Should we just use int64 everywhere instead of int ??
	minWaste := 1<<31 - 1

	go layoutCanvas(ch, byWidth, LayoutByWidth)
	go layoutCanvas(ch, byHeight, LayoutByHeight)
	go layoutCanvas(ch, byArea, LayoutByArea)
	go layoutCanvas(ch, byMax, LayoutByMax)

	var bestCanvas *Canvas

	for i := 0; i < numCanvi; i++ {
		c := <-ch
		waste := (c.Root.Width * c.Root.Height) - blockArea
		//fmt.Printf("%s <%dx%d> has wasted %d pixels\n", c.layout, c.Root.Width, c.Root.Height, waste)
		if waste < minWaste {
			minWaste = waste
			bestCanvas = c
		}
	}
	//fmt.Println("USING ", bestCanvas.layout)

	return bestCanvas
}

func layoutCanvas(ch chan<- *Canvas, blocks Blocks, layout Layout) {
	switch layout {
	case LayoutByWidth:
		sort.Sort(BlocksByWidth(blocks))
	case LayoutByHeight:
		sort.Sort(BlocksByHeight(blocks))
	case LayoutByArea:
		sort.Sort(BlocksByArea(blocks))
	default:
		sort.Sort(BlocksByMax(blocks))
	}

	canvas := fit(blocks)
	canvas.layout = layout
	ch <- canvas
}

// Fit blocks in a rectangle.  Blocks must be sorted before calling Fit.  It's
// easiest to call BestFit which calls this method with 4 different sorts to
// determine the tightest packing.
func fit(blocks Blocks) *Canvas {

	root := &Block{Width: blocks[0].Width, Height: blocks[0].Height}
	canvas := &Canvas{Root: root}

	for _, block := range blocks {
		w := block.Width
		h := block.Height
		if node := canvas.findNode(canvas.Root, w, h); node != nil {
			block.fit = canvas.splitNode(node, w, h)
		} else {
			block.fit = canvas.growNode(w, h)
		}
		block.fit.Name = block.Name
	}

	return canvas.dup(blocks)
}

func (c *Canvas) dup(nodes Blocks) *Canvas {
	r := c.Root
	root := &Block{Name: "#root#", X: r.X, Y: r.Y, Width: r.Width, Height: r.Height}
	blocks := make(Blocks, len(nodes))
	for i, s := range nodes {
		blocks[i] = s.fit
	}

	return &Canvas{Root: root, Blocks: blocks}
}

func (c *Canvas) findNode(node *Block, w, h int) *Block {
	if node.used {
		if r := c.findNode(node.right, w, h); r != nil {
			return r
		}
		return c.findNode(node.down, w, h)
	} else if w <= node.Width && h <= node.Height {
		return node
	}

	return nil
}

func (c *Canvas) splitNode(node *Block, w, h int) *Block {
	node.used = true
	node.down = &Block{X: node.X, Y: node.Y + h, Width: node.Width, Height: node.Height - h}
	node.right = &Block{X: node.X + w, Y: node.Y, Width: node.Width - w, Height: h}
	node.Width = w
	node.Height = h

	return node
}

func (c *Canvas) growNode(w, h int) *Block {

	rw := c.Root.Width
	rh := c.Root.Height

	canGrowDown := (w <= rw)
	canGrowRight := (h <= rh)

	// attempt to keep square-ish by growing right when height is much greater than width
	shouldGrowRight := canGrowRight && (rh >= (rw + w))

	// attempt to keep square-ish by growing down  when width  is much greater than height
	shouldGrowDown := canGrowDown && (rw >= (rh + h))

	if shouldGrowRight {
		return c.growRight(w, h)
	} else if shouldGrowDown {
		return c.growDown(w, h)
	} else if canGrowRight {
		return c.growRight(w, h)
	} else if canGrowDown {
		return c.growDown(w, h)
	}

	// need to ensure sensible root starting size to avoid this happening
	return nil
}

// duplicate block ... is a deep copy needed here????
func dup(block *Block) *Block {
	b := *block
	return &b
}

func (c *Canvas) growRight(w, h int) *Block {
	newRoot := &Block{Width: c.Root.Width + w, Height: c.Root.Height}
	newRoot.used = true
	newRoot.down = dup(c.Root)
	newRoot.right = &Block{X: c.Root.Width, Y: 0, Width: w, Height: c.Root.Height}

	c.Root = newRoot

	if node := c.findNode(c.Root, w, h); node != nil {
		return c.splitNode(node, w, h)
	}

	return nil
}

func (c *Canvas) growDown(w, h int) *Block {
	newRoot := &Block{Width: c.Root.Width, Height: c.Root.Height + h}
	newRoot.used = true
	newRoot.down = &Block{X: 0, Y: c.Root.Height, Width: c.Root.Width, Height: h}
	newRoot.right = dup(c.Root)

	c.Root = newRoot

	if node := c.findNode(c.Root, w, h); node != nil {
		return c.splitNode(node, w, h)
	}

	return nil
}
