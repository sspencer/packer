# Packer

Go implementation of a Binary Tree Algorithm for bin packing.  Based on
Jake Gordon's [bin-packing](https://github.com/jakesgordon/bin-packing)
project on github.

There is still lots to do and structs I want to rename, but the base
algorithm is working.

    $ cd tester
    $ go run main.go

The above will pack a number of different sets of blocks and display
the results in series of PNGs.  On OSX, the images open in the
Preview app, other OSes will get a Exit 1 code.

Eventual goal is to turn this into a sprite packer utility, that takes
a directory of PNGs and outputs a sprite image and stylesheet, packing
the images as tight as possible.

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
