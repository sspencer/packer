package packer

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func getImageMap(config *SpriteConfig) (images map[string]*image.Image, err error) {

	imagePaths, err := getImagePaths(config)
	if err != nil {
		return nil, err
	}

	return getImageData(imagePaths)
}

// getImageData returns map of normalized file name (becomes css name) to image data.
func getImageData(filenames []string) (images map[string]*image.Image, err error) {

	images = make(map[string]*image.Image)

	// image css class name gets normalized
	re := regexp.MustCompile("([^-a-zA-Z0-9])")

	for _, fn := range filenames {
		log.Println(fn)
		base := path.Base(fn)
		ext := strings.ToLower(filepath.Ext(base))
		name := re.ReplaceAllLiteralString(base[:len(base)-len(ext)], "-")

		// filtering twice, but why open file needlessly if ext does not match
		if ext != ".png" && ext != ".jpeg" && ext != ".jpg" {
			continue
		}

		f, err := os.Open(fn)
		if err != nil {
			return nil, fmt.Errorf("Could not open file: %q", fn)
		}

		defer f.Close()

		if ext == ".png" {
			img, err := png.Decode(f)
			if err != nil {
				return nil, fmt.Errorf("Problem decoding PNG image: %q", fn)
			}

			images[name] = &img

		} else if ext == ".jpg" || ext == ".jpeg" {
			img, err := jpeg.Decode(f)
			if err != nil {
				return nil, fmt.Errorf("Problem decoding JPEG image: %q", fn)
			}

			images[name] = &img
		}
	}

	return images, nil
}

func getImagePaths(config *SpriteConfig) ([]string, error) {

	var files []string

	for _, s := range config.Includes {
		j, err := filepath.Abs(s)
		if err != nil {
			return nil, fmt.Errorf("ERROR: Could not determine absolute path of includes:", err)
		}

		matches, err := filepath.Glob(j)
		if err != nil {
			return nil, fmt.Errorf("ERROR: Could not glob matching files:", err)
		}

		// add file to list if not excluded
		if len(config.Excludes) == 0 {
			files = append(files, matches...)
		} else {
			for _, m := range matches {
				for _, x := range config.Excludes {
					if strings.Index(m, x) == -1 {
						files = append(files, m)
					}
				}
			}
		}

	}

	return files, nil
}
