package ast

// Miguel kinds
const (
	MiguelKindCreate = "create"
	MiguelKindUpdate = "update"
	MiguelKindDelete = "delete"
)

const (
	valueKindString = iota + 1
	valueKindBytes
	valueKindInt
)

// schemaKeyField is the JSON key that holds the variant discriminator of option and or values.
const schemaKeyField = "schemaKey"
