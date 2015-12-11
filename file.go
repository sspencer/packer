package packer

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ReadImages returns map of normalized file name (becomes css name) to image data.
func readImages(dir string) (images map[string]*image.Image, err error) {

	entries, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	images = make(map[string]*image.Image)

	// image css class name gets normalized
	re := regexp.MustCompile("([^-a-zA-Z0-9])")

	for _, entry := range entries {

		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))

		// filtering twice, but why open file needlessly if ext does not match
		if ext != ".png" && ext != ".jpeg" && ext != ".jpg" {
			continue
		}

		fullName := filepath.Join(dir, name)
		imageName := re.ReplaceAllLiteralString(name[:len(name)-len(ext)], "-")
		imgFile, err := os.Open(fullName)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open file: %s", fullName)
			continue
		}

		defer imgFile.Close()

		if ext == ".png" {
			img, err := png.Decode(imgFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Problem decoding PNG image: %s", fullName)
				continue
			}

			images[imageName] = &img

		} else if ext == ".jpg" || ext == ".jpeg" {
			img, err := jpeg.Decode(imgFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Problem decoding JPEG image: %s", fullName)
				continue
			}

			images[imageName] = &img
		}
	}

	return images, nil
}
