package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Reads all files in the 'gen' folder and encodes them as strings literals
// in assets/assets.go.
func main() {
	fs, _ := ioutil.ReadDir("gen/")
	out, _ := os.Create("assets/assets.go")
	defer out.Close()
	out.Write([]byte("package assets \n\nconst (\n"))
	for _, f := range fs {
		if f.Name() == "assets.go" {
			continue // we do not want to include the assets file itself
		}
		out.Write([]byte(strings.Title(strings.Replace(f.Name(), ".", "", -1)) + " = `"))
		f, _ := os.Open("gen/" + f.Name())
		io.Copy(out, f)
		out.Write([]byte("`\n"))
	}
	out.Write([]byte(")\n"))
}
