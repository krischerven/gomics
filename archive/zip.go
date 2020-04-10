package archive

import (
	"archive/zip"
	"errors"
	"github.com/gotk3/gotk3/gdk"
	"path"
	"sort"
)

type Zip struct {
	files  []*zip.File // File elements sorted by their Names
	reader *zip.ReadCloser
	name   string // Name of the Zip file
}

type zipfile []*zip.File

func (p zipfile) Len() int           { return len(p) }
func (p zipfile) Less(i, j int) bool { return strcmp(p[i].Name, p[j].Name, true) }
func (p zipfile) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

/* Reads filenames from a given zip archive, and sorts them */
func NewZip(name string) (*Zip, error) {
	var err error

	ar := new(Zip)

	ar.name = path.Base(name)
	ar.files = make([]*zip.File, 0, MaxArchiveEntries)
	ar.reader, err = zip.OpenReader(name)
	if err != nil {
		return nil, err
	}

	for _, f := range ar.reader.File {
		if ExtensionMatch(f.Name, ImageExtensions) == false {
			continue
		}
		ar.files = append(ar.files, f)
	}

	if len(ar.files) == 0 {
		return nil, errors.New(ar.name + ": no images in the zip file")
	}

	sort.Sort(zipfile(ar.files))

	return ar, nil
}

func (ar *Zip) checkbounds(i int) error {
	if i < 0 || i >= len(ar.files) {
		return ErrBounds
	}
	return nil
}

func (ar *Zip) Load(i int, autorotate bool) (*gdk.Pixbuf, error) {
	if err := ar.checkbounds(i); err != nil {
		return nil, err
	}

	f, err := ar.files[i].Open()
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return LoadPixbuf(f, autorotate)
}

func (ar *Zip) Name(i int) (string, error) {
	if err := ar.checkbounds(i); err != nil {
		return "", err
	}

	return ar.files[i].Name, nil
}

func (ar *Zip) Len() int {
	return len(ar.files)
}

func (ar *Zip) Close() error {
	return ar.reader.Close()
}
