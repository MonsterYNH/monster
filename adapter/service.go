package adapter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"monster/config"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

// ServiceAdapter 服务适配器
type ServiceAdapter interface {
	// 启动服务
	Run(*config.Config) error
	// 关闭服务
	Close() error
	// GetName get service name
	GetName() string
}

// ServiceAdapterPool service pool
type ServiceAdapterPool struct {
	pool  map[string]ServiceAdapter
	lock  *sync.Mutex
	group *errgroup.Group
	done  context.Context
}

// NewService create new service
func NewService() *ServiceAdapterPool {
	group, done := errgroup.WithContext(context.Background())
	return &ServiceAdapterPool{
		pool:  make(map[string]ServiceAdapter),
		lock:  &sync.Mutex{},
		done:  done,
		group: group,
	}
}

// Register register a service
func (adapter *ServiceAdapterPool) Serve(services ...ServiceAdapter) error {
	adapter.lock.Lock()
	defer adapter.lock.Unlock()

	for _, service := range services {
		name := service.GetName()
		if _, exist := adapter.pool[name]; exist {
			return fmt.Errorf("ERROR: service %s is already exist", name)
		}
		config, err := config.GetConfigEntry(name)
		if err != nil {
			return err
		}
		adapter.pool[name] = service
		adapter.group.Go(convertAdapter(service, config))
		log.Println("INFO: register service named", name)
	}
	return adapter.group.Wait()
}

// StopRegistedService stop all the service
func (adapter *ServiceAdapterPool) StopRegistedService() error {
	adapter.lock.Lock()
	defer adapter.lock.Unlock()

	if len(adapter.pool) == 0 {
		return errors.New("no registed service to stop")
	}
	for name, service := range adapter.pool {
		if err := service.Close(); err != nil {
			return err
		}
		log.Println(fmt.Sprintf("INFO: service %s close success", name))
	}
	return nil
}

func convertAdapter(s ServiceAdapter, c *config.Config) func() error {
	return func() error {
		log.Println(fmt.Sprintf("INFO: service %s start success", s.GetName()))
		splitEndpoint := strings.Split(c.Endpoint, ":")
		if len(splitEndpoint) != 2 {
			return fmt.Errorf("format endpoint: %s", c.Endpoint)
		}
		// 启动服务
		return s.Run(c)
	}
}
