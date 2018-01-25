package packer

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

//import

const retinaTag = "@2x"
const hoverTrigger = "_hover"

// Config is the configuration structure needed to build sprites.
type Config struct {
	Base64     bool
	Retina     bool
	Output     string
	ImgURL     string
	Format     string
	Name       string
	Prefix     string
	Margin     int
	Background string
}

// Sprite is result containing image(s) and stylesheet computed from packing blocks into a canvas.
type Sprite struct {
	Image       *image.RGBA
	RetinaImage *image.RGBA
	Stylesheet  string
}

// CreateSprite creates a sprite and stylesheet for the config data.
func (c *Config) CreateSprite(files []string) (*Sprite, error) {

	// Validate configuration parameters
	if err := c.validate(); err != nil {
		return nil, err
	}

	// Convert file paths into image data
	images, err := c.getImages(files)
	if err != nil {
		return nil, err
	}

	var retinaImage *image.RGBA
	retinaImage = nil

	var normalImage *image.RGBA
	normalImage = nil

	//	var
	if c.Retina {
		retinaImage = c.createImage(images)
		// resize images in image map
		resized := make(map[string]*image.Image)
		for name, img := range images {
			w := (*img).Bounds().Max.X - (*img).Bounds().Min.X
			m := resize.Resize(uint(w/2), 0, *img, resize.Lanczos3)
			resized[name] = &m
		}

		images = resized
	}

	normalImage = c.createImage(images)

	return &Sprite{Image: normalImage, RetinaImage: retinaImage, Stylesheet: "/* TBD - run sprite template */"}, nil

}

func (c *Config) createImage(images map[string]*image.Image) *image.RGBA {
	// create proxy 'block' for each image
	blocks := make(Blocks, len(images))
	i := 0
	for name, image := range images {
		img := *image
		w := img.Bounds().Max.X - img.Bounds().Min.X
		h := img.Bounds().Max.Y - img.Bounds().Min.Y
		blocks[i] = &Block{Name: name, Width: w, Height: h}

		i++
	}

	canvas := Fit(blocks)

	rgba := image.NewRGBA(image.Rect(0, 0, canvas.Root.Width, canvas.Root.Height))
	draw.Draw(rgba, rgba.Bounds(), ColorToUniform(c.Background), image.ZP, draw.Src)

	for _, b := range canvas.Blocks {
		img, ok := images[b.Name]
		if ok {
			src := *img
			dp := image.Pt(b.X, b.Y)
			r := image.Rectangle{dp, dp.Add(image.Pt(b.Width, b.Height))}
			draw.Draw(rgba, r, src, image.ZP, draw.Src)
		}
	}

	return rgba
}

func (c *Config) String() string {
	return fmt.Sprintf("CONFIG: base64=%t retina=%t output=%s imgurl=%s format=%s name=%s prefix=%s bg=%s margin=%d>",
		c.Base64,
		c.Retina,
		c.Output,
		c.ImgURL,
		c.Format,
		c.Name,
		c.Prefix,
		c.Background,
		c.Margin)
}

// validate config parameters
func (c *Config) validate() error {
	if c.Output == "" {
		c.Output = "."
	}

	if c.Margin < 0 || c.Margin > 100 {
		return fmt.Errorf("margin must have a value between 0 and 100")
	}

	if c.Format != "jpg" && c.Format != "png" {
		return fmt.Errorf("illegal option %q for format (only 'png' or 'jpg' allowed)", c.Format)
	}

	if c.Format == "jpg" && strings.ToLower(c.Background) == "transparent" {
		c.Background = "white"
	}

	return nil
}

// Save saves stylesheet and image(s) to disk.
func (c *Config) Save(sprite *Sprite) error {
	fn, err := filepath.Abs(path.Join(c.Output, fmt.Sprintf("%s.%s", c.Name, c.Format)))
	if err != nil {
		return err
	}

	err = c.saveImage(fn, sprite.Image)
	if err != nil {
		return err
	}

	if sprite.RetinaImage != nil {
		fn, err = filepath.Abs(path.Join(c.Output, fmt.Sprintf("%s%s.%s", c.Name, retinaTag, c.Format)))
		if err != nil {
			return err
		}

		err = c.saveImage(fn, sprite.RetinaImage)
	}

	return err
}

// save given image to disk
func (c *Config) saveImage(fn string, img *image.RGBA) error {
	w, _ := os.Create(fn)
	defer w.Close()
	if c.Format == "png" {
		return png.Encode(w, img)
	}

	return jpeg.Encode(w, img, nil)
}
