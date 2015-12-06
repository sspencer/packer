package packer

import (
	"sort"
	"testing"
)

func getSprites() Sprites {
	sprites := Sprites{
		NewSprite("", 5, 30),
		NewSprite("", 20, 15),
		NewSprite("", 25, 10),
	}

	return sprites
}

func TestSortByWidth(t *testing.T) {
	sprites := getSprites()
	sort.Sort(SpritesByWidth(sprites))
	if sprites[0].Width != 25 {
		t.Error("Sort failed - first item")
	}

	if sprites[2].Width != 5 {
		t.Error("Sort failed - last item")
	}
}

func TestSortByHeight(t *testing.T) {
	sprites := getSprites()
	sort.Sort(SpritesByHeight(sprites))
	if sprites[0].Height != 30 {
		t.Error("Sort failed - first item")
	}

	if sprites[2].Height != 10 {
		t.Error("Sort failed - last item")
	}
}

func TestSortByArea(t *testing.T) {
	sprites := getSprites()
	sort.Sort(SpritesByArea(sprites))
	if sprites[0].Width != 20 {
		t.Error("Sort failed - first item")
	}

	if sprites[2].Width != 5 {
		t.Error("Sort failed - last item")
	}
}

func TestSortByMax(t *testing.T) {
	sprites := getSprites()
	sort.Sort(SpritesByMax(sprites))
	if sprites[0].Height != 30 {
		t.Error("Sort failed - first item")
	}

	if sprites[1].Width != 25 {
		t.Error("Sort failed - last item")
	}
}
