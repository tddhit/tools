package bloom

type Options struct {
	file    string
	maxSize int64
}

type Option func(*Options)

func WithMmap(file string, maxSize int64) Option {
	return func(o *Options) {
		o.file = file
		o.maxSize = maxSize
	}
}
