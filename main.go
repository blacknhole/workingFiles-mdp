package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const defaultTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html"; charset="utf-8">
    <title>{{ .Title }}</title>
  </head>
  <body>
{{ .File }}
{{ .Body }}
  </body>
</html>
`

type content struct {
	Title string
	File  template.HTML
	Body  template.HTML
}

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")
	flag.Parse()

	if *tFname == "" {
		*tFname = os.Getenv("TEMPLATE_NAME")
	}

	if err := run(os.Stdin, *filename, *tFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(r io.Reader, filename, tFname string, out io.Writer, skipPreview bool) error {
	input := []byte{}
	if filename != "" {
		var err error
		input, err = os.ReadFile(filename)
		if err != nil {
			return err
		}
	} else {
		s := bufio.NewScanner(r)
		for s.Scan() {
			input = append(input, []byte(s.Text()+"\n")...)
		}
		if s.Err() != nil {
			return s.Err()
		}
	}

	htmlData, err := parseContent(input, filename, tFname)
	if err != nil {
		return err
	}

	tf, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}

	if err := tf.Close(); err != nil {
		return err
	}

	outName := tf.Name()
	fmt.Fprintln(out, outName)
	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)
	return preview(outName)
}

func parseContent(input []byte, filename, tFname string) ([]byte, error) {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	if filename != "" {
		filename = `<p>The file <b>` + filename +
			`</b> is being previewed...</p>
<hr>
<br>
`
	}
	c := content{
		Title: "Markdown Preview Tool",
		File:  template.HTML(filename),
		Body:  template.HTML(body),
	}

	var buffer bytes.Buffer

	if err = t.Execute(&buffer, c); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func saveHTML(outFname string, data []byte) error {
	return os.WriteFile(outFname, data, 0644)
}

func preview(filename string) error {
	cName := ""
	cParams := []string{}
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	cParams = append(cParams, filename)

	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	err = exec.Command(cPath, cParams...).Run()
	time.Sleep(2 * time.Second)
	return err
}
