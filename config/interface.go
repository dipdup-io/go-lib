package config

// Configurable -
type Configurable interface {
	Validate() error
	Substitute() error
}
