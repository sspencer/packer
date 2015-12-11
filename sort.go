package packer

// BlocksByWidth is used to sort images by width (then height if width is same)
type BlocksByWidth Blocks

func (s BlocksByWidth) Len() int      { return len(s) }
func (s BlocksByWidth) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s BlocksByWidth) Less(i, j int) bool {
	// Sort by (1)Width, (2)Height
	diff := s[i].Width - s[j].Width
	if diff == 0 {
		diff = s[i].Height - s[j].Height
	}

	return diff > 0
}

// BlocksByHeight is used to sort images by height (then width if height is same)
type BlocksByHeight Blocks

func (s BlocksByHeight) Len() int      { return len(s) }
func (s BlocksByHeight) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s BlocksByHeight) Less(i, j int) bool {
	// Sort by (1)Height, (2)Width
	diff := s[i].Height - s[j].Height
	if diff == 0 {
		diff = s[i].Width - s[j].Width
	}

	return diff > 0
}

// BlocksByArea is used to sort images by area (then height, then width)
type BlocksByArea Blocks

func (s BlocksByArea) Len() int      { return len(s) }
func (s BlocksByArea) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s BlocksByArea) Less(i, j int) bool {
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

// BlocksByMax is used to sort images by max dimension.
type BlocksByMax Blocks

func (s BlocksByMax) Len() int      { return len(s) }
func (s BlocksByMax) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s BlocksByMax) Less(i, j int) bool {
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
