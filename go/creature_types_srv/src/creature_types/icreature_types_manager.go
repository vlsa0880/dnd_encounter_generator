package creature_types

type CreatureTypeManager interface {
	Init() error
	Load() error
	GetData() Types
}
