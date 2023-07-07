package ui

type Config struct {
	Width     int
	Height    int
	CachePath string
	AssetPath string
	Name      string
	Args      []string
	Port      int
}

type Option interface {
	Apply(config *Config) error
}

type WithName string

func (a WithName) Apply(config *Config) error {
	config.Name = string(a)
	return nil
}

type WithSize [2]int

func (s WithSize) Apply(config *Config) error {
	config.Width = s[0]
	config.Height = s[1]
	return nil
}

type WithCachePath string

func (l WithCachePath) Apply(config *Config) error {
	config.CachePath = string(l)
	return nil
}

type WithAssetPath string

func (l WithAssetPath) Apply(config *Config) error {
	config.AssetPath = string(l)
	return nil
}

type WithArgs []string

func (a WithArgs) Apply(config *Config) error {
	config.Args = append(config.Args, a...)
	return nil
}

type WithPort int

func (p WithPort) Apply(config *Config) error {
	config.Port = int(p)
	return nil
}
