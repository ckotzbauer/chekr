package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"

	"github.com/sirupsen/logrus"
)

type Printer struct {
	Type string
	File string
}

func (p Printer) Print(list PrintableList) {
	var out string

	if p.Type == "json" {
		out = list.ToJson()
	} else if p.Type == "html" {
		out = list.ToHtml()
	} else {
		out = list.ToTable()
	}

	if p.File != "" {
		data := []byte(out)
		err := os.WriteFile(p.File, data, 0)

		if err != nil {
			logrus.WithError(err).Fatalf("Could write output to file!")
		}
	} else {
		fmt.Fprint(os.Stdout, out)
	}
}

func ToJson(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		logrus.WithError(err).Fatalf("Could not marshal object!")
	}

	return string(b)
}

func ToHtml(html string, v interface{}) string {
	buf := new(bytes.Buffer)
	tpl := template.New("page")
	tpl, err := tpl.Parse(html)

	if err != nil {
		logrus.WithError(err).Fatalf("Could create html-page!")
	}

	err = tpl.Execute(buf, v)

	if err != nil {
		logrus.WithError(err).Fatalf("Could create html-page!")
	}

	content := buf.String()
	return fmt.Sprintf(HtmlContainer, content)
}
