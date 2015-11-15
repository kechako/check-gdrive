package local

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kechako/check-gdrive/md5sum"
)

type File struct {
	Path string
	os.FileInfo
}

func NewFile(path string) (*File, error) {
	var err error

	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	file := &File{
		Path:     path,
		FileInfo: info,
	}

	return file, nil
}

func (f *File) GetFiles() ([]*File, error) {
	if !f.IsDir() {
		return nil, fmt.Errorf("%s is not a directory.", f.Path)
	}

	fis, err := ioutil.ReadDir(f.Path)
	if err != nil {
		return nil, err
	}

	files := make([]*File, 0, 20)
	for _, fi := range fis {
		file := &File{
			Path:     f.Join(fi.Name()),
			FileInfo: fi,
		}
		files = append(files, file)
	}

	return files, nil
}

func (f *File) Md5Checksum() (string, error) {
	return md5sum.SumFile(f.Path)
}

func (f *File) Join(name string) string {
	return filepath.Join(f.Path, name)
}
