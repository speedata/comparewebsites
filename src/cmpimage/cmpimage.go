package cmpimage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	startpath  string
	destpath   string
	diffdir    string
	buf        bytes.Buffer
	images     []string
	wg         sync.WaitGroup
	curwd      string
	maxclients int
	sema       chan struct{}
)

// Run compare from imagemagick
func compare(path, cmppath, diffpath, relpath string) error {
	defer wg.Done()
	sema <- struct{}{}
	defer func() { <-sema }()
	p := longestCommonPrefix(path, cmppath)
	fmt.Println("compare", strings.TrimPrefix(path, p), "↔", strings.TrimPrefix(cmppath, p))
	err := exec.Command("compare", "-metric", "AE", path, cmppath, diffpath).Run()
	if err != nil {
		images = append(images, relpath)
	}
	return err
}

// Recurse into the first directory and assume that the second directory has the
// same structure
func recurse(path string, info os.FileInfo, err error) error {
	relpath, err := filepath.Rel(startpath, path)
	if err != nil {
		return err
	}
	// Target directory
	cmppath := filepath.Join(destpath, relpath)

	// Filename / file path for the diff
	diffpath := filepath.Join(diffdir, relpath)
	if cmppathinfo, err := os.Stat(cmppath); err != nil {
		fmt.Println("Path does not exist:", cmppath)
	} else {
		if cmppathinfo.IsDir() {
			os.MkdirAll(diffpath, 0755)
		} else {
			if strings.HasSuffix(path, ".png") {
				wg.Add(1)
				go compare(path, cmppath, diffpath, relpath)
			}
		}
	}

	return nil
}

func render(outfile string) error {
	tmpl := `<!DOCTYPE html>
<html>
<head>
	<title>Bilder</title>
	<style type="text/css">
		img { width: 30% ; }
		tr.img td	{ border-bottom: 1pt solid black; }
		tr  {vertical-align: top;}
	</style>
</head>
<body>
	<table>
	{{ range .Images -}}
	<tr>
		<td colspan="3">{{ .}}</td>
	</tr>
	<tr class="img">
		<td><a target="_blank" href="{{$.Diffpath}}/{{ . }}"><img src="{{$.Diffpath}}/{{ . }}"></a></td>
		<td><a target="_blank" href="{{$.Startpath}}/{{ . }}"><img src="{{$.Startpath}}/{{ . }}"></a></td>
		<td><a target="_blank" href="{{$.Destpath}}/{{ . }}"><img src="{{$.Destpath}}/{{ . }}"></a></td>
	</tr>
	{{- end }}
</table>
</body>
</html>
`

	t := template.Must(template.New("html").Parse(tmpl))
	data := struct {
		Diffpath  string
		Startpath string
		Destpath  string
		Images    []string
	}{
		Diffpath:  diffdir,
		Startpath: startpath,
		Destpath:  destpath,
		Images:    images,
	}
	err := t.Execute(&buf, data)
	if err != nil {
		return err
	}
	f, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = buf.WriteTo(f)
	if err != nil {
		return err
	}
	return nil
}

// The “main” function
func startcompare() error {
	var err error
	curwd, err = os.Getwd()
	if err != nil {
		return err
	}

	if !filepath.IsAbs(startpath) {
		startpath = filepath.Join(curwd, startpath)
	}
	if !filepath.IsAbs(destpath) {
		destpath = filepath.Join(curwd, destpath)
	}

	diffdir = filepath.Join(os.TempDir(), "diff")
	os.RemoveAll(diffdir)

	filepath.Walk(startpath, recurse)

	wg.Wait()
	outhtml := filepath.Join(diffdir, "out.html")
	err = render(outhtml)
	if err != nil {
		return err
	}
	fmt.Println("html:", outhtml)
	if runtime.GOOS == "darwin" {
		exec.Command("open", outhtml).Run()
	} else {
		fmt.Println(runtime.GOOS, "no open command provided")
	}

	return nil
}

// Create a json file of the differences and write it into the
// current directory
func jsonit() error {
	guessHostname := filepath.Base(destpath)
	var jsonbuf bytes.Buffer
	var tojson []string
	for _, img := range images {
		p := strings.TrimSuffix(img, ".png")
		if strings.HasSuffix(p, "index.html") {
			p = strings.TrimSuffix(p, "index.html")
		}
		tojson = append(tojson, "http://"+guessHostname+"/"+p)

	}
	b, err := json.Marshal(tojson)
	if err != nil {
		return err
	}
	json.Indent(&jsonbuf, b, "", " ")
	jsonpath := filepath.Join(curwd, "diff.json")
	outfile, err := os.Create(jsonpath)
	if err != nil {
		return err
	}
	defer outfile.Close()
	jsonbuf.WriteTo(outfile)
	fmt.Println("json:", jsonpath)
	return nil
}

// Start comparing two directories
func Dothings(paths []string) error {
	if len(paths) != 2 {
		return fmt.Errorf("Two paths needed for Dothings")
	}

	maxclients = runtime.NumCPU() * 3
	sema = make(chan struct{}, maxclients)

	startpath = paths[0]
	destpath = paths[1]
	fmt.Println("Comparing", startpath, destpath)
	err := startcompare()
	err = jsonit()
	return err
}
