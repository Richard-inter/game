package grpc

import (
	"fmt"

	"github.com/Richard-inter/game/internal/discovery"
)

type ServiceDiscovery interface {
	GetService(serviceName string) (string, error)
	Close() error
}

type DiscoveryManager struct {
	discovery discovery.ServiceDiscovery
}

func NewDiscoveryManager(etcdEndpoints []string) (*DiscoveryManager, error) {
	etcdDiscovery, err := discovery.NewEtcdDiscovery(etcdEndpoints)
	if err != nil {
		return nil, fmt.Errorf("failed to create service discovery: %w", err)
	}

	return &DiscoveryManager{
		discovery: etcdDiscovery,
	}, nil
}

func (dm *DiscoveryManager) GetService(serviceName string) (string, error) {
	return dm.discovery.GetService(serviceName)
}

func (dm *DiscoveryManager) Close() error {
	return dm.discovery.Close()
}
