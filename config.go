package packer

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alecthomas/template"
	"github.com/nfnt/resize"
)

//import

const retinaTag = "@2x"
const hoverTrigger = "_hover"
const hoverCSS = ":hover"

// Config is the configuration structure needed to build sprites.
type Config struct {
	Base64     bool
	Retina     bool
	HTML       bool
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

type spriteimage struct {
	Name   string
	Hover  string
	X      int
	Y      int
	Width  int
	Height int
}

type stylesheet struct {
	CSSPath string
	ImgPath string
	Retina  bool
	ImgURL  string
	Format  string
	Prefix  string
	Name    string
	Images  []spriteimage
}

// CreateSprite creates a sprite and stylesheet for the config data.
func (c *Config) CreateSprite(files []string) (*Sprite, error) {

	if len(files) == 0 {
		return nil, fmt.Errorf("One or more images must be specified")
	}
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

	//	var
	if c.Retina {
		_, retinaImage = c.createImage(images)
		// resize images in image map
		resized := make(map[string]*image.Image)
		for name, img := range images {
			w := (*img).Bounds().Max.X - (*img).Bounds().Min.X
			m := resize.Resize(uint(w/2), 0, *img, resize.Lanczos3)
			resized[name] = &m
		}

		images = resized
	}

	css, normalImage := c.createImage(images)

	return &Sprite{Image: normalImage, RetinaImage: retinaImage, Stylesheet: css}, nil
}

func (c *Config) createImage(images map[string]*image.Image) (string, *image.RGBA) {
	// create proxy 'block' for each image
	blocks := make(Blocks, len(images))
	i := 0
	for name, image := range images {
		img := *image
		w := img.Bounds().Max.X - img.Bounds().Min.X + c.Margin*2
		h := img.Bounds().Max.Y - img.Bounds().Min.Y + c.Margin*2
		blocks[i] = &Block{Name: name, Width: w, Height: h}

		i++
	}

	canvas := Fit(blocks)

	ss := stylesheet{
		CSSPath: path.Join(c.CSSPath, fmt.Sprintf("%s.css", c.Name)),
		ImgPath: path.Join(c.ImgPath, fmt.Sprintf("%s.%s", c.Name, c.Format)),
		Retina:  c.Retina,
		Format:  c.Format,
		ImgURL:  c.ImgURL,
		Name:    c.Name,
		Prefix:  c.Prefix,
	}

	var sprites []spriteimage

	rgba := image.NewRGBA(image.Rect(0, 0, canvas.Root.Width, canvas.Root.Height))
	draw.Draw(rgba, rgba.Bounds(), colorToUniform(c.Background), image.ZP, draw.Src)

	for _, b := range canvas.Blocks {
		img, ok := images[b.Name]
		if ok {
			src := *img
			x := b.X + c.Margin
			y := b.Y + c.Margin
			dp := image.Pt(x, y)
			r := image.Rectangle{dp, dp.Add(image.Pt(b.Width, b.Height))}
			draw.Draw(rgba, r, src, image.ZP, draw.Src)

			hasHover := strings.HasSuffix(b.Name, hoverTrigger)
			name := b.Name
			if hasHover {
				name = name[0 : len(name)-len(hoverTrigger)]
			}

			si := spriteimage{
				Name:   fmt.Sprintf("%s_%s", c.Prefix, name),
				X:      -x,
				Y:      -y,
				Width:  b.Width,
				Height: b.Height,
			}

			if hasHover {
				si.Hover = hoverCSS
			}

			sprites = append(sprites, si)
		}
	}

	ss.Images = sprites

	var doc bytes.Buffer
	// CSS Template
	tmpl, err := template.New("css").Parse(string(CSSTemplate))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tmpl.Execute(&doc, &ss)

	// Optional Test template (to stdout)
	if c.HTML {
		tmpl, err := template.New("html").Parse(string(HTMLTemplate))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		tmpl.Execute(os.Stdout, &ss)
	}

	return doc.String(), rgba
}

func (c *Config) String() string {
	return fmt.Sprintf("CONFIG: base64=%t retina=%t html=%t css=%s img=%s imgurl=%s format=%s name=%s prefix=%s bg=%s margin=%d>",
		c.Base64,
		c.Retina,
		c.HTML,
		c.CSSPath,
		c.ImgPath,
		c.ImgURL,
		c.Format,
		c.Name,
		c.Prefix,
		c.Background,
		c.Margin)
}

// validate config parameters
func (c *Config) validate() error {

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
	fn, err := filepath.Abs(path.Join(c.CSSPath, fmt.Sprintf("%s.css", c.Name)))
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(fn, []byte(sprite.Stylesheet), 0644); err != nil {
		return err
	}

	fn, err = filepath.Abs(path.Join(c.ImgPath, fmt.Sprintf("%s.%s", c.Name, c.Format)))
	if err != nil {
		return err
	}

	err = c.saveImage(fn, sprite.Image)
	if err != nil {
		return err
	}

	if sprite.RetinaImage != nil {
		fn, err = filepath.Abs(path.Join(c.ImgPath, fmt.Sprintf("%s%s.%s", c.Name, retinaTag, c.Format)))
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
