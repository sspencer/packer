package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sspencer/packer"
)

type blockdef []string

var (
	pastels = []color.Color{
		color.RGBA{255, 247, 165, 255}, // yellow
		color.RGBA{255, 165, 224, 255}, // red
		color.RGBA{165, 179, 255, 255}, // blue
		color.RGBA{191, 255, 165, 255}, // green
		color.RGBA{255, 203, 165, 255}, // orange
	}

	examples = map[string]blockdef{
		"simple": {
			"500x200",
			"250x200",
			"50x50x20",
		},

		"square": {
			"50x50x100",
		},

		"tall": {
			"50x400x2",
			"50x300x5",
			"50x200x10",
			"50x100x20",
			"50x50x40",
		},

		"wide": {
			"400x50x2",
			"300x50x5",
			"200x50x10",
			"100x50x20",
			"50x50x40",
		},

		"tallAndWide": {
			"400x100",
			"100x400",
			"400x100",
			"100x400",
			"400x100",
			"100x400",
		},

		"powersOf2": {
			"2x2x256",
			"4x4x128",
			"8x8x64",
			"16x16x32",
			"32x32x16",
			"64x64x8",
			"128x128x4",
			"256x256x2",
		},

		"oddAndEven": {
			"50x50x20",
			"47x31x20",
			"23x17x20",
			"109x42x20",
			"42x109x20",
			"17x33x20",
		},

		"complex": {
			"100x100x3",
			"60x60x3",
			"50x20x20",
			"20x50x20",
			"250x250",
			"250x100",
			"100x250",
			"400x80",
			"80x400",
			"10x10x100",
			"5x5x500",
		},
	}
)

// parse blockdef "(width)x(height)x(num)" or "(width)x(height)"
func parse(block string) (int, int, int) {
	s := strings.Split(block, "x")
	w, _ := strconv.Atoi(s[0])
	h, _ := strconv.Atoi(s[1])
	n := 1
	if len(s) > 2 {
		n, _ = strconv.Atoi(s[2])
	}

	return w, h, n
}

func pack(name string, blocks []string) {
	sprites := make(packer.Sprites, 0)
	for _, block := range blocks {
		w, h, n := parse(block)
		for i := 0; i < n; i++ {
			sprites = append(sprites, packer.NewSprite("", 0, 0, w, h))
		}
	}

	canvas := packer.BestFit(sprites)

	render(name, canvas)
}

func render(name string, c *packer.Canvas) {

	img := image.NewRGBA(image.Rect(0, 0, c.Root.Width, c.Root.Height))

	for n, s := range c.Sprites {
		rect := image.Rect(s.X, s.Y, s.X+s.Width, s.Y+s.Height)
		col := pastels[n%len(pastels)]
		draw.Draw(img, rect, &image.Uniform{col}, image.ZP, draw.Src)
	}

	w, _ := os.Create(name + ".png")
	defer w.Close()
	png.Encode(w, img)

}

// open all the images in Preview app on OSX.
func show(names ...string) {
	command := "open"
	arg1 := "-a"
	arg2 := "/Applications/Preview.app"
	var args []string
	args = append(args, arg1)
	args = append(args, arg2)
	args = append(args, names...)
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	names := make([]string, len(examples))
	i := 0
	for name, blocks := range examples {
		pack(name, blocks)
		names[i] = fmt.Sprintf("%s.png", name)
	}

	show(names...)
}
