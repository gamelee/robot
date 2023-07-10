package client

type Config struct {
	ID   uint32
	Name string
	Addr string
}

type Option interface {
	Apply(conf *Config)
}

type WithID uint32

func (id WithID) Apply(conf *Config) {
	conf.ID = uint32(id)
}

type WithName string

func (name WithName) Apply(conf *Config) {
	conf.Name = string(name)
}

type WithAddr string

func (addr WithAddr) Apply(conf *Config) {
	conf.Addr = string(addr)
}

func NewConfig(opts ...Option) *Config {
	this := new(Config)
	for _, opt := range opts {
		opt.Apply(this)
	}

	return this
}

func (conf *Config) Apply(opts ...Option) *Config {
	for _, opt := range opts {
		opt.Apply(conf)
	}
	return conf
}
