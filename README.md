# Compare two websites

Compare two websites by taking screenshots of pages and visually comparing each.

## Prerequisites

You need a [Go compiler](https://golang.org/), [PhantomJS](http://phantomjs.org/) and [ImageMagick](https://www.imagemagick.org/script/index.php) (preferably version 7 and above).

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


## Building

	git clone https://github.com/speedata/comparewebsites
	cd comparewebsites
	make

Now you should have a file called `cmpimages`. Try to run it:

    ./cmpimages example/example.json screenshots/www.speedata.de/de/ screenshots/www.speedata.de/en/

Which does two things:

1. Create a directory called `screenshots` in the current working directory and save screenshots of the given web pages (the json files).
1. Create a directory $TMPDIR/cmpimage-diff and writes a web page called `out.html`. This page shows the differences. Another file called `diff.json` in the current working directory lists all the pages which have differences.




## Caveats

This is tested only on a Mac, feedback etc. is welcome. First version. ImageMagick < version 7 does not compare files with different sizes, so expect strange output in the HTML file.


## Contact

gundlach@speedata.de
