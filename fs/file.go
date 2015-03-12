package fs

import (
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type File struct {
	atime   uint64
	mtime   uint64
	name    string
	size    uint64
	content []byte
}

func NewFile(name string) *File {
	return &File{
		atime: uint64(time.Now().Unix()),
		mtime: uint64(time.Now().Unix()),
		name: name,
		size: 0,
	}
}

func (f *File) Allocate(off uint64, size uint64, mode uint32) (code fuse.Status) {
	return fuse.OK
}

func (f *File) String() string {
	return f.name
}

func (f *File) Read(buf []byte, off int64) (fuse.ReadResult, fuse.Status) {
	f.atime = uint64(time.Now().Unix())
	return fuse.ReadResultData(f.content), fuse.OK
}

func (f *File) Write(content []byte, off int64) (uint32, fuse.Status) {
	f.atime = uint64(time.Now().Unix())
	f.mtime = f.atime

	f.content = content
	return uint32(len(content)), fuse.OK
}

func (f *File) Flush() fuse.Status {
	return fuse.OK
}

func (f *File) Fsync(flags int) (code fuse.Status) {
	return fuse.OK
}

func (f *File) Truncate(size uint64) (code fuse.Status) {
	f.content = []byte{}
	return fuse.OK
}

func (f *File) Chmod(mode uint32) fuse.Status {
	return fuse.OK
}

func (f *File) Chown(uid uint32, gid uint32) fuse.Status {
	return fuse.OK
}

func (f *File) GetAttr(out *fuse.Attr) fuse.Status {
	out.Atime = f.atime
	out.Mtime = f.mtime
	out.Mode = fuse.S_IFREG | 0644
	out.Size = uint64(len(f.content))
	return fuse.OK
}

func (f *File) InnerFile() nodefs.File {
	return nil
}

func (f *File) Release() {
}

func (f *File) SetInode(n *nodefs.Inode) {
}

func (f *File) Utimens(a *time.Time, m *time.Time) fuse.Status {
	f.atime = uint64(a.Unix())
	f.mtime = f.atime
	return fuse.OK
}
