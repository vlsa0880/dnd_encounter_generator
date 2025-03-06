package kafka

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/kafka/consumers"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/kafka/interfaces"
	ikafka "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/kafka/interfaces"
	msghandlers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/kafka/msg_handlers"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
)

type Manager struct {
	ctx            context.Context
	ctxCancel      context.CancelFunc
	settingsLoader settings.SettingsLoader
	consumers      []interfaces.Consumer
}

func New(ctx context.Context, settingsLoader settings.SettingsLoader, creatureTypesDB dbinterfaces.CreatureTypesRepository) (*Manager, error) {
	if settingsLoader == nil {
		return nil, fmt.Errorf("bad settings loader")
	}

	var manager Manager
	manager.settingsLoader = settingsLoader

	manager.ctx, manager.ctxCancel = context.WithCancel(ctx)
	if err := manager.setupDataHandler(creatureTypesDB); err != nil {
		return nil, fmt.Errorf("can't setup data handler: %w", err)
	}
	return &manager, nil
}

func (manager *Manager) Run() error {
	var eg errgroup.Group
	for _, consumer := range manager.consumers {
		eg.Go(consumer.Run)
	}
	return eg.Wait()
}

func (manager *Manager) Stop() error {
	var eg errgroup.Group
	for _, consumer := range manager.consumers {
		eg.Go(consumer.Stop)
	}
	return eg.Wait()
}

func (manager *Manager) setupDataHandler(creatureTypesDB dbinterfaces.CreatureTypesRepository) error {
	if err := manager.isErrors(); err != nil {
		return err
	}

	var handler ikafka.MsgHandler
	var err error
	handler, err = msghandlers.NewCreatureTypeHandler(
		manager.settingsLoader,
		creatureTypesDB,
	)
	if err != nil {
		return err
	}

	consumer, err := consumers.NewGetCreatureType(
		manager.ctx,
		manager.settingsLoader,
		handler,
	)
	if err != nil {
		return err
	}

	manager.consumers = append(manager.consumers, consumer)
	return nil
}

func (manager *Manager) isErrors() error {
	if manager.settingsLoader == nil {
		return fmt.Errorf("bad settings loader")
	}
	return nil
}
