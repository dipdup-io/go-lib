package config

// Configurable -
type Configurable interface {
	Substitute() error
}
