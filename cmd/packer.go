package main

import (
	"fmt"
	"os"

	"github.com/sspencer/packer"
)

// default filename, if not supplied on command line
const configFile = "./packer.toml"

func main() {

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
