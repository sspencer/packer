package main

import "fmt"
import "gopkg.in/alecthomas/kingpin.v2"

var (
	base64     = kingpin.Flag("base64", "Create css with base64 encoded sprite.").Short('b').Bool()
	cssout     = kingpin.Flag("css", "CSS output path (defaults to stdout).").Short('c').String()
	imgout     = kingpin.Flag("img", "Image output path (if not set, base64 set to true).").Short('o').String()
	imgpath    = kingpin.Flag("imgpath", "http path to images on the web server").Short('i').Default("../img").String()
	format     = kingpin.Flag("format", "Output format of the sprite (png or jpg)  [png].").Short('f').Default("png").String()
	name       = kingpin.Flag("name", "Name of sprite file without file extension (image and css).").Short('n').Default("sprite").String()
	retina     = kingpin.Flag("retina", "Generate retina and normal sprite. Source images must be in retina resolution.").Short('r').Bool()
	background = kingpin.Flag("bg", "Background color of the sprite in hex (or 'transparent')").Short('g').Default("transparent").String()
	margin     = kingpin.Flag("margin", "Margin in px between tiles.").Default("4").Short('m').Int()
	prefix     = kingpin.Flag("prefix", "Prefix for the class name used in css.").Short('p').Default("sprite_").String()

	src = kingpin.Arg("src", "input images").Required().Strings()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	if *imgout == "" {
		*base64 = true
	}

	if *margin < 0 {
		*margin = 0
	}

	if *margin > 100 {
		*margin = 100
	}

	if *format != "jpg" && *format != "png" {
		kingpin.Fatalf("Illegal option %q for format. Must be 'png' or 'jpg' instead.", *format)
	}

	fmt.Printf("Base64: %t\n", *base64)
	fmt.Printf("CSS: %s\n", *cssout)
	fmt.Printf("Image: %s\n", *imgout)
	fmt.Printf("ImgPath: %s\n", *imgpath)
	fmt.Printf("Format: %s\n", *format)
	fmt.Printf("Name: %s\n", *name)
	fmt.Printf("Retina: %t\n", *retina)
	fmt.Printf("Background: %s\n", *background)
	fmt.Printf("Margin: %d\n", *margin)
	fmt.Printf("Prefix: %s\n", *prefix)
	fmt.Printf("Source: %s\n", *src)
}
