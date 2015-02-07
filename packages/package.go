package packages

import (
	"bytes"
	"io"
	"path"
)

type Kind string

const (
	// will be imported at the top level
	// may cause clashes
	KindFile = Kind("file")

	// will be imported at local level
	KindPackage = Kind("package")
)

func DetermineKind(packagename string) Kind {
	switch path.Ext(packagename) {
	case ".js":
		return KindFile
	}
	return KindPackage
}

type Package struct {
	Kind  Kind
	Name  string
	Files []string

	Deps []string `json:",omitempty"`

	root    *Packages
	content *bytes.Buffer
}

func NewPackage(root *Packages, packagename string) (*Package, error) {
	p := &Package{
		Kind:    DetermineKind(packagename),
		Name:    packagename,
		root:    root,
		content: new(bytes.Buffer),
	}

	var err error

	switch p.Kind {
	case KindFile:
		p.Files = append(p.Files, p.Name)
		if err := p.transpileSingle(); err != nil {
			return nil, err
		}
	case KindPackage:
		if p.Files, err = p.root.Dir.List(p.Name); err != nil {
			return nil, err
		}
		if err := p.transpileModule(); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *Package) WriteTo(w io.Writer) error {
	_, err := p.content.WriteTo(w)
	return err
}
