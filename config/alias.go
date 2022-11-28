package config

// Alias - type for aliasing in config
type Alias[Type any] struct {
	name   string
	entity Type
}

// UnmarshalYAML -
func (a *Alias[Type]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&a.name); err == nil {
		return nil
	}

	var typ Type
	if err := unmarshal(&typ); err != nil {
		return err
	}
	a.entity = typ
	return nil
}

// Name - returns alias name if it exists
func (a *Alias[Type]) Name() string {
	return a.name
}

// Struct - returns substitute struct
func (a *Alias[Type]) Struct() Type {
	return a.entity
}

// SetStruct - set entity. Use in Substitute
func (a *Alias[Type]) SetStruct(entity Type) {
	a.entity = entity
}
