package packer

import (
	"fmt"
	"sort"
)

// BestFit packs sprites into a rectangle using 4 different sorting algorithms,
// modifying the x/y of the sprite to give the tighest pack in a rectangle.
func BestFit(sprites Sprites) *Canvas {

	// compute area of the shapes to determine best layout below
	spriteArea := 0

	// Copy each list of sprites so they can be packed most efficiently.
	// Since we're dealing with a slice of pointers, seems like we have
	// to copy each struct.  Is there a better way??
	byWidth := make(Sprites, len(sprites))
	byHeight := make(Sprites, len(sprites))
	byArea := make(Sprites, len(sprites))
	byMax := make(Sprites, len(sprites))

	for i, s := range sprites {
		byWidth[i] = NewSprite(s.Name, s.Width, s.Height)
		byHeight[i] = NewSprite(s.Name, s.Width, s.Height)
		byArea[i] = NewSprite(s.Name, s.Width, s.Height)
		byMax[i] = NewSprite(s.Name, s.Width, s.Height)

		spriteArea += s.Width * s.Height
	}

	// Try to layout Sprites 4 different ways.  What we have here
	// is an "embarrassingly parallel" problem, the easiest kind
	// to perform concurrently
	ch := make(chan *Canvas)

	go doit(ch, byWidth, LayoutByWidth)
	go doit(ch, byHeight, LayoutByHeight)
	go doit(ch, byArea, LayoutByArea)
	go doit(ch, byMax, LayoutByMax)

	// Canvi ... canvases
	numCanvi := 4

	// TODO DANGER what if we're laying huge, int64 range area here
	// Should we just use int64 everywhere instead of int ??
	minWaste := 1<<31 - 1

	var bestCanvas *Canvas

	for i := 0; i < numCanvi; i++ {
		c := <-ch
		waste := (c.Root.Width * c.Root.Height) - spriteArea
		fmt.Printf("%s <%dx%d> has wasted %d pixels\n", c.layout, c.Root.Width, c.Root.Height, waste)
		if waste < minWaste {
			minWaste = waste
			bestCanvas = c
		}
	}
	fmt.Println(">>>> RETURNING ", bestCanvas.layout)

	return bestCanvas
}

func doit(ch chan<- *Canvas, sprites Sprites, layout Layout) {
	switch layout {
	case LayoutByWidth:
		sort.Sort(SpritesByWidth(sprites))
	case LayoutByHeight:
		sort.Sort(SpritesByHeight(sprites))
	case LayoutByArea:
		sort.Sort(SpritesByArea(sprites))
	default:
		sort.Sort(SpritesByMax(sprites))
	}

	canvas := Fit(sprites)
	canvas.layout = layout
	ch <- canvas
}

// Fit sprites in a rectangle.  Sprites must be sorted before calling Fit.  It's
// easiest to call BestFit which calls this method with 4 different sorts to
// determine the tightest packing.
func Fit(sprites Sprites) *Canvas {

	root := newPosSprite("", 0, 0, sprites[0].Width, sprites[0].Height)
	canvas := &Canvas{Root: root}

	for _, sprite := range sprites {
		w := sprite.Width
		h := sprite.Height
		if node := canvas.findNode(canvas.Root, w, h); node != nil {
			sprite.fit = canvas.splitNode(node, w, h)
		} else {
			sprite.fit = canvas.growNode(w, h)
		}
		sprite.fit.Name = sprite.Name
	}

	return canvas.dup(sprites)
}

func (c *Canvas) dup(nodes Sprites) *Canvas {
	r := c.Root
	root := newPosSprite("#root#", r.X, r.Y, r.Width, r.Height)
	sprites := make(Sprites, len(nodes))
	for i, s := range nodes {
		sprites[i] = s.fit
	}

	return &Canvas{Root: root, Sprites: sprites}
}

func (c *Canvas) findNode(node *Sprite, w, h int) *Sprite {
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

func (c *Canvas) splitNode(node *Sprite, w, h int) *Sprite {
	node.used = true
	node.down = newPosSprite("", node.X, node.Y+h, node.Width, node.Height-h)
	node.right = newPosSprite("", node.X+w, node.Y, node.Width-w, h)
	node.Width = w
	node.Height = h

	return node
}

func (c *Canvas) growNode(w, h int) *Sprite {

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

// duplicate sprite ... is a deep copy needed here????
func dup(s *Sprite) *Sprite {
	newSprite := *s
	return &newSprite
}

func (c *Canvas) growRight(w, h int) *Sprite {
	newRoot := newPosSprite("", 0, 0, c.Root.Width+w, c.Root.Height)
	newRoot.used = true
	newRoot.down = dup(c.Root)
	newRoot.right = newPosSprite("", c.Root.Width, 0, w, c.Root.Height)

	c.Root = newRoot

	if node := c.findNode(c.Root, w, h); node != nil {
		return c.splitNode(node, w, h)
	}

	return nil
}

func (c *Canvas) growDown(w, h int) *Sprite {
	newRoot := newPosSprite("", 0, 0, c.Root.Width, c.Root.Height+h)
	newRoot.used = true
	newRoot.down = newPosSprite("", 0, c.Root.Height, c.Root.Width, h)
	newRoot.right = dup(c.Root)

	c.Root = newRoot

	if node := c.findNode(c.Root, w, h); node != nil {
		return c.splitNode(node, w, h)
	}

	return nil
}
