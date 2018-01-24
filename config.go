package packer

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/nfnt/resize"
)

//import

//const retinaTag = "@2x"
const hoverTrigger = "_hover"

// Config is the configuration structure needed to build sprites.
type Config struct {
	Base64     bool
	Retina     bool
	CSSPath    string
	ImgPath    string
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
	return fmt.Sprintf("CONFIG: base64=%t retina=%t csspath=%s imgpath=%s imgurl=%s format=%s name=%s prefix=%s bg=%s margin=%d>",
		c.Base64,
		c.Retina,
		c.CSSPath,
		c.ImgPath,
		c.ImgURL,
		c.Format,
		c.Name,
		c.Prefix,
		c.Background,
		c.Margin)
}

func (c *Config) validate() error {
	if c.ImgPath == "" {
		c.Base64 = true
	}

	if c.Margin < 0 {
		c.Margin = 0
	}

	if c.Margin > 100 {
		c.Margin = 100
	}

	if c.Format != "jpg" && c.Format != "png" {
		return fmt.Errorf("illegal option %q for format (only 'png' or 'jpg' allowed)", c.Format)
	}

	return nil
}

// SaveStylesheet saves the css style sheet to the CSSPath specified in the config Sprite value.
func (c *Config) SaveStylesheet() error {

	fmt.Println("SAVE STYLESHEET")
	var p string
	var err error
	stylesheet := path.Join(c.CSSPath, fmt.Sprintf("%s.css", c.Name))

	if p, err = filepath.Abs(stylesheet); err != nil {
		return fmt.Errorf("ERROR: Can not write to file %q: %s", p, err)
	}

	if err = ioutil.WriteFile(p, []byte(stylesheet), 0644); err != nil {
		return fmt.Errorf("ERROR: Failed writing to file %q: %s", p, err)
	}

	return nil
}

// SaveSprite saves the sprite image to file specified in the config Sprite value.
func (c *Config) SaveSprite(img *image.RGBA) error {
	fmt.Println("SAVE SPRITE")
	var p string
	var err error
	if p, err = filepath.Abs(c.ImgPath); err != nil {
		return fmt.Errorf("ERROR: Can not write to file %q: %s", p, err)
	}

	w, _ := os.Create(p)
	defer w.Close()
	png.Encode(w, img)
	return nil
}

// SaveRetinaSprite saves the retina sprite image, if retina config value was set..
func (c *Config) SaveRetinaSprite(img *image.RGBA) error {
	/*
		if !c.Retina {
			return nil
		}
		fmt.Println("SAVE RETINA SPRITE")

		fn := path.Base(c.Sprite)
		ext := path.Ext(c.Sprite)
		fn = path.Dir(c.Sprite) + fn[:len(fn)-len(ext)] + retinaTag + ext
		fmt.Printf("  Retina filename: %q", fn)

		var err error
		if fn, err = filepath.Abs(fn); err != nil {
			return fmt.Errorf("ERROR: Can not write to file %q: %s", fn, err)
		}

		w, _ := os.Create(fn)
		defer w.Close()
		png.Encode(w, img)
	*/
	return nil
}
