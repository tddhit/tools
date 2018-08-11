package mmap

import (
	"encoding/binary"
	"errors"
	"os"
	"syscall"
	"unsafe"

	"github.com/tddhit/tools/log"
)

const (
	RDONLY = iota
	CREATE
	APPEND

	ALLOCSIZE = 1 << 24 // 16M
)

type MmapFile struct {
	*os.File
	buf       []byte
	fileSize  int64
	maxBufOff int64
}

func New(path string, size, mode int) (*MmapFile, error) {
	var (
		flag      int
		prot      int
		exclusive bool = true
	)
	switch mode {
	case RDONLY:
		flag = os.O_RDONLY
		prot = syscall.PROT_READ
		exclusive = false
	case CREATE:
		flag = os.O_CREATE | os.O_RDWR | os.O_TRUNC
		prot = syscall.PROT_READ | syscall.PROT_WRITE
	case APPEND:
		flag = os.O_CREATE | os.O_RDWR | os.O_APPEND
		prot = syscall.PROT_READ | syscall.PROT_WRITE
	}
	file, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if err = flock(int(file.Fd()), exclusive); err != nil {
		return nil, err
	}
	buf, err := syscall.Mmap(int(file.Fd()), 0, size,
		prot, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	return &MmapFile{
		File:      file,
		buf:       buf,
		fileSize:  info.Size(),
		maxBufOff: info.Size(),
	}, nil
}

func (m *MmapFile) Size() int64 {
	return m.fileSize
}

func (m *MmapFile) WriteAt(b []byte, off int64) error {
	if err := m.tryGrow(off, len(b)); err != nil {
		return err
	}
	copy(m.buf[off:off+int64(len(b))], b)
	return nil
}

func (m *MmapFile) PutUint32At(off int64, v uint32) error {
	if err := m.tryGrow(off, 4); err != nil {
		return err
	}
	binary.LittleEndian.PutUint32(m.buf[off:off+4], v)
	return nil
}

func (m *MmapFile) PutUint64At(off int64, v uint64) error {
	if err := m.tryGrow(off, 8); err != nil {
		return err
	}
	binary.LittleEndian.PutUint64(m.buf[off:off+8], v)
	return nil
}

func (m *MmapFile) tryGrow(off int64, n int) error {
	if off+int64(n) >= m.fileSize {
		err := syscall.Ftruncate(int(m.File.Fd()), m.fileSize+ALLOCSIZE)
		if err != nil {
			log.Error(err)
			return err
		}
		m.fileSize += ALLOCSIZE
	}
	if off+int64(n) > m.maxBufOff {
		m.maxBufOff = off + int64(n)
	}
	return nil
}

func (m *MmapFile) Uint64At(off int64) uint64 {
	return binary.LittleEndian.Uint64(m.buf[off : off+8])
}

func (m *MmapFile) Uint32At(off int64) uint32 {
	return binary.LittleEndian.Uint32(m.buf[off : off+8])
}

func (m *MmapFile) Buf(off int64) unsafe.Pointer {
	return unsafe.Pointer(&m.buf[off])
}

func (m *MmapFile) Close() error {
	if m.buf == nil || m.File == nil {
		return errors.New("buf/file is nil.")
	}
	if err := syscall.Ftruncate(int(m.File.Fd()), m.maxBufOff); err != nil {
		log.Error(err)
		return err
	}
	if err := syscall.Munmap(m.buf); err != nil {
		log.Error(err)
		return err
	}
	if err := m.File.Sync(); err != nil {
		log.Error(err)
		return err
	}
	if err := funlock(int(m.File.Fd())); err != nil {
		log.Error(err)
		return err
	}
	if err := m.File.Close(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func flock(fd int, exclusive bool) error {
	flag := syscall.LOCK_SH
	if exclusive {
		flag = syscall.LOCK_EX
	}
	return syscall.Flock(fd, flag|syscall.LOCK_NB)
}

func funlock(fd int) error {
	return syscall.Flock(fd, syscall.LOCK_UN)
}
