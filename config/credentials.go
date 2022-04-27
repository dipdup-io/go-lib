package config

// Credentials -
type Credentials struct {
	User *UserCredentials `yaml:"user,omitempty" validate:"omitempty"`
}

// UserCredentials -
type UserCredentials struct {
	Name     string `yaml:"name" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}
