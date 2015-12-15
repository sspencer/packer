package packer

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
)

type SpritesConfig struct {
	Defaults SpriteDefaults
	Sprites  []SpriteConfig `toml:"sprite"`
}

type SpriteDefaults struct {
	Classname  string
	Background string
	Hover      string
	Padding    int
}

type SpriteConfig struct {
	Includes   []string
	Excludes   []string
	Stylesheet string
	Sprite     string
	Classname  string
	Background string
	Retina     bool
	Hover      string
	Padding    int
}

func (c SpriteDefaults) String() string {
	return fmt.Sprintf(`<SpriteDefaults
	classname=%s
	background=%s
	hover=%s
	padding=%d />`,
		c.Classname,
		c.Background,
		c.Hover,
		c.Padding)
}

func (c SpriteConfig) String() string {
	return fmt.Sprintf(`<SpriteConfig
	includes=%s
	excludes=%s
	stylesheet=%s
	sprite=%s
	classname=%s
	background=%s
	retina=%t
	hover=%s
	padding=%d />`,
		c.Includes,
		c.Excludes,
		c.Stylesheet,
		c.Sprite,
		c.Classname,
		c.Background,
		c.Retina,
		c.Hover,
		c.Padding)
}

func (c *SpriteConfig) setDefaults(d *SpriteDefaults) {

	if c.Classname == "" {
		c.Classname = d.Classname
	}
	if c.Background == "" {
		c.Background = d.Background
	}
	if c.Hover == "" {
		c.Hover = d.Hover
	}
	if c.Padding == 0 {
		c.Padding = d.Padding
	}
}

func (c *SpriteConfig) verify() error {
	if len(c.Includes) == 0 {
		return errors.New("At least one 'include' path must be specified.")
	}

	// verify each path
	if c.Stylesheet == "" {
		return errors.New("Config value 'stylesheet' must be specified.")
	}

	if c.Sprite == "" {
		return errors.New("Config value 'sprite' must be specified.")
	}
	return nil
}

// NewConfig returns sprite configuration wth defaults applied.
func NewConfig(filename string, appDefaults *SpriteDefaults) (*SpritesConfig, error) {
	var config SpritesConfig

	if _, err := toml.DecodeFile(filename, &config); err != nil {
		err = fmt.Errorf("ERROR: Could not read configuration file %q: %s", filename, err)
		return nil, err
	}

	for i, c := range config.Sprites {
		c.setDefaults(&config.Defaults)
		c.setDefaults(appDefaults)

		if err := c.verify(); err != nil {
			err = fmt.Errorf("ERROR: Configuration problem in sprite config #%d: %s", i+1, err)
			return nil, err
		}
	}

	return &config, nil
}
