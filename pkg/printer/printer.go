package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
)

type Printer struct {
	Type string
	File string
}

func (p Printer) Print(list PrintableList) error {
	var out string
	var err error

	if p.Type == "json" {
		out, err = list.ToJson()
	} else if p.Type == "html" {
		out, err = list.ToHtml()
	} else {
		out, err = list.ToTable()
	}

	if err != nil {
		return err
	}

	if p.File != "" {
		data := []byte(out)
		err = os.WriteFile(p.File, data, 0)
	} else {
		fmt.Fprint(os.Stdout, out)
	}

	return err
}

func ToJson(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func ToHtml(html string, v interface{}) (string, error) {
	buf := new(bytes.Buffer)
	tpl := template.New("page")
	tpl, err := tpl.Parse(html)

	if err != nil {
		return "", nil
	}

	err = tpl.Execute(buf, v)
	content := buf.String()
	return fmt.Sprintf(HtmlContainer, content), err
}
