package bloom

type Options struct {
	file    string
	maxSize int
}

type Option func(*Options)

func WithMmap(file string, maxSize int) Option {
	return func(o *Options) {
		o.file = file
		o.maxSize = maxSize
	}
}
