package packages

import (
	"fmt"
	"io"
)

type Packages struct {
	Dir  Dir
	List map[string]*Package

	Order []string
}

func New(dir string) *Packages {
	return &Packages{
		Dir:  Dir(dir),
		List: make(map[string]*Package),
	}
}

func (ps *Packages) Load(pkgs ...string) error {
	var current string
	unloaded := pkgs
	for len(unloaded) > 0 {
		current, unloaded = unloaded[len(unloaded)-1], unloaded[:len(unloaded)-1]
		if _, loaded := ps.List[current]; loaded {
			continue
		}

		p, err := NewPackage(ps, current)
		if err != nil {
			return err
		}

		ps.List[current] = p
		unloaded = append(unloaded, p.Deps...)
	}
	return nil
}

func (ps *Packages) WriteTo(w io.Writer) error {
	for _, pkgname := range ps.Order {
		pkg := ps.List[pkgname]
		if err := pkg.WriteTo(w); err != nil {
			return err
		}
		fmt.Fprintln(w)
	}
	return nil
}
