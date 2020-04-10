package archive

import (
	"bytes"
	"github.com/krischerven/gomics/natsort"
	"errors"
	"github.com/gotk3/gotk3/gdk"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Loader interface {
	Load(i int) (*gdk.Pixbuf, error)
	Name(i int) (string, error)
	Len() int
}

// TODO(utkan): check rar support

//var ArchiveExtensions = []string{".zip", ".cbz", ".7z", ".rar", ".tar", ".tgz", ".tbz2", ".cb7", ".cbr", ".cbt"}
var ArchiveExtensions = []string{".zip", ".cbz"}
var ImageExtensions = []string{ // FIXME(utkan): Use gdk_pixbuf_get_formats()
	".jpg", ".jpeg", ".gif", ".png", ".tif", ".bmp", ".pcx", ".xv", ".xpm",
	".xcf", ".tif", ".tga", ".pnm", ".lbm", ".cur", ".ico",
	".jp2", ".j2k", ".jpf", ".jpx", ".jpm",
}

func ExtensionMatch(p string, extensions []string) bool {
	pext := strings.ToLower(filepath.Ext(p))
	for _, ext := range extensions {
		if pext == ext {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func wrap(val, low, mod int) int {
	val %= mod
	if val < low {
		val = mod + val
	}
	return val
}

type stringArray []string

func (p stringArray) Len() int           { return len(p) }
func (p stringArray) Less(i, j int) bool { return strings.ToLower(p[i]) < strings.ToLower(p[j]) }
func (p stringArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func ListArchives(dir string) (anames []string, err error) {
	file, err := os.Open(dir)
	if err != nil {
		return
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return
	}

	if !fi.Mode().IsDir() {
		err = errors.New(dir + " is not a directory!")
		return
	}

	names, err := file.Readdirnames(-1)
	if err != nil {
		return
	}

	anames = make([]string, 0, len(names))
	for _, name := range names {
		if !ExtensionMatch(name, ArchiveExtensions) {
			continue
		}
		anames = append(anames, name)
	}

	sort.Sort(stringArray(anames))

	return
}

func LoadPixbuf(r io.Reader, autorotate bool) (*gdk.Pixbuf, error) {
	w, _ := gdk.PixbufLoaderNew()
	defer w.Close()
	_, err := io.Copy(w, r)
	if err != nil {
		return nil, err
	}

	pixbuf, err := w.GetPixbuf()
	if err != nil {
		return nil, err
	}

	if autorotate == false {
		return pixbuf, nil
	}

	return pixbuf.ApplyEmbeddedOrientation()
}

type File struct {
	*os.File
}

func NewFile(f *os.File) *File {
	return &File{f}
}

func (r *File) Size() (int64, error) {
	fi, err := r.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func (r *File) SetSize(n int64) error {
	return r.Truncate(n)
}

func (r *File) Ext() string {
	ext := filepath.Ext(r.Name())
	if len(ext) <= 1 || ext[0] != '.' {
		return ""
	}

	return ext[1:]
}

type Buffer struct {
	bytes.Buffer
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{*bytes.NewBuffer(data)}
}

func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	return offset, nil
}

func (b *Buffer) SetSize(int64) error {
	return nil
}

func (b *Buffer) Size() (int64, error) {
	return int64(b.Len()), nil
}

func strcmp(a, b string, nat bool) bool {
	if nat {
		return natsort.Less(a, b)
	}
	return a < b
}
