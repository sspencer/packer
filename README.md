# Packer

Go implementation of a Binary Tree Algorithm for bin packing.  Based on
Jake Gordon's [bin-packing](https://github.com/jakesgordon/bin-packing)
project on github.

There are not a lot of options yet and a stylesheet is not created yet,
but a CSS Sprite image is created.

## Usage

To generate a series of test images created with pastel colored blocks,
run the tester program.  On OSX, the images open in the
Preview app, other OSes will get a Exit 1 code.

    $ cd tester
    $ go run main.go

<img src="misc/screenshot.png?raw=true" />

To generate an actual css sprite from your own directory of images, invoke
the pack command with a configuration file:

    $ cd cmd
    $ go run pack.go some/path/to/packer.toml

## Configuration

The configuration file is evolving, but uses [TOML](https://github.com/toml-lang/toml),
Tom's Obvious, Minimal Language.  All paths specified in the config file are relative to
its location.

```toml
# packer.toml
includes=["./icons/*.png"]
retina=true
stylesheet="./sprites.css"
sprite="./sprites.png"
csspath="/img/sprites.png"
hover="_hover"
background="transparent"
prefix="sprite"
margin=0
```

TBD - explain configuration options and default values.


## About the Code

Channels are used for an "embarrassingly parallel" problem ... pack the
images in 4 different ways by sorting the images differently:

* sort by width
* sort by height
* sort by area
* sort by max side (width or height)

By changing the sort order of the images, an occasional advantage can
be realized.

For example:

	==== Packing complex ====
	LayoutByWidth <650x730> has wasted 194700 pixels
	LayoutByArea <650x650> has wasted 142700 pixels
	LayoutByHeight <530x530> has wasted 1100 pixels
	LayoutByMax <730x400> has wasted 12200 pixels
	>>>> RETURNING  LayoutByHeight

## Todo

* Generate stylesheet
* Integrate CLI commander like [cobra](https://github.com/spf13/cobra)
