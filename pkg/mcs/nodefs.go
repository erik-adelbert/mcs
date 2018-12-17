/* nodefs.go
erik adelbert - 2018 - erik _ adelbert _ fr
*/

package mcs

/*
import (
	"fmt"
	"golang.org/x/tools/godoc/vfs"
	"os"
	"path"
	"path/filepath"
	"time"
)


type nodeFI struct {
	name string
	file *Node
}

func (fi nodeFI) Name() string {
	if fi.name != "" {
		return fi.name
	}
	return fmt.Sprintf("%p", fi.file)
}

func (fi nodeFI) Size() int64 {
	return 0
}

func (fi nodeFI) ModeTime() time.Time {
	return time.Time{}
}

func (fi nodeFI) Mode() os.FileMode {
	return os.ModeDir | os.ModePerm
}

func (fi nodeFI) IsDir() bool {
	return fi.name == ""
}

func (fi nodeFI) Sys() interface{} {
	return nil
}


func NewNameSpace() vfs.NameSpace {
	ns := vfs.NameSpace{}
	ns.Bind("/", &nodeFS{}, "/", vfs.BindReplace)
	return ns
}

func New(n *Node) vfs.FileSystem {
	return &nodeFS{n}
}

type nodeFS struct {
	root *Node
}

func (fs *nodeFS) Lstat(path string) (os.FileInfo, error) {
	return fs.Stat(path)
}

func (fs *nodeFS) Open(path string) (vfs.ReadSeekCloser, error) {

	node := fs.root
	err := filepath.Walk()
}

func (fs *nodeFS) ReadDir(path string) ([]os.FileInfo, error) {

}

func (fs *nodeFS) RootType(path string) vfs.RootType {
	return ""
}

func (fs *nodeFS) Stat(path string) (os.FileInfo, error) {
	if fs == nil {
		return nil, os.ErrInvalid
	}
	return nil, os.ErrNotExist
}

func (fs *nodeFS) String() string {
	return fmt.Sprintf("node(@%p)", fs.root)
}

func (fs *nodeFS) RootType(path string) vfs.RootType {
	return vfs.RootTypeGoPath
}

*/
