package main

import (
	"fmt"
	"os"

	"github.com/sspencer/packer"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("packer", "CSS Sprite generator")
	//base64     = app.Flag("base64", "Create css with base64 encoded sprite.").Short('b').Bool()
	format     = app.Flag("format", "Output format of the sprite (png or jpg)  [png].").Short('f').Default("png").String()
	cssout     = app.Flag("css", "destination path for sprite image file").Default("css/").String()
	imgout     = app.Flag("img", "destination path for sprite image file").Default("img/").String()
	imgurl     = app.Flag("imgpath", "HTTP path to images on the web server").Short('i').Default("../img").String()
	margin     = app.Flag("margin", "Margin in px between tiles.").Default("4").Short('m').Int()
	name       = app.Flag("name", "Name of sprite file without file extension (image and css).").Short('n').Default("sprite").String()
	prefix     = app.Flag("prefix", "Prefix for the class name used in css.").Short('p').Default("sprite").String()
	retina     = app.Flag("retina", "Generate retina and normal sprite. Source images must be in retina resolution.").Short('r').Bool()
	background = app.Flag("background", "Background color of the sprite in hex (or 'transparent')").Default("transparent").String()
	html       = app.Flag("html", "Output test HTML file to stdout").Bool()
	showCSS    = app.Flag("show-css-template", "Print CSS Template to <stdout> and exit").Bool()
	showHTML   = app.Flag("show-html-template", "Print HTML Template to <stdout> and exit").Bool()

	images = app.Arg("src", "source images").Strings()
)

func main() {
	app.Version("0.0.1")
	//kingpin.Parse()
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *showCSS {
		fmt.Println(packer.CSSTemplate)
		os.Exit(0)
	}

	if *showHTML {
		fmt.Println(packer.HTMLTemplate)
		os.Exit(0)
	}

	c := packer.Config{
		Base64:     false, /**base64,*/
		Retina:     *retina,
		HTML:       *html,
		CSSPath:    *cssout,
		ImgPath:    *imgout,
		ImgURL:     *imgurl,
		Format:     *format,
		Name:       *name,
		Prefix:     *prefix,
		Margin:     *margin,
		Background: *background,
	}

	sprite, err := c.CreateSprite(*images)
	if err != nil {
		app.Fatalf("%s\n", err)
	}

	if err := c.Save(sprite); err != nil {
		app.Fatalf("%s\n", err)
	}
}
