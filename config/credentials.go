package config

// Credentials -
type Credentials struct {
	User   *UserCredentials `validate:"omitempty" yaml:"user,omitempty"`
	ApiKey *ApiKey          `validate:"omitempty" yaml:"api_key,omitempty"`
}

// UserCredentials -
type UserCredentials struct {
	Name     string `validate:"required" yaml:"name"`
	Password string `validate:"required" yaml:"password"` //nolint:gosec
}

type ApiKey struct {
	Header string `validate:"required" yaml:"header"`
	Key    string `validate:"required" yaml:"key"`
}
