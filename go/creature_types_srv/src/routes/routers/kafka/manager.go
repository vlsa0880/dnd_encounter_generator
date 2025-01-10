package kafka

import (
	"context"
	"fmt"

	data_handlers_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/kafka/consumers"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/kafka/interfaces"
	msghandlers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/kafka/msg_handlers"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type Manager struct {
	ctx            context.Context
	ctxCancel      context.CancelFunc
	settingsLoader settings.SettingsLoader
	consumers      []interfaces.Consumer
}

func New(ctx context.Context, settingsLoader settings.SettingsLoader) *Manager {
	if settingsLoader == nil {
		panic("bad settings loader")
	}

	manager := &Manager{
		settingsLoader: settingsLoader,
	}

	manager.ctx, manager.ctxCancel = context.WithCancel(ctx)
	return manager
}

func (manager *Manager) Run() {
	for _, consumer := range manager.consumers {
		go consumer.Run()
	}
}

func (manager *Manager) Stop() {
	for _, consumer := range manager.consumers {
		go consumer.Stop()
	}
}

func (manager *Manager) SetupDataHandler(data_handler data_handlers_interfaces.DataHandler) error {
	if err := manager.isErrors(); err != nil {
		return err
	}
	handlerGetCreatureType := msghandlers.NewCreatureTypeHandler(
		manager.settingsLoader,
		data_handler,
	)
	var consumerGetCreatureTypes interfaces.Consumer = consumers.NewGetCreatureType(
		manager.ctx,
		manager.settingsLoader,
	)
	if err := consumerGetCreatureTypes.SetupMsgHandler(handlerGetCreatureType); err != nil {
		return err
	}
	manager.consumers = append(manager.consumers, consumerGetCreatureTypes)
	return nil
}

func (manager *Manager) isErrors() error {
	if manager.settingsLoader == nil {
		return fmt.Errorf("bad settings loader")
	}
	return nil
}
