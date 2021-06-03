package printer

type PrintableList interface {
	ToJson() string
	ToHtml() string
	ToTable() string
}

type Printable interface {
}
