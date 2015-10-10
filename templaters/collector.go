package main

import "bufio"
import "fmt"
import "os"
import "strings"
import "unicode"
import "unicode/utf8"
import "text/template"
import "github.com/Sirupsen/logrus"
import "github.com/ogier/pflag"

// Data is an exported type that
// contains information used for the template
type Data struct {
	Name string
}

func main() {
	names := pflag.String("names", "", "Comma-separated list of collectors")
	pflag.Parse()
	if *names == "" {
		logrus.Fatal("names of collector not specified")
	}

	for _, name := range strings.Split(*names, ",") {
		data := Data{Name: upperFirst(name)}
		collectorName := strings.ToLower(fmt.Sprintf("%s_collector.go", name))
		f, err := os.Create(fmt.Sprintf("../collectors/%s", collectorName))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		t := template.New("collector.tmpl")
		t, err = t.ParseFiles("templates/collector.tmpl")
		err = t.Execute(w, data)
		if err != nil {
			panic(err)
		}

		w.Flush()
	}
}

func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}
