package mmap

import (
	"encoding/binary"
	"errors"
	"os"
	"sync"
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

const (
	RANDOM = iota
	SEQUENTIAL
	WILLNEED
)

type MmapFile struct {
	sync.RWMutex
	*os.File
	buf       []byte
	maxSize   int64
	fileSize  int64
	maxBufOff int64
}

func New(path string, size int64, mode, advise int) (*MmapFile, error) {
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
	switch advise {
	case RANDOM:
		advise = syscall.MADV_RANDOM
	case SEQUENTIAL:
		advise = syscall.MADV_SEQUENTIAL
	case WILLNEED:
		advise = syscall.MADV_WILLNEED
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
	buf, err := syscall.Mmap(int(file.Fd()), 0, int(size),
		prot, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	if _, _, err := syscall.Syscall(syscall.SYS_MADVISE,
		uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)),
		uintptr(advise)); err != 0 {

		return nil, err
	}
	return &MmapFile{
		File:      file,
		buf:       buf,
		maxSize:   size,
		fileSize:  info.Size(),
		maxBufOff: info.Size(),
	}, nil
}

func (m *MmapFile) Size() int64 {
	return m.fileSize
}

func (m *MmapFile) WriteAt(b []byte, off int64) error {
	if err := m.tryGrow(off, int64(len(b))); err != nil {
		return err
	}
	copy(m.buf[off:off+int64(len(b))], b)
	return nil
}

func (m *MmapFile) OrAt(b []byte, off int64) error {
	if err := m.tryGrow(off, int64(len(b))); err != nil {
		return err
	}
	for i := range b {
		m.buf[off+int64(i)] |= b[int64(i)]
	}
	return nil
}

func (m *MmapFile) ReadAt(off, n int64) ([]byte, error) {
	if err := m.tryGrow(off, n); err != nil {
		return nil, err
	}
	return m.buf[off : off+n], nil
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

func (m *MmapFile) tryGrow(off, n int64) error {
	if off+n > m.maxSize {
		return errors.New("oversize")
	}
	m.Lock()
	defer m.Unlock()
	for off+n >= m.fileSize {
		err := syscall.Ftruncate(int(m.File.Fd()), m.fileSize+ALLOCSIZE)
		if err != nil {
			log.Error(err)
			return err
		}
		m.fileSize += ALLOCSIZE
	}
	if off+n > m.maxBufOff {
		m.maxBufOff = off + int64(n)
	}
	return nil
}

func (m *MmapFile) Uint64At(off int64) (uint64, error) {
	if err := m.tryGrow(off, 8); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(m.buf[off : off+8]), nil
}

func (m *MmapFile) Uint32At(off int64) (uint32, error) {
	if err := m.tryGrow(off, 4); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(m.buf[off : off+8]), nil
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
