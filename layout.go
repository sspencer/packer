package packer

import "sort"

// Canvas contains location information for all the sprites
type Canvas struct {
	Root    *Sprite
	Sprites Sprites
}

// BestFit packs sprites into a rectangle using 4 different sorting algorithms,
// modifying the x/y of the sprite to give the tighest pack in a rectangle.
func BestFit(sprites Sprites) *Canvas {
	sort.Sort(SpritesByHeight(sprites))
	return fit(sprites)
}

// Fit sprites
func fit(sprites Sprites) *Canvas {

	root := NewSprite("", 0, 0, sprites[0].Width, sprites[0].Height)
	canvas := &Canvas{root, nil}

	for _, sprite := range sprites {
		w := sprite.Width
		h := sprite.Height
		if node := canvas.findNode(canvas.Root, w, h); node != nil {
			sprite.Fit = canvas.splitNode(node, w, h)
		} else {
			sprite.Fit = canvas.growNode(w, h)
		}
		sprite.Fit.Name = sprite.Name
	}

	return canvas.dup(sprites)
}

func (c *Canvas) dup(nodes Sprites) *Canvas {
	r := c.Root
	root := NewSprite("#root#", r.X, r.Y, r.Width, r.Height)
	sprites := make(Sprites, len(nodes))
	for i, s := range nodes {
		sprites[i] = s.Fit
	}

	return &Canvas{root, sprites}
}

func (c *Canvas) findNode(node *Sprite, w, h int) *Sprite {
	if node.Used {
		if r := c.findNode(node.Right, w, h); r != nil {
			return r
		}
		return c.findNode(node.Down, w, h)
	} else if w <= node.Width && h <= node.Height {
		return node
	}

	return nil
}

func (c *Canvas) splitNode(node *Sprite, w, h int) *Sprite {
	node.Used = true
	node.Down = NewSprite("", node.X, node.Y+h, node.Width, node.Height-h)
	node.Right = NewSprite("", node.X+w, node.Y, node.Width-w, h)
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
	newRoot := NewSprite("", 0, 0, c.Root.Width+w, c.Root.Height)
	newRoot.Used = true
	newRoot.Down = dup(c.Root)
	newRoot.Right = NewSprite("", c.Root.Width, 0, w, c.Root.Height)

	c.Root = newRoot

	if node := c.findNode(c.Root, w, h); node != nil {
		return c.splitNode(node, w, h)
	}

	return nil
}

func (c *Canvas) growDown(w, h int) *Sprite {
	newRoot := NewSprite("", 0, 0, c.Root.Width, c.Root.Height+h)
	newRoot.Used = true
	newRoot.Down = NewSprite("", 0, c.Root.Height, c.Root.Width, h)
	newRoot.Right = dup(c.Root)

	c.Root = newRoot

	if node := c.findNode(c.Root, w, h); node != nil {
		return c.splitNode(node, w, h)
	}

	return nil
}
