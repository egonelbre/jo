package packages

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Dir string

type File struct {
	Name     string // it's the relative path from Dir
	Filename string // it's the actual filepath
	ModTime  time.Time

	content    bytes.Buffer
	transpiled bytes.Buffer
}

type byFilename []*File

func (a byFilename) Len() int           { return len(a) }
func (a byFilename) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFilename) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (f *File) Load() error {
	stat, err := os.Stat(f.Filename)
	if err != nil {
		return err
	}

	f.ModTime = stat.ModTime()
	data, err := ioutil.ReadFile(f.Filename)
	if err != nil {
		return err
	}

	f.content.Reset()
	f.content.Write(data)
	return nil
}

func (d Dir) LoadFile(name string) (*File, error) {
	filename := filepath.Join(string(d), filepath.FromSlash(name))
	f := &File{Name: name, Filename: filename}
	return f, f.Load()
}

func (d Dir) LoadFiles(name string) ([]*File, error) {
	dirname := filepath.Join(string(d), filepath.FromSlash(name))

	entries, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	var fs []*File
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".js" {
			continue
		}

		f, err := d.LoadFile(path.Join(name, e.Name()))
		if err != nil {
			return fs, err
		}

		fs = append(fs, f)
	}

	return fs, nil
}
