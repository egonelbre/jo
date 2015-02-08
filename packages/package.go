package packages

type Package struct {
	Name  string
	Files []*File

	Deps []string `json:",omitempty"`

	root *Packages
}

func NewPackage(root *Packages, packagename string) (*Package, error) {
	p := &Package{
		Name: packagename,
		root: root,
	}

	var err error
	p.Files, err = p.root.Dir.LoadFiles(p.Name)
	if err != nil {
		return nil, err
	}
	if err := p.Transpile(); err != nil {
		return nil, err
	}
	return p, nil
}
