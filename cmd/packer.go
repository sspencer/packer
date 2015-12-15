package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sspencer/packer"
)

// default filename, if not supplied on command line
const configFile = "./packer.toml"

func main() {

	var defaults = packer.SpriteDefaults{
		Classname:  "icon-",
		Background: "transparent",
		Hover:      "",
		Padding:    0,
	}

	if len(os.Args) >= 2 {
		for _, fn := range os.Args[1:] {
			pack(fn, &defaults)
		}
	} else {
		pack(configFile, &defaults)
	}
}

func pack(filename string, defaults *packer.SpriteDefaults) {
	var config *packer.SpritesConfig
	var err error
	if config, err = packer.NewConfig(filename, defaults); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// var Discard io.Writer = devNull(0)

	for _, c := range config.Sprites {
		img, stylesheet, err := packer.CreateSprite(&c)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		saveStylesheet(stylesheet, &c)
		saveSprite(img, &c)
	}
}

func saveStylesheet(stylesheet string, c *packer.SpriteConfig) {
	fmt.Println("SAVE STYLESHEET")
	var p string
	var err error

	if p, err = filepath.Abs(c.Stylesheet); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Can not write to file %q: %s", p, err)
		os.Exit(1)
	}

	if err = ioutil.WriteFile(p, []byte(stylesheet), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed writing to file %q: %s", p, err)
	}
}

func saveSprite(img *image.RGBA, c *packer.SpriteConfig) {
	fmt.Println("SAVE SPRITE")
	var p string
	var err error
	if p, err = filepath.Abs(c.Sprite); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Can not write to file %q: %s", p, err)
		os.Exit(1)
	}

	w, _ := os.Create(p)
	defer w.Close()
	png.Encode(w, img)
}
