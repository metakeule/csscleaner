package lib

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
csstidy input-file [options] [output-file]

--template=highest
*/

func Cleanup(in string) (out string, err error) {

	var (
		dir     string
		infile  string
		outfile string
		csstidy string
	)

steps:
	for jump := 1; err == nil; jump++ {
		switch jump - 1 {
		default:
			break steps
		// count a number up for each following step
		case 0:
			cmd := exec.Command("which", "csstidy")
			var o []byte
			o, err = cmd.Output()
			csstidy = strings.TrimSpace(string(o))
		case 1:
			dir, err = ioutil.TempDir("/tmp", "csstidy_")
		case 2:
			defer os.RemoveAll(dir)
		case 3:
			infile = filepath.Join(dir, "input.css")
			outfile = filepath.Join(dir, "output.css")
			err = ioutil.WriteFile(infile, []byte(in), 0644)
		case 4:
			cmd := exec.Command(csstidy, infile, "--template=highest", outfile)
			err = cmd.Run()
		case 5:
			var o []byte
			o, err = ioutil.ReadFile(outfile)
			out = string(o)
		}
	}

	return
}
