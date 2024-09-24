package config

// Credentials -
type Credentials struct {
	User   *UserCredentials `yaml:"user,omitempty"    validate:"omitempty"`
	ApiKey *ApiKey          `yaml:"api_key,omitempty" validate:"omitempty"`
}

// UserCredentials -
type UserCredentials struct {
	Name     string `yaml:"name" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

type ApiKey struct {
	Header string `yaml:"header" validate:"required"`
	Key    string `yaml:"key"    validate:"required"`
}
