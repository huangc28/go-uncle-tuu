package deps

import (
	"huangc28/go-ios-iap-vendor/internal/app/inventory"
	"huangc28/go-ios-iap-vendor/internal/app/users"
	"sync"

	"github.com/golobby/container"
	cinternal "github.com/golobby/container/pkg/container"
)

type DepContainer struct {
	Container cinternal.Container
}

type DepRegistrar func() error
type ServiceProvider func(cinternal.Container) DepRegistrar

var (
	_depContainer *DepContainer

	once sync.Once
)

func Get() *DepContainer {
	once.Do(func() {
		_depContainer = &DepContainer{
			Container: container.NewContainer(),
		}
	})

	return _depContainer
}

func (dep *DepContainer) Run() error {
	depRegistrars := []DepRegistrar{
		users.UserDaoServiceProvider(dep.Container),
		inventory.InventoryDaoServiceProvider(dep.Container),
	}

	for _, depRegistrar := range depRegistrars {
		if err := depRegistrar(); err != nil {
			return err
		}
	}

	return nil
}
