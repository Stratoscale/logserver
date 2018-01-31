package filesystem

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/go-systemd/sdjournal"
	"github.com/kr/fs"
)

const (
	tmpDir = "/tmp"
	sep    = string(os.PathSeparator)
)

var log = logrus.StandardLogger().WithField("pkg", "journal")

type journal struct {
	inner   FileSystem
	dirName string
	// copyDir holds a copy of a journal directory in case the inner filesystem is not a local filesystem
	// the copy is deleted when the filesystem is closed
	copyDir string
	sync.Mutex
}

func (j *journal) Join(elem ...string) string {
	return j.inner.Join(elem...)
}

// WrapTar wraps a filesystem, and show a journalctl directory on journalDirName
// as a log file and not as a directory.
func WrapJournal(inner FileSystem, journalDirName string) FileSystem {
	return &journal{
		inner:   inner,
		dirName: journalDirName,
	}
}

func (j *journal) ReadDir(dirname string) ([]os.FileInfo, error) {
	files, err := j.inner.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	for i, f := range files {
		path := j.Join(dirname, f.Name())
		files[i] = j.lStat(path, f)
	}
	return files, nil
}

func (j *journal) Lstat(name string) (os.FileInfo, error) {
	stat, err := j.inner.Lstat(name)
	if err != nil {
		return nil, err
	}
	return j.lStat(name, stat), nil
}

func (j *journal) Open(name string) (File, error) {
	if j.isJournalDir(name) {
		path, err := j.cpJournalDir(name)
		if err != nil {
			return nil, err
		}
		r, err := sdjournal.NewJournalReader(sdjournal.JournalReaderConfig{Path: path})
		if err != nil {
			return nil, err
		}
		log.Debugf("Serving journal from %s", path)
		return &journalFile{JournalReader: r}, nil
	}
	return j.inner.Open(name)
}

// lStat fixes internal stat information to show journal dir as a journal file
func (j *journal) lStat(path string, stat os.FileInfo) os.FileInfo {
	if j.isJournalDir(path) {
		return file{
			isDir: false,
			name:  stat.Name(),
			size:  stat.Size(),
			time:  stat.ModTime(),
		}
	}
	return stat
}

// cpJournalDir copies the journal to a local filesystem so it could be used.
func (j *journal) cpJournalDir(path string) (string, error) {
	// for local file system, we don't need to copy the journal
	if local, ok := j.inner.(*Local); ok {
		path = j.Join(local.basePath, path)
		log.Debugf("Journal from local path: %s", path)
		return path, nil
	}

	j.Lock()
	defer j.Unlock()
	if j.copyDir != "" {
		log.Debugf("Journal from copy directory: %s", j.copyDir)
		return j.copyDir, nil
	}

	copyDir, err := ioutil.TempDir(tmpDir, "journal-")
	if err != nil {
		return "", fmt.Errorf("create temp directory: %s", err)
	}

	log.Infof("Copying journal to %s", copyDir)

	// walk the journal filesystem
	for w := fs.WalkFS(path, j.inner); w.Step(); {
		if err := w.Err(); err != nil {
			return "", fmt.Errorf("walk dir: %s", err)
		}
		if w.Stat().IsDir() {
			continue
		}
		// copy file
		remotePath := w.Path()
		localPath := filepath.Join(copyDir, remotePath[len(path):])
		err := j.copyFile(remotePath, localPath)
		if err != nil {
			return "", fmt.Errorf("copy file: %s", err)
		}
	}
	j.copyDir = copyDir
	return copyDir, nil
}

func (j *journal) isJournalDir(path string) bool {
	return strings.Trim(path, sep) == strings.Trim(j.dirName, sep)
}

// copyFile copies a file from a given FileSystem to a local filesystem
func (j *journal) copyFile(remotePath, localPath string) error {

	// open remote file for reading
	r, err := j.Open(remotePath)
	if err != nil {
		return fmt.Errorf("open file %s for reading: %s", remotePath, err)
	}
	defer r.Close()

	// create file directory if not exists, omit the error
	os.MkdirAll(filepath.Dir(localPath), 0755)

	// open local file for writing
	w, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("open local file %s for writing: %s", localPath, err)
	}
	defer w.Close()

	// copy file content
	n, err := io.Copy(w, r)
	if err != nil {
		return err
	}
	log.Debugf("Copied %s -> %s (%d B)", remotePath, localPath, n)
	return err
}

func (j *journal) Close() error {
	if j.copyDir != "" {
		os.RemoveAll(j.copyDir)
	}
	j.copyDir = ""
	return j.inner.Close()
}

// journalFile a file like object for a journal reader.
// The Seek interface is implementing only seek to start
type journalFile struct {
	*sdjournal.JournalReader
}

func (j *journalFile) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		return 0, j.Rewind()
	}
	return 0, nil
}
