package packages

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/egonelbre/bundlejs/indenter"
)

func (p *Package) transpileSingle() error {
	data, err := p.root.Dir.ReadFile(p.Files[0])
	if err != nil {
		return err
	}

	p.content.Reset()
	p.content.Write(data)
	return nil
}

var (
	packageHeader = "/* %v */\nvar %v = {};\n(function(Ɛ){\n"
	packageFooter = "\n})(%v);\n"

	fileHeader = "\n/* %v */\n"

	rpPkgName = `([a-zA-Z\/\.\-\__0-9]+)`
	rxPkgName = regexp.MustCompile(`"` + rpPkgName + `"`)

	rxImport  = regexp.MustCompile(`import\s+"` + rpPkgName + `"[ \t;]*\n?`)
	rxImports = regexp.MustCompile(`(?s)import\s*\((?:\s*"` + rpPkgName + `"\s*)+\s*\)[ \t;]*\n?`)
	rxExport  = regexp.MustCompile(`export\s+([a-zA-Z_0-9$]+)\s*[ \t;]`)
)

func pkgvar(pkgname string) string {
	return "ノ" + strings.Replace(pkgname, "/", "ノ", -1)
}

func importStatement(pkgname string) string {
	if DetermineKind(pkgname) == KindPackage {
		name := path.Base(pkgname)
		return fmt.Sprintf("var %s = %s;\n", name, pkgvar(pkgname))
	}
	return ""
}

func (p *Package) exportStatement(name string) string {
	return fmt.Sprintf("Ɛ.%s = %s;", name, name)
}

func (p *Package) importFile(file string) error {
	data, err := p.root.Dir.ReadFile(file)
	if err != nil {
		return fmt.Errorf(`error in read file: %v`, err)
	}

	data = replaceAllSubmatchFunc(rxImport, data, func(packagename []byte) (r []byte) {
		p.Deps = append(p.Deps, string(packagename))
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

	indent := indenter.New(p.content, []byte{'\t'})
	fmt.Fprintf(indent, fileHeader, file)
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
