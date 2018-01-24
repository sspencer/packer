package main

import (
	"fmt"
	"os"

	"github.com/sspencer/packer"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	base64     = kingpin.Flag("base64", "Create css with base64 encoded sprite.").Short('b').Bool()
	csspath    = kingpin.Flag("css", "CSS output path (defaults to stdout).").Short('c').String()
	imgpath    = kingpin.Flag("img", "Image output path (if not set, base64 set to true).").Short('o').String()
	imgurl     = kingpin.Flag("imgpath", "http path to images on the web server").Short('u').Default("../img").String()
	format     = kingpin.Flag("format", "Output format of the sprite (png or jpg)  [png].").Short('f').Default("png").String()
	name       = kingpin.Flag("name", "Name of sprite file without file extension (image and css).").Short('n').Default("sprite").String()
	prefix     = kingpin.Flag("prefix", "Prefix for the class name used in css.").Short('p').Default("sprite_").String()
	retina     = kingpin.Flag("retina", "Generate retina and normal sprite. Source images must be in retina resolution.").Short('r').Bool()
	margin     = kingpin.Flag("margin", "Margin in px between tiles.").Default("4").Short('m').Int()
	background = kingpin.Flag("bg", "Background color of the sprite in hex (or 'transparent')").Short('g').Default("transparent").String()

	images = kingpin.Arg("src", "input images").Required().Strings()
)

/*
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
*/
func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	c := packer.Config{
		Base64:     *base64,
		Retina:     *retina,
		CSSPath:    *csspath,
		ImgPath:    *imgpath,
		ImgURL:     *imgurl,
		Format:     *format,
		Name:       *name,
		Prefix:     *prefix,
		Margin:     *margin,
		Background: *background,
	}

	sprite, err := c.CreateSprite(*images)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := sprite.SaveSprite(sprite.Image); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// default filename, if not supplied on command line
const configFile = "./packer.toml"

func mai2n() {

	if len(os.Args) >= 2 {
		pack(os.Args[1])
	} else {
		pack(configFile)
	}
}

func pack(filename string) {
	var c *packer.SpriteConfig
	var err error
	if c, err = packer.NewConfig(filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// var Discard io.Writer = devNull(0)

	r, err := packer.CreateSprite(c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = c.SaveStylesheet(r.Stylesheet); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = c.SaveSprite(r.Sprite); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = c.SaveRetinaSprite(r.Sprite2x); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
