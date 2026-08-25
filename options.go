package mimekit

type Options struct {
	MaxHeaderBytes int
	MaxBodyBytes   int
	StrictFold     bool
}

func (o Options) withDefaults() Options {
	if o.MaxHeaderBytes <= 0 {
		o.MaxHeaderBytes = 1 << 20
	}
	if o.MaxBodyBytes <= 0 {
		o.MaxBodyBytes = 32 << 20
	}
	return o
}
