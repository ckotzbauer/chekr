package printer

type PrintableList interface {
	ToJson() (string, error)
	ToHtml() (string, error)
	ToTable() (string, error)
}

type Printable interface {
}

type PrintableResult struct {
	Item  Printable
	Error error
}
