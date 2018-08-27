# Compare two websites

Compare two websites by taking screenshots of pages and visually comparing each.

## Prerequisites

You need a Go compiler, phantomjs and imagemagick (compare).

## Usage

You can supply zero or more JSON files that lists URLs to take screenshots and optionally two directories with screenshots to compare.

    Usage: cmpimages [.json*] [dir1 dir2]

The JSON files must have this format:

    [
     	"https://www.example.com/",
     	"https://www.example.com/path/to/page.html",
     	"https://www.example.com/anotherpath/"
    ]

The screenshots are crated in a sub directory of the current working directory called `screenshots`.


## Caveats

This is tested only on a Mac, feedback etc. is welcome. First version.


## Contact

gundlach@speedata.de
