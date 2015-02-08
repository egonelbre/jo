package packages

import (
	"fmt"
	"io"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/egonelbre/jo/indenter"
)

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
	return fmt.Sprintf("var %s = %s;\n", path.Base(pkgname), pkgvar(pkgname))
}

func (p *Package) exportStatement(name string) string {
	return fmt.Sprintf("Ɛ.%s = %s;", name, name)
}

func (p *Package) transpileFile(file *File) error {
	data := file.content.Bytes()

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

	file.transpiled.Reset()
	fmt.Fprintf(&file.transpiled, fileHeader, file.Name)
	file.transpiled.Write(data)
	return nil
}

func (p *Package) Transpile() error {
	for _, file := range p.Files {
		if err := p.transpileFile(file); err != nil {
			if err != nil {
				return fmt.Errorf(`importing "%v" failed: %v`, file.Name, err)
			}
		}
	}
	return nil
}

func (p *Package) WriteTo(w io.Writer) error {
	fmt.Fprintf(w, packageHeader, p.Name, pkgvar(p.Name))
	defer fmt.Fprintf(w, packageFooter, pkgvar(p.Name))

	sort.Sort(byFilename(p.Files))

	indent := indenter.New(w, []byte{'\t'})
	for _, file := range p.Files {
		_, err := file.transpiled.WriteTo(indent)
		if err != nil {
			return err
		}
	}

	return nil
}
