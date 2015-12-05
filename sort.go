package packer

// SpritesByWidth is used to sort images by width (then height if width is same)
type SpritesByWidth Sprites

func (s SpritesByWidth) Len() int      { return len(s) }
func (s SpritesByWidth) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SpritesByWidth) Less(i, j int) bool {
	// Sort by (1)Width, (2)Height
	diff := s[i].Width - s[j].Width
	if diff == 0 {
		diff = s[i].Height - s[j].Height
	}

	return diff > 0
}

// SpritesByHeight is used to sort images by height (then width if height is same)
type SpritesByHeight Sprites

func (s SpritesByHeight) Len() int      { return len(s) }
func (s SpritesByHeight) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SpritesByHeight) Less(i, j int) bool {
	// Sort by (1)Height, (2)Width
	diff := s[i].Height - s[j].Height
	if diff == 0 {
		diff = s[i].Width - s[j].Width
	}

	return diff > 0
}

// SpritesByArea is used to sort images by area (then height, then width)
type SpritesByArea Sprites

func (s SpritesByArea) Len() int      { return len(s) }
func (s SpritesByArea) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SpritesByArea) Less(i, j int) bool {
	// Sort by (1)Area, (2)Height, (3)Width
	diff := (s[i].Width * s[i].Height) - (s[j].Width * s[j].Height)

	if diff == 0 {
		diff = s[i].Height - s[j].Height
	}

	if diff == 0 {
		diff = s[i].Width - s[j].Width
	}

	return diff > 0
}

// SpritesByMax is used to sort images by max dimension.
type SpritesByMax Sprites

func (s SpritesByMax) Len() int      { return len(s) }
func (s SpritesByMax) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SpritesByMax) Less(i, j int) bool {
	diff := max(s[i].Width, s[i].Height) - max(s[j].Width, s[j].Height)
	if diff == 0 {
		diff = min(s[i].Width, s[i].Height) - min(s[j].Width, s[j].Height)
	}

	if diff == 0 {
		diff = s[i].Height - s[j].Height
	}

	if diff == 0 {
		diff = s[i].Width - s[j].Width
	}

	return diff > 0
}
