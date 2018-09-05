package bloom

import (
	"math"
	"os"
	"sync"

	"github.com/spaolacci/murmur3"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/tools/mmap"
)

var pool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1)
	},
}

type Bloom struct {
	opts Options
	m    int64
	k    int
	buf  []byte
	file *mmap.MmapFile
}

var defaultOption = Options{}

func New(n int, FP float64, opts ...Option) (*Bloom, error) {
	opt := defaultOption
	for _, o := range opts {
		o(&opt)
	}
	m, k := EstimateParameters(n, FP)
	b := &Bloom{
		opts: opt,
		m:    m,
		k:    k,
	}
	if b.opts.file != "" && b.opts.maxSize > 0 {
		f, err := mmap.New(b.opts.file, b.opts.maxSize, mmap.CREATE, mmap.RANDOM)
		if err != nil {
			return nil, err
		}
		b.file = f
	} else {
		b.buf = make([]byte, b.m/8+1)
	}
	return b, nil
}

func (b *Bloom) Add(key []byte) error {
	h := hash(key)
	delta := h>>17 | h<<15
	for i := 0; i < b.k; i++ {
		bitPos := h % uint64(b.m)
		if b.file != nil {
			off := int64(bitPos / 8)
			buf := pool.Get().([]byte)
			buf[0] = 1 << (bitPos % 8)
			if err := b.file.OrAt(buf, off); err != nil {
				return err
			}
			buf[0] = 0
			pool.Put(buf)
		} else {
			b.buf[bitPos/8] |= 1 << (bitPos % 8)
		}
		h += delta
	}
	return nil
}

func (b *Bloom) MayContain(key []byte) bool {
	h := hash(key)
	delta := h>>17 | h<<15
	for i := 0; i < b.k; i++ {
		bitPos := h % uint64(b.m)
		if b.file != nil {
			off := int64(bitPos / 8)
			buf, err := b.file.ReadAt(off, 1)
			if err != nil {
				log.Fatal(err)
			}
			if buf[0]&(1<<(bitPos%8)) == 0 {
				return false
			}
		} else {
			if b.buf[bitPos/8]&(1<<(bitPos%8)) == 0 {
				return false
			}
		}
		h += delta
	}
	return true
}

func (b *Bloom) Sync() error {
	if b.file != nil {
		return b.file.Sync()
	}
	return nil
}

func (b *Bloom) Close() error {
	if b.file != nil {
		if err := b.file.Close(); err != nil {
			return err
		}
		b.file = nil
	}
	return nil
}

func (b *Bloom) Delete() error {
	if b.file != nil {
		if err := b.Close(); err != nil {
			return err
		}
		if err := os.Remove(b.opts.file); err != nil {
			return err
		}
		b.file = nil
	}
	return nil
}

func EstimateParameters(n int, p float64) (m int64, k int) {
	m = int64(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	k = int(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	return
}

func hash(key []byte) uint64 {
	hasher := murmur3.New64()
	hasher.Write(key)
	return hasher.Sum64()
}
