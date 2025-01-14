package creature_types_tests

import creature_types "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"

type CreatureTypeManagerMockValid struct {
	types creature_types.Data
}

func (mock *CreatureTypeManagerMockValid) Init() error {
	mock.setDefaultData()
	return nil
}

func (mock *CreatureTypeManagerMockValid) GetData() (*creature_types.Data, error) {
	return &mock.types, nil
}

func (mock *CreatureTypeManagerMockValid) setDefaultData() {
	mock.types = creature_types.Data{
		Types: []creature_types.Type{
			{Name: "тест1"},
			{Name: "тест2"},
		},
	}
}
