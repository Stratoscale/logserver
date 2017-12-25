package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
)


type SshFS struct{
    BaseFS
}

func (f *SshFS) ReadDir(dirname string) ([]os.FileInfo, error) {
    return ioutil.ReadDir(filepath.Join(f.Url.Path, dirname)) }

func (f *SshFS) Lstat(name string) (os.FileInfo, error) {
    return os.Lstat(filepath.Join(f.Url.Path, name)) }

func (f *SshFS) Join(elem ...string) string {
    return filepath.Join(f.Url.Path, filepath.Join(elem...)) }
