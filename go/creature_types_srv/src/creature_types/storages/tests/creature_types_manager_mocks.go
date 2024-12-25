package creature_types_tests

import creature_types "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"

type CreatureTypeManagerMockValid struct {
}

func (mock *CreatureTypeManagerMockValid) Init() error {
	return nil
}

func (mock *CreatureTypeManagerMockValid) Load() error {
	return nil
}

func (mock *CreatureTypeManagerMockValid) GetData() creature_types.Types {
	return creature_types.Types{
		creature_types.Type{Name: "тест1"},
		creature_types.Type{Name: "тест2"},
	}
}
