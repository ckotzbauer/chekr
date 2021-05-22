package kubernetes

type KindVersions map[string][]KindVersion

type KindVersion struct {
	Group     string
	Version   string
	Preferred bool
	Name      string
}
