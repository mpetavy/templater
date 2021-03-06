package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"os"
	"strings"
	"text/template"
)

type data struct {
	Type string
	Name string
}

var (
	inputFile     *string
	outputFile    *string
	searchReplace *string
	d             data
)

func init() {
	common.Init(false, "1.0.0", "", "2018", "GO code generator by template", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, nil, run, 0)

	inputFile = flag.String("i", "", "The file to be parsed")
	outputFile = flag.String("o", "", "The file to be generated")
	searchReplace = flag.String("sr", "", "Search")
	flag.StringVar(&d.Type, "t", "", "The subtype used for the queue being generated")
	flag.StringVar(&d.Name, "n", "", "The name used for the queue being generated. This should start with a capital letter so that it is exported.")
}

func run() error {
	b, err := os.ReadFile(common.CleanPath(*inputFile))
	if err != nil {
		return err
	}

	t := template.Must(template.New(".").Parse(string(b)))

	var buf bytes.Buffer

	err = t.Execute(&buf, d)
	if err != nil {
		return err
	}

	code := string(buf.Bytes())

	for i := 5; i > 1; i-- {
		code = strings.Replace(code, "    ", "\t", -1)
	}

	items := strings.Split(*searchReplace, ";")

	for _, v := range items {
		p := strings.Index(v, "=")
		search := v[:p]
		replace := v[p+1:]

		code = strings.Replace(code, search, replace, -1)
	}

	if *outputFile != "" {
		err := os.WriteFile(common.CleanPath(*outputFile), []byte(code), common.DefaultFileMode)
		if common.Error(err) {
			return err
		}
	}

	return nil
}

func main() {
	defer common.Done()

	common.Run(nil)
}
