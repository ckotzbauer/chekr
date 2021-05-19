package internal

type FileGroup struct {
	registerFile string
	typesFile    string
}

type GroupVersion struct {
	Group     string     `json:"group"`
	Version   string     `json:"version"`
	Resources []Resource `json:"resources"`
}

type GroupVersionKind struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

type Resource struct {
	Name        string           `json:"name"`
	Introduced  string           `json:"introduced"`
	Deprecated  string           `json:"deprecated"`
	Removed     string           `json:"removed"`
	Replacement GroupVersionKind `json:"replacement"`
}
