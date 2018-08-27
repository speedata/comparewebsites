package screenshot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Write a javascript file to some temporary location
// that is the input of the phantomjs call
func preparePhantom() (string, error) {
	mkscreenshot_js := `
var page = require('webpage').create();
var system = require('system');
var address, output

if (system.args.length < 3) {
	phantom.exit(1);
}

address = system.args[1];
output = system.args[2];
page.viewportSize = { width: 1200, height: 600 };

page.open(address, function(status) {
	if(status === "success") {
	  page.render(output);
	};
  phantom.exit();
});
`
	f, err := ioutil.TempFile("", "compscreenr")
	if err != nil {
		return "", err
	}
	defer f.Close()
	jsfile := f.Name()
	f.WriteString(mkscreenshot_js)
	return jsfile, nil
}

// Run phantomjs and wait for it to finish
func callPhantom(url string, destpath string, wg *sync.WaitGroup, sema chan struct{}, jsfile string) error {
	sema <- struct{}{}
	defer func() { <-sema }()
	defer wg.Done()
	cmd := exec.Command("phantomjs", jsfile, url, destpath)
	fmt.Println("screenshot", url, "→", destpath)
	err := cmd.Run()
	return err
}

// Remove the temporary javascript file for phantomjs
func finishPhantom(jsfile string) {
	os.RemoveAll(jsfile)
}

// Remove all paths used in the domain list to have an empty
// directory for the screenshots.
func clearScreenshotsPath(screeshotpath string, filelist []string) error {
	domainlist := make(map[string]bool)

	for _, v := range filelist {
		u, err := url.Parse(v)
		if err != nil {
			return err
		}
		domainlist[u.Hostname()] = true
	}
	var err error
	for k, _ := range domainlist {
		err = os.RemoveAll(filepath.Join(screeshotpath, k))
		if err != nil {
			return err
		}
	}
	return nil
}

// Read one of the JSON files given on the command line
func readJsonURLList(jsonpath string) ([]string, error) {
	var localurllist []string
	fmt.Println("Trying to read json file", jsonpath)

	jsonfile, err := os.Open(jsonpath)
	if err != nil {
		return localurllist, err
	}
	defer jsonfile.Close()

	dec := json.NewDecoder(jsonfile)
	err = dec.Decode(&localurllist)
	if err != nil {
		return localurllist, err
	}

	return localurllist, nil
}

// The “main” function for the screenshot part.
func Dothings(jsonfiles []string) error {
	curwd, err := os.Getwd()
	if err != nil {
		return err
	}
	urllist := []string{}
	for _, thisjsonfile := range jsonfiles {
		localurls, err := readJsonURLList(thisjsonfile)
		if err != nil {
			return err
		}

		urllist = append(urllist, localurls...)

	}
	screeshotpath := filepath.Join(curwd, "screenshots")

	clearScreenshotsPath(screeshotpath, urllist)

	// Create a js file to be used for phantomjs
	jsfile, err := preparePhantom()
	if err != nil {
		return err
	}

	maxClients := runtime.NumCPU() * 2
	var wg sync.WaitGroup
	sema := make(chan struct{}, maxClients)

	for _, thisfile := range urllist {
		u, err := url.Parse(thisfile)
		if err != nil {
			return err
		}
		destpath := filepath.Join(u.Hostname(), u.Path)
		if strings.HasSuffix(u.Path, "/") || u.Path == "" {
			destpath = destpath + "/index.html"
		}
		destpath = destpath + ".png"

		wg.Add(1)
		go callPhantom(thisfile, filepath.Join(screeshotpath, destpath), &wg, sema, jsfile)

	}
	wg.Wait()
	finishPhantom(jsfile)
	return nil
}
