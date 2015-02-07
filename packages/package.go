package packages

import (
	"bytes"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
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

	switch p.Kind {
	case KindFile:
		if err := p.addFile(filepath.FromSlash(p.Name)); err != nil {
			return nil, err
		}
		if err := p.transpileSingle(); err != nil {
			return nil, err
		}
	case KindPackage:
		if err := p.addFolder(filepath.FromSlash(p.Name)); err != nil {
			return nil, err
		}
		if err := p.transpileModule(); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *Package) addFile(filename string) error {
	if filepath.IsAbs(filename) {
		p.Files = append(p.Files, filename)
		return nil
	}
	abs, err := filepath.Abs(filepath.Join(p.root.Dir, filename))
	if err != nil {
		return err
	}
	p.Files = append(p.Files, abs)
	return nil
}

func (p *Package) addFolder(foldername string) error {
	dirname := filepath.Join(p.root.Dir, foldername)
	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		if filepath.Ext(info.Name()) == ".js" {
			err := p.addFile(filepath.Join(dirname, info.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Package) WriteTo(w io.Writer) error {
	_, err := p.content.WriteTo(w)
	return err
}
