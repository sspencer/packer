package packer

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// Given list of file paths, return map of css name to (path/extension removed) to image data
func (c *Config) getImages(files []string) (map[string]*image.Image, error) {
	images := make(map[string]*image.Image)
	re := regexp.MustCompile("([^_a-zA-Z0-9])")

	for _, fn := range files {
		base := path.Base(fn)
		ext := strings.ToLower(filepath.Ext(base))
		name := re.ReplaceAllLiteralString(base[:len(base)-len(ext)], "_")

		if ext != ".png" && ext != ".jpeg" && ext != ".jpg" {
			return nil, fmt.Errorf("Unrecognized file extension, %q", ext)
		}

		f, err := os.Open(fn)
		if err != nil {
			return nil, fmt.Errorf("Could not open file, %q", fn)
		}

		defer f.Close()

		if ext == ".png" {
			img, err := png.Decode(f)
			if err != nil {
				return nil, fmt.Errorf("Problem decoding PNG image, %q", fn)
			}

			images[name] = &img

		} else if ext == ".jpg" || ext == ".jpeg" {
			img, err := jpeg.Decode(f)
			if err != nil {
				return nil, fmt.Errorf("Problem decoding JPEG image, %q", fn)
			}

			images[name] = &img
		}
	}

	return images, nil
}
