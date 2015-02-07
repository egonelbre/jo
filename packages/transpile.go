package packages

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/egonelbre/bundlejs/indenter"
)

func (p *Package) transpileSingle() error {
	data, err := ioutil.ReadFile(p.Files[0])
	if err != nil {
		return fmt.Errorf("error in read file: %s", err)
	}

	p.content.Reset()
	p.content.Write(data)
	return nil
}

var (
	packageHeader = "/* %v */\nvar %v = {};\n(function($pkg){"
	packageFooter = "\n\n})(%v);\n"

	fileHeader = "\n\n/* %v */\n"

	pathslash = "ノ"
	pkgname   = `\s*"([a-zA-Z\/\.\-\__0-9]+)"`

	rxPkgName = regexp.MustCompile(pkgname)
	rxImport  = regexp.MustCompile(`import ` + pkgname + `[ \t;]*`)
	rxImports = regexp.MustCompile(`(?s)import\s*\((?:` + pkgname + `)+\s*\)[ \t;]*`)
	rxExport  = regexp.MustCompile(`export\s+([a-zA-Z_0-9$]+)\s*[ \t;]*`)
)

func pkgvar(pkgname string) string {
	return pathslash + strings.Replace(pkgname, "/", pathslash, -1)
}

func importStatement(pkgname string) string {
	if DetermineKind(pkgname) == KindPackage {
		name := path.Base(pkgname)
		return fmt.Sprintf("var %s = %s;\n", name, pkgvar(pkgname))
	}
	return ""
}

func (p *Package) exportStatement(name string) string {
	return fmt.Sprintf("$pkg.%s = %s;", name, name)
}

func (p *Package) importFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf(`error in read file: %v`, err)
	}

	data = replaceAllSubmatchFunc(rxImport, data, func(packagename []byte) (r []byte) {
		return []byte(importStatement(string(packagename)))
	})

	data = rxImports.ReplaceAllFunc(data, func(imports []byte) (r []byte) {
		pkgs := rxPkgName.FindAllSubmatch(imports, -1)
		for _, pkg := range pkgs {
			dep := string(pkg[1])
			p.Deps = append(p.Deps, dep)
			r = append(r, []byte(importStatement(dep))...)
		}
		return r
	})

	data = replaceAllSubmatchFunc(rxExport, data, func(name []byte) (r []byte) {
		return []byte(p.exportStatement(string(name)))
	})

	rel, err := filepath.Rel(p.root.Dir, filename)
	if err != nil {
		rel = filename
	}
	rel = filepath.ToSlash(rel)

	indent := indenter.New(p.content, []byte{'\t'})
	fmt.Fprintf(indent, fileHeader, rel)
	indent.Write(data)
	return nil
}

func (p *Package) transpileModule() error {
	p.content.Reset()

	fmt.Fprintf(p.content, packageHeader, p.Name, pkgvar(p.Name))
	defer fmt.Fprintf(p.content, packageFooter, pkgvar(p.Name))

	for _, filename := range p.Files {
		if err := p.importFile(filename); err != nil {
			if err != nil {
				return fmt.Errorf(`importing "%v" failed: %v`, filename, err)
			}
		}
	}
	return nil
}
