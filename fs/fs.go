package fs

import (
	"sync"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

type Fs struct {
	pathfs.FileSystem
	Files     map[string]*File
	filesLock sync.Mutex
}

func (f *Fs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	f.filesLock.Lock()
	defer f.filesLock.Unlock()
	if name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	fsFile, ok := f.Files[name]
	if ok {
		return &fuse.Attr{
			Atime: fsFile.atime,
			Mtime: fsFile.mtime,
			Mode:  fuse.S_IFREG | 0644, Size: uint64(len(fsFile.content)),
		}, fuse.OK
	}

	return nil, fuse.ENOENT
}

func (f *Fs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	f.filesLock.Lock()
	defer f.filesLock.Unlock()
	if name == "" {
		c = []fuse.DirEntry{}
		for _, fsFile := range f.Files {
			c = append(c, fuse.DirEntry{Name: fsFile.name, Mode: fuse.S_IFREG})
		}

		return c, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (f *Fs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	f.filesLock.Lock()
	defer f.filesLock.Unlock()
	fsFile, ok := f.Files[name]
	if !ok {
		fsFile = NewFile(name)
		f.Files[name] = fsFile
	}
	return fsFile, fuse.OK
}

func (f *Fs) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	f.filesLock.Lock()
	defer f.filesLock.Unlock()
	fsFile := NewFile(name)
	f.Files[name] = fsFile

	return fsFile, fuse.OK
}

func (f *Fs) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	f.filesLock.Lock()
	defer f.filesLock.Unlock()
	delete(f.Files, name)
	return fuse.OK
}

func (f *Fs) Rename(oldPath string, newPath string, context *fuse.Context) (codee fuse.Status) {
	f.filesLock.Lock()
	defer f.filesLock.Unlock()
	if oldFile, ok := f.Files[oldPath]; ok {
		oldFile.name = newPath
		return fuse.OK
	}
	return fuse.ENOENT
}
