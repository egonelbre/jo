package packages

import (
	"io/ioutil"
	"path"
	"path/filepath"
)

type Dir string

func (d Dir) ReadFile(name string) ([]byte, error) {
	filename := filepath.FromSlash(name)
	return ioutil.ReadFile(filepath.Join(string(d), filename))
}

func (d Dir) List(name string) ([]string, error) {
	dirname := filepath.Join(string(d), filepath.FromSlash(name))
	entries, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	fs := []string{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".js" {
			continue
		}

		fs = append(fs, path.Join(name, e.Name()))
	}
	return fs, nil
}
