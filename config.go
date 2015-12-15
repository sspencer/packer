package packer

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const retinaTag = "@2x"

type SpriteConfig struct {
	Includes   []string
	Retina     bool
	Stylesheet string
	Sprite     string
	CSSPath    string
	Prefix     string
	Background string
	Margin     int
	Hover      string
}

func (c *SpriteConfig) String() string {
	return fmt.Sprintf(`<SpriteConfig
	includes=%s
	retina=%t
	stylesheet=%s
	sprite=%s
	csspath=%s
	prefix=%s
	background=%s
	margin=%d
	hover=%s />`,
		c.Includes,
		c.Retina,
		c.Stylesheet,
		c.Sprite,
		c.CSSPath,
		c.Prefix,
		c.Background,
		c.Margin,
		c.Hover)
}

// NewConfig returns sprite configuration wth defaults applied.
func NewConfig(filename string) (*SpriteConfig, error) {
	var c SpriteConfig

	if _, err := toml.DecodeFile(filename, &c); err != nil {
		err = fmt.Errorf("ERROR: Could not read configuration file %q: %s", filename, err)
		return nil, err
	}

	// set default values if not already set
	if len(c.Includes) == 0 {
		c.Includes = []string{"./"}
	}
	if c.Sprite == "" {
		c.Sprite = "./sprite.png"
	}
	if c.CSSPath == "" {
		c.CSSPath = "/img/"
	}
	if c.Prefix == "" {
		c.Prefix = "icon"
	}
	if c.Background == "" {
		c.Background = "transparent"
	}

	return &c, nil
}

func (c *SpriteConfig) SaveStylesheet(stylesheet string) error {

	// stylesheet is optional, so just return
	if stylesheet == "" {
		return nil
	}

	fmt.Println("SAVE STYLESHEET")
	var p string
	var err error

	if p, err = filepath.Abs(c.Stylesheet); err != nil {
		return fmt.Errorf("ERROR: Can not write to file %q: %s", p, err)
	}

	if err = ioutil.WriteFile(p, []byte(stylesheet), 0644); err != nil {
		return fmt.Errorf("ERROR: Failed writing to file %q: %s", p, err)
	}

	return nil
}

// SaveSprite saves the sprite image to file specified in the config Sprite value.
func (c *SpriteConfig) SaveSprite(img *image.RGBA) error {
	fmt.Println("SAVE SPRITE")
	var p string
	var err error
	if p, err = filepath.Abs(c.Sprite); err != nil {
		return fmt.Errorf("ERROR: Can not write to file %q: %s", p, err)
	}

	w, _ := os.Create(p)
	defer w.Close()
	png.Encode(w, img)
	return nil
}

// SaveRetinaSprite saves the retina sprite image, if retina config value was set..
func (c *SpriteConfig) SaveRetinaSprite(img *image.RGBA) error {
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
	return nil
}
