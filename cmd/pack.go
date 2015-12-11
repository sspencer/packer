package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/sspencer/packer"
)

var (
	prefix     = flag.String("prefix", "icon", "css class prefix")
	imgSrc     = flag.String("img", ".", "sprite image source")
	cssDest    = flag.String("css", ".", "destination folder for sprite stylesheet")
	spriteDest = flag.String("sprite", ".", "destination folder for sprite image")
)

func main() {
	flag.Parse()

	// var Discard io.Writer = devNull(0)
	img, stylesheet, err := packer.CreateSprite(*imgSrc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Stylesheet: %s\n", stylesheet)

	w, _ := os.Create("sprite.png")
	defer w.Close()
	png.Encode(w, img)
}
