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
the images as tight as possible.  Images can be packed tighter simply
by sorting them in differnt orders (height first, width first, etc).

