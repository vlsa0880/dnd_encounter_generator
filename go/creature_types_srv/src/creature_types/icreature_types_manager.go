package creature_types

type ICreatureTypeManager interface {
	Init() error
	Load() error
	GetData() Types
}
